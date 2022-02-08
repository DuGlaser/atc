package scraper

import (
	"errors"
	"io"

	"github.com/PuerkitoBio/goquery"
)

type SubmitPage struct {
	doc *goquery.Document
}

func NewSubmitPage(r io.Reader) (*SubmitPage, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	s := &SubmitPage{
		doc: doc,
	}

	return s, nil
}

func (sp *SubmitPage) GetCSRFToken() (string, error) {
	value, exists := sp.doc.Find("input[name='csrf_token']").First().Attr("value")
	if !exists {
		return "", errors.New("csrf_token is not found")
	}

	return value, nil
}

type Language struct {
	Value string
	Name  string
}

func (sp *SubmitPage) GetLanguageIds() []*Language {
	ls := []*Language{}

	sp.doc.Find("div#select-lang > div:first-child > select > option").Each(func(i int, s *goquery.Selection) {
		value, exists := s.Attr("value")
		if !exists {
			return
		}

		name := s.Text()
		ls = append(ls, &Language{
			Value: value,
			Name:  name,
		})
	})

	return ls
}

type Task struct {
	Value string
	Name  string
}

func (sp *SubmitPage) GetTasks() []*Task {
	ts := []*Task{}

	sp.doc.Find("select#select-task > option").Each(func(i int, s *goquery.Selection) {
		value, exists := s.Attr("value")
		if !exists {
			return
		}

		name := s.Text()
		ts = append(ts, &Task{
			Value: value,
			Name:  name,
		})
	})

	return ts
}
