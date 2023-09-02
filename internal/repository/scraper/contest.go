package scraper

import (
	"errors"
	"io"
	"net/url"
	"strings"
	"time"

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
	ID string `json:"id"`
	// 問題ページで表示される一意なID
	DisplayID string `json:"display_id"`
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

func (cp *ContestPage) GetStartAt() (time.Time, error) {
	u, ok := cp.doc.Find(".contest-duration a").First().Attr("href")
	if !ok {
		return time.Time{}, errors.New("target url is not found")
	}

	parsed, err := url.Parse(u)
	if err != nil {
		return time.Time{}, errors.New("start time is not found")
	}

	t := parsed.Query().Get("iso")
	return time.Parse("20060102T1504", t)
}

func (cp *ContestPage) GetCSRFToken() (string, error) {
	value, exists := cp.doc.Find("input[name='csrf_token']").First().Attr("value")
	if !exists {
		return "", errors.New("csrf_token is not found")
	}

	return value, nil
}
