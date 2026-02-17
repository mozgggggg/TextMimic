// pixellist.go
package pixellist

import (
	"image"
	"image/color"
	_ "image/color"
	_ "image/png"
)

func Getlist(img image.Image) [][][]int {

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	//Making a matrix with [Pixel Count] arrays with lengths of 3
	matrix := make([][][]int, height)
	for y := 0; y < height; y++ {
		matrix[y] = make([][]int, width)
		for x := 0; x < width; x++ {
			matrix[y][x] = make([]int, 4)
			r, g, b, a := img.At(x, y).RGBA()
			matrix[y][x][0] = int(r)
			matrix[y][x][1] = int(g)
			matrix[y][x][2] = int(b)
			matrix[y][x][3] = int(a)
		}
	}
	// Getting Colors of ALL Pixels in the image and sending them to matrix

	return matrix
}

func Savelist(matrix [][][]int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, len(matrix[0]), len(matrix)))

	for y := 0; y < len(matrix); y++ {
		for x := 0; x < len(matrix[0]); x++ {
			img.Set(x, y, color.RGBA64{R: uint16(matrix[y][x][0]), G: uint16(matrix[y][x][1]), B: uint16(matrix[y][x][2]), A: uint16(matrix[y][x][3])})
		}
	}
	return img
}
