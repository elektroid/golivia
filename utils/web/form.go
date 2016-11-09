package web

import (
    "github.com/gin-gonic/gin"
    "github.com/elektroid/golivia/constants"
    "os"
    "log"
    "errors"
    "io"
)

func ExtractFileFromForm(c *gin.Context, name string) (string, error ){
	file, header , err := c.Request.FormFile(name)
	if err!=nil{
		log.Print(err)
		return "", err
	}
    filename := header.Filename
    out, err := os.Create(constants.TmpDir+"/"+filename)
    if err != nil {
        log.Fatal(err)
		return "", errors.New("extraction failure, not able to create")
    }
    defer out.Close()
    _, err = io.Copy(out, file)
    if err != nil {
        log.Fatal(err)
	    return "", errors.New("extraction failure, not able to copy content")
    }   
	return constants.TmpDir+"/"+filename, nil
}