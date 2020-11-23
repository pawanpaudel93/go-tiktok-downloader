package tiktok

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"strings"
	"time"

	goquery "github.com/PuerkitoBio/goquery"
)

var headers map[string]string
var cookies []*http.Cookie

// Profile - Tiktok Profile
type Profile struct {
	URL        string
	filePath   string
	data       PageProps
	httpClient *http.Client
	Proxy      string
	BaseDIR    string
}

// Video - Tiktok Video
type Video struct {
	URL        string
	filePath   string
	data       VideoData
	Proxy      string
	httpClient *http.Client
	BaseDIR    string
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

func (video *Video) setProxy() {
	if !(strings.Contains(video.Proxy, "http://") || strings.Contains(video.Proxy, "https://")) {
		video.Proxy = "http://" + string(video.Proxy)
	}
}

func (video *Video) setClient(jar *cookiejar.Jar) {
	video.httpClient = &http.Client{
		Jar: jar,
	}
	if video.Proxy != "" {
		if proxyURL, err := url.Parse(video.Proxy); err == nil {
			transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
			video.httpClient = &http.Client{
				Jar:       jar,
				Transport: transport,
			}
		} else {
			fmt.Println(err)
			fmt.Println("Not Using Proxy")
		}
	}
}

func (profile *Profile) setProxy() {
	if !(strings.Contains(profile.Proxy, "http://") || strings.Contains(profile.Proxy, "https://")) {
		profile.Proxy = "http://" + string(profile.Proxy)
	}
}

func (profile *Profile) setClient(jar *cookiejar.Jar) {
	profile.httpClient = &http.Client{
		Jar: jar,
	}
	if profile.Proxy != "" {
		if proxyURL, err := url.Parse(profile.Proxy); err == nil {
			transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
			profile.httpClient = &http.Client{
				Jar:       jar,
				Transport: transport,
			}
		} else {
			fmt.Println(err)
			fmt.Println("Not Using Proxy")
		}
	}
}

// Download - Download Tiktok video
func (video *Video) Download() (string, error) {
	jar, _ := cookiejar.New(nil)
	video.filePath = path.Join(video.BaseDIR, video.data.Author.UniqueID+"_"+video.data.ItemStruct.VideoID+".mp4")
	URL := video.data.Video.URL
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return "", err
	}
	parsedURL, _ := url.Parse(URL)

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	jar.SetCookies(parsedURL, cookies)
	video.setClient(jar)
	resp, err := video.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	err = saveTiktok(video.filePath, resp)
	if err != nil {
		return "", err
	}
	return video.filePath, err
}

// FetchInfo - Get Tiktok video Information.
func (video *Video) FetchInfo() error {
	jar, _ := cookiejar.New(nil)
	videoURL := video.URL
	if video.Proxy != "" {
		video.setProxy()
	}

	req, err := http.NewRequest("GET", videoURL, nil)
	if err != nil {
		return err
	}
	parsedURL, _ := url.Parse(videoURL)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	jar.SetCookies(parsedURL, cookies)
	video.setClient(jar)
	resp, err := video.httpClient.Do(req)
	if err != nil {
		return err
	}
	VideoData := VideoData{}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}
	doc.Find("#__NEXT_DATA__").Each(func(i int, s *goquery.Selection) {
		err = json.Unmarshal([]byte(s.Text()), &VideoData)
		video.data = VideoData
	})
	return err
}

// GetInfo returns Tiktok video information
func (video *Video) GetInfo() (string, error) {
	tiktokData := map[string]interface{}{
		"video": map[string]interface{}{
			"ID":          video.data.ItemStruct.VideoID,
			"URL":         video.URL,
			"likes":       video.data.Likes,
			"shares":      video.data.Shares,
			"comments":    video.data.Comments,
			"played":      video.data.Played,
			"createdTime": video.data.CreatedTime,
			"description": video.data.Description,
			"cover":       video.data.ItemStruct.Video.Cover,
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
	return string(data), err
}

// FetchInfo - Get Tiktok Profile information.
func (profile *Profile) FetchInfo() error {
	jar, _ := cookiejar.New(nil)
	profileURL := profile.URL
	if profile.Proxy != "" {
		profile.setProxy()
	}

	req, err := http.NewRequest("GET", profileURL, nil)
	if err != nil {
		return err
	}
	parsedURL, _ := url.Parse(profileURL)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	jar.SetCookies(parsedURL, cookies)
	profile.setClient(jar)
	resp, err := profile.httpClient.Do(req)
	if err != nil {
		return err
	}
	VideoData := VideoData{}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}
	doc.Find("#__NEXT_DATA__").Each(func(i int, s *goquery.Selection) {
		err = json.Unmarshal([]byte(s.Text()), &VideoData)
		profile.data = VideoData.Props.PageProps
	})
	return err
}

// GetPPInfo returns Tiktok profile picture information.
func (profile *Profile) GetPPInfo() (string, error) {
	photoData := map[string]string{
		"Thumbnail": replaceUnicode(profile.data.AvatarThumb),
		"Medium":    replaceUnicode(profile.data.AvatarMedium),
		"Larger":    replaceUnicode(profile.data.AvatarLarger),
	}
	data, err := json.Marshal(photoData)
	return string(data), err
}

// GetProfileInfo returns Tiktok profile information.
func (profile *Profile) GetProfileInfo() (string, error) {
	profileData := map[string]interface{}{
		"user":          profile.data.User,
		"userStats":     profile.data.UserStats,
		"userMetaStats": profile.data.UserMetaParams,
	}
	data, err := json.Marshal(profileData)
	return string(data), err
}

// DownloadPhoto - Download Tiktok Profile Picture.
func (profile *Profile) DownloadPhoto(PhotoType string) (string, error) {
	var photoURL string
	jar, _ := cookiejar.New(nil)
	switch PhotoType {
	case "thumbnail":
		photoURL = profile.data.AvatarThumb
		profile.filePath = profile.data.UserInfo.UniqueID + "_thumbnail.jpg"
	case "medium":
		photoURL = profile.data.AvatarMedium
		profile.filePath = profile.data.UserInfo.UniqueID + "_medium.jpg"
	default:
		photoURL = profile.data.AvatarLarger
		profile.filePath = profile.data.UserInfo.UniqueID + "_large.jpg"
	}
	profile.filePath = path.Join(profile.BaseDIR, profile.filePath)
	parsedURL, err := url.Parse(photoURL)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("GET", photoURL, nil)
	if err != nil {
		return "", err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	jar.SetCookies(parsedURL, cookies)
	profile.setClient(jar)
	resp, err := profile.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	err = saveTiktok(profile.filePath, resp)
	if err != nil {
		return "", err
	}
	return profile.filePath, err
}
