package api

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	database "github.com/dqrk0jeste/letscube-backend/database/sqlc"
	"github.com/dqrk0jeste/letscube-backend/token"
	"github.com/dqrk0jeste/letscube-backend/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func makeUserResponse(user database.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,printascii,min=1,max=20"`
	Password string `json:"password" binding:"required,printascii,min=6,max=60"`
}

func (server *Server) createUser(context *gin.Context) {
	var req CreateUserRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	passwordHash, err := util.GeneratePasswordHash(req.Password)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := database.CreateUserParams{
		ID:           id,
		Username:     strings.Trim(req.Username, " "),
		PasswordHash: passwordHash,
	}

	user, err := server.database.CreateUser(context, arg)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				context.JSON(http.StatusConflict, errorResponse(err))
				return
			}
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusCreated, makeUserResponse(user))
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,printascii"`
	Password string `json:"password" binding:"required,printascii"`
}

type loginUserResponse struct {
	Token string       `json:"token"`
	User  userResponse `json:"user"`
}

func (server *Server) loginUser(context *gin.Context) {
	var req loginUserRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.database.GetUserByUsername(context, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.VerifyPassword(req.Password, user.PasswordHash)
	if err != nil {
		context.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	token, err := server.tokenMaker.CreateToken(user.ID, server.config.TokenDuration)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, loginUserResponse{
		Token: token,
		User:  makeUserResponse(user),
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

	context.JSON(http.StatusOK, makeUserResponse(user))
}

type GetUsersByUsernameRequest struct {
	Page     int32  `form:"page_number" binding:"required"`
	PageSize int32  `form:"page_size" binding:"required,min=1,max=20"`
	Input    string `form:"input" binding:"required,min=3,printascii"`
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

	res := make([]userResponse, 0)

	for _, user := range users {
		res = append(res, makeUserResponse(user))
	}

	context.JSON(http.StatusOK, res)
}

type UpdateUsersUsernameRequest struct {
	Username string `json:"username" binding:"required,printascii,min=1,max=20"`
}

func (server *Server) updateUsersUsername(context *gin.Context) {
	var req UpdateUsersUsernameRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authorizationPayload := context.MustGet("authorization_payload").(*token.Payload)

	arg := database.UpdateUsersUsernameParams{
		ID:       authorizationPayload.UserID,
		Username: req.Username,
	}

	user, err := server.database.UpdateUsersUsername(context, arg)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				context.JSON(http.StatusConflict, errorResponse(err))
				return
			}
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, makeUserResponse(user))
}

type UpdateUsersPasswordRequest struct {
	Password string `json:"password" binding:"required,printascii,min=6,max=60"`
}

func (server *Server) updateUsersPassword(context *gin.Context) {
	var req UpdateUsersPasswordRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authorizationPayload := context.MustGet("authorization_payload").(*token.Payload)

	passwordHash, err := util.GeneratePasswordHash(req.Password)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := database.UpdateUsersPasswordParams{
		ID:           authorizationPayload.UserID,
		PasswordHash: passwordHash,
	}

	user, err := server.database.UpdateUsersPassword(context, arg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, makeUserResponse(user))
}

// type GetUserByUsernameRequest struct {
// 	Username string `uri:"username" binding:"required"`
// }

// func (server *Server) getUserByUsername(context *gin.Context) {
// 	var req GetUserByUsernameRequest
// 	if err := context.ShouldBindUri(&req); err != nil {
// 		context.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	user, err := server.database.GetUserByUsername(context, req.Username)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			context.JSON(http.StatusNotFound, errorResponse(err))
// 			return
// 		}
// 		context.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	context.JSON(http.StatusOK, makeUserResponse(user))
// }

// type UpdateUserRequest struct {
// 	Username string `json:"username" binding:"printascii"`
// 	Password string `json:"password" binding:"printascii"`
// }

// func (server *Server) updateUser(context *gin.Context) {
// 	var req UpdateUserRequest
// 	if err := context.ShouldBindJSON(&req); err != nil {
// 		context.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	authorizationPayload := context.MustGet("authorization_payload").(*token.Payload)

// 	passwordHash, err := util.GeneratePasswordHash(req.Password)
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	arg := database.UpdateUserParams{
// 		ID:           authorizationPayload.UserID,
// 		Username:     req.Username,
// 		PasswordHash: passwordHash,
// 	}

// 	user, err := server.database.UpdateUser(context, arg)
// 	if err != nil {
// 		if err, ok := err.(*pq.Error); ok {
// 			if err.Code.Name() == "unique_violation" {
// 				context.JSON(http.StatusConflict, errorResponse(err))
// 				return
// 			}
// 		}
// 		context.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	context.JSON(http.StatusOK, makeUserResponse(user))
// 	return
// }
