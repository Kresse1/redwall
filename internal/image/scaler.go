package image

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"

	xdraw "golang.org/x/image/draw"
)

func Scale(screenW int, screenH int, img io.Reader) ([]byte, error) {
	wallpaper, _, err := image.Decode(img)
	if err != nil {
		return nil, err
	}

	wallpaperRec := wallpaper.Bounds()
	scaledImage := image.NewRGBA(image.Rect(0, 0, screenW, screenH))

	imgW := wallpaperRec.Dx()
	imgH := wallpaperRec.Dy()

	var cropRec image.Rectangle

	if imgW*screenH > imgH*screenW {
		cropW := imgH * screenW / screenH
		offsetX := (imgW - cropW) / 2
		cropRec = image.Rect(offsetX, 0, offsetX+cropW, imgH)
	} else {
		cropH := imgW * screenH / screenW
		offsetY := (imgH - cropH) / 2
		cropRec = image.Rect(0, offsetY, imgW, offsetY+cropH)
	}

	xdraw.CatmullRom.Scale(scaledImage, scaledImage.Bounds(), wallpaper, cropRec, xdraw.Src, nil)

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, scaledImage, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil

}
