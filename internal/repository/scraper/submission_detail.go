package scraper

import (
	"fmt"
	"io"
	"strconv"
	"strings"

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

	sd.doc.Find("div#main-container h4").Each(func(i int, s *goquery.Selection) {
		if s.Text() != "ジャッジ結果" {
			return
		}

		s.Next().Find("table > tbody > tr:nth-child(3) > td tr").Each(func(i int, s *goquery.Selection) {
			// NOTE: td:nth-child(2)が使えないので最初の要素は無視する
			if i == 0 {
				return
			}

			nds := []string{}
			s.Find("td").Each(func(i int, s *goquery.Selection) { nds = append(nds, strings.TrimSpace(s.Text())) })

			if status, ok := statusCodeMap[nds[0]]; ok {
				ss := strings.Split(nds[1], " ")
				count, err := strconv.Atoi(ss[len(ss)-1])

				if err != nil {
					return
				}

				m[status] += count
			} else {
				err = fmt.Errorf("invalid status. got=%s", s.Text())
				return
			}
		})
	})

	return m, err
}
