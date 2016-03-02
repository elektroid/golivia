package main

import (
	"github.com/gin-gonic/gin"
	"github.com/elektroid/golivia/models"
)

type NewAlbumIn struct {
	Title string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	ViewType string `json:"view_type" binding:"required"`
}

func NewAlbum(c *gin.Context, in *NewAlbumIn) (*models.Album, error) {

	a, err := models.CreateAlbum(db, in.Title, in.Description, in.ViewType)
	if err != nil {
		return nil, err
	}

	return a, nil
}

