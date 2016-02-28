package qrcode

import (
	"github.com/shunde/rsc/qr"
	"image"
)

func Encode(msg string) image.Image {
	code, err := qr.Encode(msg, qr.Q)
	if err != nil {
		return nil
	}
	return code.Image()
}
