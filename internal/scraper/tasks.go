package scraper

import (
	"io"
	"strings"

	"github.com/DuGlaser/atc/internal"
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

func (tp *TasksPage) GetProblemIds() []internal.Problem {
	ps := []internal.Problem{}
	tp.doc.Find("table tbody tr td:first-child a").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if !exists {
			return
		}

		// /contests/abc001/tasks/abc001_1
		ls := strings.Split(link, "/")

		// abc001_1
		problem := ls[len(ls)-1]

		ids := strings.Split(problem, "_")
		// 1
		id := strings.ToLower(ids[len(ids)-1])

		t := strings.ToLower(s.Text())
		ps = append(ps, internal.Problem{ID: id, DisplayedID: t})
	})

	return ps
}
