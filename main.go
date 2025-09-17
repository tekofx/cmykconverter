package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Image struct {
	Filename  string
	Name      string
	Extension string
}

func downloadCMYKProfile() *error {
	file, err := os.Create("USWebCoatedSWOP.icc")
	if err != nil {
		return &err
	}
	defer file.Close()

	// Send HTTP GET request
	resp, err := http.Get("www.color.org/registry/profiles/SWOP2006_Coated3v2.icc")
	if err != nil {
		return &err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return &err
	}

	// Copy the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return &err
	}

	return nil
}

func getImagesInCurrentDir() ([]Image, error) {
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

	var images []Image
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
			newImage := Image{
				Filename:  entry.Name(),
				Extension: ext,
				Name:      strings.Split(entry.Name(), ".")[0],
			}
			images = append(images, newImage)
		}
	}
	return images, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {

	magickCommand := "magick"

	if runtime.GOOS == "windows" {
		magickCommand = "magick.exe"
	}

	images, err := getImagesInCurrentDir()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if !fileExists("USWebCoatedSWOP.icc") {
		downloadCMYKProfile()
	}
	if len(images) == 0 {
		fmt.Println("No se encontraron imágenes. Arrastra imágenes a esta carpeta")
		os.Exit(0)
	}

	for _, img := range images {
		fmt.Println(img)
		//convert 02\ PRINT\ ARTORIAS\ imprimir.png -colorspace CMYK -profile USWebCoatedSWOP.icc image_CMYK.png
		cmd := exec.Command(magickCommand, img.Filename, "-colorspace", "CMYK", "-profile", "USWebCoatedSWOP.icc", "cmyk_"+img.Name+".jpg")
		_, err := cmd.Output()
		if err != nil {
			fmt.Println("Error ", err)
			os.Exit(0)
		}
	}

}
