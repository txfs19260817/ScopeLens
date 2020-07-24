package file

import (
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/image/draw"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
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

// resize image
// https://text.baldanders.info/golang/resize-image/
func Rescale(filePath string) (err error) {
	//open original image file
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return
	}
	defer file.Close()

	//decode image
	imgSrc, ext, err := image.Decode(file)
	if err != nil {
		return
	}

	//rectange of image
	rctSrc := imgSrc.Bounds()

	//rescale to 1280 * 720
	imgDst := image.NewRGBA(image.Rect(0, 0, 1280, 720))
	draw.CatmullRom.Scale(imgDst, imgDst.Bounds(), imgSrc, rctSrc, draw.Over, nil)

	//create resized image file
	dst, err := os.Create(filePath) //dst file path is same as src
	if err != nil {
		return
	}
	defer dst.Close()

	//encode resized image
	switch ext {
	case "jpeg":
		if err = jpeg.Encode(dst, imgDst, &jpeg.Options{Quality: 80}); err != nil {
			return
		}
	case "gif":
		if err = gif.Encode(dst, imgDst, nil); err != nil {
			return
		}
	case "png":
		enc := png.Encoder{CompressionLevel: png.BestCompression}
		if err = enc.Encode(dst, imgDst); err != nil {
			return
		}
	default:
		err = errors.New("format error")
	}
	return nil
}