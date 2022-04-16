package scraper

import (
	"errors"
	"io"
	"strings"

	"github.com/DuGlaser/atc/internal"
	"github.com/PuerkitoBio/goquery"
)

type ContestPage struct {
	doc *goquery.Document
}

func NewContestPage(r io.Reader) (*ContestPage, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	s := &ContestPage{
		doc: doc,
	}

	return s, nil
}

func (cp *ContestPage) GetProblemIds() []internal.Problem {
	ps := []internal.Problem{}
	cp.doc.Find("div#contest-statement h3").Each(func(i int, s *goquery.Selection) {
		if s.Text() == "配点" {
			s.Next().Find("table tbody tr > td:first-child").Each(func(i int, s *goquery.Selection) {
				t := strings.ToLower(s.Text())
				ps = append(ps, internal.Problem{ID: "", DisplayID: t})
			})
		}
	})

	return ps
}

func (cp *ContestPage) GetCSRFToken() (string, error) {
	value, exists := cp.doc.Find("input[name='csrf_token']").First().Attr("value")
	if !exists {
		return "", errors.New("csrf_token is not found")
	}

	return value, nil
}
