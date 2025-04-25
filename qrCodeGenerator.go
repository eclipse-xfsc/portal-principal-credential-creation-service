package main

import (
	"fmt"

	qrcode "github.com/skip2/go-qrcode"
)

func generateQR(data string) ([]byte, error) {
	var png []byte
	png, err := qrcode.Encode(data, qrcode.Medium, 256)
	if err != nil {
		Logger.Error(err)
		err = fmt.Errorf("Error generating QR code")
		return png, err
	}

	return png, nil
}