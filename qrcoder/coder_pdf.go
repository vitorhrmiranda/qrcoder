package qrcoder

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type CoderPDF struct {
	CoderPNG
}

func NewCoderPDF(content string, size int) *CoderPDF {
	cp := NewCoderPNG(content, size)
	return &CoderPDF{CoderPNG: *cp}
}

func (cs CoderPDF) Encode() (raw []byte, err error) {
	var tempPNG string = fmt.Sprintf("./_temp-%d.png", time.Now().Unix())

	if raw, err = cs.CoderPNG.Encode(); err != nil {
		return raw, fmt.Errorf("PDF Encode: %w", err)
	}

	// use temp image to add in pdf
	ioutil.WriteFile(tempPNG, raw, 0644)
	defer os.Remove(tempPNG)

	buffer := new(bytes.Buffer)
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	imgOpt := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}
	pdf.ImageOptions(tempPNG, 0, 0, 0, 0, false, imgOpt, 0, "")

	pdf.Output(buffer)
	return buffer.Bytes(), err
}
