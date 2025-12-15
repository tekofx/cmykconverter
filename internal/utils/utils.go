package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bodgit/sevenzip"
	"github.com/minio/selfupdate"
	"github.com/tekofx/cmykconverter/internal/models"
)

const VERSION = "0.5"

func CheckCmykConverterUpdates() error {
	resp, err := http.Get("https://api.github.com/repos/tekofx/cmykconverter/releases/latest")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var release models.GithubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return err
	}

	if VERSION != release.TagName {
		err = doUpdate(release.Assets[0].Url)
		if err != nil {
			return err
		}

		fmt.Printf("Se ha actualizado cymkconverter a la version %s. Pulse ENTER para cerrar el programa y vuelva a abrirlo.", release.TagName)
		fmt.Scanln() // Waits for Enter key press
		os.Exit(0)
	}

	// Now you can use release.TagName, release.Url, etc.
	return nil

}

func doUpdate(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = selfupdate.Apply(resp.Body, selfupdate.Options{})
	return err
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
func ExtractFile(zipPath, destDir string) error {
	r, err := sevenzip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	for _, f := range r.File {
		// Skip the root directory entry itself (if any)
		if isRootDir(f.Name) {
			continue
		}

		// Remove the first directory component (e.g., "project/" from "project/file.txt")
		trimmedName := trimFirstDir(f.Name)

		filePath := filepath.Join(destDir, trimmedName)

		if !strings.HasPrefix(filePath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			continue // Skip unsafe paths
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, f.Mode()); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return err
		}

		srcFile, err := f.Open()
		if err != nil {
			return err
		}

		dstFile, err := os.Create(filePath)
		if err != nil {
			srcFile.Close()
			return err
		}

		_, err = io.Copy(dstFile, srcFile)
		srcFile.Close()
		dstFile.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func FolderExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func DownloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// isRootDir checks if the file entry is a root directory (ends with / and has no subpath)
func isRootDir(name string) bool {
	return strings.Contains(name, "/") && strings.Count(strings.Trim(name, "/"), "/") == 0 && strings.HasSuffix(name, "/")
}

// trimFirstDir removes the first directory part from a path
func trimFirstDir(p string) string {
	if i := strings.Index(p, "/"); i != -1 {
		return p[i+1:]
	}
	return p
}
