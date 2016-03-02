package models_test

import (
	"testing"
	"github.com/elektroid/golivia/models"
)

func TestAlbum(t *testing.T) {

	// we don't want duplicated albums
	// we want to be able to list albums
	// we want to attach photos to an album
	// we want to list photos in an album ....
	// and so on ... let's see if we make methods in the models or write code in the handlers

	_, err := models.CreateAlbum(nil, "random title", "a test album", models.PUBLIC)
	if err== nil {
		t.Error("Failed to detect nil db")
		t.Fatal(err)
	}

}


