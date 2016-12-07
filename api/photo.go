package main

import (
	//"errors"
	"github.com/elektroid/golivia/constants"
	"github.com/elektroid/golivia/models"
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"

	"io"
	"log"
	"os"
)

type NewPhotoIn struct {
	AlbumId int64 `path:"album_id, required"`
}

func NewPhoto(c *gin.Context, in *NewPhotoIn) (*models.Photo, error) {

	log.Println("in new photo")

	path, err := ExtractFileFromForm(c)
	if err != nil {
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

type GetPhotoIn struct {
	Md5 string `path:"md5, required"`
}

func LoadPhotoFromMD5(c *gin.Context, in *GetPhotoIn) (*models.Photo, error) {
	p, err := models.LoadPhotoFromMd5(db, in.Md5)
	return p, err
}

func ExtractFileFromForm(c *gin.Context) (string, error) {
	file, header, err := c.Request.FormFile("upload")
	if err != nil {
		return "", errors.BadRequestf("unable to find upload file: %s", err.Error())
	}
	filename := header.Filename
	out, err := os.Create(constants.TmpDir + "/" + filename)
	if err != nil {
		log.Fatal(err)
		return "", errors.New("extraction failure, not able to create")
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
		return "", errors.New("extraction failure, not able to copy content")
	}
	return constants.TmpDir + "/" + filename, nil
}
