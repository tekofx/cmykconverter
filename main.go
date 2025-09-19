package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/tekofx/cmykconverter/internal/utils"
)

func main() {

	images, err := utils.GetImagesInCurrentDir()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if !utils.FileExists("USWebCoatedSWOP.icc") {
		utils.DownloadCMYKProfile()
	}
	if !utils.FolderExists("imagemagick") {
		utils.DownloadFile("https://imagemagick.org/archive/binaries/ImageMagick-7.1.2-3-portable-Q16-x64.zip", "imagemagick.zip")
		utils.ExtractFile("imagemagick.zip", "imagemagick")
		os.Remove("imagemagick.zip")
	}

	if len(images) == 0 {
		fmt.Println("No se encontraron imágenes. Arrastra imágenes a esta carpeta")
		os.Exit(0)
	}

	for _, img := range images {
		cmd := exec.Command("imagemagick/magick.exe", img.Filename, "-colorspace", "CMYK", "-profile", "USWebCoatedSWOP.icc", "cmyk_"+img.Name+".jpg")
		_, err := cmd.Output()
		if err != nil {
			fmt.Println("Error ", err)
			os.Exit(0)
		}
	}

}
