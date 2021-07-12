package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"videoproject/controller"
	"videoproject/middleware"
	"videoproject/repository"
	"videoproject/service"
	"videoproject/util"

	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	videoRepository repository.VideoRepository = repository.NewVideoRepository()
	videoService    service.VideoService       = service.NewVideoService(videoRepository)
	loginService    service.LoginService       = service.NewLoginService()
	jwtService      service.JWTService         = service.NewJWTService()

	videoController controller.VideoController = controller.New(videoService)
	loginController controller.LoginController = controller.NewLoginController(loginService, jwtService)
)

/*func setupLogOutput() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

}*/

func main() {
	//close database connection
	defer videoRepository.CloseDB()

	//write GIN HTTP output to a log file so we can read stdout better
	logfile, err := os.Create("./gin_http.log")
	if err != nil {
		fmt.Println("Could not create log file")
	}
	gin.SetMode(gin.DebugMode)
	gin.DefaultWriter = io.MultiWriter(logfile)

	server := gin.New()

	//server.Use(gin.Recovery(), middleware.Logger(), middleware.BasicAuth(), gindump.Dump())

	server.Use(gin.Recovery(), gin.Logger())

	server.Static("/css", "./templates/css")

	server.LoadHTMLGlob("./templates/*.html")

	// Login Endpoint: Authentication + Token creation
	server.POST("/login", func(ctx *gin.Context) {
		token := loginController.Login(ctx)
		if token != "" {
			ctx.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		} else {
			ctx.JSON(http.StatusUnauthorized, nil)
		}
	})

	//JWT Authorization Middleware applies to */api* only.
	apiRoutes := server.Group("/api", middleware.AuthorizeJWT())
	{
		apiRoutes.GET("/videos", func(ctx *gin.Context) {
			ctx.JSON(200, videoController.FindAll())
		})

		apiRoutes.POST("/videos", func(ctx *gin.Context) {
			err := videoController.Save(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"message": "Video Input is Saved"})
			}

		})

		apiRoutes.PUT("/videos/:id", func(ctx *gin.Context) {
			err := videoController.Update(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"message": "Video is Updated"})
			}

		})

		apiRoutes.DELETE("/videos/:id", func(ctx *gin.Context) {
			err := videoController.Delete(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"message": "Video is Deleted"})
			}

		})

	}

	// The "/view" endpoints are public (no Authorization required)
	viewRoutes := server.Group("/view")
	{
		viewRoutes.GET("/videos", videoController.ShowAll)
	}

	config, err := util.LoadConfig(".")

	fmt.Println(config)
	if err != nil {
		log.Fatal("cannot load configuration \n", err)
	}

	port := config.HttpPort

	if port == "" {
		port = "8080"
		fmt.Printf("HTTPPORT not passed from env variable, using default %s \n", port)
	}
	server.Run(":" + port)

	//server.Run("127.0.0.1:8080")
	fmt.Println("Server is Running!")
}
