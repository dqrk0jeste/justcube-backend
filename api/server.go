package api

import (
	database "github.com/dqrk0jeste/letscube-backend/database/sqlc"
	"github.com/dqrk0jeste/letscube-backend/token"
	"github.com/dqrk0jeste/letscube-backend/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     util.Config
	database   *database.Queries
	tokenMaker token.PasetoMaker
	router     *gin.Engine
}

func CreateServer(config util.Config, database *database.Queries) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSecret)
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:     config,
		database:   database,
		tokenMaker: *tokenMaker,
	}

	server.addRouter()

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) addRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	router.GET("/users/:id", server.getUserById)
	router.GET("/users", server.authMiddleware, server.getUsersByUsername)

	server.router = router
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
