package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-gorp/gorp"
	"github.com/loopfz/gadgeto/zesty"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/elektroid/golivia/constants"
	"github.com/elektroid/golivia/utils/errors"
	"github.com/elektroid/golivia/db/initdb"
	"flag"
	"fmt"
	"time"
	"log"
	"github.com/vharitonsky/iniflags"
)


func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        t := time.Now()

        // Set example variable
        c.Set("example", "12345")

        // before request

        c.Next()

        // after request
        latency := time.Since(t)
        log.Print(latency)

        // access the status we are sending
        status := c.Writer.Status()
        log.Println(status)
    }
}


var db *gorp.DbMap // TODO remove

func main() {

	// let's read config
	var (
    	dataBaseDir   = flag.String("dataBaseDir", "/tmp", "Directory where sqlite db is")
    	photosBaseDir = flag.String("photosBaseDir", "/tmp", "Directory where photos are stored")
    	tmpDir		  = flag.String("tmpDir", "/tmp", "Directory where photos are stored")
    	adminPassword = flag.String("adminPassword", "adminPassword", "Admin password to create and modify albums/photos")
    	usersPassword = flag.String("usersPassword", "usersPassword", "User password to view albums/photos")
    	serverPort    = flag.Int("serverPort", 8080, "Port to listen to")
	)
	iniflags.Parse()
	constants.PhotosDir = *photosBaseDir
	constants.TmpDir = *tmpDir


	// start db
	tdb, err := initdb.InitSqlite(*dataBaseDir)
	if err != nil {
		panic(err)
	}
	db = tdb

	tonic.SetErrorHook(goofs.ErrHook)
	zesty.RegisterDB(zesty.NewDB(tdb), constants.DBName)

	router := gin.Default()
  	router.Use(Logger())
  	router.Use(gin.Recovery())

  	// set up admin routes
  	administrated := router.Group("/", gin.BasicAuth(gin.Accounts{
  		"admin": *adminPassword,
  	}))
	administrated.POST("/album", tonic.Handler(NewAlbum, 201))
	administrated.POST("/album/:album_id/photo", tonic.Handler(NewPhoto, 201))
	administrated.GET("/photo/:md5", tonic.Handler(CheckPhotoMd5, 201))

	// set up client routes, might become 100 % static
	viewers := router.Group("/", gin.BasicAuth(gin.Accounts{
        "amigos": *usersPassword,
    }))
	viewers.GET("/album1/:album_id", tonic.Handler(GetAlbumHtml, 201))
	viewers.GET("/gal/:album_id", tonic.Handler(GetGalleriaHtml, 201))
	
	viewers.GET("/by_dates", tonic.Handler(GetPopulatedDates, 201))
	viewers.GET("/date/:year/:month", tonic.Handler(GetAlbumByDate, 201))



	router.Static("/photos", constants.PhotosDir)
  //  router.StaticFS("/more_static", http.Dir("my_file_system"))
	router.Static("/resources", "./resources")
	router.StaticFile("/", "./generator/templates/assets/welcome.html")
    //router.StaticFile("/favicon.ico", "./resources/favicon.ico")


	panic(router.Run(fmt.Sprintf(":%d", *serverPort)))
}
