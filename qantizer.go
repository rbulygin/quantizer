package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
)

const using = `
quant [--cfg <config>] <src> [<dst>] [<level> [<level> [...]] 
<config> is an optional config file with levels divided by space
<src> source file name (JPEG, PNG, GIF is accepted)
<dst> optional destination file, by default <src>_out.png 
<level is one of quantize levels, by default used 64 128 and 192 to quantize on 4 colors 
`

type level struct {
	Level int `json:"level"`
	Color int `json:"color"`
}

type config struct {
	Left   int     `json:"left"`
	Levels []level `json:"levels"`
	Right  int     `json:"right"`
}

func (cfg *config) convert(in color.Color) color.Color {
	r, g, b, _ := in.RGBA()
	y := uint8((299*r + 587*g + 114*b) / 1000)
	for i := range cfg.Levels {
		if y < uint8(cfg.Levels[i].Level) {
			y = uint8(cfg.Levels[i].Color)
			return color.Gray{Y: y}
		}
	}
	return color.Gray{Y: 255}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf(using)
		return
	}

	cfg := config{
		Left: 0,
		Levels: []level{
			{Level: 50, Color: 0},
			{Level: 100, Color: 50},
			{Level: 150, Color: 100},
		},
		Right: 256,
	}

	if len(os.Args) > 1 {
		img, ext, err := getImageFromFilePath(os.Args[1])
		if err != nil {
			fmt.Printf("Error on start: %v\n", err)
			return
		}
		fmt.Printf("Image extension: %v\n", ext)
		if img != nil {
			fmt.Printf("Image extension loaded: %v\n", ext)

		}

		fmt.Println("bounds", img.Bounds())
		fmt.Println("color model", img.ColorModel())
		fmt.Println("color at", img.At(100, 100))

		out := image.NewGray(img.Bounds())
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
				out.Set(x, y, cfg.convert(img.At(x, y)))
			}
			fmt.Printf("x: %v\r", x)
		}
		if err := storeImageToFile("D:\\out.png", out); err != nil {
			fmt.Printf("Error on store result: %v\n", err)
		}

	}
	fmt.Println(cfg)
}

func getImageFromFilePath(filePath string) (image.Image, string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("open file %q: %w", filePath, err)
	}
	defer f.Close()
	image, kind, err := image.Decode(f)
	if err != nil {
		return nil, "", fmt.Errorf("decode image: %w", err)
	}
	return image, kind, nil
}

func storeImageToFile(filePath string, img image.Image) error {
	outf, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer outf.Close()
	if err := png.Encode(outf, img); err != nil {
		return fmt.Errorf("encode image to png: %v\n", err)
	}
	return nil
}
