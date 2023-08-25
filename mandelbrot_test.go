package github.com/leorighi/goFractal

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
)

func TestGeneratePNG(t *testing.T) {
	m := &Mandelbrot{Width: 200, Height: 200}
	filename := "test_output.png"
	err := m.GeneratePNG(filename, 0.0, 0.0)
	if err != nil {
		t.Errorf("Failed to generate PNG: %v", err)
	}

	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		t.Errorf("PNG file was not created")
	}
	defer file.Close()

	actualFilename := filepath.Base(file.Name())
	if actualFilename != filename {
		t.Errorf("File name does not match: got %v, want %v", actualFilename, filename)
	}

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		t.Errorf("Failed to decode PNG: %v", err)
	}
	if config.Width != 200 || config.Height != 200 {
		t.Errorf("PNG dimensions are incorrect: got %dx%d, want 200x200", config.Width, config.Height)
	}

	os.Remove(filename)
}

func TestGenerateGIF(t *testing.T) {
	m := &Mandelbrot{Width: 200, Height: 200}
	filename := "test_output.gif"
	err := m.GenerateGIF(filename, 10, 0.0, 0.0, 0.1)
	if err != nil {
		t.Errorf("Failed to generate GIF: %v", err)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("GIF file was not created")
	}
	os.Remove(filename)
}

func TestMandelbrotGenerator(t *testing.T) {
	inSet := MandelbrotGenerator(0 + 0i)
	if inSet != (color.Black) {
		t.Errorf("Expected black color, got: %v", inSet)
	}
	outSet := MandelbrotGenerator(2 + 2i)
	if outSet == (color.Black) {
		t.Errorf("Expected color, got black")
	}
}

func TestGeneratePNG_Failure(t *testing.T) {
	m := &Mandelbrot{Width: 200, Height: 200}
	filename := "////fail_file.png"

	err := m.GeneratePNG(filename, 0.0, 0.0)
	if err == nil {
		t.Errorf("Expected to fail file creation, got success")
	} else if !os.IsPermission(err) {
		t.Errorf("Expected a permission error, got: %v", err)
	}
}

func TestGenerateGif_Failure(t *testing.T) {
	m := &Mandelbrot{Width: 200, Height: 200}
	filename := "////fail_file.gif"

	err := m.GenerateGIF(filename, 10, 0.0, 0.0, 0.1)
	if err == nil {
		t.Errorf("Expected to fail file creation, got success")
	} else if !os.IsPermission(err) {
		t.Errorf("Expected a permission error, got: %v", err)
	}
}
