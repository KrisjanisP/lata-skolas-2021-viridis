// platform/router/router.go

package router

import (
	"encoding/gob"

	"github.com/KrisjanisP/viridis/platform/authenticator"

	"github.com/KrisjanisP/viridis/web/app/callback"
	"github.com/KrisjanisP/viridis/web/app/index"
	"github.com/KrisjanisP/viridis/web/app/karte"
	"github.com/KrisjanisP/viridis/web/app/login"
	"github.com/KrisjanisP/viridis/web/app/logout"
	"github.com/KrisjanisP/viridis/web/app/profile"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// New registers the routes and returns the router.
func New(auth *authenticator.Authenticator) *gin.Engine {
	router := gin.Default()

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

	router.Static("/assets", "./web/assets") // serve html
	router.Static("/dist", "./web/dist")     // serve javascript, css
	router.LoadHTMLGlob("web/template/*")

	router.GET("/", index.Handler)
	router.GET("/login", login.Handler(auth))
	router.GET("/callback", callback.Handler(auth))
	router.GET("/profile.html", profile.Handler)
	router.GET("/map.html", karte.Handler)
	router.GET("/logout", logout.Handler)
	/*
		router.GET("/tiles", getTiles)
		router.POST("/tiles", postTiles)
		config := cors.DefaultConfig()
		config.AllowAllOrigins = true
		router.Use(cors.New(config))
	*/
	return router
}

/*
func getTiles(c *gin.Context) {
	c.File("./data/tiles.geojson")
}

func postTiles(c *gin.Context) {
	var tile_names []string

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&tile_names); err != nil {
		return
	}

	if len(tile_names) >= 10 {
		c.AbortWithStatus(400)
		return
	}

	l := log.New(os.Stdout, "[API] ", log.Ldate|log.Ltime)
	l.Println("Received tiles: " + strings.Join(tile_names, " "))

	for _, tile_name := range tile_names {
		if tile_name == "" {
			c.AbortWithStatus(400)
			return
		}
		tile, err := dbapi.
		if err != nil {
			log.Fatal(err)
			c.AbortWithStatus(500)
			return
		}
		if tile.Name == "" {
			c.AbortWithStatus(400)
			return
		}
		time := time.Now().Format("2006-01-02 15:04:05")
		stmt, err := db.Prepare("INSERT INTO tasks_queue(tile_name, req_date, user_id) values(?,?,?)")
		if err != nil {
			log.Fatal(err)
			c.AbortWithStatus(500)
			return
		}

		_, err = stmt.Exec(tile_name, time, 1)
		if err != nil {
			log.Fatal(err)
			c.AbortWithStatus(500)
			return
		}
	}

	c.IndentedJSON(http.StatusCreated, tile_names)
}
*/
