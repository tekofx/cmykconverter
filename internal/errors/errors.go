package errors

import (
	"fmt"
	"os"
)

func ManagerError(err error) {
	fmt.Println("Error:", err)
	fmt.Println("Pulsa Enter para continuar...")
	fmt.Scanln() // Waits for Enter key press
	os.Exit(0)
}
