// package api

// import "github.com/gin-gonic/gin"

// func (server *Server) authUser(context *gin.Context) {

// }

package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func (server *Server) authMiddleware(context *gin.Context) {
	authorizationHeader := context.GetHeader(authorizationHeaderKey)

	if len(authorizationHeader) == 0 {
		err := errors.New("authorization header is not provided")
		context.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	fields := strings.Fields(authorizationHeader)
	if len(fields) != 2 {
		err := errors.New("invalid authorization header format")
		context.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != authorizationTypeBearer {
		err := fmt.Errorf("unsupported authorization type %s", authorizationType)
		context.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	token := fields[1]
	payload, err := server.tokenMaker.VerifyToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	context.Set(authorizationPayloadKey, payload)
	context.Next()
}
