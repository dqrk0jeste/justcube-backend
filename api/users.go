package api

import (
	"database/sql"
	"net/http"
	"strings"

	database "github.com/dqrk0jeste/letscube-backend/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Username     string         `json:"username" binding:"required"`
	Password     string         `json:"password" binding:"required"`
	ProfileImage sql.NullString `json:"profile_image"`
}

func (server *Server) createUser(context *gin.Context) {
	var req CreateUserRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := database.CreateUserParams{
		ID:           uuid.New(),
		Username:     strings.Trim(req.Username, " "),
		PasswordHash: req.Password,
		ProfileImage: req.ProfileImage,
	}

	user, err := server.database.CreateUser(context, arg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"id":            user.ID,
		"username":      user.Username,
		"profile_image": user.ProfileImage.String,
	})
}

type GetUserByIdRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func (server *Server) getUserById(context *gin.Context) {
	var req GetUserByIdRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.database.GetUserById(context, id)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"id":            user.ID,
		"username":      user.Username,
		"profile_image": user.ProfileImage.String,
	})
}

type GetUsersByUsernameRequest struct {
	Page     int32  `form:"page_number" binding:"required"`
	PageSize int32  `form:"page_size" binding:"required"`
	Input    string `form:"input" binding:"required"`
}

func (server *Server) getUsersByUsername(context *gin.Context) {
	var req GetUsersByUsernameRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := database.GetUsersByUsernameParams{
		Input:  strings.Trim(req.Input, " ") + "%",
		Offset: (req.Page - 1) * req.PageSize,
		Limit:  req.PageSize,
	}

	users, err := server.database.GetUsersByUsername(context, arg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := make([]gin.H, 0)

	for _, user := range users {
		res = append(res, gin.H{
			"id":            user.ID,
			"username":      user.Username,
			"profile_image": user.ProfileImage.String,
		})
	}

	context.JSON(http.StatusOK, res)
}
