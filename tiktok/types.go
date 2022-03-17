package tiktok

import "net/http"

type ProfileData struct {
	ItemInfo struct {
		ItemStruct struct {
			Video struct {
				URL   string `json:"playAddr"`
				Cover string `json:"dynamicCover"`
			} `json:"video"`
			VideoID string `json:"id"`
			Author  struct {
				ID           string `json:"id"`
				UniqueID     string `json:"uniqueId"`
				Nickname     string `json:"nickname"`
				AvatarThumb  string `json:"avatarThumb"`
				AvatarMedium string `json:"avatarMedium"`
				AvatarLarger string `json:"avatarLarger"`
				Signature    string `json:"signature"`
				Verified     bool   `json:"verified"`
			} `json:"author"`
			CreatedTime int    `json:"createTime"`
			Description string `json:"desc"`
			VideoStats  struct {
				Likes    int `json:"diggCount"`
				Shares   int `json:"shareCount"`
				Comments int `json:"commentCount"`
				Played   int `json:"playCount"`
			} `json:"stats"`
			AuthorStats struct {
				Followings int `json:"followingCount"`
				Followers  int `json:"followerCount"`
				Hearts     int `json:"heartCount"`
				Videos     int `json:"videoCount"`
			} `json:"authorStats"`
		} `json:"itemStruct"`
	} `json:"itemInfo"`
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
		UserStats struct {
			Followings int `json:"followingCount"`
			Followers  int `json:"followerCount"`
			Hearts     int `json:"heartCount"`
			Videos     int `json:"videoCount"`
		} `json:"stats"`
	} `json:"userInfo"`
	UserMetaParams struct {
		Title       string `json:"title"`
		Keywords    string `json:"keywords"`
		Description string `json:"description"`
		URL         string `json:"canonicalHref"`
	} `json:"metaParams"`
}

type UserModule struct {
	Users struct {
		User struct {
			ID             string `json:"id"`
			ShortID        string `json:"shortId"`
			UniqueID       string `json:"uniqueId"`
			Nickname       string `json:"nickname"`
			AvatarLarger   string `json:"avatarLarger"`
			AvatarMedium   string `json:"avatarMedium"`
			AvatarThumb    string `json:"avatarThumb"`
			Signature      string `json:"signature"`
			CreateTime     int    `json:"createTime"`
			Verified       bool   `json:"verified"`
			SecUID         string `json:"secUid"`
			Ftc            bool   `json:"ftc"`
			Relation       int    `json:"relation"`
			OpenFavorite   bool   `json:"openFavorite"`
			CommentSetting int    `json:"commentSetting"`
			DuetSetting    int    `json:"duetSetting"`
			StitchSetting  int    `json:"stitchSetting"`
			PrivateAccount bool   `json:"privateAccount"`
			Secret         bool   `json:"secret"`
			IsADVirtual    bool   `json:"isADVirtual"`
			RoomID         string `json:"roomId"`
		} `json:"user"`
	} `json:"users"`
	Stats struct {
		User struct {
			FollowerCount  int  `json:"followerCount"`
			FollowingCount int  `json:"followingCount"`
			Heart          int  `json:"heart"`
			HeartCount     int  `json:"heartCount"`
			VideoCount     int  `json:"videoCount"`
			DiggCount      int  `json:"diggCount"`
			NeedFix        bool `json:"needFix"`
		} `json:"user"`
	} `json:"stats"`
}

type SEO struct {
	MetaParams struct {
		Title            string `json:"title"`
		Keywords         string `json:"keywords"`
		Description      string `json:"description"`
		CanonicalHref    string `json:"canonicalHref"`
		RobotsContent    string `json:"robotsContent"`
		ApplicableDevice string `json:"applicableDevice"`
	} `json:"metaParams"`
	JsonldList [][]interface{} `json:"jsonldList"`
	PageType   int             `json:"pageType"`
}

type VideoData struct {
	Props struct {
		ProfileData `json:"pageProps"`
	} `json:"props"`
}

type VideoDataV2 struct {
	SEO        `json:"SEO"`
	ItemModule struct {
		Item struct {
			ID           string `json:"id"`
			Desc         string `json:"desc"`
			CreateTime   string `json:"createTime"`
			ScheduleTime int    `json:"scheduleTime"`
			Video        struct {
				ID            string   `json:"id"`
				Height        int      `json:"height"`
				Width         int      `json:"width"`
				Duration      int      `json:"duration"`
				Ratio         string   `json:"ratio"`
				Cover         string   `json:"cover"`
				OriginCover   string   `json:"originCover"`
				DynamicCover  string   `json:"dynamicCover"`
				PlayAddr      string   `json:"playAddr"`
				DownloadAddr  string   `json:"downloadAddr"`
				ShareCover    []string `json:"shareCover"`
				ReflowCover   string   `json:"reflowCover"`
				Bitrate       int      `json:"bitrate"`
				EncodedType   string   `json:"encodedType"`
				Format        string   `json:"format"`
				VideoQuality  string   `json:"videoQuality"`
				EncodeUserTag string   `json:"encodeUserTag"`
				CodecType     string   `json:"codecType"`
				Definition    string   `json:"definition"`
			} `json:"video"`
			Author string `json:"author"`
			Music  struct {
				ID                 string `json:"id"`
				Title              string `json:"title"`
				PlayURL            string `json:"playUrl"`
				CoverLarge         string `json:"coverLarge"`
				CoverMedium        string `json:"coverMedium"`
				CoverThumb         string `json:"coverThumb"`
				AuthorName         string `json:"authorName"`
				Original           bool   `json:"original"`
				Duration           int    `json:"duration"`
				Album              string `json:"album"`
				ScheduleSearchTime int    `json:"scheduleSearchTime"`
			} `json:"music"`
			Challenges []struct {
				ID            string `json:"id"`
				Title         string `json:"title"`
				Desc          string `json:"desc"`
				ProfileLarger string `json:"profileLarger"`
				ProfileMedium string `json:"profileMedium"`
				ProfileThumb  string `json:"profileThumb"`
				CoverLarger   string `json:"coverLarger"`
				CoverMedium   string `json:"coverMedium"`
				CoverThumb    string `json:"coverThumb"`
				IsCommerce    bool   `json:"isCommerce"`
			} `json:"challenges"`
			Stats struct {
				DiggCount    int `json:"diggCount"`
				ShareCount   int `json:"shareCount"`
				CommentCount int `json:"commentCount"`
				PlayCount    int `json:"playCount"`
			} `json:"stats"`
			IsActivityItem bool `json:"isActivityItem"`
			DuetInfo       struct {
				DuetFromID string `json:"duetFromId"`
			} `json:"duetInfo"`
			WarnInfo     []interface{} `json:"warnInfo"`
			OriginalItem bool          `json:"originalItem"`
			OfficalItem  bool          `json:"officalItem"`
			TextExtra    []struct {
				AwemeID      string `json:"awemeId"`
				Start        int    `json:"start"`
				End          int    `json:"end"`
				HashtagID    string `json:"hashtagId"`
				HashtagName  string `json:"hashtagName"`
				Type         int    `json:"type"`
				SubType      int    `json:"subType"`
				UserID       string `json:"userId"`
				IsCommerce   bool   `json:"isCommerce"`
				UserUniqueID string `json:"userUniqueId"`
				SecUID       string `json:"secUid"`
			} `json:"textExtra"`
			Secret            bool          `json:"secret"`
			ForFriend         bool          `json:"forFriend"`
			Digged            bool          `json:"digged"`
			ItemCommentStatus int           `json:"itemCommentStatus"`
			ShowNotPass       bool          `json:"showNotPass"`
			Vl1               bool          `json:"vl1"`
			TakeDown          int           `json:"takeDown"`
			ItemMute          bool          `json:"itemMute"`
			EffectStickers    []interface{} `json:"effectStickers"`
			AuthorStats       struct {
				FollowerCount  int `json:"followerCount"`
				FollowingCount int `json:"followingCount"`
				Heart          int `json:"heart"`
				HeartCount     int `json:"heartCount"`
				VideoCount     int `json:"videoCount"`
				DiggCount      int `json:"diggCount"`
			} `json:"authorStats"`
			PrivateItem           bool          `json:"privateItem"`
			DuetEnabled           bool          `json:"duetEnabled"`
			StitchEnabled         bool          `json:"stitchEnabled"`
			StickersOnItem        []interface{} `json:"stickersOnItem"`
			IsAd                  bool          `json:"isAd"`
			ShareEnabled          bool          `json:"shareEnabled"`
			Comments              []interface{} `json:"comments"`
			DuetDisplay           int           `json:"duetDisplay"`
			StitchDisplay         int           `json:"stitchDisplay"`
			IndexEnabled          bool          `json:"indexEnabled"`
			DiversificationLabels []string      `json:"diversificationLabels"`
			Nickname              string        `json:"nickname"`
			AuthorID              string        `json:"authorId"`
			AuthorSecID           string        `json:"authorSecId"`
			AvatarThumb           string        `json:"avatarThumb"`
		} `json:"item"`
	} `json:"ItemModule"`
	UserModule `json:"UserModule"`
}

type ProfileDataV2 struct {
	SEO        `json:"SEO"`
	UserModule `json:"UserModule"`
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
