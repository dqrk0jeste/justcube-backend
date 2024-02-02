package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) refreshAccessToken(context *gin.Context) {
	refreshToken, err := context.Cookie("refresh_token")
	if err != nil {
		context.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	payload, err := server.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		context.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := server.database.GetSessionById(context, payload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked || session.UserID != payload.UserID {
		context.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	user, err := server.database.GetUserById(context, session.UserID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	newAccessToken, _, err := server.tokenMaker.CreateToken(session.UserID, server.config.AccessTokenDuration)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
		"user":         user.MakeResponse(),
	})
}
