package main

import (
	"github.com/gin-gonic/gin"
	"github.com/elektroid/golivia/models"
	"github.com/elektroid/golivia/constants"
	"os"
	"log"
	"errors"
	"io"
)

type NewPhotoIn struct {
	AlbumId int64 `path:"album_id, required"`
}

func NewPhoto(c *gin.Context, in *NewPhotoIn) (*models.Photo, error) {

	path, err := ExtractFileFromForm(c)
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


func ExtractFileFromForm(c *gin.Context) (string, error ){
	file, header , err := c.Request.FormFile("upload")
	if err!=nil{
		log.Print(err)
		return "", err
	}
    filename := header.Filename
    out, err := os.Create(constants.TmpDir+"/"+filename)
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
	return constants.TmpDir+"/"+filename, nil
}

