package downloader

type Video struct {
	URL string `json:"playAddr"`
}

type Author struct {
	UniqueID string `json:"uniqueId"`
	Nickname string `json:"nickname"`
}
type VideoStats struct {
	Likes    int `json:"diggCount"`
	Shares   int `json:"shareCount"`
	Comments int `json:"commentCount"`
	Played   int `json:"playCount"`
}

type AuthorStats struct {
	Followings int `json:"followingCount"`
	Followers  int `json:"followerCount"`
	Hearts     int `json:"heartCount"`
	Videos     int `json:"videoCount"`
}

type ItemStruct struct {
	Video       `json:"video"`
	Author      `json:"author"`
	CreatedTime int    `json:"createTime"`
	Description string `json:"desc"`
	VideoStats  `json:"stats"`
	AuthorStats `json:"authorStats"`
}

type ItemInfo struct {
	ItemStruct `json:"itemStruct"`
}

type PageProps struct {
	ItemInfo `json:"itemInfo"`
}

type Props struct {
	PageProps `json:"pageProps"`
}

type VideoData struct {
	Props `json:"props"`
}
