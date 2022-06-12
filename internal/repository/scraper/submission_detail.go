package scraper

import (
	"fmt"
	"io"

	"github.com/PuerkitoBio/goquery"
)

type SubmissionDetailPage struct {
	doc *goquery.Document
}

func NewSubmissionDetailPage(r io.Reader) (*SubmissionDetailPage, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	s := &SubmissionDetailPage{
		doc: doc,
	}

	return s, nil
}

func (sd *SubmissionDetailPage) GetSubmissionStatusMap() (map[StatusCode]int, error) {
	var err error
	m := map[StatusCode]int{}
	sd.doc.Find("table").Last().Find("tbody > tr > td:nth-child(2)").Each(func(i int, s *goquery.Selection) {
		if status, ok := statusCodeMap[s.Text()]; ok {
			m[status] += 1
		} else {
			err = fmt.Errorf("invalid status. got=%s", s.Text())
			return
		}
	})

	return m, err
}
