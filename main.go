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
	video := tiktok.Video{URL: "https://www.tiktok.com/@manjuxettri07/video/6886732628280593666", UseProxy: false, BaseDIR: baseDIR}
	err := video.FetchInfo()
	if err == nil {
		err = video.Download()
		if err == nil {
			fmt.Println("Video Downloaded Successfully!!!")
		} else {
			logError(err)
		}
	} else {
		fmt.Println("error")
		logError(err)
	}

	profile := tiktok.Profile{URL: "https://www.tiktok.com/@manjuxettri07", UseProxy: false, BaseDIR: baseDIR}
	err = profile.FetchInfo()
	if err == nil {
		err = profile.DownloadPhoto("large")
		if err == nil {
			fmt.Println("Profile Image Downloaded Successfully!!!")
		} else {
			logError(err)
		}
	} else {
		logError(err)
	}
}
