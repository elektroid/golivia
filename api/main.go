package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-gorp/gorp"
	"github.com/loopfz/gadgeto/zesty"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/loopfz/gadgeto/tonic/jujuerrhook"
	"github.com/elektroid/golivia/constants"
	"github.com/elektroid/golivia/db/initdb"
	"flag"
	"fmt"
	"github.com/vharitonsky/iniflags"
)





var db *gorp.DbMap // TODO remove

func main() {

	// let's read config
	var (
    	dataBaseDir   = flag.String("dataBaseDir", "/tmp", "Directory where sqlite db is")
    	photoBaseDir  = flag.String("photosBaseDir", "/tmp", "Directory where photos are stored")
    	adminPassword = flag.String("adminPassword", "adminPassword", "Admin password to create and modify albums/photos")
    	usersPassword = flag.String("usersPassword", "usersPassword", "User password to view albums/photos")
    	serverPort    = flag.Int("serverPort", 8080, "Port to listen to")
	)
	iniflags.Parse()
	constants.PhotosDir = *photoBaseDir


	// start db
	tdb, err := initdb.InitSqlite(*dataBaseDir)
	if err != nil {
		panic(err)
	}
	db = tdb

	tonic.SetErrorHook(jujuerrhook.ErrHook)
	zesty.RegisterDB(zesty.NewDB(tdb), constants.DBName)

	router := gin.Default()
  	router.Use(gin.Logger())
  	router.Use(gin.Recovery())

  	// set up admin routes
  	administrated := router.Group("/", gin.BasicAuth(gin.Accounts{
  		"admin": *adminPassword,
  	}))
	administrated.POST("/album", tonic.Handler(NewAlbum, 201))
	administrated.POST("/album/:album_id/photo", tonic.Handler(NewPhoto, 201))

	// set up client routes, might become 100 % static
	viewers := router.Group("/", gin.BasicAuth(gin.Accounts{
        "amigos": *usersPassword,
    }))
	viewers.POST("/album/:album_id", tonic.Handler(GenerateAlbum, 201))
	viewers.GET("/album/:album_id", tonic.Handler(GetAlbumHtml, 201))


	router.Static("/photos", constants.PhotosDir)
  //  router.StaticFS("/more_static", http.Dir("my_file_system"))
    router.StaticFile("/favicon.ico", "./resources/favicon.ico")


	router.Run(fmt.Sprintf(":%d", *serverPort))
}
