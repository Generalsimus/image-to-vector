package vector

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"sort"
)

type VectorPath struct {
	id    *uint8
	color color.Color
	Move  *map[int]*[2]int
}

func (p *VectorPath) Equal(el *VectorPath) bool {
	return *p.id == *el.id
}
func (p VectorPath) ColorEqual(color color.Color) bool {
	return p.color == color
}

var incIndex uint8 = 0

func NewVectorPath(color color.Color) *VectorPath {
	incIndex = incIndex + 1
	move := make(map[int]*[2]int)
	// fmt.Println("IDDDD", incIndex)
	return &VectorPath{
		id:    &incIndex,
		color: color,
		Move:  &move,
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
func (p *VectorPath) AddMove(x int, y int) {
	ok, YVector := p.GetYVectorXPos(y)
	// fmt.Println("YVector: ", ok, YVector)
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
	// fmt.Println("YVector: ", ok, YVector)
}

func (p *VectorPath) Assign(p2 *VectorPath) {
	move := *p2.Move
	for y, XVector := range move {
		p.AddMove(XVector[0], y)
		p.AddMove(XVector[1], y)
	}
	p2.id = p.id
	p2.Move = p.Move

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

func (v VectorImage) ImageVector() (image.Image, map[uint8]*VectorPath) {
	bounds := v.Img.Bounds()
	widget := bounds.Max.X
	height := bounds.Max.Y
	img := v.Img
	colorDiffNum := float64(255 * v.ColorDiffPercent)
	paths := map[uint8]*VectorPath{}
	// []*VectorPath{}

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

			prevOk, previous := PathInclude(pathShapes, column-1)
			curOk, current := PathInclude(pathShapes, column)

			isColorUpper := curOk && current.ColorEqual(pixelColor)
			isColorLeft := prevOk && previous.ColorEqual(pixelColor)
			equal := curOk && prevOk && current.Equal(previous)

			if isColorLeft && isColorUpper && !equal {
				// delete(paths, previous.id)
				// current.Assign(previous)
			}

			if !isColorLeft && !isColorUpper {
				pathShape := NewVectorPath(pixelColor)

				pathShape.AddMove(column, row)
				pathShapes[column] = pathShape
				paths[*pathShape.id] = pathShape
			}

			if prevOk && !isColorLeft {
				previous.AddMove(column-1, row)
			}

			if curOk && !isColorUpper {
				current.AddMove(column, row-1)
			}

			if isColorUpper {
				current.AddMove(column, row)
			}
			if isColorLeft {
				previous.AddMove(column, row)
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

func (v VectorImage) SavePathsToSVGFile(paths map[uint8]*VectorPath, fileName string) {
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
		move := *path.Move
		// size := len(move)
		keys := make([]int, 0, len(move))

		for k := range move {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		d := ""

		// d1 := ""
		// d2 := ""
		for _, YVector := range keys {
			XVectors := move[YVector]
			// d1 = d1 + fmt.Sprintf("L%v %v ", XVectors[0], YVector)
			// d2 = fmt.Sprintf("L%v %v ", XVectors[1], YVector) + d2
			d = fmt.Sprintf("L%v %v ", XVectors[1], YVector) + d + fmt.Sprintf("L%v %v ", XVectors[0], YVector)

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
