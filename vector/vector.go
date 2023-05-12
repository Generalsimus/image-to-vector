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
	index     *uint8
	color     color.Color
	Move      *map[int]*[2]int
	MoveLines *[]*[2]int
}

func (p *VectorPath) Equal(el *VectorPath) bool {
	return *p.index == *el.index
}
func (p VectorPath) ColorEqual(color color.Color) bool {
	return p.color == color
}

func NewVectorPath(color color.Color, index uint8) *VectorPath {
	// move := make(map[int]*[2]int)
	moveLines := []*[2]int{}
	return &VectorPath{
		index:     &index,
		color:     color,
		MoveLines: &moveLines,
	}
}

func (p *VectorPath) GetYVectorXPos(y int) (bool, *[2]int) {
	move := *p.Move
	yMove, ok := move[y]
	if ok {
		return ok, yMove
	}
	XVectors := [2]int{0, 0}
	yMove = &XVectors
	move[y] = yMove

	return ok, yMove
}
func (p *VectorPath) AddMoveLeft(x int, y int) {
	move := [2]int{x, y}
	moveLines := append(*p.MoveLines, &move)
	p.MoveLines = &moveLines
}
func (p *VectorPath) AddMoveRight(x int, y int) {
	move := [2]int{x, y}
	moveLines := append([]*[2]int{&move}, *p.MoveLines...)
	p.MoveLines = &moveLines
}

func (p *VectorPath) Assign(p2 *VectorPath) {
	moveLines := append(*p2.MoveLines, *p.MoveLines...)

	p.MoveLines = &moveLines
	p2.MoveLines = &moveLines
	// move := *p.Move
	fmt.Println("ASSIGN")
	p2.index = p.index
	// p2.Move = p.Move
	// for y, XVector := range move {
	// 	p.AddMove(XVector[0], y)
	// 	p.AddMove(XVector[1], y)
	// }

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
	var index uint8 = 0
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

			prevOk, previous := PathInclude(pathShapes, column-1)
			curOk, current := PathInclude(pathShapes, column)

			isColorUpper := curOk && current.ColorEqual(pixelColor)
			isColorLeft := prevOk && previous.ColorEqual(pixelColor)
			equal := curOk && prevOk && current.Equal(previous)
			if isColorLeft && isColorUpper && !equal {
				index := *previous.index
				current.Assign(previous)
				paths[index] = nil
			}

			if !isColorLeft && !isColorUpper {
				pathShape := NewVectorPath(pixelColor, index)
				index++
				pathShape.AddMoveLeft(column, row)
				// pathShape.AddMoveRight(column, row)
				pathShapes[column] = pathShape
				paths = append(paths, pathShape)
			}

			if prevOk && !isColorLeft {
				previous.AddMoveRight(column-1, row)
			}

			if curOk && !isColorUpper {
				current.AddMoveLeft(column, row-1)
			}

			if isColorUpper && !isColorLeft {
				current.AddMoveLeft(column, row)
			}
			if isColorLeft && !isColorUpper {
				previous.AddMoveRight(column, row)
				pathShapes[column] = previous

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
		if path == nil {
			continue
		}
		moveLines := *path.MoveLines

		d := ""
		for _, XYPoint := range moveLines {
			// XVectors := move[YVector]
			// if /
			d += fmt.Sprintf("L%v %v ", XYPoint[0], XYPoint[1])

		}
		// fmt.Println("D: ", d)
		r, g, b, a := path.color.RGBA()
		color := fmt.Sprintf("rgba(%v,%v,%v,%v)", r>>8, g>>8, b>>8, a>>8)

		if _, err := f.Write([]byte(fmt.Sprintf("<path fill=\"%v\" d=\"M%vZ\" />\n", color, d[1:]))); err != nil {
			log.Fatal(err)
		}

	}
	if _, err := f.Write([]byte("</svg>")); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

}
