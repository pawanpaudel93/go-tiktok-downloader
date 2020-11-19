package downloader

type Video struct {
	URL   string `json:"playAddr"`
	cover string `json:"dynamicCover"`
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

type UserMetaParams struct {
	Title       string `json:"title"`
	keywords    string `json:"keywords"`
	description string `json:"description"`
	URL         string `json:"canonicalHref"`
}

type UserInfo struct {
	Author      `json:"user"`
	AuthorStats `json:"stats"`
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
