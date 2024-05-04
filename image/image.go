package image

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

const (
	width      int = 35
	height     int = 35
	marginHor  int = 5
	marginVert int = 5
)

var (
	yellow_correctPlace = color.RGBA{249, 234, 0, 255}
	white_anotherPlace  = color.RGBA{255, 255, 255, 255}
	grey_notPresent     = color.RGBA{95, 95, 95, 255}
	background          = color.RGBA{34, 33, 35, 255}
	// border                  = color.RGBA{103, 107, 48, 255}

	hardLetter = [...]string{"ж", "җ", "м", "ш", "щ", "ы", "ю"}
)

type Painter struct {
	fontfile string
}

func New(fontfile string) (*Painter, error) {
	// check if file exists
	if _, err := os.Stat(fontfile); os.IsNotExist(err) {
		return nil, fmt.Errorf("font file does not exist: %s", fontfile)
	}
	return &Painter{
		fontfile: fontfile,
	}, nil
}

func (p Painter) RectOnly(answer string, words []string) []byte {
	return p.createImage(answer, words, false)
}

func (p Painter) FullImage(answer string, words []string) []byte {
	return p.createImage(answer, words, true)
}

func (p Painter) createImage(answer string, words []string, withLetters bool) []byte {

	myimage := image.NewRGBA(image.Rect(0, 0, marginHor+5*(width+marginHor), marginVert+6*(height+marginVert)))
	draw.Draw(myimage, myimage.Bounds(), &image.Uniform{background}, image.Point{}, draw.Src)
	for line := 0; line != len(words); line++ {
		responce := []rune(words[line])
		for i := 0; i != 5; i++ {
			colorRect := definceColor([]rune(answer), responce, i)
			if !withLetters {
				p.drawRectOnly(myimage, colorRect, marginHor+i*(marginHor+width), marginVert+line*(marginVert+height))
			}
			if withLetters {
				p.drawRectWithLetter(myimage, string(responce[i]), colorRect, marginHor+i*(marginHor+width), marginVert+line*(marginVert+height))
			}
		}
	}

	for n := len(words); n < 6; n++ {
		p.drawEmptyLine(myimage, marginHor, marginVert+n*(marginVert+height))
	}

	buf := new(bytes.Buffer)
	png.Encode(buf, myimage)
	send_s3 := buf.Bytes()

	return send_s3
}

func (p Painter) drawRectOnly(myimage *image.RGBA, colorRect, x, y int) image.Rectangle {
	red_rect := image.Rect(x, y, x+width, y+height)
	var rectColor color.RGBA
	switch colorRect {
	case 0:
		rectColor = grey_notPresent
	case 1:
		rectColor = white_anotherPlace
	case 2:
		rectColor = yellow_correctPlace
	}
	draw.Draw(myimage, red_rect, &image.Uniform{rectColor}, image.Point{}, draw.Src)
	return red_rect
}

func (p Painter) drawRectWithLetter(myimage *image.RGBA, letter string, colorRect, x, y int) {
	red_rect := p.drawRectOnly(myimage, colorRect, x, y)
	fontColor := image.Black
	if colorRect == 0 {
		fontColor = image.White
	}

	var fonSize float64 = 18
	c := freetype.NewContext()
	c.SetClip(red_rect)
	c.SetDPI(90)
	c.SetDst(myimage)
	c.SetFontSize(fonSize)
	c.SetSrc(fontColor)
	fontBytes, _ := os.ReadFile(p.fontfile)
	f, _ := freetype.ParseFont(fontBytes)
	c.SetHinting(font.HintingFull)
	c.SetFont(f)

	leftEdge := x + width/3
	if isHardLetter(letter) {
		leftEdge = x + width/4
	}

	pt := freetype.Pt(leftEdge, y+height/2+int(c.PointToFixed(fonSize)>>6)/2)

	_, err := c.DrawString(strings.ToUpper(letter), pt)
	if err != nil {
		log.Println(err)
		return
	}
}

func (p Painter) drawEmptyLine(myimage *image.RGBA, x, y int) {
	rectColor := grey_notPresent
	for i := 0; i != 5; i++ {
		red_rect := image.Rect(x, y, x+width, y+height)
		draw.Draw(myimage, red_rect, &image.Uniform{rectColor}, image.Point{}, draw.Src)
		x += marginHor + width
	}
}

func isHardLetter(s string) bool {
	for i := 0; i < len(hardLetter); i++ {
		if s == hardLetter[i] {
			return true
		}
	}
	return false
}

func definceColor(riddle []rune, responce []rune, i int) int {
	if riddle[i] == responce[i] {
		return 2
	}
	for n := 0; n != 5; n++ {
		if riddle[n] == responce[i] && riddle[n] != responce[n] {
			return 1
		}
	}
	return 0
}
