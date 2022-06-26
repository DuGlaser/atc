package scraper

import (
	"errors"
	"io"
	"strings"

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

type Problem struct {
	// URLに使われる、コンテスト内の問題を一意に識別するID
	ID string
	// 問題ページで表示される一意なID
	DisplayID string
}

func (cp *ContestPage) GetProblemIds() []Problem {
	ps := []Problem{}
	cp.doc.Find("div#contest-statement h3").Each(func(i int, s *goquery.Selection) {
		if s.Text() == "配点" {
			s.Next().Find("table tbody tr > td:first-child").Each(func(i int, s *goquery.Selection) {
				t := strings.ToLower(s.Text())
				ps = append(ps, Problem{ID: "", DisplayID: t})
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
