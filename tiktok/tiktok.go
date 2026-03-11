package tiktok

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"regexp"
	"time"

	goquery "github.com/PuerkitoBio/goquery"
)

var headers map[string]string
var cookies []*http.Cookie

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

// fetchPage fetches a TikTok page and parses the embedded universal data.
func fetchPage(pageURL, proxy string) (*goquery.Document, *http.Client, error) {
	jar, _ := cookiejar.New(nil)
	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return nil, nil, err
	}
	parsedURL, _ := url.Parse(pageURL)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	jar.SetCookies(parsedURL, cookies)
	client := newClient(proxy, jar)
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("page returned status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return doc, client, nil
}

// parseUniversalData extracts the __UNIVERSAL_DATA_FOR_REHYDRATION__ JSON from the page.
func parseUniversalData(doc *goquery.Document) (*UniversalData, error) {
	var result *UniversalData
	var parseErr error
	doc.Find("#__UNIVERSAL_DATA_FOR_REHYDRATION__").Each(func(i int, s *goquery.Selection) {
		ud := UniversalData{}
		if err := json.Unmarshal([]byte(s.Text()), &ud); err != nil {
			parseErr = err
		} else {
			result = &ud
		}
	})
	if result != nil {
		return result, nil
	}
	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse page data: %w", parseErr)
	}
	return nil, errors.New("page data not found")
}

var videoURLRegex = regexp.MustCompile(`(?m)tiktok\.com\/@([\w.-]+)\/video\/(\d+)/?`)

// FetchInfo fetches TikTok video information from the page.
func (video *Video) FetchInfo() error {
	match := videoURLRegex.FindStringSubmatch(video.URL)
	if len(match) < 3 {
		return errors.New("invalid TikTok video URL")
	}
	username, videoId := match[1], match[2]
	video.filePath = path.Join(video.BaseDIR, username+"_"+videoId+".mp4")

	doc, client, err := fetchPage(video.URL, video.Proxy)
	if err != nil {
		return err
	}
	video.httpClient = client

	ud, err := parseUniversalData(doc)
	if err != nil {
		return err
	}
	if ud.DefaultScope.VideoDetail == nil {
		return errors.New("video data not found")
	}
	video.data = *ud.DefaultScope.VideoDetail
	return nil
}

// Download downloads the TikTok video. Must call FetchInfo first.
func (video *Video) Download() (string, error) {
	data, ok := video.data.(VideoDataV3)
	if !ok {
		return "", errors.New("call FetchInfo before Download")
	}
	URL := data.ItemInfo.ItemStruct.Video.PlayAddr
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return "", err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := video.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	if err := saveTiktok(video.filePath, resp); err != nil {
		return "", err
	}
	return video.filePath, nil
}

// GetInfo returns TikTok video information as JSON.
func (video *Video) GetInfo() (string, error) {
	v3data, ok := video.data.(VideoDataV3)
	if !ok {
		return "", errors.New("video data not found")
	}
	item := v3data.ItemInfo.ItemStruct
	tiktokData := map[string]interface{}{
		"video": map[string]interface{}{
			"ID":          item.ID,
			"URL":         video.URL,
			"likes":       item.Stats.DiggCount,
			"shares":      item.Stats.ShareCount,
			"comments":    item.Stats.CommentCount,
			"played":      item.Stats.PlayCount,
			"createdTime": item.CreateTime,
			"description": item.Desc,
			"cover":       item.Video.Cover,
		},
		"author": map[string]interface{}{
			"uniqueID":  item.Author.UniqueID,
			"nickname":  item.Author.Nickname,
			"url":       "https://tiktok.com/@" + item.Author.UniqueID,
			"followers": item.AuthorStats.FollowerCount,
			"following": item.AuthorStats.FollowingCount,
			"hearts":    item.AuthorStats.HeartCount,
			"videos":    item.AuthorStats.VideoCount,
		},
	}
	data, err := json.Marshal(tiktokData)
	return string(data), err
}

// FetchInfo fetches TikTok profile information from the page.
func (profile *Profile) FetchInfo() error {
	doc, client, err := fetchPage(profile.URL, profile.Proxy)
	if err != nil {
		return err
	}
	profile.httpClient = client

	ud, err := parseUniversalData(doc)
	if err != nil {
		return err
	}
	if ud.DefaultScope.UserDetail == nil {
		return errors.New("profile data not found")
	}
	profile.data = *ud.DefaultScope.UserDetail
	return nil
}

// GetPPInfo returns TikTok profile picture URLs as JSON.
func (profile *Profile) GetPPInfo() (string, error) {
	v3data, ok := profile.data.(ProfileDataV3)
	if !ok {
		return "", errors.New("profile data not found")
	}
	photoData := map[string]string{
		"Thumbnail": replaceUnicode(v3data.UserInfo.User.AvatarThumb),
		"Medium":    replaceUnicode(v3data.UserInfo.User.AvatarMedium),
		"Larger":    replaceUnicode(v3data.UserInfo.User.AvatarLarger),
	}
	data, err := json.Marshal(photoData)
	return string(data), err
}

// GetProfileInfo returns TikTok profile information as JSON.
func (profile *Profile) GetProfileInfo() (string, error) {
	v3data, ok := profile.data.(ProfileDataV3)
	if !ok {
		return "", errors.New("profile data not found")
	}
	profileData := map[string]interface{}{
		"user":      v3data.UserInfo.User,
		"userStats": v3data.UserInfo.Stats,
		"userMetaStats": map[string]interface{}{
			"title":       v3data.ShareMeta.Title,
			"description": v3data.ShareMeta.Desc,
		},
	}
	data, err := json.Marshal(profileData)
	return string(data), err
}

// DownloadPhoto downloads a TikTok profile picture. Must call FetchInfo first.
func (profile *Profile) DownloadPhoto(PhotoType string) (string, error) {
	v3data, ok := profile.data.(ProfileDataV3)
	if !ok {
		return "", errors.New("call FetchInfo before DownloadPhoto")
	}
	var photoURL string
	switch PhotoType {
	case "thumbnail":
		photoURL = v3data.UserInfo.User.AvatarThumb
		profile.filePath = v3data.UserInfo.User.UniqueID + "_thumbnail.jpg"
	case "medium":
		photoURL = v3data.UserInfo.User.AvatarMedium
		profile.filePath = v3data.UserInfo.User.UniqueID + "_medium.jpg"
	default:
		photoURL = v3data.UserInfo.User.AvatarLarger
		profile.filePath = v3data.UserInfo.User.UniqueID + "_larger.jpg"
	}
	profile.filePath = path.Join(profile.BaseDIR, profile.filePath)
	req, err := http.NewRequest("GET", photoURL, nil)
	if err != nil {
		return "", err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := profile.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	if err := saveTiktok(profile.filePath, resp); err != nil {
		return "", err
	}
	return profile.filePath, nil
}
