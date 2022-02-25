package scraper

import (
	"errors"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ProblemPage struct {
	doc *goquery.Document
}

func NewProblemPage(r io.Reader) (*ProblemPage, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	p := &ProblemPage{
		doc: doc,
	}

	return p, nil
}

type Sample struct {
	In  string
	Out string
}

func (cp *ProblemPage) GetProblemSamples() ([]*Sample, error) {
	sms := []*Sample{}

	ins := []string{}
	outs := []string{}

	cp.doc.Find("div#task-statement div.part").Each(func(i int, s *goquery.Selection) {
		t := s.Find("h3").Text()

		p := s.Find("pre").Text()
		p = strings.TrimSpace(p)
		p = strings.TrimRight(p, "\n")

		switch {
		case strings.Contains(t, "入力例"):
			ins = append(ins, p)
		case strings.Contains(t, "出力例"):
			outs = append(outs, p)
		}
	})

	if len(ins) != len(outs) {
		return nil, errors.New("The number of input and output examples is different.")
	}

	for i := range ins {
		sms = append(sms, &Sample{In: ins[i], Out: outs[i]})
	}

	return sms, nil
}
