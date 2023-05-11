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
	index     uint8
	color     color.Color
	MoveLines *[][]uint8
	// MoveLinesRight []*[]uint8
	// MoveNames      []*string
	// name *string
	//[Y][2]uint8{LeftX,RightX}
	Move *map[int][2]int
	// Y      int
	// LeftX  int
	// RightX int
}

func (p *VectorPath) Equal(el *VectorPath) bool {
	return p.index == el.index
}
func (p VectorPath) ColorEqual(color color.Color) bool {
	return p.color == color
}

var incIndex uint8 = 0

func NewVectorPath(color color.Color) *VectorPath {
	incIndex = incIndex + 1
	move := make(map[int][2]int)
	return &VectorPath{
		index: incIndex,
		color: color,
		Move:  &move,
	}
}

func (p *VectorPath) GetYVectorXPos(y int) (bool, [2]int) {
	move := *p.Move
	yMove, ok := move[y]
	if ok {
		return ok, yMove
	}
	yMove = [2]int{0, 0}
	move[y] = yMove
	return ok, yMove
}
func (p *VectorPath) AddMove(x int, y int) {
	ok, YVector := p.GetYVectorXPos(y)
	if ok {
		if YVector[0] > x {
			YVector[0] = x
		}
		if YVector[1] < x {
			YVector[1] = x
		}
	} else {
		YVector[0] = x
		YVector[1] = x
	}
}

// func (p *VectorPath) AddMoveLeft(x int, y int) {
// 	ok, YVector := p.GetYVectorXPos()
// 	if !ok {

// 	}
// 	// mve := []uint8{uint8(x), uint8(y)}
// 	// yMove := p.GetYVectorXPos(move[1])
// 	// LeftX := yMove[0]
// 	// if LeftX == 0 || move[0] < LeftX {
// 	// 	yMove[0] = move[0]
// 	// 	move := []uint8{uint8(x), uint8(y)}
// 	// 	moveLines := *p.MoveLines

// 	// 	moveLines = append(moveLines, move)
// 	// 	p.MoveLines = &moveLines
// 	// }

// }
// func (p *VectorPath) AddMoveRight(x int, y int) {
// 	move := []uint8{uint8(x), uint8(y)}
// 	yMove := p.GetYVectorXPos(move[1])
// 	moveLines := *p.MoveLines
// 	RightX := yMove[1]
// 	if move[0] > RightX && RightX != 0 {
// 		yMove[1] = move[0]
// 		moveLines[0] = move
// 	} else {
// 		yMove[1] = move[0]
// 		// p.Y = y

// 		moveLines = append([][]uint8{move}, moveLines...)
// 		p.MoveLines = &moveLines

// 	}
// }

// func (p *VectorPath) AssignMoves(arg2 ...[]uint8) *[][]uint8 {
// 	moveLines := append([][]uint8{}, arg2...)
// 	return &moveLines
// }

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
			isColorLeft := prevOk && previous.ColorEqual(pixelColor)
			// equal := curOk && prevOk && current.Equal(previous)

			if !isColorLeft && !isColorUpper {
				pathShape := NewVectorPath(pixelColor)

				pathShape.AddMove(column, row)
				pathShapes[column] = pathShape
				paths = append(paths, pathShape)
			}
			if isColorUpper {
				current.AddMove(column, row)
			}
			// if isColorLeft && !
			// fmt.Println("COLOR: ", color.RGBA{5, 9, 10, 2} == color.RGBA{5, 9, 10, 2})
			// if !isColorLeft && !isColorUpper {
			// 	pathShape := NewVectorPath(pixelColor)

			// 	pathShape.AddMoveLeft(column, row)
			// 	pathShapes[column] = pathShape
			// 	paths = append(paths, pathShape)

			// }
			// if prevOk && !isColorLeft {
			// 	previous.AddMoveRight(column-1, row)
			// }

			// if curOk && !isColorUpper {
			// 	current.AddMoveLeft(column, row-1)
			// }

			// if isColorLeft && !isColorUpper {
			// 	// if curOk {
			// 	// 	current.AddMoveRight(column, row-1)
			// 	// }
			// 	pathShapes[column] = previous
			// 	previous.AddMoveRight(column, row)

			// }
			// if !isColorLeft && isColorUpper {
			// 	// if prevOk {
			// 	// 	previous.AddMoveRight(column-1, row)
			// 	// }
			// 	current.AddMoveLeft(column, row)
			// }

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
		moveLines := *path.MoveLines
		size := len(moveLines)
		fmt.Println("SIZE: ", size)
		if size > 0 {
			d := ""
			for index, move := range moveLines {

				if index == 0 {
					d += fmt.Sprintf("M%v %v", move[0], move[1])
				} else {
					d += fmt.Sprintf(" L%v %v", move[0], move[1])
				}
			}

			r, g, b, a := path.color.RGBA()

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
