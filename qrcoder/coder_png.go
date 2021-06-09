package qrcoder

import (
	"fmt"

	"github.com/skip2/go-qrcode"
)

// Dimension limits
const (
	MAX_SIZE int = 1 << 11
	MIN_SIZE int = 1 << 7
)

type CoderPNG struct {
	content string
	size    int
}

func NewCoderPNG(content string, size int) *CoderPNG {

	if size > MAX_SIZE {
		size = MAX_SIZE
	}
	if size < MIN_SIZE {
		size = MIN_SIZE
	}

	return &CoderPNG{content: content, size: size}
}

func (cs CoderPNG) Encode() (raw []byte, err error) {
	if raw, err = qrcode.Encode(cs.content, qrcode.Highest, cs.size); err != nil {
		return raw, fmt.Errorf("PNG Encode: %w", err)
	}
	return
}
