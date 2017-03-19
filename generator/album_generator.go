package generator

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/elektroid/golivia/models"
	"html/template"
	"path/filepath"
	"strings"
)

var BaseTemplate = "./generator/templates/album_template.html"
var GalleriaTemplate = "./generator/templates/galleria_template.html"
var BrowseTemplate = "./generator/templates/browse.html"

func getContent(album interface{}, template_file string) (string, error) {

	funcMap := template.FuncMap{
		"DashToSlash": func(s string) string { return strings.Replace(s, "-", "/", -1) },
	}

	t := template.New("album").Funcs(funcMap) //create a new template
	var err error
	t, err = t.ParseFiles(template_file) //open and parse a template text file
	if err != nil {
		return "", err
	}

	var w bytes.Buffer

	_, file := filepath.Split(template_file)
	if err = t.ExecuteTemplate(&w, file, album); err != nil {
		w.WriteString(err.Error())
		log.Print(err)
		return "", err
	}

	return w.String(), nil
}

func GetAlbumHtml(album *models.Album) (string, error) {
	return getContent(album, BaseTemplate)
}

func GetGalleriaHtml(album *models.Album) (string, error) {
	return getContent(album, GalleriaTemplate)
}

func GetDatesListHtml(dates []*models.Populated) (string, error) {
	return getContent(dates, BrowseTemplate)
}
