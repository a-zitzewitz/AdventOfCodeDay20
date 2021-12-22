package main

import "testing"

func Test_enhanceImage(t *testing.T) {
	image, err := loadImage("test-image.txt")
	if err != nil {
		t.Fatal(err)
	}
	algo, err := initAlgorithmFromFile("test-algo.txt")
	if err != nil {
		t.Fatal(err)
	}
	image.print()
	if image.countPixels() != 10 {
		t.Errorf("Expected 10 pixels, but got %d", image.countPixels())
	}
	image = image.enhanceImage(algo)
	image.print()
	if image.countPixels() != 24 {
		t.Errorf("Expected 24 pixels, but got %d", image.countPixels())
	}
	image = image.enhanceImage(algo)
	image.print()
	if image.countPixels() != 35 {
		t.Errorf("Expected 35 pixels, but got %d", image.countPixels())
	}
}

func Test_50Times(t *testing.T) {
	image, err := loadImage("test-image.txt")
	if err != nil {
		t.Fatal(err)
	}
	algo, err := initAlgorithmFromFile("test-algo.txt")
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 50; i++ {
		image = image.enhanceImage(algo)
	}
	if image.countPixels() != 3351 {
		t.Errorf("Expected 3351 pixels, but got %d", image.countPixels())
	}
}
