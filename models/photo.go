package models

import (
	"errors"
	"os"
	"io"
	"fmt"
	"crypto/md5"
	"github.com/go-gorp/gorp"
	"github.com/elektroid/golivia/imgResize"
	"github.com/elektroid/golivia/constants"
	"github.com/Masterminds/squirrel"
	"github.com/elektroid/golivia/utils/sqlgenerator"
	"github.com/rwcarlsen/goexif/exif"
	//"github.com/go-sql-driver/mysql"

	"path/filepath"
	"time"
	"log"
)

type Photo struct {
	ID int64 `json:"id" db:"id"`
	AlbumId int64 `json:"album_id" db:"album_id"`  // no multi album for now
	LocalPath string `json:"local_path" db:"local_path"`
	Description string `json:"description" db:"description"`
	Md5Sum string `json:"md5sum" db:"md5sum"`
	//Time mysql.NullTime `json:"time" db:"time"`
	Time *time.Time `json:"time" db:"time"`
}
var PathSeparator = fmt.Sprintf("%c", os.PathSeparator)

// Create a photo
func CreatePhoto(db *gorp.DbMap, A *Album, Path string) (*Photo, error){
	if db == nil {
		return nil, errors.New("Missing db parameter to create photo")
	}
	p := &Photo{
		AlbumId: A.ID,
	}

	md5sum, err := md5sum(Path)
	if err != nil{
		return nil, err
	}
 	p.Md5Sum = md5sum

 	d, err := getExifDate(Path)
 	if err == nil {
 		//p.Time=mysql.NullTime{Time: *d, Valid: true}
 		p.Time=d
 	}else{
 		log.Print(err)
 	}

	existing, err := LoadPhotoFromMd5(db, p.Md5Sum)
	if err==nil{
		return nil, fmt.Errorf("Image already exists under album %d ", existing.AlbumId)		
	}


	newPath := constants.PhotosDir+PathSeparator+md5sum+filepath.Ext(Path)
	p.LocalPath=md5sum+filepath.Ext(Path)
	err = os.Rename(Path, newPath)
	if err != nil{
		return nil, err
	}

	err = p.Valid()
	if err != nil {
		return nil, err
	}

	err = p.setMini(A.MiniatureWidth, A.MiniatureHeight, constants.MiniatureQuality, constants.PhotosDir+PathSeparator+constants.MiniSubDir, newPath)
	if err != nil {
		return nil, err		
	}

	err = db.Insert(p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Photo) setMini(targetWidth uint, targetHeight uint, quality int, miniDirPath string, originalPath string) error{
		var miniP=miniDirPath+PathSeparator+p.LocalPath
		err:=imgResize.MakeMini(targetWidth, targetHeight, quality, originalPath, miniP)
		if err!=nil{
			return err
		}
		return nil
}


func getExifDate(filePath string) (*time.Time, error){

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	x, err := exif.Decode(f)
    if err != nil {
        return nil, err
    }
    t, err := x.DateTime()
    if err != nil{
    	return nil, err
    }
    return &t, err
}

func md5sum(filePath string) (string, error){
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// Load album by ID
func LoadPhotoFromMd5(db *gorp.DbMap, Md5Sum string) (*Photo, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to list elements")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"photo"`).Where(
		squirrel.Eq{`md5sum`: Md5Sum},
	)

	query, args, err := selector.ToSql()
	if err != nil {
		return nil, err
	}

	var photo Photo
	err = db.SelectOne(&photo, query, args...)
	if err != nil {
		return nil, err
	}

	return &photo, nil
}

// Verify that a photo object is valid before creating/updating it
func (p *Photo) Valid() error {
	// TODO coherency checks
	return nil
}
