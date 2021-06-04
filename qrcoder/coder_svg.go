package qrcoder

import (
	"bytes"
	"fmt"

	"github.com/aaronarduino/goqrsvg"
	svg "github.com/ajstarks/svgo"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

// BLOCK_SIZE is a default size of pixels in QRCODE
const BLOCK_SIZE uint16 = 5

type CoderSVG struct {
	content   string
	Blocksize uint16
}

func NewCoderSVG(content string) *CoderSVG {
	return &CoderSVG{content: content, Blocksize: BLOCK_SIZE}
}

func (cs CoderSVG) Encode() (raw []byte, err error) {
	var buf bytes.Buffer
	var qrCode barcode.Barcode
	defer buf.Reset()

	if qrCode, err = qr.Encode(cs.content, qr.H, qr.Auto); err != nil {
		return raw, fmt.Errorf("SVG Encode: %w", err)
	}

	s := svg.New(&buf)

	// Write QR code to SVG
	qs := goqrsvg.NewQrSVG(qrCode, 5)
	qs.StartQrSVG(s)

	if err = qs.WriteQrSVG(s); err != nil {
		return raw, fmt.Errorf("SVG Encode: %w", err)
	}

	s.End()
	return buf.Bytes(), nil
}
