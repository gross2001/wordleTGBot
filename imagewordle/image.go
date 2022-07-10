package imgwordle

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

/*
var (
	dpi      float64 = 72
	fontfile string  = "./luxisr.ttf"
	size     float64 = 12
)
*/
var (
	width      int = 30
	height     int = 30
	marginHor  int = 10
	marginVert int = 10
)

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
	var myred color.RGBA
	var fontColor *image.Uniform
	switch {
	case colorRect == 0:
		myred = color.RGBA{128, 128, 128, 255}
		fontColor = image.White
	case colorRect == 1:
		myred = color.RGBA{0, 255, 255, 255}
		fontColor = image.Black
	case colorRect == 2:
		myred = color.RGBA{255, 255, 0, 255}
		fontColor = image.Black

	}
	draw.Draw(myimage, red_rect, &image.Uniform{myred}, image.Point{}, draw.Src)

	c := freetype.NewContext()
	c.SetClip(red_rect)
	c.SetDPI(72)
	c.SetDst(myimage)
	c.SetFontSize(12)
	c.SetSrc(fontColor)
	fontfile := "./Roboto-Black.ttf"
	fontBytes, _ := ioutil.ReadFile(fontfile)
	f, _ := freetype.ParseFont(fontBytes)
	c.SetHinting(font.HintingNone)
	c.SetFont(f)

	pt := freetype.Pt(x+width/3, y+int(c.PointToFixed(12)>>6)+height/3)
	_, err := c.DrawString(letter, pt)
	if err != nil {
		log.Println(err)
		return
	}

}

func CreateImage(answer string, words []string) []byte {

	myimage := image.NewRGBA(image.Rect(0, 0, 220, len(words)*(height+marginVert)+marginVert))
	mygreen := color.RGBA{0, 100, 0, 255} //  R, G, B, Alpha
	// backfill entire background surface with color mygreen
	draw.Draw(myimage, myimage.Bounds(), &image.Uniform{mygreen}, image.Point{}, draw.Src)

	for line := 0; line != len(words); line++ {
		runes := []rune(words[line])
		for i := 0; i != 5; i++ {
			colorRect := definceColor([]rune(answer), runes[i], i)
			drawRectWithLetter(myimage, string(runes[i]), colorRect, marginHor+i*(marginHor+width), marginVert+line*(marginVert+height))
		}
	}

	buf := new(bytes.Buffer)
	png.Encode(buf, myimage)
	send_s3 := buf.Bytes()

	return send_s3
}
