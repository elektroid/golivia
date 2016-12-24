package generator

import (
	"bytes"
	"github.com/elektroid/golivia/models"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

var BaseTemplate = "./generator/templates/album_template.html"
var GalleriaTemplate = "./generator/templates/galleria_template.html"

func getContent(album *models.Album, template_file string) (string, error) {
	t := template.New("album") //create a new template
	var err error
	t, err = t.ParseFiles(template_file) //open and parse a template text file
	if err != nil {
		return "", err
	}

	var w bytes.Buffer

	_, file := filepath.Split(template_file)
	if err = t.ExecuteTemplate(&w, file, album); err != nil {
		w.WriteString(err.Error())
		log.Fatal(err)
		return "", err
	}

	return w.String(), nil
}

func GenerateAlbum(album *models.Album) error {
	content, err := getContent(album, BaseTemplate)
	if err != nil {
		return err
	}
	f, err := os.Create("/tmp/album.html")
	if err != nil {
		return err
	}

	f.WriteString(content)

	f.Close()
	return nil
}

func GetAlbumHtml(album *models.Album) (string, error) {
	return getContent(album, BaseTemplate)
}

func GetGalleriaHtml(album *models.Album) (string, error) {
	return getContent(album, GalleriaTemplate)
}
