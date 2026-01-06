package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v82"
	"schej.it/server/db"
	"schej.it/server/logger"
	"schej.it/server/routes"
	"schej.it/server/services/gcloud"
	"schej.it/server/slackbot"
	"schej.it/server/utils"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "schej.it/server/docs"
)

// @title Schej.it API
// @version 1.0
// @description This is the API for Schej.it!

// @host localhost:3002/api

func main() {
	// Set release flag
	release := flag.Bool("release", false, "Whether this is the release version of the server")
	flag.Parse()
	if *release {
		os.Setenv("GIN_MODE", "release")
		gin.SetMode(gin.ReleaseMode)
	} else {
		os.Setenv("GIN_MODE", "debug")
	}

	// Init logfile
	logFile, err := os.OpenFile("logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)

	// Init logger
	logger.Init(logFile)

	// Load .env variables
	loadDotEnv()

	// Init router
	router := gin.New()
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}

		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}
		return fmt.Sprintf("%v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
			param.TimeStamp.Format("2006/01/02 15:04:05"),
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Path,
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())

	// Cors
	// Get allowed origins from environment variable or use defaults
	allowedOrigins := []string{"http://localhost:3002", "http://localhost:8080"}
	if envOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); envOrigins != "" {
		// Split by comma and trim spaces
		for _, origin := range strings.Split(envOrigins, ",") {
			trimmedOrigin := strings.TrimSpace(origin)
			if trimmedOrigin != "" {
				allowedOrigins = append(allowedOrigins, trimmedOrigin)
			}
		}
		logger.StdOut.Printf("CORS: Added custom origins from CORS_ALLOWED_ORIGINS: %v\n", envOrigins)
	} else {
		// Add production domains if not using custom origins
		allowedOrigins = append(allowedOrigins, "https://www.schej.it", "https://schej.it", "https://www.timeful.app", "https://timeful.app")
	}
	
	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Init database
	closeConnection := db.Init()
	defer closeConnection()

	// Init google cloud stuff
	closeTasks := gcloud.InitTasks()
	defer closeTasks()

	// Session
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("session", store))

	// Init routes
	apiRouter := router.Group("/api")
	routes.InitAuth(apiRouter)
	routes.InitUser(apiRouter)
	routes.InitEvents(apiRouter)
	routes.InitAnalytics(apiRouter)
	routes.InitStripe(apiRouter)
	routes.InitFolders(apiRouter)
	slackbot.InitSlackbot(apiRouter)

	// Serve frontend static files only if the directory exists
	// In Docker deployments with separated containers, frontend is served by Nginx
	frontendDistPath := "../frontend/dist"
	if _, err := os.Stat(frontendDistPath); err == nil {
		// Frontend dist directory exists, serve static files
		err = filepath.WalkDir(frontendDistPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && d.Name() != "index.html" {
				split := splitPath(path)
				newPath := filepath.Join(split[3:]...)
				router.StaticFile(fmt.Sprintf("/%s", newPath), path)
			}
			return nil
		})
		if err != nil {
			log.Fatalf("failed to walk directories: %s", err)
		}

		router.LoadHTMLFiles("../frontend/dist/index.html")
		router.NoRoute(noRouteHandler())
		logger.StdOut.Println("Frontend static files loaded from ../frontend/dist")
	} else {
		// Frontend not found, likely running in separated container mode
		logger.StdOut.Println("Frontend directory not found, skipping static file serving (expected in separated container deployments)")
		// In separated container mode, only serve API routes
		router.NoRoute(func(c *gin.Context) {
			c.JSON(404, gin.H{"error": "Not found. This is the API server - frontend is served separately."})
		})
	}

	// Init swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Run server
	router.Run(":3002")
}

// Load .env variables
func loadDotEnv() {
	err := godotenv.Load(".env")

	// Load stripe key
	stripe.Key = os.Getenv("STRIPE_API_KEY")

	if err != nil {
		// In Docker or production environments, .env file might not exist
		// Environment variables should be set directly
		logger.StdOut.Println("No .env file found, using environment variables directly")
	}
}

func noRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		params := gin.H{}
		path := c.Request.URL.Path

		// Determine meta tags based off URL
		if match := regexp.MustCompile(`\/e\/(\w+)`).FindStringSubmatchIndex(path); match != nil {
			// /e/:eventId
			eventId := path[match[2]:match[3]]
			event := db.GetEventByEitherId(eventId)

			if event != nil {
				title := fmt.Sprintf("%s - Timeful (formerly Schej)", event.Name)
				params = gin.H{
					"title":   title,
					"ogTitle": title,
				}

				if len(utils.Coalesce(event.When2meetHref)) > 0 {
					params["ogImage"] = "/img/when2meetOgImage2.png"
				}
			}
		}

		c.HTML(http.StatusOK, "index.html", params)
	}
}

func splitPath(path string) []string {
	dir, last := filepath.Split(path)
	if dir == "" {
		return []string{last}
	}
	return append(splitPath(filepath.Clean(dir)), last)
}
