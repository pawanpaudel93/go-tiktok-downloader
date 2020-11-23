package main

import (
	"fmt"

	"github.com/pawanpaudel93/go-tiktok-downloader/tiktok"
)

func logError(err error) {
	if err != nil {
		fmt.Println("ERROR", err)
	}
}

func main() {
	baseDIR := "./downloads"
	video := tiktok.Video{URL: "https://www.tiktok.com/@manjuxettri07/video/6886732628280593666", BaseDIR: baseDIR}
	err := video.FetchInfo()
	if err == nil {
		_, err = video.Download()
		if err == nil {
			fmt.Println("Video Downloaded Successfully!!!")
		} else {
			logError(err)
		}
	} else {
		logError(err)
	}

	profile := tiktok.Profile{URL: "https://www.tiktok.com/@manjuxettri07", BaseDIR: baseDIR}
	err = profile.FetchInfo()
	if err == nil {
		_, err = profile.DownloadPhoto("large")
		if err == nil {
			fmt.Println("Profile Image Downloaded Successfully!!!")
		} else {
			logError(err)
		}
	} else {
		logError(err)
	}
}
