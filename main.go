package main

import (
	"fmt"

	downloader "github.com/pawanpaudel93/go-tiktok-downloader/downloader"
)

func main() {
	tiktokDownloader := downloader.TikTok{URL: "https://www.tiktok.com/@manjuxettri07/video/6886732628280593666", FilePath: "temp.mp4", UseProxy: true}
	tiktokDownloader.Download()
	fmt.Println(tiktokDownloader.GetTiktokInfo())
}
