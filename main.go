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
	video := tiktok.Video{URL: "https://www.tiktok.com/@berywambeatbox/video/6897238157025086721", BaseDIR: baseDIR}
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

	if info, err := video.GetInfo(); err == nil {
		fmt.Println("Video Info:", info)
	} else {
		logError(err)
	}

	fmt.Println("------------------------------------------------------")

	profile := tiktok.Profile{URL: "https://www.tiktok.com/@berywambeatbox", BaseDIR: baseDIR}
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
	if ppInfo, err := profile.GetPPInfo(); err == nil {
		fmt.Println("Profile Info:", ppInfo)
	} else {
		logError(err)
	}

	if pInfo, err := profile.GetProfileInfo(); err == nil {
		fmt.Println("Profile Info:", pInfo)
	} else {
		logError(err)
	}
}
