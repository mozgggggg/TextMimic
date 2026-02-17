package main

import (
	"bufio"
	"fmt"
	_ "image/color"
	"image/draw"

	"image/png"
	"os"
	"strconv"
	"strings"

	"image"
	_ "image/color"
	_ "image/draw"
	_ "image/png"
	"math/rand"
	"slices"
	pxllist "textmimic/pixellist"
)

func main() {
	var userinputraw string

	scanner := bufio.NewScanner(os.Stdin)

	// prompt the user for input
	fmt.Print("Enter maximum width and height, e.g. 2000 2000: ")

	if scanner.Scan() {
		userinputraw = scanner.Text()
	}

	rawcoords := strings.Fields(userinputraw)

	var maxwidth int
	var maxheight int

	if len(rawcoords) == 2 {
		maxwidth, _ = strconv.Atoi(rawcoords[0])
		maxheight, _ = strconv.Atoi(rawcoords[1])
		if maxheight <= 0 || maxwidth <= 0 {
			fmt.Print("Impossible Coordinates")
		}
	} else {
		fmt.Print("Impossible Coordinates")
	}

	fmt.Print("Enter Text Line Count: ")
	if scanner.Scan() {
		userinputraw = scanner.Text()
	}

	linecount, err := strconv.Atoi(userinputraw)
	if err != nil {
		fmt.Println(err)
	}
	pxperline := maxheight / linecount

	matrix := make([][][]int, maxheight)
	for y := 0; y < maxheight; y++ {
		matrix[y] = make([][]int, maxwidth)
		for x := 0; x < maxwidth; x++ {
			matrix[y][x] = make([]int, 4)
		}
	}

	rawimg, err := os.Open("Lettersheet.png")
	if err != nil {
		panic(err)
	}
	defer rawimg.Close()
	img, _, err := image.Decode(rawimg)
	if err != nil {
		panic(err)
	}

	rawlettersheet := pxllist.Getlist(img)

	var separatorColor = []int{65535, 0, 0, 65535}
	//the whole lettersheet parsed into individual images, separated by separatorcolor

	lastsplit, letterrow := 0, 0
	var lettersheet [][]image.Image
	for y := 0; y < len(rawlettersheet); y++ {
		if slices.Equal(rawlettersheet[y][0], separatorColor) {
			lettersheetstrips := rawlettersheet[lastsplit+1 : y] //finding the strip
			lastsplit = y
			lastsplit2 := 0
			lettersheet = append(lettersheet, make([]image.Image, 0))
			for x := 0; x < len(lettersheetstrips[0]); x++ {
				if slices.Equal(lettersheetstrips[0][x], separatorColor) {
					Letterimage := make([][][]int, len(lettersheetstrips))
					for i, lett := range lettersheetstrips {
						Letterimage[i] = lett[lastsplit2+1 : x]
					}
					//getting a singular image from the strips
					if len(Letterimage) <= pxperline {
						lettersheet[letterrow] = append(lettersheet[letterrow], pxllist.Savelist(Letterimage))
					} else {
						lettersheet[letterrow] = append(lettersheet[letterrow], pxllist.Savelist(Letterimage[len(Letterimage)-pxperline:]))
					}
					//appending it in its row alongside its copies
					lastsplit2 = x
				}
			}
			letterrow++
		}
	}
	fmt.Println("Lettersheet Parsed Successfully")
	var usertext string
	fmt.Println("Enter Text:")
	if scanner.Scan() {
		usertext = scanner.Text()
	}
	words := strings.Fields(usertext)

	writerx, writery := 0, 0

	canvas := image.NewRGBA(image.Rect(0, 0, maxwidth, maxheight))

	wordlistlen := len(words)
	for i := 0; i < wordlistlen; i++ {
		word := words[i]
		if i < wordlistlen-1 {
			word = word + " "
		}
		// converting a word into an image containing that word (assembling from letters from lettersheet)
		wordimg0 := image.NewRGBA(image.Rect(0, 0, maxwidth, pxperline))
		wordlen := 0
		for l, ascii := range word {
			var imagebuffer image.Image
			switch {
			case ascii > 64 && ascii < 91:
				imagebuffer = lettersheet[ascii-65][rand.Intn(5)] //rand.Intn(5)
			case ascii > 96 && ascii < 123:
				imagebuffer = lettersheet[ascii-97][rand.Intn(15)+6] //rand.Intn(15)+6
			case ascii == 32:
				imagebuffer = lettersheet[rand.Intn(len(lettersheet))][5]
			default:
				fmt.Println("Symbol ["+string(rune(ascii))+"] is not accepted, ascii code:", int(ascii))
				imagebuffer = lettersheet[rand.Intn(len(lettersheet))][5]
			}
			boundaries := imagebuffer.Bounds()
			if wordlen+boundaries.Dx() <= maxwidth {
				draw.Draw(
					wordimg0,
					image.Rect(wordlen, pxperline-boundaries.Dy(), wordlen+boundaries.Dx(), pxperline),
					imagebuffer,
					image.Point{},
					draw.Over,
				)
				wordlen += boundaries.Dx()
			} else {
				fmt.Println("Word [" + word + "] is too long, it will be split") //use slice tricks and make it overlap to the next line if the word doesnt fit
				words = append(words, "ERRORWORD")
				copy(words[i+1:], words[i+1:])
				words[i+1] = word[l:]
				wordlistlen++
				break
			}
		}

		wordimg := wordimg0.SubImage(image.Rect(0, 0, wordlen, pxperline))

		if wordlen > maxwidth-writerx {
			if writery+pxperline < maxheight {
				writerx = 0
				writery = writery + pxperline
			} else {
				fmt.Println("reached the end of the given coordinates")
				break
			}
		}

		draw.Draw(
			canvas,
			image.Rect(writerx, writery+pxperline-wordimg.Bounds().Dy(), writerx+wordlen, writery+pxperline*2-wordimg.Bounds().Dy()), // where to draw
			wordimg,
			image.Point{},
			draw.Over,
		)
		writerx = writerx + wordlen

	}

	result, err := os.Create("result.png")
	if err != nil {
		panic(err)
	}
	defer result.Close()
	if err := png.Encode(result, canvas); err != nil {
		panic(err)
	}

	fmt.Println("Successfully Edited Image!")

}
