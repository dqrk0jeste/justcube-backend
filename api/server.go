package api

import (
	database "github.com/dqrk0jeste/letscube-backend/database/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	database *database.Queries
	router   *gin.Engine
}

func CreateServer(database *database.Queries) *Server {
	router := gin.Default()

	server := &Server{
		database: database,
		router:   router,
	}

	router.POST("/users", server.createUser)
	router.GET("/users/:id", server.getUserById)
	router.GET("/users", server.getUsersByUsername)

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
