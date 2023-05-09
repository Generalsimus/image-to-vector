package vector

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
)

type VectorPath struct {
	color          color.Color
	MoveLinesLeft  []*[]uint8
	MoveLinesRight []*[]uint8
	// MoveNames      []*string
	// name *string
	// X    int
	// Y    int
}

func (p VectorPath) Equal(el *VectorPath) bool {
	return p.color == el.color
}
func (p VectorPath) ColorEqual(color color.Color) bool {
	// fmt.Println("EQ:", p.color, color, p.color == color)
	return p.color == color
}
func NewVectorPath(color color.Color) VectorPath {
	// name := ""
	return VectorPath{
		color:          color,
		MoveLinesLeft:  []*[]uint8{},
		MoveLinesRight: []*[]uint8{},
		// MoveNames: []*string{},
		// name: &name,
		// MoveLines: &[]uint8{},
		// X: 0,
		// Y: 0,
	}
}

func (p *VectorPath) AddMoveLeft(x int, y int) {
	moveLines := []uint8{uint8(x), uint8(y)}
	p.MoveLinesLeft = append(p.MoveLinesLeft, &moveLines)
}
func (p *VectorPath) AddMoveRight(x int, y int) {
	moveLines := []uint8{uint8(x), uint8(y)}
	p.MoveLinesRight = append(p.MoveLinesRight, &moveLines)
}

func NewVectorPaths(count int) []*VectorPath {
	vectors := []*VectorPath{}
	for i := 0; i < count; i++ {
		path := NewVectorPath(color.RGBA{0, 0, 0, 0})
		vectors = append(vectors, &path)
	}
	return vectors
}
func PathInclude(vectorPaths []*VectorPath, index int) (bool, *VectorPath) {
	count := len(vectorPaths)

	if index > -1 && index < count {
		el := vectorPaths[index]

		return el != nil, el
	}

	return false, nil
}

type VectorImage struct {
	ColorDiffPercent float64
	Img              image.Image
}

// func (v VectorImage) RGBToNum(r uint8, g uint8, b uint8) uint8 {
// 	return (r << 16) + (g << 8) + (b)
// }
// func (v VectorImage) NumToRGB(num int) (uint8, uint8, uint8) {

//		return uint8((num & 0xff0000) >> 16), uint8((num & 0x00ff00) >> 8), uint8(num & 0x0000ff)
//	}
func isColorEqual(color1 color.Color, color2 color.Color) bool {
	r1, g1, b1, a1 := color1.RGBA()
	r2, g2, b2, a2 := color2.RGBA()
	return r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2
}

//	func PixelColorIsLonely(img image.Image, x, y  ) {
//		if x
//
// mainPixel:= := img.At(x, y)
// }

func (v VectorImage) ImageVector() (image.Image, []*VectorPath) {
	bounds := v.Img.Bounds()
	widget := bounds.Max.X
	height := bounds.Max.Y
	img := v.Img
	colorDiffNum := float64(255 * v.ColorDiffPercent)
	paths := []*VectorPath{}

	newImage := image.NewRGBA(image.Rect(0, 0, widget, height))
	pathShapes := make([]*VectorPath, widget)
	// NewVectorPaths(widget)
	// pathShapes := []VectorPath{}
	for row := 0; row < height; row++ {
		for column := 0; column < widget; column++ {
			r, g, b, a := img.At(column, row).RGBA()

			red := math.Round(float64((r>>8))/colorDiffNum) * colorDiffNum
			green := math.Round(float64((g>>8))/colorDiffNum) * colorDiffNum
			blue := math.Round(float64((b>>8))/colorDiffNum) * colorDiffNum

			pixelColor := color.RGBA{
				uint8(red),
				uint8(green),
				uint8(blue),
				uint8(a),
			}
			newImage.Set(column, row, pixelColor)

			//
			prevOk, previous := PathInclude(pathShapes, column-1)
			curOk, current := PathInclude(pathShapes, column)
			// current := pathShapes[column]
			isColorUpper := curOk && current.ColorEqual(pixelColor)
			isColorLeft := prevOk && curOk && previous.ColorEqual(current.color)

			if !isColorLeft && isColorUpper {
				current.AddMoveLeft(column, row)

			}
			if !isColorUpper && isColorLeft {
				current.AddMoveRight(column, row)

			}
			fmt.Println("pixelColor: ", pixelColor, column, isColorUpper, isColorLeft)
			if !isColorUpper && !isColorLeft {
				pathShape := NewVectorPath(pixelColor)
				pathShapeAddr := &pathShape

				pathShapes[column] = pathShapeAddr
				paths = append(paths, pathShapeAddr)

				pathShapeAddr.AddMoveLeft(column, row)
				// fmt.Println("pathShapeAddr: ", pixelColor, pathShapeAddr.color)
				if prevOk {
					previous.AddMoveRight(column-1, row)
				}
			}

		}
	}
	// fmt.Println("ENDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD")
	return newImage, paths
}

///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////

func (v VectorImage) SavePathsToSVGFile(paths []*VectorPath, fileName string) {
	os.Remove(fileName)

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	bounds := v.Img.Bounds()
	widget := bounds.Max.X
	height := bounds.Max.Y
	// if _, err := f.Write([]byte(fmt.Printf(""))); err != nil {
	// 	log.Fatal(err)
	// }
	// viewBox="0 0 %v %v"
	if _, err := f.Write([]byte(fmt.Sprintf(
		"<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"%vpx\" height=\"%vpx\">\n",
		widget,
		height,
	))); err != nil {
		log.Fatal(err)
	}
	// <svg height="210" width="400">
	// fmt.Println("paths: ", paths)</svg></svg>

	for _, path := range paths {
		size := len(path.MoveLinesLeft)
		// if size > 1 {
		// 	fmt.Println("EEEE: ", path.color)
		// }
		if size > 0 {
			// fmt.Println("EEEE: ", path.color)
			d := ""
			for index, moveAddr := range path.MoveLinesLeft {
				move := *moveAddr
				if index == 0 {
					d += fmt.Sprintf("M%v %v", move[0], move[1])
				} else {
					d += fmt.Sprintf(" L%v %v", move[0], move[1])
				}
			}
			for _, moveAddr := range path.MoveLinesRight {
				move := *moveAddr
				d += fmt.Sprintf(" L%v %v", move[0], move[1])
			}
			r, g, b, a := path.color.RGBA()
			// color.NRGBA
			if _, err := f.Write([]byte(fmt.Sprintf("<path fill=\"%v\" d=\"%v Z\" />\n", fmt.Sprintf("rgba(%v,%v,%v,%v)", r>>8, g>>8, b>>8, a>>8), d))); err != nil {
				log.Fatal(err)
			}

		}
	}
	if _, err := f.Write([]byte("</svg>")); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

}
