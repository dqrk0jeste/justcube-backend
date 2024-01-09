package api

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	database "github.com/dqrk0jeste/letscube-backend/database/sqlc"
	"github.com/dqrk0jeste/letscube-backend/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TODO: find some image processing library for compressing images

type CreatePostRequest struct {
	ImageContent []*multipart.FileHeader `form:"image_content[]" binding:"required,max=5"`
	TextContent  string                  `form:"text_content" binding:"required,max=500"`
}

var supportedImageTypes = []string{
	"image/jpeg",
	"image/png",
}

func (server *Server) CreatePost(context *gin.Context) {
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
		if !slices.Contains(supportedImageTypes, image.Header.Get("Content-Type")) {
			context.AbortWithStatus(http.StatusUnsupportedMediaType)
			return
		}
	}

	for index, image := range files {
		// this is only for development, should delete in production
		context.SaveUploadedFile(image, "images/"+image.Filename)
		openedImage, err := image.Open()
		if err != nil {
			context.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		imageNameSeparated := strings.Split(image.Filename, ".")
		imageExtention := imageNameSeparated[len(imageNameSeparated)-1]
		imageNameToSave := id.String() + "_" + strconv.Itoa(index) + "." + imageExtention

		uploadedImage, err := server.uploader.Upload(context, &s3.PutObjectInput{
			Bucket: aws.String("letscube"),
			Key:    aws.String(imageNameToSave),
			Body:   openedImage,
		})
		if err != nil {
			context.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		// this is only for development, should delete in production
		fmt.Println(uploadedImage.Location)
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

type GetPostByIdRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// in the future, we will have public and private accounts so we will have to check if user that requested the post follows user whose post it is
func (server *Server) GetPostById(context *gin.Context) {

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

	context.JSON(http.StatusCreated, post)
}

type GetPostsByUserRequest struct {
	UserID   string `form:"user_id" binding:"required,uuid"`
	Page     int32  `form:"page_number" binding:"required"`
	PageSize int32  `form:"page_size" binding:"required,max=20"`
}

func (server *Server) GetPostsByUser(context *gin.Context) {

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

	context.JSON(http.StatusCreated, posts)
}
