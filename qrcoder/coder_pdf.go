package qrcoder

import (
	"bytes"
	"fmt"
	"image"

	"github.com/signintech/gopdf"
)

type CoderPDF struct {
	CoderPNG
}

func NewCoderPDF(content string, size int) *CoderPDF {
	cp := NewCoderPNG(content, size)
	return &CoderPDF{CoderPNG: *cp}
}

func (cs CoderPDF) Encode() (raw []byte, err error) {
	if raw, err = cs.CoderPNG.Encode(); err != nil {
		return raw, fmt.Errorf("PDF Encode: %w", err)
	}

	r := bytes.NewReader(raw)
	i, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("PDF Encode: %w", err)
	}

	buffer := new(bytes.Buffer)
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	pdf.ImageFrom(i, 0, 0, nil)

	pdf.Write(buffer)
	return buffer.Bytes(), err
}
