package tiktok

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"regexp"
	"strings"
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
	URL := video.data.Props.ProfileData.ItemInfo.ItemStruct.Video.URL
	if URL == "" {
		URL = video.dataV2.ItemModule.Item.Video.PlayAddr
	}
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
	regex := regexp.MustCompile(`(?m)tiktok\.com\/@([\w.-]+)\/video\/(\d+)/?`)
	regexResult := regex.FindStringSubmatch(videoURL)
	username, videoId := regexResult[1], regexResult[2]
	video.filePath = path.Join(video.BaseDIR, username+"_"+videoId+".mp4")
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
	VideoDataV2 := VideoDataV2{}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}
	doc.Find("#__NEXT_DATA__").Each(func(i int, s *goquery.Selection) {
		err = json.Unmarshal([]byte(s.Text()), &VideoData)
		video.data = VideoData
	})
	doc.Find("#sigi-persisted-data").Each(func(i int, s *goquery.Selection) {
		data := strings.Replace(s.Text(), "window['SIGI_STATE']=", "", -1)
		data = strings.Replace(data, fmt.Sprintf("\"%s\":", videoId), "\"item\":", -1)
		data = strings.Replace(data, fmt.Sprintf("\"%s\":", username), "\"user\":", -1)
		err = json.Unmarshal([]byte(strings.Split(data, ";window['SIGI_RETRY']")[0]), &VideoDataV2)
		video.dataV2 = VideoDataV2
	})

	return err
}

// GetInfo returns Tiktok video information
func (video *Video) GetInfo() (string, error) {
	var tiktokData map[string]interface{}
	if video.data.Props.ProfileData.ItemInfo.ItemStruct.Video.URL != "" {
		tiktokData = map[string]interface{}{
			"video": map[string]interface{}{
				"ID":          video.data.Props.ProfileData.ItemInfo.ItemStruct.VideoID,
				"URL":         video.URL,
				"likes":       video.data.Props.ProfileData.ItemInfo.ItemStruct.VideoStats.Likes,
				"shares":      video.data.Props.ProfileData.ItemInfo.ItemStruct.VideoStats.Shares,
				"comments":    video.data.Props.ProfileData.ItemInfo.ItemStruct.VideoStats.Comments,
				"played":      video.data.Props.ProfileData.ItemInfo.ItemStruct.VideoStats.Played,
				"createdTime": video.data.Props.ProfileData.ItemInfo.ItemStruct.CreatedTime,
				"description": video.data.Props.ProfileData.ItemInfo.ItemStruct.Description,
				"cover":       video.data.Props.ProfileData.ItemInfo.ItemStruct.Video.Cover,
			},
			"author": map[string]interface{}{
				"uniqueID":  video.data.Props.ProfileData.ItemInfo.ItemStruct.Author.UniqueID,
				"nickname":  video.data.Props.ProfileData.ItemInfo.ItemStruct.Author.Nickname,
				"url":       "https://tiktok.com/@" + video.data.Props.ProfileData.ItemInfo.ItemStruct.Author.UniqueID,
				"followers": video.data.Props.ProfileData.ItemInfo.ItemStruct.AuthorStats.Followers,
				"following": video.data.Props.ProfileData.ItemInfo.ItemStruct.AuthorStats.Followings,
				"hearts":    video.data.Props.ProfileData.ItemInfo.ItemStruct.AuthorStats.Hearts,
				"videos":    video.data.Props.ProfileData.ItemInfo.ItemStruct.AuthorStats.Videos,
			},
		}
	} else {
		tiktokData = map[string]interface{}{
			"video": map[string]interface{}{
				"ID":          video.dataV2.ItemModule.Item.Video.ID,
				"URL":         video.URL,
				"likes":       video.dataV2.ItemModule.Item.Stats.DiggCount,
				"shares":      video.dataV2.ItemModule.Item.Stats.ShareCount,
				"comments":    video.dataV2.ItemModule.Item.Stats.CommentCount,
				"played":      video.dataV2.ItemModule.Item.Stats.PlayCount,
				"createdTime": video.dataV2.ItemModule.Item.CreateTime,
				"description": video.dataV2.ItemModule.Item.Desc,
				"cover":       video.dataV2.ItemModule.Item.Video.Cover,
			},
			"author": map[string]interface{}{
				"uniqueID":  video.dataV2.ItemModule.Item.AuthorID,
				"nickname":  video.dataV2.ItemModule.Item.Nickname,
				"url":       "https://tiktok.com/@" + video.dataV2.ItemModule.Item.AuthorID,
				"followers": video.dataV2.ItemModule.Item.AuthorStats.FollowerCount,
				"following": video.dataV2.ItemModule.Item.AuthorStats.FollowingCount,
				"hearts":    video.dataV2.ItemModule.Item.AuthorStats.HeartCount,
				"videos":    video.dataV2.ItemModule.Item.AuthorStats.VideoCount,
			},
		}
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
	regex := regexp.MustCompile(`(?m)tiktok\.com\/@([\w.-]+)\/?`)
	regexResult := regex.FindStringSubmatch(profileURL)
	username := regexResult[1]
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
	profileDataV2 := ProfileDataV2{}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}
	doc.Find("#__NEXT_DATA__").Each(func(i int, s *goquery.Selection) {
		err = json.Unmarshal([]byte(s.Text()), &VideoData)
		profile.data = VideoData.Props.ProfileData
	})
	doc.Find("#sigi-persisted-data").Each(func(i int, s *goquery.Selection) {
		data := strings.Replace(s.Text(), "window['SIGI_STATE']=", "", -1)
		data = strings.Replace(data, fmt.Sprintf("\"%s\":", username), "\"user\":", -1)
		err = json.Unmarshal([]byte(strings.Split(data, ";window['SIGI_RETRY']")[0]), &profileDataV2)
		profile.dataV2 = profileDataV2
	})
	return err
}

// GetPPInfo returns Tiktok profile picture information.
func (profile *Profile) GetPPInfo() (string, error) {
	var photoData map[string]string
	if profile.data.UserInfo.User.ID != "" {
		photoData = map[string]string{
			"Thumbnail": replaceUnicode(profile.data.ItemInfo.ItemStruct.Author.AvatarThumb),
			"Medium":    replaceUnicode(profile.data.ItemInfo.ItemStruct.Author.AvatarMedium),
			"Larger":    replaceUnicode(profile.data.ItemInfo.ItemStruct.Author.AvatarLarger),
		}
	} else {
		photoData = map[string]string{
			"Thumbnail": replaceUnicode(profile.dataV2.UserModule.Users.User.AvatarThumb),
			"Medium":    replaceUnicode(profile.dataV2.UserModule.Users.User.AvatarMedium),
			"Larger":    replaceUnicode(profile.dataV2.UserModule.Users.User.AvatarLarger),
		}
	}
	data, err := json.Marshal(photoData)
	return string(data), err
}

// GetProfileInfo returns Tiktok profile information.
func (profile *Profile) GetProfileInfo() (string, error) {
	var profileData map[string]interface{}
	if profile.data.UserInfo.User.ID != "" {
		profileData = map[string]interface{}{
			"user":          profile.data.UserInfo.User,
			"userStats":     profile.data.UserInfo.UserStats,
			"userMetaStats": profile.data.UserMetaParams,
		}
	} else {
		profileData = map[string]interface{}{
			"user":      profile.dataV2.UserModule.Users.User,
			"userStats": profile.dataV2.UserModule.Stats.User,
			"userMetaStats": map[string]interface{}{
				"title":         profile.dataV2.SEO.MetaParams.Title,
				"keywords":      profile.dataV2.SEO.MetaParams.Keywords,
				"description":   profile.dataV2.SEO.MetaParams.Description,
				"canonicalHref": profile.dataV2.SEO.MetaParams.CanonicalHref,
			},
		}
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
		photoURL = profile.data.ItemInfo.ItemStruct.Author.AvatarThumb
		profile.filePath = profile.data.UserInfo.User.UniqueID + "_thumbnail.jpg"
		if photoURL == "" {
			photoURL = profile.dataV2.UserModule.Users.User.AvatarThumb
			profile.filePath = profile.dataV2.UserModule.Users.User.UniqueID + "_thumbnail.jpg"
		}
	case "medium":
		photoURL = profile.data.ItemInfo.ItemStruct.Author.AvatarMedium
		profile.filePath = profile.data.UserInfo.User.UniqueID + "_medium.jpg"
		if photoURL == "" {
			photoURL = profile.dataV2.UserModule.Users.User.AvatarMedium
			profile.filePath = profile.dataV2.UserModule.Users.User.UniqueID + "_medium.jpg"
		}

	default:
		photoURL = profile.data.ItemInfo.ItemStruct.Author.AvatarLarger
		profile.filePath = profile.data.UserInfo.User.UniqueID + "_large.jpg"
		if photoURL == "" {
			photoURL = profile.dataV2.UserModule.Users.User.AvatarLarger
			profile.filePath = profile.dataV2.UserModule.Users.User.UniqueID + "_larger.jpg"
		}
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
