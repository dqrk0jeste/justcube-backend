package api

import (
	database "github.com/dqrk0jeste/letscube-backend/database/sqlc"
	"github.com/dqrk0jeste/letscube-backend/s3_bucket"
	"github.com/dqrk0jeste/letscube-backend/token"
	"github.com/dqrk0jeste/letscube-backend/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config       util.Config
	database     *database.Queries
	tokenMaker   token.PasetoMaker
	router       *gin.Engine
	s3Controller *s3_bucket.S3Controller
}

func CreateServer(config util.Config, database *database.Queries) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSecret)
	if err != nil {
		return nil, err
	}

	s3Controller, err := s3_bucket.ControllerMaker()
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:       config,
		database:     database,
		tokenMaker:   *tokenMaker,
		s3Controller: s3Controller,
	}

	server.addRouter()

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) addRouter() {
	router := gin.Default()

	router.MaxMultipartMemory = 5 << 20

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	router.PUT("/users/username", server.authMiddleware, server.updateUsersUsername)
	router.PUT("/users/password", server.authMiddleware, server.updateUsersPassword)

	router.GET("/users/:id", server.getUserById)
	router.GET("/users", server.authMiddleware, server.getUsersByUsername)

	router.POST("/posts", server.authMiddleware, server.CreatePost)
	router.DELETE("/posts/:id", server.authMiddleware, server.DeletePost)
	router.GET("/posts/:id", server.GetPostById)
	router.GET("/posts", server.GetPostsByUser)

	server.router = router
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
