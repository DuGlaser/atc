package scraper

import (
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type TasksPage struct {
	doc *goquery.Document
}

func NewTasksPage(r io.Reader) (*TasksPage, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	t := &TasksPage{
		doc: doc,
	}

	return t, nil
}

func (tp *TasksPage) GetProblemIds() []Problem {
	ps := []Problem{}
	tp.doc.Find("table tbody tr td:first-child a").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if !exists {
			return
		}

		// /contests/abc001/tasks/abc001_1
		ls := strings.Split(link, "/")

		// abc001_1
		id := ls[len(ls)-1]

		t := strings.ToLower(s.Text())
		ps = append(ps, Problem{ID: id, DisplayID: t})
	})

	return ps
}

func (tp *TasksPage) GetProblemId(displayID string) *Problem {
	var ps Problem
	tp.doc.Find("table tbody tr td:first-child a").Each(func(i int, s *goquery.Selection) {
		t := strings.ToLower(s.Text())
		if t != displayID {
			return
		}

		link, exists := s.Attr("href")
		if !exists {
			return
		}

		// /contests/abc001/tasks/abc001_1
		ls := strings.Split(link, "/")

		// abc001_1
		id := ls[len(ls)-1]

		ps = Problem{ID: id, DisplayID: t}
	})

	return &ps
}
