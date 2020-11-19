package tiktok

type VideoInfo struct {
	URL   string `json:"playAddr"`
	Cover string `json:"dynamicCover"`
}

type Author struct {
	ID           string `json:"id"`
	UniqueID     string `json:"uniqueId"`
	Nickname     string `json:"nickname"`
	AvatarThumb  string `json:"avatarThumb"`
	AvatarMedium string `json:"avatarMedium"`
	AvatarLarger string `json:"avatarLarger"`
	Signature    string `json:"signature"`
	Verified     bool   `json:"verified"`
}

type User Author

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

type UserStats AuthorStats

type ItemStruct struct {
	Video       VideoInfo `json:"video"`
	VideoID     string    `json:"id"`
	Author      `json:"author"`
	CreatedTime int    `json:"createTime"`
	Description string `json:"desc"`
	VideoStats  `json:"stats"`
	AuthorStats `json:"authorStats"`
}

type ItemInfo struct {
	ItemStruct `json:"itemStruct"`
}

type UserMetaParams struct {
	Title       string `json:"title"`
	Keywords    string `json:"keywords"`
	Description string `json:"description"`
	URL         string `json:"canonicalHref"`
}

type UserInfo struct {
	User      `json:"user"`
	UserStats `json:"stats"`
}

type PageProps struct {
	ItemInfo       `json:"itemInfo"`
	UserInfo       `json:"userInfo"`
	UserMetaParams `json:"metaParams"`
}

type Props struct {
	PageProps `json:"pageProps"`
}

type VideoData struct {
	Props `json:"props"`
}
