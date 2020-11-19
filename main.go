package main

import (
	"fmt"

	"github.com/pawanpaudel93/go-tiktok-downloader/tiktok"
)

func main() {
	baseDIR := "./downloads"
	video := tiktok.Video{URL: "https://www.tiktok.com/@manjuxettri07/video/6886732628280593666", UseProxy: false, BaseDIR: baseDIR}
	_, err := video.GetVideoInfo()
	if err != nil {
		err = video.Download()
		if err == nil {
			fmt.Println("Video Downloaded Successfully!!!")
		} else {
			fmt.Println("ERROR: ", err)
		}
	}

	profile := tiktok.Profile{URL: "https://www.tiktok.com/@manjuxettri07", UseProxy: false, BaseDIR: baseDIR}
	_, err = profile.GetProfileInfo()
	if err != nil {
		err = profile.DownloadPhoto("large")
		if err == nil {
			fmt.Println("Profile Image Downloaded Successfully!!!")
		} else {
			fmt.Println("ERROR: ", err)
		}
	}
}
