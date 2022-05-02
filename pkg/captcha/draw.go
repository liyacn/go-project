package captcha

import (
	"bytes"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math/rand/v2"
	"os"
	"path"
	"unicode/utf8"
)

const (
	unitSize = 6
	adjust   = unitSize * 2
	fontSize = unitSize * 4
	height   = fontSize * 2
)

type Drawer struct {
	fontList []*opentype.Font
	fontNum  int
}

// NewDrawer fontPath为ttf或otf格式字体文件所在目录
func NewDrawer(fontPath string) *Drawer {
	dirs, err := os.ReadDir(fontPath)
	if err != nil {
		log.Fatal(err)
	}
	var fonts []*opentype.Font
	for _, entry := range dirs {
		if !entry.IsDir() {
			if b, _ := os.ReadFile(path.Join(fontPath, entry.Name())); len(b) > 0 {
				if f, err := opentype.Parse(b); err == nil {
					fonts = append(fonts, f)
				}
			}
		}
	}
	num := len(fonts)
	if num == 0 {
		log.Fatal("no fonts available")
	}
	return &Drawer{
		fontList: fonts,
		fontNum:  num,
	}
}

// Draw s字符长度2~8为宜
func (d *Drawer) Draw(s string) []byte {
	width := fontSize*utf8.RuneCountInString(s) + 3*unitSize
	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	gray := image.NewUniform(color.Gray{Y: uint8(rand.IntN(256))})
	draw.Draw(canvas, canvas.Bounds(), gray, image.Point{}, draw.Src)

	face, _ := opentype.NewFace(d.fontList[rand.IntN(d.fontNum)], &opentype.FaceOptions{
		Size: fontSize,
		DPI:  72,
	})
	c := &font.Drawer{Dst: canvas, Face: face}
	seed := rand.Perm(256)
	pos := 0
	x := unitSize
	for _, char := range s {
		c.Src = image.NewUniform(color.RGBA{
			R: uint8(seed[pos&255]),
			G: uint8(seed[(pos+1)&255]),
			B: uint8(seed[(pos+2)&255]),
		})
		c.Dot = fixed.P(x+rand.IntN(adjust), fontSize+rand.IntN(adjust))
		c.DrawString(string(char))
		x += fontSize
		pos += 3
	}

	for x = unitSize; x < width; x += unitSize {
		x0, y0 := x+rand.IntN(unitSize), rand.IntN(height)
		cl := color.RGBA{
			R: uint8(seed[pos&255]),
			G: uint8(seed[(pos+1)&255]),
			B: uint8(seed[(pos+2)&255]),
		}
		pos += 3
		switch rand.IntN(6) {
		case 0, 2, 4:
			w := rand.IntN(3) + 2
			draw.Draw(canvas, image.Rect(x0-w, y0, x0, y0+6-w), image.NewUniform(cl), image.Point{}, draw.Src)
		case 1, 3, 5:
			canvas.Set(x0, y0, cl)
			canvas.Set(x0-1, y0, cl)
			canvas.Set(x0+1, y0, cl)
			canvas.Set(x0, y0-1, cl)
			canvas.Set(x0, y0+1, cl)
		}
	}

	buf := bytes.NewBuffer(nil)
	_ = jpeg.Encode(buf, canvas, nil)
	return buf.Bytes()
}
