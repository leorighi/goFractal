package goFractal

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"math"
	"math/cmplx"
	"os"
	"sync"
)

type Mandelbrot struct {
	Width  int
	Height int
}

func (m *Mandelbrot) GeneratePNG(filename string, focusX float64, focusY float64) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	img := m.generateImage(focusX, focusY)
	png.Encode(f, img)
	return nil
}

func (m *Mandelbrot) GenerateGIF(filename string, frames int, focusX float64, focusY float64, zoomSpeed float64) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	var images []*image.Paletted
	var delays []int

	for frame := 0; frame < frames; frame++ {
		zoom := math.Exp(float64(frame) * zoomSpeed)
		img := m.generateImageWithZoom(zoom, focusX, focusY)
		palImg := image.NewPaletted(img.Bounds(), palette.Plan9)
		draw.FloydSteinberg.Draw(palImg, img.Bounds(), img, image.Point{})

		images = append(images, palImg)
		delays = append(delays, 10)
	}

	return gif.EncodeAll(f, &gif.GIF{
		Image:     images,
		Delay:     delays,
		LoopCount: 0,
	})
}

func (m *Mandelbrot) generateImage(focusX float64, focusY float64) *image.RGBA {
	return m.generateImageWithZoom(1.0, focusX, focusY)
}

func (m *Mandelbrot) generateImageWithZoom(zoom float64, focusX float64, focusY float64) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, m.Width, m.Height))
	var wg sync.WaitGroup

	for py := 0; py < m.Height; py++ {
		y := (float64(py)/float64(m.Height)*(2.5/zoom) - 1.25/zoom) + focusY
		wg.Add(1)

		go func(y float64, py int) {
			row := make([]color.Color, m.Width)
			for px := 0; px < m.Width; px++ {
				x := (float64(px)/float64(m.Width)*(3.5/zoom) - 2.5/zoom) + focusX
				z := complex(x, y)
				row[px] = MandelbrotGenerator(z)
			}
			for px, c := range row {
				img.Set(px, py, c)
			}
			wg.Done()
		}(y, py)
	}

	wg.Wait()
	return img
}

func MandelbrotGenerator(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.Gray{255 - contrast*n}
		}
	}
	return color.Black

}
