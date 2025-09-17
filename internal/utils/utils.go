package utils

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/tekofx/cmykconverter/internal/models"
)

func DownloadCMYKProfile() *error {
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

func GetImagesInCurrentDir() ([]models.Image, error) {
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

	var images []models.Image
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
			newImage := models.Image{
				Filename:  entry.Name(),
				Extension: ext,
				Name:      strings.Split(entry.Name(), ".")[0],
			}
			images = append(images, newImage)
		}
	}
	return images, nil
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
