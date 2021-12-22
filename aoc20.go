package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

func readLinesFromFile(fileName string, cap int) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines = make([]string, 0, cap)

	for scanner.Scan() {
		var line = scanner.Text()
		if len(line) > 0 {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

type Algorithm struct {
	mapping [512]byte
}

func (algo *Algorithm) getMappingFor(index int) byte {
	return algo.mapping[index]
}

func initAlgorithmFromFile(fileName string) (*Algorithm, error) {
	lines, err := readLinesFromFile(fileName, 1)
	if err != nil {
		return nil, err
	}
	line := lines[0]
	result := new(Algorithm)
	for i, c := range line {
		if c == '#' {
			result.mapping[i] = 1
		} else {
			result.mapping[i] = 0
		}
	}
	return result, nil
}

type Image struct {
	x, y         int // Coordinates of origin
	w, h         int // Width and height
	defaultPixel byte
	image        [][]byte
}

func (image *Image) getPixelAt(x int, y int) byte {
	x -= image.x
	y -= image.y
	if x < 0 || x >= image.w {
		return image.defaultPixel
	}
	if y < 0 || y >= image.h {
		return image.defaultPixel
	}
	return image.image[y][x]
}

func (image *Image) getSurroundingPixels(x int, y int) int {
	result := 0
	result += int(image.getPixelAt(x-1, y-1)) * 256
	result += int(image.getPixelAt(x, y-1)) * 128
	result += int(image.getPixelAt(x+1, y-1)) * 64
	result += int(image.getPixelAt(x-1, y)) * 32
	result += int(image.getPixelAt(x, y)) * 16
	result += int(image.getPixelAt(x+1, y)) * 8
	result += int(image.getPixelAt(x-1, y+1)) * 4
	result += int(image.getPixelAt(x, y+1)) * 2
	result += int(image.getPixelAt(x+1, y+1)) * 1
	return result
}

func (image *Image) countPixels() int {
	xLimit := image.x + image.w
	yLimit := image.y + image.h
	result := 0
	for y := image.y; y < yLimit; y++ {
		for x := image.x; x < xLimit; x++ {
			if image.getPixelAt(x, y) == 1 {
				result++
			}
		}
	}
	return result
}

func loadImage(fileName string) (*Image, error) {
	lines, err := readLinesFromFile(fileName, 100)
	if err != nil {
		return nil, err
	}
	resultImage := new(Image)
	image := make([][]byte, 0, 100)
	for _, line := range lines {
		imageLine := make([]byte, 0, 100)
		for _, c := range line {
			if c == '#' {
				imageLine = append(imageLine, 1)
			} else if c == '.' {
				imageLine = append(imageLine, 0)
			}
		}
		image = append(image, imageLine)
	}
	resultImage.x = 0
	resultImage.y = 0
	resultImage.w = len(lines[0])
	resultImage.h = len(image)
	resultImage.defaultPixel = 0
	resultImage.image = image
	return resultImage, nil
}

func (image *Image) print() {
	for y := 0; y < image.h; y++ {
		for x := 0; x < image.w; x++ {
			if image.image[y][x] == 1 {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
}

func (image *Image) enhanceImage(algo *Algorithm) *Image {
	result := new(Image)
	result.x = image.x - 2
	result.y = image.y - 2
	result.w = image.w + 4
	result.h = image.h + 4
	result.image = make([][]byte, result.h)
	if image.defaultPixel == 0 {
		result.defaultPixel = algo.getMappingFor(0)
	} else {
		result.defaultPixel = algo.getMappingFor(511)
	}

	var wg sync.WaitGroup

	wg.Add(result.h)
	yLimit := result.y + result.h
	for y := result.y; y < yLimit; y++ {
		go func(y int) {
			line := make([]byte, 0, result.w)
			limit := result.x + result.w
			for x := result.x; x < limit; x++ {
				index := image.getSurroundingPixels(x, y)
				line = append(line, algo.getMappingFor(index))
			}
			result.image[y-result.y] = line
			wg.Done()
		}(y)
	}
	wg.Wait()
	return result
}

func main() {
	image, err := loadImage("image.txt")
	if err != nil {
		log.Fatal(err)
	}
	algo, err := initAlgorithmFromFile("algo.txt")
	if err != nil {
		log.Fatal(err)
	}
	image1 := image.enhanceImage(algo)
	image2 := image1.enhanceImage(algo)
	log.Printf("Pixels in twice enhanced image: %d\n", image2.countPixels())
	for i := 2; i < 50; i++ {
		image2 = image2.enhanceImage(algo)
	}
	log.Printf("Pixels in 50 times enhanced image: %d\n", image2.countPixels())
}
