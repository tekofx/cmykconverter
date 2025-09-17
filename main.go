package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/tekofx/cmykconverter/internal/utils"
)

func main() {

	magickCommand := "magick"

	if runtime.GOOS == "windows" {
		magickCommand = "magick.exe"
	}

	images, err := utils.GetImagesInCurrentDir()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if !utils.FileExists("USWebCoatedSWOP.icc") {
		utils.DownloadCMYKProfile()
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
