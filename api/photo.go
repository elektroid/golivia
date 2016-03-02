package main

import (
	"github.com/gin-gonic/gin"
	"github.com/elektroid/golivia/models"
	"fmt"
	"os"
	"log"
	"errors"
	"io"
)

type NewPhotoIn struct {
	AlbumId int64 `json:"album_id" binding:"required"`
}

func NewPhoto(c *gin.Context, in *NewAlbumIn) (*models.Photo, error) {

	path, err := ExtractFileFromForm(c)
	if err!=nil{
		return nil, err	
	}

	s, err := models.CreatePhoto(db, path)
	if err != nil {
		return nil, err
	}

	return s, nil
}


func ExtractFileFromForm(c *gin.Context) (string, error ){
	file, header , err := c.Request.FormFile("upload")
        filename := header.Filename
        fmt.Println(header.Filename)
        out, err := os.Create("./tmp/"+filename+".png")
        if err != nil {
            log.Fatal(err)
	    return "", errors.New("extraction failure")
        }
        defer out.Close()
        _, err = io.Copy(out, file)
        if err != nil {
            log.Fatal(err)
	    return "", errors.New("extraction failure")
        }   
	return "./tmp/"+filename+".png", nil
}

