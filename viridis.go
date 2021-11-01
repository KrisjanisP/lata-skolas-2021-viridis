package main

import (
	"log"
	"net/http"

	"github.com/KrisjanisP/viridis/api"
	"github.com/KrisjanisP/viridis/database"
	"github.com/KrisjanisP/viridis/platform/authenticator"
	"github.com/KrisjanisP/viridis/platform/middleware"
	"github.com/KrisjanisP/viridis/platform/router"
	"github.com/KrisjanisP/viridis/utils"
	"github.com/KrisjanisP/viridis/web/app/callback"
	"github.com/KrisjanisP/viridis/web/app/index"
	"github.com/KrisjanisP/viridis/web/app/karte"
	"github.com/KrisjanisP/viridis/web/app/login"
	"github.com/KrisjanisP/viridis/web/app/logout"
	"github.com/KrisjanisP/viridis/web/app/profile"
	"github.com/joho/godotenv"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbapi, err := database.NewDB()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	rtr := router.New()

	rtr.Static("/assets", "./web/assets") // serve html
	rtr.Static("/dist", "./web/dist")     // serve javascript, css
	rtr.LoadHTMLGlob("web/template/*")

	rtr.GET("/", index.Handler)
	rtr.GET("/login", login.Handler(auth))
	rtr.GET("/callback", callback.Handler(auth))
	rtr.GET("/profile.html", middleware.IsAuthenticated, profile.Handler)
	rtr.GET("/map.html", middleware.IsAuthenticated, karte.Handler)
	rtr.GET("/logout", logout.Handler)

	rtr.GET("/tiles", api.DBAPI{DBAPI: dbapi}.GetTiles)
	rtr.POST("/tiles", api.DBAPI{DBAPI: dbapi}.PostTiles)

	go utils.StartWorker(dbapi)

	if err := http.ListenAndServe(":8080", rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
