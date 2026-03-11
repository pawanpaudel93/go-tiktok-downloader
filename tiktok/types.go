package tiktok

import (
	"encoding/json"
	"net/http"
)

// UniversalData wraps the __UNIVERSAL_DATA_FOR_REHYDRATION__ JSON structure
type UniversalData struct {
	DefaultScope struct {
		VideoDetail *VideoDataV3   `json:"webapp.video-detail"`
		UserDetail  *ProfileDataV3 `json:"webapp.user-detail"`
	} `json:"__DEFAULT_SCOPE__"`
}

type VideoDataV3 struct {
	ItemInfo struct {
		ItemStruct struct {
			ID         string      `json:"id"`
			Desc       string      `json:"desc"`
			CreateTime json.Number `json:"createTime"`
			Video      struct {
				ID           string `json:"id"`
				PlayAddr     string `json:"playAddr"`
				DownloadAddr string `json:"downloadAddr"`
				Cover        string `json:"cover"`
				DynamicCover string `json:"dynamicCover"`
			} `json:"video"`
			Author struct {
				ID           string `json:"id"`
				UniqueID     string `json:"uniqueId"`
				Nickname     string `json:"nickname"`
				AvatarThumb  string `json:"avatarThumb"`
				AvatarMedium string `json:"avatarMedium"`
				AvatarLarger string `json:"avatarLarger"`
				Signature    string `json:"signature"`
				Verified     bool   `json:"verified"`
			} `json:"author"`
			Stats struct {
				DiggCount    int `json:"diggCount"`
				ShareCount   int `json:"shareCount"`
				CommentCount int `json:"commentCount"`
				PlayCount    int `json:"playCount"`
			} `json:"stats"`
			AuthorStats struct {
				FollowerCount  int `json:"followerCount"`
				FollowingCount int `json:"followingCount"`
				HeartCount     int `json:"heartCount"`
				VideoCount     int `json:"videoCount"`
			} `json:"authorStats"`
		} `json:"itemStruct"`
	} `json:"itemInfo"`
	ShareMeta struct {
		Title string `json:"title"`
		Desc  string `json:"desc"`
	} `json:"shareMeta"`
	StatusCode int `json:"statusCode"`
}

type ProfileDataV3 struct {
	UserInfo struct {
		User struct {
			ID           string `json:"id"`
			UniqueID     string `json:"uniqueId"`
			Nickname     string `json:"nickname"`
			AvatarThumb  string `json:"avatarThumb"`
			AvatarMedium string `json:"avatarMedium"`
			AvatarLarger string `json:"avatarLarger"`
			Signature    string `json:"signature"`
			Verified     bool   `json:"verified"`
		} `json:"user"`
		Stats struct {
			FollowerCount  int `json:"followerCount"`
			FollowingCount int `json:"followingCount"`
			Heart          int `json:"heart"`
			HeartCount     int `json:"heartCount"`
			VideoCount     int `json:"videoCount"`
		} `json:"stats"`
	} `json:"userInfo"`
	ShareMeta struct {
		Title string `json:"title"`
		Desc  string `json:"desc"`
	} `json:"shareMeta"`
	StatusCode int `json:"statusCode"`
}

type Base struct {
	URL        string
	filePath   string
	data       interface{}
	httpClient *http.Client
	Proxy      string
	BaseDIR    string
}

type Profile Base

type Video Base
