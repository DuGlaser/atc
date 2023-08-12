package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/DuGlaser/atc/internal/repository/fetcher"
	"github.com/DuGlaser/atc/internal/repository/scraper"
)

type Language struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func main() {

	numArgs := len(os.Args)

	if numArgs <= 1 {
		fmt.Println("Please input target contest name.")
		os.Exit(1)
	}

	contestName := os.Args[1]

	res, err := fetcher.FetchSubmitPage(contestName)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	sp, err := scraper.NewSubmitPage(res.Body)
	if err != nil {
		panic(err)
	}

	_ls := sp.GetLanguageIds()

	ls := make([]Language, len(_ls))

	for i, l := range _ls {
		ls[i] = Language{
			Id:   l.Value,
			Name: l.Name,
		}
	}

	data, err := json.Marshal(ls)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
