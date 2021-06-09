package qrcoder

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"

	p "image/color/palette"

	"github.com/skip2/go-qrcode"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// Dimension limits
const (
	MAX_SIZE int = 1 << 11
	MIN_SIZE int = 1 << 7
)

const logoBase64 = `PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiIHN0YW5kYWxvbmU9Im5vIj8+CjxzdmcKICAgd2lkdGg9IjkwMCIKICAgaGVpZ2h0PSI5MDAiCiAgIHZpZXdCb3g9IjAgMCAyMzguMTI0OTkgMjM4LjEyNSIKICAgdmVyc2lvbj0iMS4xIgogICBpZD0ic3ZnNSIKCXhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIKCXhtbG5zOnN2Zz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgoJPGRlZnMgaWQ9ImRlZnMyIiAvPgoJPGcgaWQ9ImxheWVyMSI+CgkJPGNpcmNsZQogICAgICAgc3R5bGU9ImZpbGw6I2ZmZmZmZjtzdHJva2Utd2lkdGg6MC4yMDEwNzciCiAgICAgICBpZD0icGF0aDEwNDAiCiAgICAgICBjeD0iMTE5LjA2MjUiCiAgICAgICBjeT0iMTE5LjA2MjUiCiAgICAgICByPSIxMTkuMDYyNSIgLz4KCQk8cGF0aAogICAgICAgZD0ibSAxMTMuMDkyMjEsMTkuNTUwNzAzIGMgLTkuMTI4MTIsMC40NDk3OTIgLTE3LjMwMzc0NCwyLjAzNzI5MiAtMjYuNTY0MTY0LDUuMTU5Mzc1IC0yOC43NjAyLDkuNjU3MjkyIC01MS44MzE4NywzMi44ODc3MDYgLTYxLjU2ODU0LDYyLjAxODMzMSAtMy45MTU4MywxMS43NDc1MDMgLTUuNzE1LDI0Ljg3MDgzMSAtNC45NzQxNiwzNi40MzMxMjEgMC4zNzA0MSw1LjkwMDIxIDAuNDQ5NzksNi4zMjM1NCAxLjEzNzcsNi45NTg1NCAwLjU4MjA5LDAuNTI5MTcgMS4wMDU0MiwwLjUyOTE3IDk3LjM2NjY3NCwwLjUyOTE3IGggOTYuODExMDMgbCAwLjYzNSwtMC42MzUgYyAwLjg5OTU5LC0wLjg5OTU5IDEuMDMxODgsLTIuNTQgMS4wNTgzNCwtMTEuOTMyNzEgMCwtOC40OTMxMiAtMC4wNzk0LC05LjQ0NTYyIC0xLjM0OTM4LC0xNy4wNjU2MiAtMi41MTM1NCwtMTQuOTc1NDE2IC04Ljk0MjkxLC0yOS44OTc5MTggLTE4LjE3Njg4LC00Mi4yMDEwNDIgLTAuODk5NTgsLTEuMTY0MTY3IC0xLjk1NzkxLC0yLjU0IC0yLjQwNzcsLTMuMDY5MTY2IC0wLjQ0OTgsLTAuNTAyNzA5IC0xLjM0OTM4LC0xLjU4NzUgLTIuMDEwODMsLTIuMzgxMjUgLTIuMTE2NjcsLTIuNTkyOTE3IC05LjA3NTIxLC05LjMxMzMzMyAtMTIuMzAzMTMsLTExLjkwNjI0OSAtMTIuNDA4OTUsLTEwLjAwMTI1IC0yNy4xMTk3OSwtMTcuMDEyNzA4IC00MS44MDQxNywtMTkuOTQ5NTg0IC01LjU4MjcsLTEuMTExMjQ5IC01LjkwMDIsLTEuMTM3NzA3IC0xMC4wNTQxNiwtMS41NjEwNDEgLTUuMDgsLTAuNTAyNzA4IC0xMC44NDc5MSwtMC42NjE0NTggLTE1Ljc5NTYzLC0wLjM5Njg3NSB6IG0gMTIuMjIzNzYsNDYuODMxMjQ4IGMgMi4yNzU0MSwwLjI5MTA0MiA3LjU5MzU0LDEuNTA4MTI1IDkuNzg5NTgsMi4yNDg5NTggOS4wNzUyMSwzLjA5NTYyNSAxOC4xNTA0MSw5LjQxOTE2NyAyMy44Mzg5NiwxNi42NDIyOSAzLjEyMjA4LDMuOTQyMjkyIDYuMDA2MDQsOC41NzI0OTkgNS42ODg1NCw5LjEwMTY2OSAtMC4xNTg3NSwwLjIzODEyNSAtOTEuODg5Nzk0LDAuMTg1MjA1IC05Mi4xNTQzNzQsLTAuMDUyOTUgLTAuMzQzOTYsLTAuMzcwNDEgMS44NTIwOCwtNC4xMjc0OTcgNC42MzAyLC03Ljg1ODEyNCAyLjk4OTgsLTQuMDIxNjY1IDEwLjI2NTg0LC0xMC42MzYyNDkgMTMuODExMjUsLTEyLjU5NDE2NSAwLjUwMjcxLC0wLjI2NDU4NCAwLjk3ODk2LC0wLjU4MjA4MyAxLjA1ODM0LC0wLjY4NzkxNyAwLjA3OTQsLTAuMTA1ODMzIDEuNzE5NzksLTAuOTUyNSAzLjY3NzcxLC0xLjkwNDk5OSA5LjM5MjcwNCwtNC41NzcyOTIgMTkuMTI5Mzc0LC02LjE2NDc5MiAyOS42NTk3OTQsLTQuODk0NzkyIHoiCiAgICAgICBpZD0icGF0aDIiCiAgICAgICBzdHlsZT0iZmlsbDojNTIyYjMyO2ZpbGwtb3BhY2l0eToxO3N0cm9rZS13aWR0aDowLjAyNjQ1ODQiIC8+CgkJPHBhdGgKICAgICAgIGQ9Im0gNzAuODExNzk2LDE0OC42NjczNiBjIC0wLjM3MDQyLDAuMTU4NzUgLTEuMDg0NzksMC41MjkxNyAtMS41ODc1LDAuODczMTMgLTAuNTAyNywwLjM0Mzk1IC0yLjUxMzU0LDEuNTYxMDQgLTQuNDQ1LDIuNjk4NzQgLTEuOTMxNDYsMS4xNjQxOCAtMy44MzY0NSwyLjI3NTQyIC00LjIzMzMzLDIuNTEzNTUgLTAuMzk2ODcsMC4yMzgxMiAtMS40Mjg3NSwwLjg3MzEyIC0yLjMwMTg3LDEuMzc1ODMgLTAuODczMTMsMC41MjkxNyAtMi4yNDg5NywxLjM0OTM3IC0zLjA0MjcxLDEuODUyMDggLTAuNzkzNzUsMC41MDI3MSAtMi43NTE2NywxLjY2Njg4IC00LjM2NTYzLDIuNTkyOTMgLTEuNTg3NSwwLjg5OTU4IC0yLjk2MzMzLDEuNzQ2MjQgLTMuMDQyNzEsMS44MjU2MiAtMC4wNzk0LDAuMTA1ODMgLTAuNzkzNzUsMC41MjkxNyAtMS41ODc0OSwwLjk3ODk2IC0yLjI0ODk3LDEuMjQzNTQgLTguNDkzMTMsNS4wMDA2MiAtOS41Nzc5Miw1Ljc2NzkxIC0xLjQwMjI5LDAuOTc4OTYgLTEuMzQ5MzgsMS43NzI3MSAwLjIzODEyLDQuMDQ4MTMgMS41MDgxMywyLjE0MzEyIDQuMTAxMDUsNS42ODg1NCA0LjI4NjI1LDUuODczNzQgMC4wNzk0LDAuMDc5NCAwLjYzNSwwLjcxNDM4IDEuMjE3MDksMS40NTUyMSAzLjE3NDk5LDQuMDIxNjcgMTIuNDA4OTYsMTMuMDQzOTYgMTUuMzcyMjksMTUuMDU0OCAwLjQ3NjI1LDAuMjkxMDMgMS43MTk3OSwxLjIxNzA4IDIuNzc4MTIsMi4wMzcyOSAxLjA4NDc5LDAuODIwMjEgMi4zMjgzMywxLjcxOTc5IDIuNzUxNjYsMS45ODQzNyAwLjQ0OTgsMC4yNjQ1OSAwLjg0NjY3LDAuNTU1NjIgMC45MjYwNSwwLjY2MTQ2IDAuMzE3NSwwLjM3MDQxIDYuNDgyMjksMy44ODkzOCA5LjUyNSw1LjQyMzk2IDMuNjc3NzEsMS44MjU2MyA5LjE1NDU4LDQuMjU5NzkgMTAuNTgzMzMsNC42NTY2NiAwLjUwMjcxLDAuMTMyMyAxLjI5NjQ2LDAuNDIzMzMgMS43MTk3OSwwLjYwODU1IDIuNTQsMS4wNTgzMyAxMC4xMzM1NCwzLjA0MjcgMTQuNzYzNzY0LDMuODM2NDYgMS43NzI3LDAuMzE3NDkgMy41MTg5NSwwLjYzNDk5IDMuODg5MzcsMC42ODc5MSA2LjA1ODk1LDEuMDU4MzMgMTkuODk2NjYsMS4wNTgzMyAyNy43ODEyNSwtMC4wMjY1IDMuMDk1NjIsLTAuNDIzMzMgMTEuMzI0MTYsLTIuMjIyNSAxMy44OTA2MiwtMy4wNDI3MSA1Ljc0MTQ2LC0xLjg1MjA5IDcuOTExMDQsLTIuNjE5MzcgMTAuMTg2NDYsLTMuNjUxMjUgNS45NTMxMiwtMi42NDU4NCAxMC43OTUsLTUuMjEyMjkgMTQuODk2MDQsLTcuODg0NTkgMS40Mjg3NSwtMC45MjYwNCAyLjY5ODc1LC0xLjY5MzMyIDIuODMxMDQsLTEuNjkzMzIgMC4xNTg3NSwwIDAuMjY0NTgsLTAuMDc5NCAwLjI2NDU4LC0wLjIxMTY4IDAsLTAuMTA1ODMgMC40NDk3OSwtMC40NzYyNSAxLjAwNTQyLC0wLjg0NjY2IDEuMjE3MDgsLTAuNzkzNzUgNC4yODYyNSwtMy4xNDg1NCA0LjgxNTQyLC0zLjcwNDE2IDAuMjExNjYsLTAuMjExNjggMC45MjYwNCwtMC44MjAyMiAxLjU4NzUsLTEuMzIyOTIgMy4xMjIwOCwtMi4zODEyNSA4LjUxOTU4LC03LjgzMTY3IDEyLjU5NDE2LC0xMi43IDMuNzgzNTQsLTQuNTUwODQgNi40NTU4NCwtOC40OTMxMiA2LjQ1NTg0LC05LjU1MTQ2IDAsLTAuODczMTIgLTAuNjA4NTUsLTEuMzc1ODMgLTMuODM2NDYsLTMuMjgwODMgLTEuNDU1MjEsLTAuODQ2NjcgLTMuMzYwMjEsLTEuOTg0MzggLTQuMjMzMzMsLTIuNTEzNTQgLTEuOTg0MzgsLTEuMjQzNTUgLTUuMTg1ODQsLTMuMTQ4NTUgLTYuNDU1ODQsLTMuODM2NDYgLTAuNTI5MTcsLTAuMjkxMDQgLTEuMDU4MzMsLTAuNjg3OTIgLTEuMTY0MTYsLTAuODczMTIgLTAuMTA1ODQsLTAuMTU4NzYgLTAuMzE3NTEsLTAuMzE3NTEgLTAuNDQ5OCwtMC4zMTc1MSAtMC4xNTg3NCwwIC0yLjExNjY2LC0xLjEzNzcgLTQuMzkyMDgsLTIuNTEzNTQgLTIuMjQ4OTYsLTEuMzc1ODMgLTQuMTgwNDIsLTIuNTEzNTQgLTQuMjU5NzksLTIuNTEzNTQgLTAuMTA1ODMsMCAtMC41NTU2MiwtMC4yOTEwNCAtMS4wMzE4NywtMC42NjE0NiAtMC40NzYyNSwtMC4zNzA0MiAtMC45MjYwNSwtMC42NjE0NiAtMC45Nzg5NiwtMC42NjE0NiAtMC4wNzk0LDAgLTEuOTg0MzgsLTEuMTM3NzEgLTQuMjU5NzksLTIuNTQgLTMuMDE2MjUsLTEuODI1NjIgLTQuMzM5MTcsLTIuNDg3MDggLTQuODQxODgsLTIuNDM0MTcgLTAuNTU1NjIsMC4wNTI5IC0xLjIxNzA5LDAuNzQwODQgLTIuOTEwNDEsMi45NjMzNCAtNC43MzYwNSw2LjI0NDE3IC0xMS41MDkzOCwxMi4xOTcyOSAtMTYuNzQ4MTMsMTQuNzYzNzQgLTAuMzcwNDIsMC4xODUyMiAtMC43MTQzOCwwLjM5Njg4IC0wLjc5Mzc1LDAuNDc2MjUgLTAuMjkxMDQsMC4zMTc1MSAtNC43ODg5NiwyLjQwNzcxIC02Ljg3OTE3LDMuMTc1MDEgLTEuMjQzNTQsMC40NDk3OSAtMi40ODcwOCwwLjkyNjA0IC0yLjc3ODEyLDEuMDU4MzMgLTAuMjkxMDUsMC4xMzIyOSAtMS4xMTEyNSwwLjM3MDQyIC0xLjg1MjA5LDAuNTI5MTYgLTAuNzE0MzcsMC4xNTg3NiAtMS43OTkxNywwLjQyMzM0IC0yLjM4MTI1LDAuNTU1NjMgLTEuMTkwNjIsMC4yOTEwNCAtMi44ODM5NSwwLjYwODU1IC01LjQyMzk1LDEuMDA1NDIgLTIuMjQ4OTYsMC4zNzA0MSAtMTEuODAwNDIsMC4zNzA0MSAtMTQuMTU1MjEsMCAtMy42NTEyNSwtMC41NTU2MyAtNC41MjQzNywtMC43MTQzOCAtNi4wODU0MiwtMS4xNjQxNyAtMC44NzMxMiwtMC4yNjQ1OCAtMS44Nzg1NCwtMC41MDI3MSAtMi4yNDg5NSwtMC41ODIwOCAtMC4zNzA0MiwtMC4wNTMgLTAuODk5NTksLTAuMjExNjcgLTEuMTkwNjMsLTAuMzQzOTYgLTAuMjkxMDQsLTAuMTMyMjkgLTEuNTM0NTksLTAuNjA4NTQgLTIuNzUxNjY0LC0xLjA1ODMzIC01LjAyNzA5LC0xLjg3ODU0IC0xMC4yMzkzOCwtNC44NjgzMyAtMTQuNzEwODQsLTguNDY2NjcgLTIuMzI4MzMsLTEuODUyMDggLTYuMjE3NzEsLTUuNDIzOTYgLTYuMjE3NzEsLTUuNzE1IDAsLTAuMDc5NCAtMC4yMTE2NiwtMC4zNDM5NSAtMC40NDk3OSwtMC41ODIwOCAtMC44NzMxMywtMC44NzMxMyAtNC4zMTI3MSwtNS4zMTgxMyAtNC4zMTI3MSwtNS41ODI3MSAwLC0wLjI2NDU4IC0xLjI3LC0xLjM0OTM3IC0xLjU2MTA0LC0xLjMyMjkxIC0wLjA3OTQsMCAtMC40NDk3OSwwLjEzMjI5IC0wLjgyMDIxLDAuMjY0NTcgeiIKICAgICAgIGlkPSJwYXRoNCIKICAgICAgIHN0eWxlPSJmaWxsOiNlMzU4Nzk7ZmlsbC1vcGFjaXR5OjE7c3Ryb2tlLXdpZHRoOjAuMDI2NDU4NCIgLz4KCTwvZz4KPC9zdmc+Cg==`

type CoderPNG struct {
	content string
	size    int
	UseLogo bool
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
	code, err := qrcode.New(cs.content, qrcode.Highest)
	if err != nil {
		return raw, fmt.Errorf("PNG Encode: %w", err)
	}

	var buf bytes.Buffer
	img := code.Image(cs.size)

	if cs.UseLogo {
		cs.addLogo(img)
	}

	if err = png.Encode(&buf, img); err != nil {
		return raw, fmt.Errorf("PNG Encode: %w", err)
	}

	return buf.Bytes(), nil
}

func (cs CoderPNG) addLogo(qrcode image.Image) (err error) {

	//define color in image
	var palette = qrcode.(*image.Paletted)
	palette.Palette = color.Palette{
		color.RGBA{0xff, 0xff, 0xff, 0xff},
		color.RGBA{0x0, 0x0, 0x0, 0xff},
		color.RGBA{0x52, 0x2b, 0x32, 0xff},
		color.RGBA{0xe3, 0x58, 0x79, 0xff},
	}
	palette.Palette = append(palette.Palette, p.WebSafe...)

	logo, err := cs.readLogoFromCSV()

	// draw logo in qrcode
	offset := qrcode.Bounds().Max.X/2 - logo.Bounds().Max.X/2

	for x := 0; x < logo.Bounds().Max.X; x++ {
		for y := 0; y < logo.Bounds().Max.Y; y++ {
			if _, _, _, alpha := logo.At(x, y).RGBA(); alpha == 0 {
				continue
			}
			palette.Set(x+offset, y+offset, logo.At(x, y))
		}
	}

	return err
}

func (cs CoderPNG) readLogoFromCSV() (img image.Image, err error) {
	l, err := base64.StdEncoding.DecodeString(logoBase64)
	if err != nil {
		return
	}
	buf := bytes.NewBuffer(l)

	var size = cs.size / 6

	icon, err := oksvg.ReadIconStream(buf)
	if err != nil {
		return
	}

	icon.SetTarget(0, 0, float64(size), float64(size))
	rgba := image.NewRGBA(image.Rect(0, 0, size, size))
	scanner := rasterx.NewScannerGV(size, size, rgba, rgba.Bounds())
	dasher := rasterx.NewDasher(size, size, scanner)

	icon.Draw(dasher, 1)
	return rgba, err
}
