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
	index          *uint8
	color          color.Color
	MoveLinesLeft  *[]*[2]int
	MoveLinesRight *[]*[2]int
	PosY           int
}

func (p *VectorPath) Equal(el *VectorPath) bool {
	return *p.index == *el.index
}
func (p *VectorPath) ColorEqual(color color.Color) bool {
	return p.color == color
}

func NewVectorPath(color color.Color, index uint8) *VectorPath {

	moveLinesLeft := []*[2]int{}
	moveLinesRight := []*[2]int{}
	return &VectorPath{
		index:          &index,
		color:          color,
		MoveLinesLeft:  &moveLinesLeft,
		MoveLinesRight: &moveLinesRight,
	}
}

func (p *VectorPath) AddMoveLeft(x int, y int) {
	move := [2]int{x, y}

	moveLinesLeft := append(*p.MoveLinesLeft, &move)
	*p.MoveLinesLeft = moveLinesLeft

	// move := [2]int{x, y}
	// moveLinesLeft := *p.MoveLinesLeft
	// lastIndex := len(moveLinesLeft) - 1

	// if lastIndex > 0 && moveLinesLeft[lastIndex][1] == y && x < moveLinesLeft[lastIndex][0] {
	// 	moveLinesLeft[lastIndex] = &move
	// } else {
	// 	moveLinesLeft = append(moveLinesLeft, &move)
	// }
	// *p.MoveLinesLeft = moveLinesLeft
	// p.PosY = y
}
func (p *VectorPath) AddMoveRight(x int, y int) {
	move := [2]int{x, y}
	moveLinesRight := append(*p.MoveLinesRight, &move)
	*p.MoveLinesRight = moveLinesRight

	// move := [2]int{x, y}
	// moveLinesRight := *p.MoveLinesRight
	// lastIndex := len(moveLinesRight) - 1

	// if lastIndex > 0 && moveLinesRight[lastIndex][1] == y && x > moveLinesRight[lastIndex][0] {
	// 	moveLinesRight[lastIndex] = &move
	// } else {
	// 	moveLinesRight = append(moveLinesRight, &move)
	// }
	// *p.MoveLinesRight = moveLinesRight
	// p.PosY = y
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

			equal := curOk && leftOk && *current == *left
			// current.Equal(left)
			isColorCurrent := curOk && current.ColorEqual(pixelColor)
			isColorLeft := leftOk && left.ColorEqual(pixelColor)
			// equal := curOk && leftOk && *current == *left
			if isColorCurrent && isColorLeft && !equal {
				// fmt.Println("ASSIGN")
				// moveLinesLeft := append(*current.MoveLinesLeft, *left.MoveLinesLeft...)
				// moveLinesLeft = append(moveLinesLeft, *left.MoveLinesRight...)

				// moveLinesRight := append(*current.MoveLinesRight, *left.MoveLinesRight...)
				// *current.MoveLinesLeft = moveLinesLeft
				// *current.MoveLinesRight = moveLinesRight
				// index := *left.index

				// *left.MoveLinesLeft = []*[2]int{}
				// *left.MoveLinesRight = []*[2]int{}

				// *left = *current
				// paths[index] = nil
				// equal = true
			}
			if isColorLeft {
				pathShapes[column] = left
				current = left
				curOk = true
				isColorCurrent = true
				equal = true
			} else if !isColorCurrent {
				current = NewVectorPath(pixelColor, index)
				index++

				pathShapes[column] = current
				paths = append(paths, current)
				curOk = true
				isColorCurrent = true
			}

			if !equal {
				if leftOk {
					left.AddMoveLeft(column, row)
				}
				if curOk {
					current.AddMoveRight(column, row)
				}
			}

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
		if path == nil {
			continue
		}

		moveLinesLeft := *path.MoveLinesLeft
		moveLinesRight := *path.MoveLinesRight
		size := len(moveLinesLeft) + len(moveLinesRight)
		if size < 3 {
			continue
		}

		d := ""
		for _, XYPoint := range moveLinesLeft {
			d = d + fmt.Sprintf("L%v %v ", XYPoint[0], XYPoint[1])
		}

		for _, XYPoint := range moveLinesRight {
			d = fmt.Sprintf("L%v %v ", XYPoint[0], XYPoint[1]) + d
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
