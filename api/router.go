package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var ALLOWED_ORIGINS = []string{
	"http://localhost:3000",
}

var ALLOWED_HEADERS = []string{
	"Origin",
	"Authorization",
	"Content-Type",
}

func (server *Server) addRouter() {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = ALLOWED_ORIGINS
	config.AllowHeaders = ALLOWED_HEADERS

	router.Use(cors.New(config))

	router.MaxMultipartMemory = 5 << 20

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	router.GET("/users/refresh", server.refreshAccessToken)

	router.PUT("/users/username", server.authMiddleware, server.updateUsersUsername)
	router.PUT("/users/password", server.authMiddleware, server.updateUsersPassword)

	router.POST("/users/follow/:id", server.authMiddleware, server.followUser)
	router.DELETE("/users/follow/:id", server.authMiddleware, server.unfollowUser)

	router.GET("/users/followers", server.getFollowers)
	router.GET("/users/following", server.getFollowing)

	router.GET("/users/followers/:id", server.getFollowersCount)
	router.GET("/users/following/:id", server.getFollowingCount)

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
	router.GET("/posts/feed", server.authMiddleware, server.getFeed)

	server.router = router
}
