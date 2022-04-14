package scraper

import (
	"io"

	"github.com/PuerkitoBio/goquery"
)

type HomePage struct {
	doc *goquery.Document
}

func NewHomePage(r io.Reader) (*HomePage, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	h := &HomePage{
		doc: doc,
	}

	return h, nil
}

func (hp *HomePage) GetUserName() string {
	return hp.doc.Find("header div.header-mypage span.bold").Text()
}
