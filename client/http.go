package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string, user, pass string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}
	//req.Header.Add("Authorization", "Basic "+basicAuth(user, pass))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.SetBasicAuth(user, pass)
	return req, err
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func checkOnServer(md5sum string, user, pass string) (bool, error) {

	req, err := http.NewRequest("GET", *serverUrl+"/photo/"+md5sum, nil)
	if err != nil {
		log.Println("Failed to build GET request:", err)
		return false, err
	}
	req.Header.Add("Authorization", "Basic "+basicAuth(user, pass))
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to execute GET request:", err)
		return false, err
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			return true, nil
		}
		if resp.StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, fmt.Errorf("Invalid response for GET : %d", resp.StatusCode)
	}

}
