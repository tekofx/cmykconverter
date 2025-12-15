package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/tekofx/cmykconverter/internal/errors"
	"github.com/tekofx/cmykconverter/internal/utils"
)

func main() {
	version, err := utils.LoadVersion()
	fmt.Printf("Cmyk Converter - Version %s\n", version)

	err = utils.CheckCmykConverterUpdates()
	if err != nil {
		errors.ManagerError(err)
	}
	images, err := utils.GetImagesInCurrentDir()
	if err != nil {
		errors.ManagerError(err)
	}
	if !utils.FileExists("USWebCoatedSWOP.icc") {
		fmt.Println("Descargando perfil de color CMYK")
		err = utils.DownloadFile("www.color.org/registry/profiles/SWOP2006_Coated3v2.icc", "USWebCoatedSWOP.icc")
		if err != nil {
			errors.ManagerError(err)
		}
		fmt.Println("Descargado!")
	}

	if !utils.FileExists("imagemagick/magick.exe") {
		fmt.Println("Descargando imagemagick")
		err := utils.DownloadFile("https://github.com/ImageMagick/ImageMagick/releases/download/7.1.2-11/ImageMagick-7.1.2-11-portable-Q16-x64.7z", "imagemagick.7z")
		if err != nil {
			errors.ManagerError(err)
		}
		err = utils.ExtractFile("imagemagick.7z", "imagemagick")
		if err != nil {
			errors.ManagerError(err)
		}
		fmt.Println("Descargado")
	}

	if len(images) == 0 {
		fmt.Println("No se encontraron imágenes. Arrastra imágenes a esta carpeta antes de iniciar el programa.")
		fmt.Println("Pulsa Enter para continuar...")
		fmt.Scanln() // Waits for Enter key press
		os.Exit(0)
	}
	fmt.Printf("Se han encontrado %d imágenes, convirtiendo...\n", len(images))
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
			errors.ManagerError(err)
		}
	}
}
