package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/elektroid/golivia/utils/checksum"
	"github.com/vharitonsky/iniflags"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
)

/*
we should : receive a base directoy
recursively list content and try to upload images
	we compute md5sum
	we try to GET this image from server
	if we get a 404 we upload the file

Nice to have:
	some logs about execution
	local db with fstat info on files to avoid computing data on them

*/

var log = logrus.New()

func VisitFile(fp string, fi os.FileInfo, err error) error {
	if err != nil {
		fmt.Println(err) // can't walk here,
		return nil       // but continue walking elsewhere
	}
	if fi.IsDir() {
		return nil // not a file.  ignore.
	}
	matchedShort, err := filepath.Match("*.jpg", fi.Name())
	if err != nil {
		fmt.Println(err) // malformed pattern
		return err       // this is fatal.
	}
	matchedLong, err := filepath.Match("*.jpeg", fi.Name())
	if err != nil {
		fmt.Println(err) // malformed pattern
		return err       // this is fatal.
	}
	if matchedShort || matchedLong {
		TreatFile(fp, fi)
	}
	return nil
}

func ListFiles(dir string) ([]string, error) {
	var s []string
	return s, nil
}

func TreatFile(file string, fi os.FileInfo) (bool, error) {
	md5sum, err := checksum.Md5sum(file)
	if err != nil {
		log.Printf("Failed to compute checksum for file '%s' : %s\n", file, err.Error())
		return false, err
	}
	log.Printf("treat file %s : %s", file, md5sum)
	present, err := FileIsOnServer(md5sum)
	if err != nil {
		log.Printf("Failed to check file '%s' / '%s' on server: %s\n", file, md5sum, err.Error())
		return false, err
	}

	if present {
		return false, nil
	}

	err = UploadFileToServer(file, md5sum)
	if err != nil {
		log.Printf("Failed to upload file '%s' / '%s' on server: %s\n", file, md5sum, err.Error())
		return false, err
	}
	return true, nil
}

func FileIsOnServer(md5sum string) (bool, error) {

	return checkOnServer(md5sum, *viewerUser, *viewerPassword)
}

func UploadFileToServer(file, md5sum string) error {

	url := fmt.Sprintf("%s%s/%d/photo", *serverUrl, "/album", albumId)
	request, err := newfileUploadRequest(url, nil, "upload", file, *adminUser, *adminPassword)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return err
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.Println(err)
			return err
		}
		resp.Body.Close()
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		fmt.Println(body)
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			log.Println("Error returned by server : ")
			log.Println(resp.Header)
			log.Println(body)
			log.Println(resp.StatusCode)
			return fmt.Errorf("Error returned by server status : %d", resp.StatusCode)
		}
	}
	return nil

}

var client http.Client
var serverUrl *string
var viewerUser *string
var viewerPassword *string
var adminUser *string
var adminPassword *string
var photosBaseDir *string
var logFile *string
var albumId int
var PathSeparator = fmt.Sprintf("%c", os.PathSeparator)

func main() {

	// let's read config
	var (
	//	dataBaseDir   = flag.String("dataBaseDir", "/tmp", "Directory where sqlite db is")
	)
	serverUrl = flag.String("serverUrl", "http://127.0.0.1:8080", "Where to reach server")
	logFile = flag.String("logFile", "/tmp/golivia.log", "Where to store execution log")
	adminUser = flag.String("admin", "admin", "http user with admin permissions")
	adminPassword = flag.String("adminPassword", "adminPassword", "http admin user password")
	viewerUser = flag.String("viewerUser", "amigos", "http regular user name")
	viewerPassword = flag.String("viewerPassword", "usersPassword", "http regular user password")
	photosBaseDir = flag.String("photosBaseDir", "/home/bozo/datolivia", "Directory where photos are stored")
	flag.IntVar(&albumId, "albumId", 1, "album id to post to")

	iniflags.Parse()
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	client.Jar = jar
	f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer f.Close()

	log.Out = f

	log.Println("Client starting with server " + *serverUrl)

	filepath.Walk(*photosBaseDir, VisitFile)

}
