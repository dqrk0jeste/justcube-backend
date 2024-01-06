package api

import (
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	database "github.com/dqrk0jeste/letscube-backend/database/sqlc"
	"github.com/dqrk0jeste/letscube-backend/s3_bucket"
	"github.com/dqrk0jeste/letscube-backend/token"
	"github.com/dqrk0jeste/letscube-backend/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     util.Config
	database   *database.Queries
	tokenMaker token.PasetoMaker
	router     *gin.Engine
	uploader   *manager.Uploader
}

func CreateServer(config util.Config, database *database.Queries) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSecret)
	if err != nil {
		return nil, err
	}

	uploader, err := s3_bucket.UploaderMaker()
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:     config,
		database:   database,
		tokenMaker: *tokenMaker,
		uploader:   uploader,
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

	router.GET("/users/:id", server.getUserById)
	router.GET("/users", server.authMiddleware, server.getUsersByUsername)

	router.POST("/posts", server.authMiddleware, server.CreatePost)
	router.GET("/posts/:id", server.GetPostById)
	router.GET("/posts", server.GetPostsByUser)

	server.router = router
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
