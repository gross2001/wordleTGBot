package imgwordle

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"strings"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

var (
	width               int = 35
	height              int = 35
	marginHor           int = 5
	marginVert          int = 5
	yellow_correctPlace     = color.RGBA{249, 234, 0, 255}
	white_anotherPlace      = color.RGBA{255, 255, 255, 255}
	grey_notPresent         = color.RGBA{95, 95, 95, 255}
	background              = color.RGBA{34, 33, 35, 255}
	border                  = color.RGBA{103, 107, 48, 255}

	hardLetter = [...]string{"ж", "җ", "м", "ш", "щ", "ы", "ю"}
)

func isHardLetter(s string) bool {
	for i := 0; i < len(hardLetter); i++ {
		if s == hardLetter[i] {
			return true
		}
	}
	return false
}

func definceColor(answer []rune, responce rune, i int) int {
	if answer[i] == responce {
		return 2
	}
	for n := 0; n != 5; n++ {
		if answer[n] == responce {
			return 1
		}
	}
	return 0
}

func drawRectWithLetter(myimage *image.RGBA, letter string, colorRect, x, y int) {
	red_rect := image.Rect(x, y, x+width, y+height)
	var rectColor color.RGBA
	fontColor := image.Black
	switch {
	case colorRect == 0:
		rectColor = grey_notPresent
		fontColor = image.White
	case colorRect == 1:
		rectColor = white_anotherPlace
	case colorRect == 2:
		rectColor = yellow_correctPlace
	}
	draw.Draw(myimage, red_rect, &image.Uniform{rectColor}, image.Point{}, draw.Src)

	var fonSize float64 = 18
	c := freetype.NewContext()
	c.SetClip(red_rect)
	c.SetDPI(90)
	c.SetDst(myimage)
	c.SetFontSize(fonSize)
	c.SetSrc(fontColor)
	fontfile := "./Roboto-Black.ttf"
	fontBytes, _ := ioutil.ReadFile(fontfile)
	f, _ := freetype.ParseFont(fontBytes)
	c.SetHinting(font.HintingFull)
	c.SetFont(f)

	//	pt := freetype.Pt(x+width/3, y+int(c.PointToFixed(12)>>6)+height/3)
	// pt - left lower point

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

func drawEmptyLine(myimage *image.RGBA, x, y int) {
	rectColor := grey_notPresent
	for i := 0; i != 5; i++ {
		red_rect := image.Rect(x, y, x+width, y+height)
		draw.Draw(myimage, red_rect, &image.Uniform{rectColor}, image.Point{}, draw.Src)
		x += marginHor + width
	}
}

func CreateImage(answer string, words []string) []byte {

	myimage := image.NewRGBA(image.Rect(0, 0, marginHor+5*(width+marginHor), marginVert+6*(height+marginVert)))

	draw.Draw(myimage, myimage.Bounds(), &image.Uniform{background}, image.Point{}, draw.Src)

	for line := 0; line != len(words); line++ {
		runes := []rune(words[line])
		for i := 0; i != 5; i++ {
			colorRect := definceColor([]rune(answer), runes[i], i)
			drawRectWithLetter(myimage, string(runes[i]), colorRect, marginHor+i*(marginHor+width), marginVert+line*(marginVert+height))
		}
	}

	for n := len(words); n < 6; n++ {
		drawEmptyLine(myimage, marginHor, marginVert+n*(marginVert+height))
	}

	buf := new(bytes.Buffer)
	png.Encode(buf, myimage)
	send_s3 := buf.Bytes()

	return send_s3
}
