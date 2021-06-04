package qrcode

type QRCoder interface {
	// Encode a QR Code and return a raw file
	Encode() ([]byte, error)
}
