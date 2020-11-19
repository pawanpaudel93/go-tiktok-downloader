package main

import (
	"fmt"

	"github.com/pawanpaudel93/go-tiktok-downloader/downloader"
)

func main() {
	videoDownloader := downloader.TiktokVideo{URL: "https://www.tiktok.com/@manjuxettri07/video/6886732628280593666", FilePath: "temp.mp4", UseProxy: false}
	videoDownloader.Download("temp.mp4")
	fmt.Println(videoDownloader.GetVideoInfo())
	photoDownloader := downloader.TiktokProfile{URL: "https://www.tiktok.com/@manjuxettri07", FilePath: "temp.jpg", UseProxy: false}
	print(photoDownloader.GetProfilePicture())
	photoDownloader.DownloadPhoto("full", "temp.jpg")
}
