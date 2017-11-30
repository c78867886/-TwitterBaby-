package server

import (
	"handler"
	"fmt"
	"time"
	"context"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

// NewServer : Instantiate a server
func NewServer(h *handler.Handler) (e *echo.Echo) {
	e = echo.New()
	e.HideBanner = true
	e.Logger.SetLevel(log.ERROR)
	e.Use(middleware.Logger())
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(handler.Key),
		Skipper: func(c echo.Context) bool {
			// Skip authentication for and signup login requests
			if c.Request().Method == "OPTIONS" || c.Path() == "/api/v1/login" || c.Path() == "/api/v1/signup" || c.Path() == "/api/v1/ws/:username" {
				return true
			}
			return false
		},
	}))
	
	// CORS config
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:4200"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowCredentials: true, 
		AllowHeaders: []string{echo.HeaderAuthorization, echo.HeaderContentType}, 
	}))

	// Routes
	e.GET("/api/v1/ws/:username", h.NotifHandler.GetConnection)
	e.POST("/api/v1/signup", h.UserHandler.Signup)
	e.POST("/api/v1/login", h.UserHandler.Login)
	e.POST("/api/v1/follow/:username", h.UserHandler.Follow)
	e.POST("/api/v1/unfollow/:username", h.UserHandler.Unfollow)
	e.GET("/api/v1/userInfo/:username", h.UserHandler.FetchUserInfo)
	e.POST("/api/v1/updateUserInfo", h.UserHandler.UpdateUserInfo)
	e.POST("/api/v1/updateProfilePic", h.UserHandler.UpdateProfilePicture)
	e.GET("/api/v1/showFollower/:username", h.UserHandler.ShowFollower)
	e.GET("/api/v1/showFollowing/:username", h.UserHandler.ShowFollowing)
	e.GET("/api/v1/tweetlist/:username", h.TweetHandler.FetchTweets)
	e.POST("/api/v1/newTweet", h.TweetHandler.NewTweet)
	e.DELETE("/api/v1/deleteTweet/:tweet", h.TweetHandler.DeleteTweet)
	e.GET("/api/v1/tweettimeline/:username", h.TweetHandler.FetchTweetTimeLine)
	e.POST("/api/v1/newcomment/:tweet", h.CommentHandler.NewComment)
	e.GET("/api/v1/fetchcomment/:tweet", h.CommentHandler.FetchComment)
	e.POST("/api/v1/reTweet", h.TweetHandler.ReTweet)

	return e;
}

// TerminalControl : Thread for terminal control
func TerminalControl(e *echo.Echo, h *handler.Handler, srvAddr string) {
	var op string
	
	for {
		fmt.Println("Listening on " + srvAddr)
		fmt.Print("Option('h' for help): ")
		fmt.Scanln(&op)
		if op == "q" {
			fmt.Println("Shutting down server.")
			ShutdownServer(e, h)
			break
		} else if op == "h" {
			fmt.Println("'q' to shutdown server")
			fmt.Println("'i' to reconstruct database to default (w/ some initial collections)")
			fmt.Println("'r' to reconstruct testing database")
		} else if op == "i" {
			dbReinsert()
			fmt.Println("Database reconstructed.")
		} else if op == "r" {
			reconstructTestDB()
			fmt.Println("Testing database reconstructed.")
		}
	}
}

// ShutdownServer : Shutdown the server
func ShutdownServer(e *echo.Echo, h *handler.Handler) {
	h.Shutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
