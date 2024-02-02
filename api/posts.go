package api

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"net/http"
	"slices"
	"strconv"
	"sync"

	database "github.com/dqrk0jeste/letscube-backend/database/sqlc"
	"github.com/dqrk0jeste/letscube-backend/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CreatePostRequest struct {
	ImageContent []*multipart.FileHeader `form:"image_content[]" binding:"max=5"`
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

	for _, file := range files {
		if !slices.Contains(SupportedImageTypes, file.Header.Get("Content-Type")) {
			context.AbortWithStatus(http.StatusUnsupportedMediaType)
			return
		}
	}

	ch := make(chan error)
	wg := sync.WaitGroup{}
	errorProcessingImagesCounter := 0

	for index, image := range files {
		wg.Add(1)
		go func(image *multipart.FileHeader, index int) {
			defer wg.Done()
			imageNameToSave := id.String() + "_" + strconv.Itoa(index) + ".jpg"

			_, err := server.s3Controller.Upload(context, image, imageNameToSave)
			ch <- err
		}(image, index)
	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	for err := range ch {
		if err != nil {
			errorProcessingImagesCounter++
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

	context.JSON(http.StatusCreated, gin.H{
		"post":             post,
		"number_of_errors": errorProcessingImagesCounter,
	})
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
		go func(i int) {
			err := server.s3Controller.Delete(context, id.String()+"_"+strconv.Itoa(i)+".jpg")
			if err != nil {
				fmt.Println("there has been an error deleting image" + id.String() + "_" + strconv.Itoa(i) + ".jpg")
			}
		}(i)
	}

	context.Status(http.StatusOK)
}

type GetPostByIdRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// in the future, we will have public and private accounts so we will have to check if user that requested the post follows the user whose post it is
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

	context.JSON(http.StatusOK, post.MakeResponse())
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

	res := make([]database.PostResponse, 0)

	for _, post := range posts {
		res = append(res, post.MakeResponse())
	}

	context.JSON(http.StatusOK, res)
}

type GetFeedRequest struct {
	Page     int32 `form:"page_number" binding:"required"`
	PageSize int32 `form:"page_size" binding:"required,max=20"`
}

func (server *Server) getFeed(context *gin.Context) {
	var req GetFeedRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authorizationPayload := context.MustGet("authorization_payload").(*token.Payload)

	arg := database.GetFeedParams{
		UserID: authorizationPayload.UserID,
		Offset: (req.Page - 1) * req.PageSize,
		Limit:  req.PageSize,
	}

	posts, err := server.database.GetFeed(context, arg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := make([]database.PostResponse, 0)

	for _, post := range posts {
		res = append(res, post.MakeResponse())
	}

	context.JSON(http.StatusOK, res)
}

func (server *Server) getGuestFeed(context *gin.Context) {
	var req GetFeedRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := database.GetGuestFeedParams{
		Offset: (req.Page - 1) * req.PageSize,
		Limit:  req.PageSize,
	}

	posts, err := server.database.GetGuestFeed(context, arg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := make([]database.PostResponse, 0)

	for _, post := range posts {
		res = append(res, post.MakeResponse())
	}

	context.JSON(http.StatusOK, res)
}

type PostCommentRequest struct {
	PostID  string `json:"post_id" binding:"required,uuid"`
	Content string `json:"content" binding:"required,max=200"`
}

func (server *Server) postComment(context *gin.Context) {
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

	arg := database.PostCommentParams{
		ID:      id,
		UserID:  userID,
		PostID:  postID,
		Content: req.Content,
	}

	comment, err := server.database.PostComment(context, arg)
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
	PostID string `uri:"id" binding:"required,uuid"`
}

func (server *Server) deleteComment(context *gin.Context) {
	var req DeleteCommentRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(req.PostID)
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

type GetCommentsRequest struct {
	PostID   string `form:"post_id" binding:"required,uuid"`
	Page     int32  `form:"page_number" binding:"required,min=1"`
	PageSize int32  `form:"page_size" binding:"required,min=1,max=20"`
}

func (server *Server) getComments(context *gin.Context) {
	var req GetCommentsRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	postID, err := uuid.Parse(req.PostID)
	if err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := database.GetCommentsByPostParams{
		PostID: postID,
		Limit:  req.PageSize,
		Offset: (req.Page - 1) * req.PageSize,
	}

	comments, err := server.database.GetCommentsByPost(context, arg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	res := make([]database.CommentResponse, 0)

	for _, comment := range comments {
		res = append(res, comment.MakeResponse())
	}

	context.JSON(http.StatusOK, res)
}

type PostReplyRequest struct {
	CommentID string `json:"comment_id" binding:"required,uuid"`
	Content   string `json:"content" binding:"required,max=200"`
}

func (server *Server) postReply(context *gin.Context) {
	var req PostReplyRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	postID, err := uuid.Parse(req.CommentID)
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

	arg := database.PostReplyParams{
		ID:        id,
		UserID:    userID,
		CommentID: postID,
		Content:   req.Content,
	}

	reply, err := server.database.PostReply(context, arg)
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

	context.JSON(http.StatusCreated, reply)
}

type DeleteReplyRequest struct {
	ReplyID string `uri:"id" binding:"required,uuid"`
}

func (server *Server) deleteReply(context *gin.Context) {
	var req DeleteReplyRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(req.ReplyID)
	if err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authorizationPayload := context.MustGet("authorization_payload").(*token.Payload)
	userID := authorizationPayload.UserID

	reply, err := server.database.GetReplyById(context, id)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if userID != reply.UserID {
		context.Status(http.StatusForbidden)
		return
	}

	err = server.database.DeleteReply(context, id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.Status(http.StatusOK)
}

type GetRepliesRequest struct {
	CommentID string `form:"comment_id" binding:"required,uuid"`
	Page      int32  `form:"page_number" binding:"required,min=1"`
	PageSize  int32  `form:"page_size" binding:"required,min=1,max=20"`
}

func (server *Server) getReplies(context *gin.Context) {
	var req GetRepliesRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	commentID, err := uuid.Parse(req.CommentID)
	if err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := database.GetRepliesByCommentParams{
		CommentID: commentID,
		Limit:     req.PageSize,
		Offset:    (req.Page - 1) * req.PageSize,
	}

	replies, err := server.database.GetRepliesByComment(context, arg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	res := make([]database.ReplyResponse, 0)

	for _, reply := range replies {
		res = append(res, reply.MakeResponse())
	}

	context.JSON(http.StatusOK, res)
}
