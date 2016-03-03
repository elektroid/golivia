package models

import (
	"errors"
	"os"
	"io"
	"crypto/md5"
	"github.com/go-gorp/gorp"
)

type Photo struct {
	ID int64 `json:"id" db:"id"`
	AlbumId int64 `json:"album_id" db:"album_id"`  // no multi album for now
	LocalPath string `json:"local_path" db:"local_path"`
	Md5Sum string `json:"md5sum" db:"md5sum"`
}

// Create a photo
func CreatePhoto(db *gorp.DbMap, AlbumId int64, Path string) (*Photo, error){
	if db == nil {
		return nil, errors.New("Missing db parameter to create photo")
	}
	p := &Photo{
		LocalPath : Path,
		AlbumId: AlbumId,
	}

	err := p.Valid()
	if err != nil {
		return nil, err
	}

	err = db.Insert(p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Photo) SetFile(filePath string) error{
  var result []byte
  file, err := os.Open(filePath)
  if err != nil {
    return err
  }
  defer file.Close()

  hash := md5.New()
  if _, err := io.Copy(hash, file); err != nil {
    return err
  }

  p.Md5Sum=string(hash.Sum(result))
  return nil

}


// Verify that a photo object is valid before creating/updating it
func (p *Photo) Valid() error {
	// TODO coherency checks
	return nil
}