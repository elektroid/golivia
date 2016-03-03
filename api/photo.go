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
	AlbumId int64 `path:"album_id, required"`
}

func NewPhoto(c *gin.Context, in *NewPhotoIn) (*models.Photo, error) {

	fmt.Println("extract file from form")
	path, err := ExtractFileFromForm(c)
	if err!=nil{
		fmt.Println("failed to extract from form")
		return nil, err	
	}

	p, err := models.CreatePhoto(db, in.AlbumId, path)
	if err != nil {
		return nil, err
	}

	fmt.Println("set file")
	p.SetFile(path)

	return p, nil
}


func ExtractFileFromForm(c *gin.Context) (string, error ){
	file, header , err := c.Request.FormFile("upload")
	if err!=nil{
		log.Print(err)
		return "", err
	}
    filename := header.Filename
    fmt.Println(header.Filename)
    out, err := os.Create("/tmp/"+filename)
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
	return "/tmp/"+filename, nil
}

