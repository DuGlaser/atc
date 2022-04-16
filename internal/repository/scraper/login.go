package scraper

import (
	"errors"
	"io"

	"github.com/PuerkitoBio/goquery"
)

type LoginPage struct {
	doc *goquery.Document
}

func NewLoginPage(r io.Reader) (*LoginPage, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	l := &LoginPage{
		doc: doc,
	}

	return l, nil
}

func (lp *LoginPage) GetCSRFToken() (string, error) {
	value, exists := lp.doc.Find("input[name='csrf_token']").First().Attr("value")
	if !exists {
		return "", errors.New("csrf_token is not found")
	}

	return value, nil
}
