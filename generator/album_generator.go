package generator

import (
	"html/template"
    "bytes"
    "log"
    "path/filepath"
	"github.com/elektroid/golivia/models"
	)


var BaseTemplate="./generator/templates/album_template.html"
var GalleriaTemplate="./generator/templates/galleria_template.html"
var ListTemplate="./generator/templates/by_dates_template.html"

func getAlbumContent(album *models.Album, template_file string) (string, error){
    t := template.New("album") //create a new template
    var err error
    t, err = t.ParseFiles(template_file) //open and parse a template text file
    if err!=nil{
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


func GetAlbumHtml(album *models.Album) ( string, error){
    return getAlbumContent(album, BaseTemplate)
}

func GetGalleriaHtml(album *models.Album) ( string, error){
    return getAlbumContent(album, GalleriaTemplate)
}

func GetDatesListHtml(dates []*models.Populated) ( string, error){
    t := template.New("list") //create a new template
    var err error
    t, err = t.ParseFiles(ListTemplate) //open and parse a template text file
    if err!=nil{
        return "", err
    }
   
    var w bytes.Buffer 

    _, file := filepath.Split(ListTemplate)
    if err = t.ExecuteTemplate(&w, file, dates); err != nil {
        w.WriteString(err.Error())
        log.Fatal(err)
        return "", err
    }
    
    return w.String(), nil
}