package models

import (
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/go-gorp/gorp"
	"github.com/elektroid/golivia/utils/sqlgenerator"
	"github.com/elektroid/golivia/constants"
	"fmt"
//	"log"
)

const (
	// Const names to be retrievable from model code
	PUBLIC  = "public"
	PRIVATE   = "private"

)

type Album struct {
	ID         int64  `json:"id" db:"id"`
	Title  	   string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	ViewType string `json:"view_type" db:"view_type"`
	MiniatureWidth uint `json:"miniature_width" db:"miniatureWidth"`
	MiniatureHeight uint `json:"miniature_height" db:"miniatureHeight"`
	Photos []*Photo `db:"-"`
}

// Create an icon
func CreateAlbum(db *gorp.DbMap, Title string, Description string, ViewType string) (*Album, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to create album")
	}

	a := &Album{
		Title : Title,
		Description: Description,
		ViewType: ViewType,
		Photos: []*Photo{},
		MiniatureWidth: constants.DEFAULT_MINIATURE_W,
		MiniatureHeight: constants.DEFAULT_MINIATURE_H,
	}

	err := a.Valid()
	if err != nil {
		return nil, err
	}

	err = db.Insert(a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Album) LoadPhotos(db *gorp.DbMap)(error){
	if db == nil {
		return errors.New("Missing db parameter to list albums")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"photo"`).Where(squirrel.Eq{`album_id`: a.ID})
	query, args, err := selector.ToSql()
	if err != nil {
		return err
	}

	var photos []*Photo
	_, err = db.Select(&photos, query, args...)
	if err != nil {
		return  err
	}
	a.Photos=photos
	return nil
}


func LoadPhotosByDate(db *gorp.DbMap, Year int64, Month int64)(*Album, error){
	if db == nil {
		return nil, errors.New("Missing db parameter to list albums")
	}



	query := "SELECT * FROM photo WHERE strftime('%Y-%m', time) = '"+fmt.Sprintf("%d-%02d", Year, Month)+"'"
	fmt.Print(query)

	a := &Album{
		Title : fmt.Sprintf("Album for %d-%02d", Year, Month),
		Description: fmt.Sprintf("Album for %d-%02d", Year, Month),
		MiniatureWidth: constants.DEFAULT_MINIATURE_W,
		MiniatureHeight: constants.DEFAULT_MINIATURE_H,
	}


	var photos []*Photo
	_, err := db.Select(&photos, query, nil)
	if err != nil {
		return nil, err
	}
	a.Photos=photos
	return a, nil

}



func ListAlbums(db *gorp.DbMap, ViewType *string)([]*Album, error){
	if db == nil {
		return nil, errors.New("Missing db parameter to list albums")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"album"`)

	if ViewType != nil{
		selector = selector.Where(squirrel.Eq{`view_type`: ViewType})
	}
		
	query, args, err := selector.ToSql()
	if err != nil {
		return nil, err
	}

	var albums []*Album
	_, err = db.Select(&albums, query, args...)
	if err != nil {
		return nil, err
	}

	return albums, nil
}


// Load album by ID
func LoadAlbumFromID(db *gorp.DbMap, ID int64) (*Album, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to list elements")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"album"`).Where(
		squirrel.Eq{`id`: ID},
	)

	query, args, err := selector.ToSql()
	if err != nil {
		return nil, err
	}

	var album Album
	err = db.SelectOne(&album, query, args...)
	if err != nil {
		return nil, err
	}

	return &album, nil
}

// Verify that an album object is valid before creating/updating it
func (a *Album) Valid() error {
	// TODO coherency checks
	return nil
}
