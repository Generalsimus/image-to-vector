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
	color     color.Color
	MoveLines []*[]uint8
	MoveNames []*string
	name      *string
	X         int
	Y         int
}

func (p VectorPath) Equal(el VectorPath) bool {
	return p.color == el.color
}
func (p VectorPath) ColorEqual(color color.Color) bool {
	// fmt.Println("EQ:", p.color, color, p.color == color)
	return p.color == color
}
func NewVectorPath(color color.Color) VectorPath {
	name := ""
	return VectorPath{
		color:     color,
		MoveLines: []*[]uint8{},
		// MoveNames: []*string{},
		name: &name,
		// MoveLines: &[]uint8{},
		X: 0,
		Y: 0,
	}
}

func (p *VectorPath) AddMove(x int, y int) {
	p.X = x
	p.Y = y
	// p.endX = x

	// fmt.Printf("X: %v, Y: %v \n", x, y)
	moveLines := []uint8{uint8(x), uint8(y)}
	p.MoveLines = append(p.MoveLines, &moveLines)
	// fmt.Println("EEE: ", p.MoveLines)
	// p.MoveLines = &moveLines
	// p.pathMoveNames = append(p.pathMoveNames, &name)
	// p.name = &name
	// return p
}

func NewVectorPaths(count int) []VectorPath {
	vectors := []VectorPath{}
	for i := 0; i < count; i++ {
		vectors = append(vectors, NewVectorPath(color.RGBA{0, 0, 0, 0}))
	}
	return vectors
}
func PathInclude(vectorPaths []VectorPath, index1 int, index2 int) (bool, VectorPath, VectorPath) {
	count := len(vectorPaths)

	if index1 > -1 && index1 < count && index2 > -1 && index2 < count {
		return true, vectorPaths[index1], vectorPaths[index2]

	}
	var element1 VectorPath
	var element2 VectorPath
	return false, element1, element2
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

func (v VectorImage) ImageVector() (image.Image, []VectorPath) {
	bounds := v.Img.Bounds()
	widget := bounds.Max.X
	height := bounds.Max.Y
	img := v.Img
	colorDiffNum := float64(255 * v.ColorDiffPercent)
	paths := []VectorPath{}

	newImage := image.NewRGBA(image.Rect(0, 0, widget, height))
	pathShapes := NewVectorPaths(widget)
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
			ok, element1, element2 := PathInclude(pathShapes, column, column-1)
			//
			fmt.Println("ELMN1: ", column, element1, element2, ok && element1.Equal(element2), element1.ColorEqual(pixelColor))
			if ok && element1.Equal(element2) {
				if !element1.ColorEqual(pixelColor) {
					// 	continue
					// } else {
					// fmt.Println("ELMN1: ", element1.MoveLines)
					element1.AddMove(column, row)
					// fmt.Println("ELMN2: ", element1.MoveLines)
					// go func() {
					// }()
				}

			} else {
				fmt.Println("NEWWWWWWWWWWWWWWWWWWWWWWW: ")

				pathShape := NewVectorPath(pixelColor)

				pathShape.AddMove(column, row)
				paths = append(paths, pathShape)
				pathShapes[column] = pathShape
				// fmt.Println("NEWWWWWWWWWWWWWWWWWWWWWWW: ", paths)
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

func (v VectorImage) SavePathsToSVGFile(paths []VectorPath, fileName string) {
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

	if _, err := f.Write([]byte(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%vpx" height="%vpx" viewBox="0 0 %v %v">`, widget, height, widget, height))); err != nil {
		log.Fatal(err)
	}
	// <svg height="210" width="400">
	// fmt.Println("paths: ", paths)</svg></svg>

	for _, path := range paths {
		size := len(path.MoveLines)
		if size > 1 {
			fmt.Println("EEEE: ", size)
		}
		if size > 0 {
			d := ""
			for index, moveAddr := range path.MoveLines {
				move := *moveAddr
				if index == 0 {
					d += fmt.Sprintf("M %v %v", move[0], move[1])
				} else {
					d += fmt.Sprintf("L %v %v", move[0], move[1])
				}
			}
			if _, err := f.Write([]byte(fmt.Sprintf(`<path d="%v Z" />`, d))); err != nil {
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
