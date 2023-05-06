package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"vectoral/vector"
)

func main() {
	dirPath := "./images/"
	entries, _ := os.ReadDir(dirPath)

	for _, file := range entries {
		fileName := file.Name()
		img, _ := getImageFromFilePath(dirPath + fileName)
		vectorImg := vector.VectorImage{
			ColorDiffPercent: 0.3,
			Img:              img,
		}
		// vectorImg.ImageVector()
		fmt.Println("FILENAME: ", fileName)
		saveImageAt(vectorImg.ImageVector(), "./save/"+fileName)
	}
	fmt.Println("PROCESS END")

}

func saveImageAt(image image.Image, path string) {
	var imageBuf bytes.Buffer
	png.Encode(&imageBuf, image)

	// Write to file.
	outfile, err := os.Create(path)
	if err != nil {
		// replace this with real error handling
		panic(err.Error())
	}
	defer outfile.Close()
	png.Encode(outfile, image)
}

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _ := png.Decode(f)
	return image, err
}
