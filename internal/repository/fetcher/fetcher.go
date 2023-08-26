package fetcher

import (
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/DuGlaser/atc/internal/auth"
	"github.com/DuGlaser/atc/internal/repository/scraper"
)

const ATCODER_URL = "https://atcoder.jp/"

func getDefaultClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second,
		},
		Timeout: 30 * time.Second}

}

func setCookie(req *http.Request) error {
	session, err := auth.GetSession()
	if err != nil {
		return err
	}

	req.Header.Add("Cookie", session)
	return nil
}

func is2xx(res *http.Response) bool {
	return res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusMultipleChoices
}

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

func FetchAuthSession(username, password string) (*http.Response, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := getDefaultClient()
	client.Jar = jar

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

	if err := setCookie(req); err != nil {
		return nil, err
	}

	res, err := getDefaultClient().Do(req)
	if err != nil {
		return nil, err
	}

	if !is2xx(res) {
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

	if err := setCookie(req); err != nil {
		return nil, err
	}

	res, err := getDefaultClient().Do(req)
	if err != nil {
		return nil, err
	}

	if !is2xx(res) {
		return nil, fmt.Errorf("Could not access %s problems.", c)
	}

	return res, nil
}

func FetchProblemPage(contest, problem string) (*http.Response, error) {
	c := strings.ToLower(contest)
	p := strings.ToLower(problem)

	req, err := http.NewRequest("GET", GetAtcoderUrl("contests", c, "tasks", p), nil)
	if err != nil {
		return nil, err
	}

	if err := setCookie(req); err != nil {
		return nil, err
	}

	res, err := getDefaultClient().Do(req)
	if err != nil {
		return nil, err
	}

	if !is2xx(res) {
		return nil, fmt.Errorf("Could not access %s problem.", p)
	}

	return res, nil
}

func FetchSubmitPage(contest string) (*http.Response, error) {
	req, err := http.NewRequest("GET", GetAtcoderUrl("contests", contest, "submit"), nil)
	if err != nil {
		return nil, err
	}

	if err := setCookie(req); err != nil {
		return nil, err
	}

	return getDefaultClient().Do(req)
}

func FetchHomePage() (*http.Response, error) {
	req, err := http.NewRequest("GET", GetAtcoderUrl("home"), nil)
	if err != nil {
		return nil, err
	}

	if err := setCookie(req); err != nil {
		return nil, err
	}

	return getDefaultClient().Do(req)
}

func FetchSubmissionsMe(contest string) (*http.Response, error) {
	req, err := http.NewRequest("GET", GetAtcoderUrl("contests", contest, "submissions", "me"), nil)
	if err != nil {
		return nil, err
	}

	if err := setCookie(req); err != nil {
		return nil, err
	}

	return getDefaultClient().Do(req)
}

func FetchSubmissionDetail(contest, submissionID string) (*http.Response, error) {
	req, err := http.NewRequest("GET", GetAtcoderUrl("contests", contest, "submissions", submissionID), nil)
	if err != nil {
		return nil, err
	}

	if err := setCookie(req); err != nil {
		return nil, err
	}

	return getDefaultClient().Do(req)
}

func PostProblemAnswer(contest, problem, lang, code string) (*http.Response, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := getDefaultClient()
	client.Jar = jar

	submitUrl := GetAtcoderUrl("contests", contest, "submit")

	req, err := http.NewRequest("GET", submitUrl, nil)
	if err != nil {
		return nil, err
	}

	if err := setCookie(req); err != nil {
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
	form.Add("data.TaskScreenName", problem)
	form.Add("data.LanguageId", lang)
	form.Add("csrf_token", csrf)
	form.Add("sourceCode", code)

	return client.PostForm(submitUrl, form)
}
