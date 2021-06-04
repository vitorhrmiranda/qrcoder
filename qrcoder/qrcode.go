package qrcoder

import "fmt"

type QRCoder interface {
	// Encode a QR Code and return a raw file
	Encode() ([]byte, error)
}

// Mimetypes enableds
const (
	MIME_SVG string = "image/svg+xml"
	MIME_PNG string = "image/png"
	MIME_PDF string = "application/pdf"
)

// Allowed extensions
const (
	SVG string = "svg"
	PNG string = "png"
	PDF string = "pdf"
)

// CreateQRCode is a factory to gerenate QRCode
func CreateQRCode(content string, extension string, size int) (qrcode []byte, mimetype string, err error) {
	switch extension {
	case SVG:
		qrcode, err = NewCoderSVG(content).Encode()
		mimetype = MIME_SVG

	case PNG:
		qrcode, err = NewCoderPNG(content, size).Encode()
		mimetype = MIME_PNG

	case PDF:
		qrcode, err = NewCoderPDF(content, size).Encode()
		mimetype = MIME_PDF

	default:
		err = fmt.Errorf("extension not recognized")
	}

	return qrcode, mimetype, err
}
