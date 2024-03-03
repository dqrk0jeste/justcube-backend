package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var ALLOWED_ORIGINS = []string{
	"http://localhost:3000",
	"http://localhost:5173",
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

	usersRouter := router.Group("/users")

	usersRouter.POST("/", server.createUser)
	usersRouter.POST("/login", server.loginUser)

	usersRouter.GET("/refresh", server.refreshAccessToken)

	usersRouter.PUT("/username", server.authMiddleware, server.updateUsersUsername)
	usersRouter.PUT("/password", server.authMiddleware, server.updateUsersPassword)

	usersRouter.POST("/follows/:id", server.authMiddleware, server.followUser)
	usersRouter.DELETE("/follows/:id", server.authMiddleware, server.unfollowUser)

	usersRouter.GET("/followers", server.getFollowers)
	usersRouter.GET("/following", server.getFollowing)

	usersRouter.GET("/followers/count/:id", server.getFollowersCount)
	usersRouter.GET("/following/count/:id", server.getFollowingCount)

	usersRouter.GET("/:id", server.getUserById)
	usersRouter.GET("/", server.getUsersByUsername)

	postsRouter := router.Group("/posts")

	postsRouter.POST("/", server.authMiddleware, server.createPost)
	postsRouter.DELETE("/:id", server.authMiddleware, server.deletePost)

	postsRouter.GET(":id", server.getPostById)
	postsRouter.GET("/", server.getPostsByUser)

	postsRouter.GET("/feed", server.authMiddleware, server.getFeed)
	postsRouter.GET("/guest-feed", server.getGuestFeed)

	commentsRouter := postsRouter.Group("/comments")

	commentsRouter.POST("/", server.authMiddleware, server.postComment)
	commentsRouter.DELETE("/:id", server.authMiddleware, server.deleteComment)
	commentsRouter.GET("/", server.getComments)

	repliesRouter := commentsRouter.Group("/replies")

	repliesRouter.POST("/", server.authMiddleware, server.postReply)
	repliesRouter.DELETE("/:id", server.authMiddleware, server.deleteReply)
	repliesRouter.GET("/", server.getReplies)

	server.router = router
}
