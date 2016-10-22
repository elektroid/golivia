package generator

import (
	"html/template"
    "os"
    "bytes"
    "log"
	"github.com/elektroid/golivia/models"
	)


func getContent(album *models.Album) (string, error){
   t := template.New("album") //create a new template
    var err error
    t, err = t.ParseFiles("./generator/templates/album_template.html") //open and parse a template text file
    if err!=nil{
        return "", err
    }
   
    var w bytes.Buffer 

    if err = t.ExecuteTemplate(&w, "album_template.html", album); err != nil {
        w.WriteString(err.Error())
        log.Fatal(err)
        return "", err
    }
    
    return w.String(), nil
}

func GenerateAlbum(album *models.Album) error{
	content, err := getContent(album)
    if err != nil{
        return err
    }
    f, err := os.Create("/tmp/album.html")
    if err!=nil{
      return err
    }

    f.WriteString(content)

    f.Close()
    return nil
}

func GetAlbumHtml(album *models.Album) ( string, error){
    return getContent(album)
}