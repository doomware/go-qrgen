package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"

	"github.com/akamensky/argparse"
	"github.com/disintegration/imaging"
	qrcode "github.com/skip2/go-qrcode"
)

func qr_gen(url *string, qr_size *int, out_file *string) {
	// Qr generator
	qr := qrcode.WriteFile(*url, qrcode.Highest, *qr_size, *out_file)
	_ = qr
}

func put_logo(out_file *string, logo *string, logo_size *int) {

	// Image 1 opening
	image1, err := imaging.Open(*out_file)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	// Image 2 opening
	// *logo_size = 500
	// *logo = "jwlogo.jpg"
	image2, err := imaging.Open(*logo)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	// Resize image 2 to logo_size
	image2 = imaging.Resize(image2, int(*logo_size), int(0), imaging.Lanczos)

	// Centering vars
	back_bounds := image1.Bounds() // Get size of image
	width := back_bounds.Dx()
	init_xy := int((width / 2) - (*logo_size / 2)) // Get x and y to put logo in center

	// Image destine
	dst := imaging.New(width, width, color.NRGBA{0, 0, 0, 0})    // Create a new image
	dst = imaging.Paste(dst, image1, image.Pt(0, 0))             // Paste image1 (qr) to new image
	dst = imaging.Paste(dst, image2, image.Pt(init_xy, init_xy)) // Paste logo in the center of qr

	// Save the result as PNG
	err = imaging.Save(dst, *out_file)
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
}

func main() {
	parser := argparse.NewParser("qr-gen", "Simple QR generator with logo in center")
	url := parser.String("u", "url", &argparse.Options{Required: true, Help: "Url to generate qr code"})
	qr_size := parser.Int("s", "size", &argparse.Options{Required: false, Help: "Size in px of png file", Default: 1850})
	out_file := parser.String("o", "output", &argparse.Options{Required: false, Help: "Output filename", Default: "qr.png"})
	logo := parser.String("l", "logo", &argparse.Options{Required: false, Help: "Filename of logo"})
	logo_size := parser.Int("S", "lsize", &argparse.Options{Required: false, Help: "Size of logo in px", Default: 400})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}

	// If pass out_file without extention, it give one
	fileExtension := filepath.Ext(*out_file)
	if fileExtension == "" {
		*out_file = fmt.Sprintf("%v.png", *out_file)
	}

	qr_gen(url, qr_size, out_file)
	if *logo != "" {
		put_logo(out_file, logo, logo_size)
	}
}
