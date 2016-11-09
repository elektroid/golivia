package checksum

import (
	"os"
	"io"
	"fmt"
	"crypto/md5"
)


func Md5sum(filePath string) (string, error){
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
