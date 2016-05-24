package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-gorp/gorp"
	"github.com/loopfz/gadgeto/zesty"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/loopfz/gadgeto/tonic/jujuerrhook"
	"github.com/elektroid/golivia/constants"
	"github.com/elektroid/golivia/db/initdb"
)

var db *gorp.DbMap // TODO remove

func main() {

	tdb, err := initdb.InitSqlite()
	if err != nil {
		panic(err)
	}

	tonic.SetErrorHook(jujuerrhook.ErrHook)

	db = tdb
	zesty.RegisterDB(zesty.NewDB(tdb), constants.DBName)

	router := gin.Default()
  	router.Use(gin.Logger())
  	router.Use(gin.Recovery())

	// Scenarios
	router.POST("/album", tonic.Handler(NewAlbum, 201))
	router.POST("/album/:album_id/photo", tonic.Handler(NewPhoto, 201))

	// call to generate
	router.GET("/album/:album_id", tonic.Handler(GenerateAlbum, 201))

	router.Run(":8080")
}
