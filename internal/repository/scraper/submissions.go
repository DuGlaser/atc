package scraper

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type SubmissionsPage struct {
	doc *goquery.Document
}

func NewSubmissionsPage(r io.Reader) (*SubmissionsPage, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	s := &SubmissionsPage{
		doc: doc,
	}

	return s, nil
}

type StatusCode string

const (
	AC  = StatusCode("AC")
	CE  = StatusCode("CE")
	IE  = StatusCode("IE")
	MLE = StatusCode("MLE")
	OLE = StatusCode("OLE")
	RE  = StatusCode("RE")
	TLE = StatusCode("TLE")
	WA  = StatusCode("WA")
	WJ  = StatusCode("WJ")
)

var statusCodeMap = map[string]StatusCode{
	"AC":  AC,
	"CE":  CE,
	"IE":  IE,
	"MLE": MLE,
	"OLE": OLE,
	"RE":  RE,
	"TLE": TLE,
	"WA":  WA,
	"WJ":  WJ,
}

type ResultMetaData struct {
	ExecTime string
	Memory   string
}

type Counter struct {
	Total   int
	Current int
}

type Submission struct {
	ID             string
	Date           string
	Task           string
	User           string
	Lang           string
	Score          int64
	CodeSize       string
	Status         StatusCode
	Counter        *Counter
	ResultMetaData *ResultMetaData
}

func (sp *SubmissionsPage) GetLatestSubmission() (*Submission, error) {
	tds := sp.doc.Find("table > tbody > tr:first-child > td")

	if tds.Length() < 10 {
		return sp.waitingJudge(tds)
	}

	return sp.completeJudge(tds)

}

func (sp *SubmissionsPage) waitingJudge(tds *goquery.Selection) (*Submission, error) {
	sm := &Submission{
		ResultMetaData: nil,
	}

	var err error

	tds.Each(func(i int, s *goquery.Selection) {

		switch i {
		case 0:
			sm.Date = s.Text()
		case 1:
			sm.Task = s.Text()
		case 2:
			sm.User = s.Text()
		case 3:
			sm.Lang = s.Text()
		case 4:
			score, err := strconv.Atoi(s.Text())
			if err != nil {
				return
			}

			sm.Score = int64(score)
		case 5:
			sm.CodeSize = s.Text()
		case 6:
			status, counter := sp.parseStatus(s.Text())
			sm.Status = status
			sm.Counter = counter
		case 7:
			href, exists := s.Find("a").Attr("href")
			if !exists {
				err = fmt.Errorf("Submission id is not found. got=%s", s.Text())
				return
			}

			paths := strings.Split(href, "/")

			sm.ID = paths[len(paths)-1]
		}
	})

	return sm, err
}

func (sp *SubmissionsPage) completeJudge(tds *goquery.Selection) (*Submission, error) {
	sm := &Submission{
		ResultMetaData: &ResultMetaData{},
	}
	var err error

	tds.Each(func(i int, s *goquery.Selection) {
		switch i {
		case 0:
			sm.Date = s.Text()
		case 1:
			sm.Task = s.Text()
		case 2:
			sm.User = s.Text()
		case 3:
			sm.Lang = s.Text()
		case 4:
			score, err := strconv.Atoi(s.Text())
			if err != nil {
				return
			}

			sm.Score = int64(score)
		case 5:
			sm.CodeSize = s.Text()
		case 6:
			if status, ok := statusCodeMap[s.Text()]; ok {
				sm.Status = status
			}
		case 7:
			sm.ResultMetaData.ExecTime = s.Text()
		case 8:
			sm.ResultMetaData.Memory = s.Text()
		case 9:
			href, exists := s.Find("a").Attr("href")
			if !exists {
				err = fmt.Errorf("Submission id is not found. got=%s", s.Text())
				return
			}

			paths := strings.Split(href, "/")

			sm.ID = paths[len(paths)-1]
		}
	})

	return sm, err
}

func (sp *SubmissionsPage) parseStatus(row string) (StatusCode, *Counter) {
	row = strings.TrimSpace(row)
	rows := strings.Split(row, " ")

	parseCounter := func(row string) *Counter {
		counter := strings.Split(rows[0], "/")
		cs := []int{}

		for _, c := range counter {
			ic, err := strconv.Atoi(c)
			if err != nil {
				return nil
			}

			cs = append(cs, ic)
		}

		return &Counter{
			Current: cs[0],
			Total:   cs[1],
		}
	}

	if len(rows) == 2 {
		return statusCodeMap[rows[1]], parseCounter(rows[0])
	}

	if strings.Contains(row, "/") {
		return WJ, parseCounter(row)
	}

	return statusCodeMap[row], nil
}
