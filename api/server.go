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
	tokenMaker   *token.PasetoMaker
	router       *gin.Engine
	s3Controller *s3_bucket.S3Controller
}

func CreateServer(config util.Config, database *database.Queries) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSecret)
	if err != nil {
		return nil, err
	}

	s3Controller, err := s3_bucket.NewController()
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:       config,
		database:     database,
		tokenMaker:   tokenMaker,
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

	router.POST("/users/follow/:id", server.authMiddleware, server.followUser)
	router.DELETE("/users/follow/:id", server.authMiddleware, server.unfollowUser)

	router.GET("/users/followers", server.getFollowers)
	router.GET("/users/following", server.getFollowing)

	router.GET("/users/followers/count/:id", server.getFollowersCount)
	router.GET("/users/following/count/:id", server.getFollowingCount)

	router.GET("/users/:id", server.getUserById)
	router.GET("/users", server.getUsersByUsername)

	router.POST("/posts", server.authMiddleware, server.createPost)
	router.DELETE("/posts/:id", server.authMiddleware, server.deletePost)

	router.POST("/posts/comments", server.authMiddleware, server.postComment)
	router.DELETE("/posts/comments/:id", server.authMiddleware, server.deleteComment)
	router.GET("/posts/comments", server.getComments)

	router.POST("/posts/comments/replies", server.authMiddleware, server.postReply)
	router.DELETE("/posts/comments/replies/:id", server.authMiddleware, server.deleteReply)
	router.GET("/posts/comments/replies", server.getReplies)

	router.GET("/posts/:id", server.getPostById)
	router.GET("/posts", server.getPostsByUser)

	server.router = router
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
