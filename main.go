package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type Image struct {
	Filename  string
	Extension string
}

func getImagesInCurrentDir() ([]string, error) {
	// Define common image file extensions
	imageExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".webp": true,
		".svg":  true,
	}

	// Read all entries in the current directory
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var images []string
	// Iterate over each entry
	for _, entry := range entries {
		// Skip directories
		if entry.IsDir() {
			continue
		}
		// Get the file extension and convert to lowercase
		ext := filepath.Ext(entry.Name())
		if len(ext) > 0 {
			ext = filepath.Ext(entry.Name())
			ext = filepath.Ext(ext)
		}
		// Check if the extension is in the list of image extensions
		if imageExtensions[ext] {
			// Append the full path to the result slice
			images = append(images, filepath.Join(".", entry.Name()))
		}
	}
	return images, nil
}
func main() {

	magickCommand := "magick"

	if runtime.GOOS == "windows" {
		magickCommand = "magick.exe"
	}

	fmt.Println(magickCommand)

	images, err := getImagesInCurrentDir()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, img := range images {
		fmt.Println(img)
		// convert 02\ PRINT\ ARTORIAS\ imprimir.png -colorspace CMYK -profile USWebCoatedSWOP.icc image_CMYK.png
		cmd := exec.Command(magickCommand, img, "-colorspace", "CMYK", "-profile", "USWebCoatedSWOP.icc", "cmyk_"+img)
		stdout, err := cmd.Output()
		if err != nil {
			fmt.Println("Error ", err)
			os.Exit(0)
		}

		fmt.Println(stdout)
	}

}
