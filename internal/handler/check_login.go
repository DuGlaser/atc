package handler

import (
	"fmt"
	"os"

	"github.com/DuGlaser/atc/internal/auth"
	"github.com/spf13/cobra"
)

func CheckLogin() {
	expired, err := auth.IsExpired()
	cobra.CheckErr(err)

	if expired {
		fmt.Println("Please login again.")
		os.Exit(1)
	}
}
