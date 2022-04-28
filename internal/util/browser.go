package util

import (
	"fmt"
	"os/exec"
	"runtime"
)

// FYI: https://gist.github.com/hyg/9c4afcd91fe24316cbf0
// Thank you!
func Openbrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}
