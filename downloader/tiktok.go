package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"time"

	goquery "github.com/PuerkitoBio/goquery"
)

var headers map[string]string
var cookies []*http.Cookie

// TikTok -
type TikTok struct {
	URL        string
	FilePath   string
	data       VideoData
	UseProxy   bool
	proxy      string
	httpClient *http.Client
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	headers = map[string]string{
		"User-Agent":      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36",
		"Accept":          "*/*",
		"Connection":      "keep-alive",
		"Pragma":          "no-cache",
		"Cache-Control":   "no-cache",
		"Sec-Fetch-Site":  "same-site",
		"Sec-Fetch-Mode":  "no-cors",
		"Sec-Fetch-Dest":  "video",
		"Referer":         "https://www.tiktok.com/",
		"Accept-Language": "en-US,en;q=0.9,bs;q=0.8,sr;q=0.7,hr;q=0.6",
		"sec-gpc":         "1",
		"DNT":             "1",
		"Range":           "bytes=0-",
	}
	webID := generateRandomNumber()
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookies = []*http.Cookie{
		{Name: "tt_webid", Value: webID, Path: "/", Expires: expiration, Domain: ".tiktok.com"},
		{Name: "tt_webid_v2", Value: webID, Path: "/", Expires: expiration, Domain: ".tiktok.com"},
	}
}

func generateRandomNumber() string {
	max := 1999999999999999999
	min := 1000000000000000000
	return strconv.Itoa(min + rand.Intn(max-min))
}

func saveTiktok(filepath string, resp *http.Response) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func (tiktok *TikTok) getProxy() {
	resp, err := http.Get("http://pubproxy.com/api/proxy?limit=1&format=txt&type=http&level=elite&last_check=1&HTTPS=true&GET=true&USER_AGENT=true&COOKIES=true&REFERER=true&country=US")
	if err == nil {
		proxy, _ := ioutil.ReadAll(resp.Body)
		tiktok.proxy = "http://" + string(proxy)
	}
}

func (tiktok *TikTok) setClient(jar *cookiejar.Jar) {
	tiktok.httpClient = &http.Client{
		Jar: jar,
	}
	if tiktok.UseProxy && tiktok.proxy != "" {
		if proxyURL, err := url.Parse(tiktok.proxy); err == nil {
			transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
			tiktok.httpClient = &http.Client{
				Jar:       jar,
				Transport: transport,
			}
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// Download -
func (tiktok *TikTok) Download() {
	jar, _ := cookiejar.New(nil)
	videoURL := tiktok.URL
	if tiktok.UseProxy {
		tiktok.getProxy()
	}

	req, err := http.NewRequest("GET", videoURL, nil)
	checkError(err)
	parsedURL, _ := url.Parse(videoURL)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	jar.SetCookies(parsedURL, cookies)
	tiktok.setClient(jar)
	if resp, err := tiktok.httpClient.Do(req); err != nil {
		fmt.Println(err)
	} else {
		VideoData := VideoData{}
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		checkError(err)
		doc.Find("#__NEXT_DATA__").Each(func(i int, s *goquery.Selection) {
			err := json.Unmarshal([]byte(s.Text()), &VideoData)
			tiktok.data = VideoData
			checkError(err)
		})
		URL := VideoData.Video.URL
		if err != nil {
			fmt.Println(err)
		}
		req, err = http.NewRequest("GET", URL, nil)
		checkError(err)
		parsedURL, _ = url.Parse(URL)

		for k, v := range headers {
			req.Header.Set(k, v)
		}
		jar.SetCookies(parsedURL, cookies)
		tiktok.setClient(jar)
		resp, err = tiktok.httpClient.Do(req)
		checkError(err)
		err = saveTiktok(tiktok.FilePath, resp)
		checkError(err)
	}
}

// GetTiktokInfo -
func (tiktok *TikTok) GetTiktokInfo() string {
	tiktokData := map[string]interface{}{
		"video": map[string]interface{}{
			"URL":         tiktok.URL,
			"likes":       tiktok.data.Likes,
			"shares":      tiktok.data.Shares,
			"comments":    tiktok.data.Comments,
			"played":      tiktok.data.Played,
			"createdTime": tiktok.data.CreatedTime,
			"description": tiktok.data.Description,
		},
		"author": map[string]interface{}{
			"uniqueID":  tiktok.data.Author.UniqueID,
			"nickname":  tiktok.data.Author.Nickname,
			"url":       "https://tiktok.com/@" + tiktok.data.Author.UniqueID,
			"followers": tiktok.data.AuthorStats.Followers,
			"following": tiktok.data.AuthorStats.Followings,
			"hearts":    tiktok.data.AuthorStats.Hearts,
			"videos":    tiktok.data.AuthorStats.Videos,
		},
	}
	data, err := json.Marshal(tiktokData)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(data)
}
