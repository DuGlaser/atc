package scraper

import (
	"errors"
	"io"
	"strings"

	"github.com/DuGlaser/atc/internal/core"
	"github.com/PuerkitoBio/goquery"
)

type TaskPage struct {
	doc *goquery.Document
}

func NewTaskPage(r io.Reader) (*TaskPage, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	p := &TaskPage{
		doc: doc,
	}

	return p, nil
}

type Sample struct {
	In  string
	Out string
}

func (tp *TaskPage) GetTaskTestCases() ([]core.TestCase, error) {
	sms := []core.TestCase{}

	ins := []string{}
	expecteds := []string{}

	tp.doc.Find("div#task-statement div.part").Each(func(i int, s *goquery.Selection) {
		t := s.Find("h3").Text()

		p := s.Find("pre").First().Text()
		p = strings.TrimSpace(p)
		p = strings.TrimRight(p, "\n")

		switch {
		case strings.Contains(t, "入力例"):
			ins = append(ins, p)
		case strings.Contains(t, "出力例"):
			expecteds = append(expecteds, p)
		}
	})

	if len(ins) != len(expecteds) {
		return nil, errors.New("The number of input and output examples is different.")
	}

	for i := range ins {
		sms = append(sms, core.TestCase{In: ins[i], Expected: expecteds[i], ID: i + 1})
	}

	return sms, nil
}
