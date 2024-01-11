package api

import (
	"net/http"

	database "github.com/dqrk0jeste/letscube-backend/database/sqlc"
	"github.com/dqrk0jeste/letscube-backend/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type FollowUserRequest struct {
	FollowedUserID string `uri:"id" binding:"required,uuid"`
}

func (server *Server) followUser(context *gin.Context) {
	var req FollowUserRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	followedUserID, err := uuid.Parse(req.FollowedUserID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authorizationPayload := context.MustGet("authorization_payload").(*token.Payload)

	if authorizationPayload.UserID == followedUserID {
		context.Status(http.StatusBadRequest)
		return
	}

	arg := database.FollowUserParams{
		FollowedUserID: followedUserID,
		UserID:         authorizationPayload.UserID,
	}

	follow, err := server.database.FollowUser(context, arg)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			switch err.Code.Name() {
			case "unique_violation":
				{
					context.JSON(http.StatusConflict, errorResponse(err))
					return
				}
			case "foreign_key_violation":
				{
					context.JSON(http.StatusBadRequest, errorResponse(err))
					return
				}
			}
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusCreated, follow)
}

type UnfollowUserRequest struct {
	FollowedUserID string `uri:"id" binding:"required,uuid"`
}

func (server *Server) unfollowUser(context *gin.Context) {
	var req UnfollowUserRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	followedUserID, err := uuid.Parse(req.FollowedUserID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authorizationPayload := context.MustGet("authorization_payload").(*token.Payload)

	arg := database.UnfollowUserParams{
		FollowedUserID: followedUserID,
		UserID:         authorizationPayload.UserID,
	}

	err = server.database.UnfollowUser(context, arg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.Status(http.StatusOK)
}

type GetFollowersRequest struct {
	Page     int32  `form:"page_number" binding:"required,min=1"`
	PageSize int32  `form:"page_size" binding:"required,min=1,max=20"`
	UserID   string `form:"user_id" binding:"required,uuid"`
}

func (server *Server) getFollowers(context *gin.Context) {
	var req GetFollowersRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := database.GetFollowersParams{
		FollowedUserID: userID,
		Offset:         (req.Page - 1) * req.PageSize,
		Limit:          req.PageSize,
	}

	followers, err := server.database.GetFollowers(context, arg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := make([]userResponse, 0)

	for _, user := range followers {
		res = append(res, makeUserResponse(user))
	}

	context.JSON(http.StatusOK, res)
}

type GetFollowingRequest struct {
	Page     int32  `form:"page_number" binding:"required,min=1"`
	PageSize int32  `form:"page_size" binding:"required,min=1,max=20"`
	UserID   string `form:"user_id" binding:"required,uuid"`
}

func (server *Server) getFollowing(context *gin.Context) {
	var req GetFollowingRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := database.GetFollowingParams{
		UserID: userID,
		Offset: (req.Page - 1) * req.PageSize,
		Limit:  req.PageSize,
	}

	following, err := server.database.GetFollowing(context, arg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := make([]userResponse, 0)

	for _, user := range following {
		res = append(res, makeUserResponse(user))
	}

	context.JSON(http.StatusOK, res)
}

type GetFollowersCountRequest struct {
	UserID string `uri:"id" binding:"required,uuid"`
}

func (server *Server) getFollowersCount(context *gin.Context) {
	var req GetFollowersCountRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	followersCount, err := server.database.GetFollowersCount(context, userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"followers_count": followersCount,
	})
}

type GetFollowingCountRequest struct {
	UserID string `uri:"id" binding:"required,uuid"`
}

func (server *Server) getFollowingCount(context *gin.Context) {
	var req GetFollowingCountRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	followingCount, err := server.database.GetFollowingCount(context, userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"following_count": followingCount,
	})
}
