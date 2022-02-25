package fetcher

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"strings"

	"github.com/DuGlaser/atc/internal/auth"
	"github.com/DuGlaser/atc/internal/scraper"
)

const ATCODER_URL = "https://atcoder.jp/"

func GetAtcoderUrl(p ...string) string {
	u, err := url.Parse(ATCODER_URL)
	if err != nil {
		panic(err)
	}

	ps := []string{u.Path}
	ps = append(ps, p...)

	u.Path = path.Join(ps...)
	return u.String() + "?lang=ja"
}

func SetCookie(req *http.Request) error {
	session, err := auth.GetSession()
	if err != nil {
		return err
	}

	req.Header.Add("Cookie", session)
	return nil
}

func FetchAuthSession(username, password string) (*http.Response, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar: jar,
	}

	loginUrl := GetAtcoderUrl("/login")
	res, err := client.Get(loginUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	lp, err := scraper.NewLoginPage(res.Body)
	if err != nil {
		return nil, err
	}

	csrf, err := lp.GetCSRFToken()
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Add("username", username)
	form.Add("password", password)
	form.Add("csrf_token", csrf)

	return client.PostForm(loginUrl, form)
}

func FetchContestPage(contest string) (*http.Response, error) {
	req, err := http.NewRequest("GET", GetAtcoderUrl("contests", contest), nil)
	if err != nil {
		return nil, err
	}

	if err := SetCookie(req); err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Could not access %s contest.", contest)
	}

	return res, nil
}

func FetchProblems(contest string) (*http.Response, error) {
	c := strings.ToLower(contest)

	req, err := http.NewRequest("GET", GetAtcoderUrl("contests", c, "tasks"), nil)
	if err != nil {
		return nil, err
	}

	if err := SetCookie(req); err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Could not access %s problems.", c)
	}

	return res, nil
}

func FetchProblemPage(contest, problem string) (*http.Response, error) {
	c := strings.ToLower(contest)
	p := strings.ToLower(problem)
	id := fmt.Sprintf("%s_%s", c, p)

	req, err := http.NewRequest("GET", GetAtcoderUrl("contests", c, "tasks", id), nil)
	if err != nil {
		return nil, err
	}

	if err := SetCookie(req); err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Could not access %s problem.", id)
	}

	return res, nil
}

func FetchSubmitPage(contest string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", GetAtcoderUrl("contests", contest, "submit"), nil)
	if err != nil {
		return nil, err
	}

	if err := SetCookie(req); err != nil {
		return nil, err
	}

	return client.Do(req)
}

func FetchHomePage() (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", GetAtcoderUrl("home"), nil)
	if err != nil {
		return nil, err
	}

	if err := SetCookie(req); err != nil {
		return nil, err
	}

	return client.Do(req)
}

func PostProblemAnswer(contest, problem, lang, code string) (*http.Response, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar: jar,
	}

	submitUrl := GetAtcoderUrl("contests", contest, "submit")

	req, err := http.NewRequest("GET", submitUrl, nil)
	if err != nil {
		return nil, err
	}

	if err := SetCookie(req); err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	sp, err := scraper.NewSubmitPage(res.Body)
	if err != nil {
		return nil, err
	}

	csrf, err := sp.GetCSRFToken()
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Add("data.TaskScreenName", fmt.Sprintf("%s_%s", strings.ToLower(contest), strings.ToLower(problem)))
	form.Add("data.LanguageId", lang)
	form.Add("csrf_token", csrf)
	form.Add("sourceCode", code)

	return client.PostForm(submitUrl, form)
}
