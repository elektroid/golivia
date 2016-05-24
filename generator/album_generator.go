package generator

import (
	"html/template"
    "bufio"
    "os"
    "log"
	"github.com/elektroid/golivia/models"
	)



func GenerateAlbum(album *models.Album) error{
	
    t := template.New("album") //create a new template
    var err error
    t, err = t.ParseFiles("../generator/album_template.html") //open and parse a template text file
    if err!=nil{
     	return err
    }
    f, err := os.Create("/tmp/album.html")
    if err!=nil{
 	  return err
    }

    w := bufio.NewWriter(f)

    if err = t.ExecuteTemplate(w, "album_template.html", album); err != nil {
        w.WriteString(err.Error())
        log.Fatal(err)
        //return err
    }
    w.Flush()
    f.Close()


    return nil

}