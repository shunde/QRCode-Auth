package qr

import (
	"github.com/qpliu/qrencode-go/qrencode"
	"image"
)

func Encode(msg string) image.Image {
	grid, err := qrencode.Encode(msg, qrencode.ECLevelH)
	if err != nil {
		return nil
	}
	return grid.Image(8)
}
