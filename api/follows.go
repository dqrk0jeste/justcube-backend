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

func (server *Server) FollowUser(context *gin.Context) {
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

func (server *Server) UnfollowUser(context *gin.Context) {
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
