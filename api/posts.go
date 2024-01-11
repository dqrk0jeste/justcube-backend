package api

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"net/http"
	"slices"
	"strconv"

	database "github.com/dqrk0jeste/letscube-backend/database/sqlc"
	"github.com/dqrk0jeste/letscube-backend/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CreatePostRequest struct {
	ImageContent []*multipart.FileHeader `form:"image_content[]" binding:"required,max=5"`
	TextContent  string                  `form:"text_content" binding:"required,max=500"`
}

var SupportedImageTypes = []string{
	"image/jpeg",
	"image/png",
}

func (server *Server) createPost(context *gin.Context) {
	var req CreatePostRequest
	if err := context.ShouldBind(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authorizationPayload := context.MustGet("authorization_payload").(*token.Payload)

	id, err := uuid.NewRandom()
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	files := req.ImageContent

	for _, image := range files {
		if !slices.Contains(SupportedImageTypes, image.Header.Get("Content-Type")) {
			context.AbortWithStatus(http.StatusUnsupportedMediaType)
			return
		}
	}

	for index, image := range files {
		imageNameToSave := id.String() + "_" + strconv.Itoa(index) + ".jpg"

		_, err := server.s3Controller.Upload(context, image, imageNameToSave)
		if err != nil {
			context.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	arg := database.CreatePostParams{
		ID:          id,
		TextContent: req.TextContent,
		ImageCount:  int32(len(files)),
		UserID:      authorizationPayload.UserID,
	}

	post, err := server.database.CreatePost(context, arg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusCreated, post)
}

type DeletePostRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func (server *Server) deletePost(context *gin.Context) {
	var req DeletePostRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authorizationPayload := context.MustGet("authorization_payload").(*token.Payload)

	id, err := uuid.Parse(req.ID)
	if err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	post, err := server.database.GetPostById(context, id)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if post.UserID != authorizationPayload.UserID {
		context.Status(http.StatusForbidden)
		return
	}

	err = server.database.DeletePost(context, id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for i := 0; i < int(post.ImageCount); i++ {
		err := server.s3Controller.Delete(context, id.String()+"_"+strconv.Itoa(i)+".jpg")
		if err != nil {
			fmt.Printf("there has been an error deleting image number %d\n", i)
		}
	}

	context.Status(http.StatusOK)
}

type GetPostByIdRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// in the future, we will have public and private accounts so we will have to check if user that requested the post follows user whose post it is
func (server *Server) getPostById(context *gin.Context) {

	var req GetPostByIdRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	post, err := server.database.GetPostById(context, id)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, post)
}

type GetPostsByUserRequest struct {
	UserID   string `form:"user_id" binding:"required,uuid"`
	Page     int32  `form:"page_number" binding:"required"`
	PageSize int32  `form:"page_size" binding:"required,max=20"`
}

func (server *Server) getPostsByUser(context *gin.Context) {
	var req GetPostsByUserRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(req.UserID)
	if err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := database.GetPostsByUserParams{
		UserID: id,
		Offset: (req.Page - 1) * req.PageSize,
		Limit:  req.PageSize,
	}

	posts, err := server.database.GetPostsByUser(context, arg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, posts)
}

type PostCommentRequest struct {
	PostID  string `json:"post_id" binding:"required,uuid"`
	Content string `json:"content" binding:"required,max=200"`
}

func (server *Server) sendComment(context *gin.Context) {
	var req PostCommentRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	postID, err := uuid.Parse(req.PostID)
	if err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authorizationPayload := context.MustGet("authorization_payload").(*token.Payload)
	userID := authorizationPayload.UserID

	id, err := uuid.NewRandom()
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := database.SendCommentParams{
		ID:      id,
		UserID:  userID,
		PostID:  postID,
		Content: req.Content,
	}

	comment, err := server.database.SendComment(context, arg)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "foreign_key_violation" {
				context.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusCreated, comment)
}

type DeleteCommentRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func (server *Server) deleteComment(context *gin.Context) {
	var req DeleteCommentRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authorizationPayload := context.MustGet("authorization_payload").(*token.Payload)
	userID := authorizationPayload.UserID

	comment, err := server.database.GetCommentById(context, id)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if userID != comment.UserID {
		context.Status(http.StatusForbidden)
		return
	}

	err = server.database.DeleteComment(context, id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.Status(http.StatusOK)
}
