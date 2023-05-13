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
	index *uint8
	color color.Color
	// Move      *map[int]*[2]int
	MoveLines *[]*[2]int
	// X         int
	// MaxY int/
}

func (p *VectorPath) Equal(el *VectorPath) bool {
	return *p.index == *el.index
}
func (p *VectorPath) ColorEqual(color color.Color) bool {
	return p.color == color
}

func NewVectorPath(color color.Color, index uint8) *VectorPath {

	moveLines := []*[2]int{}
	return &VectorPath{
		index:     &index,
		color:     color,
		MoveLines: &moveLines,
	}
}

func (p *VectorPath) AddMoveLeft(x int, y int) {
	move := [2]int{x, y}
	// fmt.Println(*p)
	moveLines := append([]*[2]int{&move}, *p.MoveLines...)
	*p.MoveLines = moveLines
}
func (p *VectorPath) AddMoveRight(x int, y int) {

	move := [2]int{x, y}
	moveLines := append(*p.MoveLines, &move)
	*p.MoveLines = moveLines
	// // count := len(moveLines)
	// indexAt := len(moveLines) - p.MaxY - y - 1
	// // fmt.Println("R", indexAt, p.Y-y)
	// if indexAt > 0 && moveLines[indexAt][0] < move[0] {
	// 	moveLines[indexAt] = &move
	// } else {

	// }
	// *p.MoveLines = moveLines

	// // p.X = x
	// if y > p.MaxY {
	// 	p.MaxY = y
	// }

}

// func (p *VectorPath) AddMove(x int, y int) {
// 	minusX = p.X - x
// 	moveLines := *p.MoveLines
// 	count := len(moveLines)

// 	if count > 0 && moveLines[count-][] < move[0] {

// 	}
// 	// move := [2]int{x, y}
// 	// moveLines := *p.MoveLines
// 	// count := len(moveLines)
// 	// if count > 0 && moveLines[0][0] < move[0] {
// 	// 	moveLines[0] = &move
// 	// } else {
// 	// 	moveLines = append([]*[2]int{&move}, *p.MoveLines...)
// 	// }
// 	*p.MoveLines = moveLines
// }

// func (p *VectorPath) SetIndex(index uint8) {
// 	p.index = index

// }

func (p *VectorPath) Assign(p2 *VectorPath) {
	moveLines := append(*p2.MoveLines, *p.MoveLines...)

	// moveLinesAddr := &moveLines
	*p.MoveLines = moveLines
	*p2.MoveLines = moveLines
	// move := *p.Move

	// fmt.Println("ASSIGN")
	// p2.SetIndex(p.index)
	// index := *p.index
	// indexAddr := &index
	// p.index = indexAddr
	// p2.index = indexAddr
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

			leftOk, left := PathInclude(pathShapes, column-1)
			curOk, current := PathInclude(pathShapes, column)

			isColorCurrent := curOk && current.ColorEqual(pixelColor)
			isColorLeft := leftOk && left.ColorEqual(pixelColor)
			equal := curOk && leftOk && current.Equal(left)

			if !curOk || !current.ColorEqual(pixelColor) {
				current = NewVectorPath(pixelColor, index)
				index++

				pathShapes[column] = current
				paths = append(paths, current)
				curOk = true
			} else if isColorCurrent && isColorLeft && !equal {
				// index := left.index

				currentLines := current.MoveLines
				leftLines := left.MoveLines

				moveLines := append(append([]*[2]int{}, *leftLines...), *currentLines...)

				*current.MoveLines = moveLines
				*left.MoveLines = moveLines
				// *currentLines = moveLines
				// *leftLines = moveLines
				*left = *current
				// *left.MoveLines = moveLines
				*left.index = *current.index
				// paths[index] = nil
				equal = true
				// fmt.Println("ASSSIGN", len(moveLines))
			}

			if !equal {
				if leftOk {
					left.AddMoveRight(column, row)
				}
				if curOk {
					current.AddMoveLeft(column, row)
				}
			}
			// if  isColorCurrent
			// if isColorLeft && isColorUpper && !equal {
			// 	index := previous.index
			// 	current.Assign(previous)
			// 	paths[index] = nil

			// }

			// if !isColorLeft && !isColorUpper {
			// 	pathShape := NewVectorPath(pixelColor, index)
			// 	index++
			// 	pathShape.AddMoveLeft(column, row)

			// 	pathShapes[column] = pathShape
			// 	paths = append(paths, pathShape)
			// }
			// if isColorLeft && !equal {
			// 	continue
			// }

			// if isColorLeft && isColorUpper && !equal {
			// 	index := left.index
			// 	moveLines := append(*left.MoveLines, *upper.MoveLines...)
			// 	*upper.MoveLines = moveLines
			// 	*left.MoveLines = moveLines
			// 	left.index = upper.index
			// 	paths[index] = nil
			// }

			// if isColorLeft {
			// 	// previous.AddMoveLeft(column, row)
			// 	pathShapes[column] = left
			// }
			//  else {
			// 	previous.AddMoveRight(column-1, row)
			// }
			// if prevOk && !isColorLeft {
			// 	previous.AddMoveRight(column-1, row)
			// }

			// if curOk && !isColorUpper {
			// 	current.AddMoveLeft(column, row-1)
			// }

			// if isColorUpper  {
			// 	current.AddMoveLeft(column, row)
			// }
			// if isColorLeft   {
			// 	previous.AddMoveRight(column, row)
			// 	pathShapes[column] = previous

			// }

		}
	}
	// fmt.Println("ENDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD")
	fmt.Println("ENDDDDDDDDDDDDDDDDDDDDDDDD")
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

		if path == nil || len(moveLines) == 0 {
			continue
		}

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
