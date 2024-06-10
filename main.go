package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

const kernel_size = 3;

func main() {
	var in_grid = read_image()

    var out_grid = make([][]color.RGBA, len(in_grid))
	for i := range out_grid {
		out_grid[i] = make([]color.RGBA, len(in_grid[i]))
	}

	// iterate through the grid
	for i := 0; i < len(in_grid); i++ {
		for j := 0; j < len(in_grid[i]); j++ {
			kernel := get_kernel(i, j, in_grid)
			out_grid[i][j] = pixel_average(kernel) 
		}
	}

	save_image(out_grid)
}

func get_kernel(i int, j int, grid [][]color.RGBA) [kernel_size][kernel_size]color.RGBA {
	var kernel = [kernel_size][kernel_size]color.RGBA{}

	// grab all surrounding pixels, if they exist
	for x := 0; x < kernel_size; x++ {
		for y := 0; y < kernel_size; y++ {

			// if in bounds:
			if i - x >= 0 && i + x <= len(grid) && j - y >= 0 && j + y <= len(grid[i]) {
				kernel[x][y] = grid[i - x][j - y]
			} else {
				kernel[x][y] = color.RGBA{0, 0, 0, 255}
			}
		}
	}

	return kernel
}

func pixel_average(kernel [kernel_size][kernel_size]color.RGBA) color.RGBA {
	// get the average of all the pixels in the kernel
	var sumR, sumG, sumB, sumA int

	for i := 0; i < kernel_size; i++ {
		for j := 0; j < kernel_size; j++ {
			sumR += int(kernel[i][j].R)
			sumG += int(kernel[i][j].G)
			sumB += int(kernel[i][j].B)
			sumA += int(kernel[i][j].A)
		}
	}

	numPixels := kernel_size * kernel_size

	return color.RGBA{
		R: uint8(sumR / numPixels),
		G: uint8(sumG / numPixels),
		B: uint8(sumB / numPixels),
		A: uint8(sumA / numPixels),
	}
}

func save_image(grid [][]color.RGBA) {
	// Create a new image with the same dimensions as the grid
    width := len(grid[0])
    height := len(grid)
    img := image.NewRGBA(image.Rect(0, 0, width, height))

    // Iterate over the grid and set the corresponding pixel in the image
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, grid[y][x])
		}
	}

    // Save the image to a file
    file, _ := os.Create("output.png")

    png.Encode(file, img)

	file.Close()
}

func read_image() [][]color.RGBA {
    // read an image and convert it to a grid

    file, _ := os.Open("in.png")

    // Decode the image
    img, _ := png.Decode(file)

    // Get the image bounds
    bounds := img.Bounds()
    width, height := bounds.Max.X, bounds.Max.Y

    // Initialize a matrix for the grid
    grid := make([][]color.RGBA, height)
    for i := range grid {
        grid[i] = make([]color.RGBA, width)
    }

    // Convert the image to a grid matrix
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            // Get the color of the pixel
            pixel := img.At(x, y)
            grid[y][x] = color.RGBAModel.Convert(pixel).(color.RGBA)
        }
    }

	file.Close()
    return grid
}