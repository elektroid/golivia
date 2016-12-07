package main

import (
	"github.com/elektroid/golivia/generator"
	"github.com/elektroid/golivia/models"
	"github.com/gin-gonic/gin"
)

type NewAlbumIn struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	ViewType    string `json:"view_type" binding:"required"`
}

func NewAlbum(c *gin.Context, in *NewAlbumIn) (*models.Album, error) {

	a, err := models.CreateAlbum(db, in.Title, in.Description, in.ViewType)
	if err != nil {
		return nil, err
	}

	return a, nil
}

type GetAlbumIn struct {
	AlbumId int64 `path:"album_id, required"`
}

func GenerateAlbum(c *gin.Context, in *GetAlbumIn) error {
	album, err := models.LoadAlbumFromID(db, in.AlbumId)
	if err != nil {
		return err
	}

	err = album.LoadPhotos(db)
	if err != nil {
		return err
	}

	return generator.GenerateAlbum(album)
}

func GetAlbumHtml(c *gin.Context, in *GetAlbumIn) error {
	album, err := models.LoadAlbumFromID(db, in.AlbumId)
	if err != nil {
		return err
	}

	err = album.LoadPhotos(db)
	if err != nil {
		return err
	}

	content, err := generator.GetAlbumHtml(album)
	if err != nil {
		return err
	}

	c.Data(200, "text/html", []byte(content))
	return nil
}

func GetGalleriaHtml(c *gin.Context, in *GetAlbumIn) error {
	album, err := models.LoadAlbumFromID(db, in.AlbumId)
	if err != nil {
		return err
	}

	err = album.LoadPhotos(db)
	if err != nil {
		return err
	}

	content, err := generator.GetGalleriaHtml(album)
	if err != nil {
		return err
	}

	c.Data(200, "text/html", []byte(content))
	return nil
}

type GetAlbumByDateIn struct {
	Year  int64 `path:"year, required"`
	Month int64 `path:"month, required"`
}

func GetAlbumByDate(c *gin.Context, in *GetAlbumByDateIn) error {
	a, err := models.LoadPhotosByDate(db, in.Year, in.Month)
	if err != nil {
		return err
	}

	content, err := generator.GetAlbumHtml(a)
	if err != nil {
		return err
	}

	c.Data(200, "text/html", []byte(content))
	return nil
}
