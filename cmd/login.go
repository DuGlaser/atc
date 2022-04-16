package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/DuGlaser/atc/internal/auth"
	"github.com/DuGlaser/atc/internal/repository/fetcher"
	"github.com/DuGlaser/atc/internal/repository/scraper"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type prompt struct {
	name     string
	password string
}

func (p *prompt) input() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	fmt.Print("Password: ***")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	fmt.Println()

	p.name = strings.TrimRight(name, "\n")
	p.password = string(password)
	return nil
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to atcoder",
	Long:  "Login to atcoder and save the session cookie locally.",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := auth.GetSession(); err == nil {
			fmt.Println("Already logged in.")
			os.Exit(0)
		}

		p := new(prompt)
		cobra.CheckErr(p.input())

		res, err := fetcher.FetchAuthSession(p.name, p.password)
		cobra.CheckErr(err)

		hp, err := scraper.NewHomePage(res.Body)
		cobra.CheckErr(err)

		name := hp.GetUserName()
		if name == "" {
			cobra.CheckErr("Login failed: Username or Password is incorrect.")
		}

		cs := []string{}
		for _, c := range res.Cookies() {
			cs = append(cs, c.String())
		}

		err = auth.StoreSession([]byte(strings.Join(cs, ";")))
		cobra.CheckErr(err)

		fmt.Printf("Success! Hi, %s!\n", name)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
