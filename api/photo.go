package main

import (
	"github.com/gin-gonic/gin"
	"github.com/elektroid/golivia/models"
	"github.com/elektroid/golivia/utils/web"

)

type NewPhotoIn struct {
	AlbumId int64 `path:"album_id, required"`
}

func NewPhoto(c *gin.Context, in *NewPhotoIn) (*models.Photo, error) {

	path, err := web.ExtractFileFromForm(c, "upload")
	if err!=nil{
		return nil, err	
	}

	a, err := models.LoadAlbumFromID(db, in.AlbumId)
	if err != nil {
		return nil, err
	}

	p, err := models.CreatePhoto(db, a, path)
	if err != nil {
		return nil, err
	}

	return p, nil
}



type CheckPhotoIn struct {
	Md5 string `path:"md5, required"`
}
func CheckPhotoMd5(c *gin.Context, in *CheckPhotoIn) (*models.Photo, error) {
	p, err := models.LoadPhotoFromMd5(db, in.Md5)
	return p, err
}



