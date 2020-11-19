package main

import (
	"github.com/pawanpaudel93/go-tiktok-downloader/tiktok"
)

func main() {
	baseDIR := "./downloads"
	video := tiktok.Video{URL: "https://www.tiktok.com/@manjuxettri07/video/6886732628280593666", UseProxy: false, BaseDIR: baseDIR}
	_, _ = video.GetVideoInfo()
	_ = video.Download()
	profile := tiktok.Profile{URL: "https://www.tiktok.com/@manjuxettri07", UseProxy: false, BaseDIR: baseDIR}
	_, _ = profile.GetProfileInfo()
	profile.DownloadPhoto("large")
}
