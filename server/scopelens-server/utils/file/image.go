package file

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"scopelens-server/config"
)

// base64
func DecodeBase64AndSave(b string) (path string, err error) {
	// check
	if len(b) < 10 {
		return "", errors.New("Empty base64 string. ")
	}
	var ext string
	if b[11] == 'j' {
		b = b[23:]
		ext = ".jpg"
	} else if b[11] == 'p' {
		b = b[22:]
		ext = ".png"
	} else if b[11] == 'g' {
		b = b[22:]
		ext = ".gif"
	} else {
		return "", errors.New("Not supported base64 image format. ")
	}

	// decode
	data, err := base64.StdEncoding.DecodeString(b)
	if err != nil {
		return "", err
	}

	// check dir
	fileDir := fmt.Sprintf("./%s", config.Path.ImageSavePath)
	if err = MkDir(fileDir); err != nil {
		return "", err
	}

	// rename
	fileName := Rename(b[:32] + ext)
	filePathStr := fileDir + fileName

	// save
	f, _ := os.OpenFile(filePathStr, os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer f.Close()
	n, err := f.Write(data)
	if err != nil {
		return "", err
	} else if n == 0 {
		return "", errors.New("Wrote an empty file. ")
	}

	return filePathStr, nil
}