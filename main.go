package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/tekofx/cmykconverter/internal/utils"
)

func main() {
	err := utils.CheckCmykConverterUpdates()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Pulsa Enter para continuar...")
		fmt.Scanln() // Waits for Enter key press
		os.Exit(0)
	}
	images, err := utils.GetImagesInCurrentDir()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Pulsa Enter para continuar...")
		fmt.Scanln() // Waits for Enter key press
		os.Exit(0)
	}
	if !utils.FileExists("USWebCoatedSWOP.icc") {
		fmt.Println("Descargando perfil de color CMYK")
		utils.DownloadFile("www.color.org/registry/profiles/SWOP2006_Coated3v2.icc", "USWebCoatedSWOP.icc")
		fmt.Println("Descargado!")
	}
	if !utils.FolderExists("imagemagick") {
		fmt.Println("Descargando imagemagick")
		utils.DownloadFile("https://imagemagick.org/archive/binaries/ImageMagick-7.1.2-3-portable-Q16-x64.zip", "imagemagick.zip")
		fmt.Println("Descargado")
		utils.ExtractFile("imagemagick.zip", "imagemagick")
		os.Remove("imagemagick.zip")
	}
	if len(images) == 0 {
		fmt.Println("No se encontraron imágenes. Arrastra imágenes a esta carpeta antes de iniciar el programa.")
		fmt.Println("Pulsa Enter para continuar...")
		fmt.Scanln() // Waits for Enter key press
		os.Exit(0)
	}
	for _, img := range images {
		cmd := exec.Command(
			"imagemagick/magick.exe",
			img.Filename,
			"-colorspace",
			"CMYK",
			"-profile",
			"USWebCoatedSWOP.icc",
			"cmyk_"+img.Name+".jpg",
		)
		_, err := cmd.Output()
		if err != nil {
			fmt.Println("Error ", err)
			fmt.Println("Pulsa Enter para continuar...")
			fmt.Scanln() // Waits for Enter key press
			os.Exit(0)
		}
	}
}
