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
	var URL string
	jar, _ := cookiejar.New(nil)
	switch video.data.(type) {
	case VideoData:
		URL = video.data.(VideoData).Props.ProfileData.ItemInfo.ItemStruct.Video.URL
	case VideoDataV2:
		URL = video.data.(VideoDataV2).ItemModule.Item.Video.PlayAddr
	default:
		return "", errors.New("Invalid Tiktok Video Data")
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
	videoData := VideoData{}
	videoDataV2 := VideoDataV2{}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}
	doc.Find("#__NEXT_DATA__").Each(func(i int, s *goquery.Selection) {
		err = json.Unmarshal([]byte(s.Text()), &videoData)
		video.data = videoData
	})
	if _, ok := video.data.(VideoData); !ok {
		doc.Find("#sigi-persisted-data").Each(func(i int, s *goquery.Selection) {
			data := strings.Replace(s.Text(), "window['SIGI_STATE']=", "", -1)
			data = strings.Replace(data, fmt.Sprintf("\"%s\":", videoId), "\"item\":", -1)
			data = strings.Replace(data, fmt.Sprintf("\"%s\":", username), "\"user\":", -1)
			err = json.Unmarshal([]byte(strings.Split(data, ";window['SIGI_RETRY']")[0]), &videoDataV2)
			video.data = videoDataV2
		})
		if _, ok := video.data.(VideoDataV2); !ok {
			return errors.New("Video Data Not Found")
		}
	}

	return err
}

// GetInfo returns Tiktok video information
func (video *Video) GetInfo() (string, error) {
	var tiktokData map[string]interface{}
	switch video.data.(type) {
	case VideoData:
		data := video.data.(VideoData)
		tiktokData = map[string]interface{}{
			"video": map[string]interface{}{
				"ID":          data.Props.ProfileData.ItemInfo.ItemStruct.VideoID,
				"URL":         video.URL,
				"likes":       data.Props.ProfileData.ItemInfo.ItemStruct.VideoStats.Likes,
				"shares":      data.Props.ProfileData.ItemInfo.ItemStruct.VideoStats.Shares,
				"comments":    data.Props.ProfileData.ItemInfo.ItemStruct.VideoStats.Comments,
				"played":      data.Props.ProfileData.ItemInfo.ItemStruct.VideoStats.Played,
				"createdTime": data.Props.ProfileData.ItemInfo.ItemStruct.CreatedTime,
				"description": data.Props.ProfileData.ItemInfo.ItemStruct.Description,
				"cover":       data.Props.ProfileData.ItemInfo.ItemStruct.Video.Cover,
			},
			"author": map[string]interface{}{
				"uniqueID":  data.Props.ProfileData.ItemInfo.ItemStruct.Author.UniqueID,
				"nickname":  data.Props.ProfileData.ItemInfo.ItemStruct.Author.Nickname,
				"url":       "https://tiktok.com/@" + data.Props.ProfileData.ItemInfo.ItemStruct.Author.UniqueID,
				"followers": data.Props.ProfileData.ItemInfo.ItemStruct.AuthorStats.Followers,
				"following": data.Props.ProfileData.ItemInfo.ItemStruct.AuthorStats.Followings,
				"hearts":    data.Props.ProfileData.ItemInfo.ItemStruct.AuthorStats.Hearts,
				"videos":    data.Props.ProfileData.ItemInfo.ItemStruct.AuthorStats.Videos,
			},
		}
	case VideoDataV2:
		data := video.data.(VideoDataV2)
		tiktokData = map[string]interface{}{
			"video": map[string]interface{}{
				"ID":          data.ItemModule.Item.Video.ID,
				"URL":         video.URL,
				"likes":       data.ItemModule.Item.Stats.DiggCount,
				"shares":      data.ItemModule.Item.Stats.ShareCount,
				"comments":    data.ItemModule.Item.Stats.CommentCount,
				"played":      data.ItemModule.Item.Stats.PlayCount,
				"createdTime": data.ItemModule.Item.CreateTime,
				"description": data.ItemModule.Item.Desc,
				"cover":       data.ItemModule.Item.Video.Cover,
			},
			"author": map[string]interface{}{
				"uniqueID":  data.ItemModule.Item.AuthorID,
				"nickname":  data.ItemModule.Item.Nickname,
				"url":       "https://tiktok.com/@" + data.ItemModule.Item.AuthorID,
				"followers": data.ItemModule.Item.AuthorStats.FollowerCount,
				"following": data.ItemModule.Item.AuthorStats.FollowingCount,
				"hearts":    data.ItemModule.Item.AuthorStats.HeartCount,
				"videos":    data.ItemModule.Item.AuthorStats.VideoCount,
			},
		}
	default:
		return "", errors.New("Video Data Not Found")
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
	videoData := VideoData{}
	profileDataV2 := ProfileDataV2{}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}
	doc.Find("#__NEXT_DATA__").Each(func(i int, s *goquery.Selection) {
		err = json.Unmarshal([]byte(s.Text()), &videoData)
		profile.data = videoData.Props.ProfileData
	})
	if _, ok := profile.data.(ProfileData); !ok {
		doc.Find("#sigi-persisted-data").Each(func(i int, s *goquery.Selection) {
			data := strings.Replace(s.Text(), "window['SIGI_STATE']=", "", -1)
			data = strings.Replace(data, fmt.Sprintf("\"%s\":", username), "\"user\":", -1)
			err = json.Unmarshal([]byte(strings.Split(data, ";window['SIGI_RETRY']")[0]), &profileDataV2)
			profile.data = profileDataV2
		})
		if _, ok := profile.data.(ProfileDataV2); !ok {
			return errors.New("Failed to fetch profile information")
		}
	}
	return err
}

// GetPPInfo returns Tiktok profile picture information.
func (profile *Profile) GetPPInfo() (string, error) {
	var photoData map[string]string
	switch profile.data.(type) {
	case ProfileData:
		data := profile.data.(ProfileData)
		photoData = map[string]string{
			"Thumbnail": replaceUnicode(data.ItemInfo.ItemStruct.Author.AvatarThumb),
			"Medium":    replaceUnicode(data.ItemInfo.ItemStruct.Author.AvatarMedium),
			"Larger":    replaceUnicode(data.ItemInfo.ItemStruct.Author.AvatarLarger),
		}
	case ProfileDataV2:
		data := profile.data.(ProfileDataV2)
		photoData = map[string]string{
			"Thumbnail": replaceUnicode(data.UserModule.Users.User.AvatarThumb),
			"Medium":    replaceUnicode(data.UserModule.Users.User.AvatarMedium),
			"Larger":    replaceUnicode(data.UserModule.Users.User.AvatarLarger),
		}
	default:
		return "", errors.New("Profile Data Not Found")
	}
	data, err := json.Marshal(photoData)
	return string(data), err
}

// GetProfileInfo returns Tiktok profile information.
func (profile *Profile) GetProfileInfo() (string, error) {
	var profileData map[string]interface{}
	switch profile.data.(type) {
	case ProfileData:
		data := profile.data.(ProfileData)
		profileData = map[string]interface{}{
			"user":          data.UserInfo.User,
			"userStats":     data.UserInfo.UserStats,
			"userMetaStats": data.UserMetaParams,
		}
	case ProfileDataV2:
		data := profile.data.(ProfileDataV2)
		profileData = map[string]interface{}{
			"user":      data.UserModule.Users.User,
			"userStats": data.UserModule.Stats.User,
			"userMetaStats": map[string]interface{}{
				"title":         data.SEO.MetaParams.Title,
				"keywords":      data.SEO.MetaParams.Keywords,
				"description":   data.SEO.MetaParams.Description,
				"canonicalHref": data.SEO.MetaParams.CanonicalHref,
			},
		}
	default:
		return "", errors.New("Profile Data Not Found")
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
		switch profile.data.(type) {
		case ProfileData:
			data := profile.data.(ProfileData)
			photoURL = data.ItemInfo.ItemStruct.Author.AvatarThumb
			profile.filePath = data.UserInfo.User.UniqueID + "_thumbnail.jpg"
		case ProfileDataV2:
			data := profile.data.(ProfileDataV2)
			photoURL = data.UserModule.Users.User.AvatarThumb
			profile.filePath = data.UserModule.Users.User.UniqueID + "_thumbnail.jpg"
		default:
			return "", errors.New("Profile Data Not Found")
		}
	case "medium":
		switch profile.data.(type) {
		case ProfileData:
			data := profile.data.(ProfileData)
			photoURL = data.ItemInfo.ItemStruct.Author.AvatarMedium
			profile.filePath = data.UserInfo.User.UniqueID + "_medium.jpg"
		case ProfileDataV2:
			data := profile.data.(ProfileDataV2)
			photoURL = data.UserModule.Users.User.AvatarMedium
			profile.filePath = data.UserModule.Users.User.UniqueID + "_medium.jpg"
		default:
			return "", errors.New("Profile Data Not Found")
		}
	default:
		switch profile.data.(type) {
		case ProfileData:
			data := profile.data.(ProfileData)
			photoURL = data.ItemInfo.ItemStruct.Author.AvatarLarger
			profile.filePath = data.UserInfo.User.UniqueID + "_larger.jpg"
		case ProfileDataV2:
			data := profile.data.(ProfileDataV2)
			photoURL = data.UserModule.Users.User.AvatarLarger
			profile.filePath = data.UserModule.Users.User.UniqueID + "_larger.jpg"
		default:
			return "", errors.New("Profile Data Not Found")
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
