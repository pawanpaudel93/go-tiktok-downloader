package downloader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	goquery "github.com/PuerkitoBio/goquery"
)

var headers map[string]string
var cookies []*http.Cookie

// TiktokProfile -
type TiktokProfile struct {
	URL        string
	FilePath   string
	data       PageProps
	httpClient *http.Client
	UseProxy   bool
	proxy      string
}

// TiktokVideo -
type TiktokVideo struct {
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

func (video *TiktokVideo) setProxy() {
	resp, err := http.Get("http://pubproxy.com/api/proxy?limit=1&format=txt&type=http&level=elite&last_check=1&HTTPS=true&GET=true&USER_AGENT=true&COOKIES=true&REFERER=true&country=US")
	if err == nil {
		proxy, _ := ioutil.ReadAll(resp.Body)
		video.proxy = "http://" + string(proxy)
	}
}

func (video *TiktokVideo) setClient(jar *cookiejar.Jar) {
	video.httpClient = &http.Client{
		Jar: jar,
	}
	if video.UseProxy && video.proxy != "" {
		if proxyURL, err := url.Parse(video.proxy); err == nil {
			transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
			video.httpClient = &http.Client{
				Jar:       jar,
				Transport: transport,
			}
		}
	}
}

func (profile *TiktokProfile) setProxy() {
	resp, err := http.Get("http://pubproxy.com/api/proxy?limit=1&format=txt&type=http&level=elite&last_check=1&HTTPS=true&GET=true&USER_AGENT=true&COOKIES=true&REFERER=true&country=US")
	if err == nil {
		proxy, _ := ioutil.ReadAll(resp.Body)
		profile.proxy = "http://" + string(proxy)
	}
}

func (profile *TiktokProfile) setClient(jar *cookiejar.Jar) {
	profile.httpClient = &http.Client{
		Jar: jar,
	}
	if profile.UseProxy && profile.proxy != "" {
		if proxyURL, err := url.Parse(profile.proxy); err == nil {
			transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
			profile.httpClient = &http.Client{
				Jar:       jar,
				Transport: transport,
			}
		}
	}
}

// Download -
func (video *TiktokVideo) Download(FilePath string) {
	jar, _ := cookiejar.New(nil)
	video.FilePath = FilePath
	videoURL := video.URL
	if video.UseProxy {
		video.setProxy()
	}

	req, err := http.NewRequest("GET", videoURL, nil)
	checkError(err)
	parsedURL, _ := url.Parse(videoURL)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	jar.SetCookies(parsedURL, cookies)
	video.setClient(jar)
	if resp, err := video.httpClient.Do(req); err != nil {
		fmt.Println(err)
	} else {
		VideoData := VideoData{}
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		checkError(err)
		doc.Find("#__NEXT_DATA__").Each(func(i int, s *goquery.Selection) {
			err := json.Unmarshal([]byte(s.Text()), &VideoData)
			video.data = VideoData
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
		video.setClient(jar)
		resp, err = video.httpClient.Do(req)
		checkError(err)
		err = saveTiktok(video.FilePath, resp)
		checkError(err)
	}
}

// GetVideoInfo -
func (video *TiktokVideo) GetVideoInfo() string {
	tiktokData := map[string]interface{}{
		"video": map[string]interface{}{
			"URL":         video.URL,
			"likes":       video.data.Likes,
			"shares":      video.data.Shares,
			"comments":    video.data.Comments,
			"played":      video.data.Played,
			"createdTime": video.data.CreatedTime,
			"description": video.data.Description,
		},
		"author": map[string]interface{}{
			"uniqueID":  video.data.Author.UniqueID,
			"nickname":  video.data.Author.Nickname,
			"url":       "https://tiktok.com/@" + video.data.Author.UniqueID,
			"followers": video.data.AuthorStats.Followers,
			"following": video.data.AuthorStats.Followings,
			"hearts":    video.data.AuthorStats.Hearts,
			"videos":    video.data.AuthorStats.Videos,
		},
	}
	data, err := json.Marshal(tiktokData)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(data)
}

// GetProfilePicture -
func (profile *TiktokProfile) GetProfilePicture() string {
	jar, _ := cookiejar.New(nil)
	profileURL := profile.URL
	if profile.UseProxy {
		profile.setProxy()
	}

	req, err := http.NewRequest("GET", profileURL, nil)
	checkError(err)
	parsedURL, _ := url.Parse(profileURL)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	jar.SetCookies(parsedURL, cookies)
	profile.setClient(jar)
	if resp, err := profile.httpClient.Do(req); err != nil {
		fmt.Println(err)
	} else {
		VideoData := VideoData{}
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		checkError(err)
		doc.Find("#__NEXT_DATA__").Each(func(i int, s *goquery.Selection) {
			err := json.Unmarshal([]byte(s.Text()), &VideoData)
			profile.data = VideoData.Props.PageProps
			checkError(err)
		})
		photoData := map[string]string{
			"Thumnail": replaceUnicode(profile.data.AvatarThumb),
			"Medium":   replaceUnicode(profile.data.AvatarMedium),
			"Larger":   replaceUnicode(profile.data.AvatarLarger),
		}
		data, err := json.Marshal(photoData)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		return string(data)
	}
	return ""
}

// DownloadPhoto -
func (profile *TiktokProfile) DownloadPhoto(PhotoType string, FilePath string) error {
	var photoURL string
	jar, _ := cookiejar.New(nil)
	profile.FilePath = FilePath
	switch PhotoType {
	case "thumbnail":
		photoURL = profile.data.AvatarThumb
	case "medium":
		photoURL = profile.data.AvatarMedium
	default:
		photoURL = profile.data.AvatarLarger
	}
	parsedURL, err := url.Parse(photoURL)
	if err != nil {
		fmt.Println(err)
		return err
	}
	req, err := http.NewRequest("GET", photoURL, nil)
	checkError(err)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	jar.SetCookies(parsedURL, cookies)
	profile.setClient(jar)
	resp, err := profile.httpClient.Do(req)
	checkError(err)
	err = saveTiktok(profile.FilePath, resp)
	checkError(err)
	return err
}
