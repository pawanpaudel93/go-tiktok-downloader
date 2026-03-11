<h1 align="center">Welcome to go-tiktok-downloader</h1>
<p>
  <img alt="Version" src="https://img.shields.io/badge/version-0.2.0-blue.svg?cacheSeconds=2592000" />
</p>

> A go package to download tiktok videos and profile pictures.

Also, Checkout the TikTok downloader webapp [govue-tiktok-downloader](https://github.com/pawanpaudel93/govue-tiktok-downloader)

### [Homepage](https://github.com/pawanpaudel93/go-tiktok-downloader)

## Install
> go get github.com/pawanpaudel93/go-tiktok-downloader

## Usage
- [Video](#video)
     - [FetchInfo](#fetchinfo)
     - [Download](#download)
	 - [GetInfo](#getinfo)
- [Profile](#profile)
     - [FetchInfo](#fetchinfo1)
     - [DownloadPhoto](#downloadphoto)
     - [GetProfileInfo](#getprofileinfo)
	 - [GetPPInfo](#getppinfo)

### Video
```go
import "github.com/pawanpaudel93/go-tiktok-downloader/tiktok"
video := tiktok.Video{URL: "https://www.tiktok.com/@berywambeatbox/video/6897238157025086721", BaseDIR: "./downloads", Proxy: ""}
```
You can also provide proxy with authentication (username:password@host:port) or simply (host:port) or leave it if you don't want to use proxy.

### ``FetchInfo``
```go
err := video.FetchInfo()
```
Fetches video metadata from the TikTok page. Must be called before `Download` or `GetInfo`. The download directory is created automatically if it doesn't exist.

### ``Download``
```go
videoPath, err := video.Download()
```
Downloads the video file and returns the path where it was saved. Requires `FetchInfo` to be called first.

### ``GetInfo``
```go
videoJSON, err := video.GetInfo()
```
Returns video information as JSON including: ID, URL, likes, shares, comments, played, createdTime, description, cover, and author details (uniqueID, nickname, url, followers, following, hearts, videos).

### Profile
```go
profile := tiktok.Profile{URL: "https://www.tiktok.com/@berywambeatbox", BaseDIR: "./downloads", Proxy: ""}
```
You can also provide proxy with authentication (username:password@host:port) or simply (host:port) or leave it if you don't want to use proxy.

### ``FetchInfo``
```go
err := profile.FetchInfo()
```
Fetches profile metadata from the TikTok page. Must be called before `DownloadPhoto`, `GetProfileInfo`, or `GetPPInfo`.

### ``DownloadPhoto``
```go
photoPath, err := profile.DownloadPhoto("medium")
```
Downloads the profile picture. PhotoType accepts `thumbnail`, `medium`, or any other string for the large version. Requires `FetchInfo` to be called first.

### ``GetProfileInfo``
```go
profileJSON, err := profile.GetProfileInfo()
```
Returns profile information as JSON including: user details (id, uniqueId, nickname, avatarThumb, avatarMedium, avatarLarger, signature, verified), user stats (followingCount, followerCount, heartCount, videoCount), and meta info (title, description).

### ``GetPPInfo``
```go
photoJSON, err := profile.GetPPInfo()
```
Returns profile picture URLs (Larger, Medium, Thumbnail) as JSON.

## Example

```go
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
}

```

## Author

**Pawan Paudel**

* Github: [@pawanpaudel93](https://github.com/pawanpaudel93)

## Contributing

Contributions, issues and feature requests are welcome!<br />Feel free to check [issues page](https://github.com/pawanpaudel93/go-tiktok-downloader/issues).

## Show your support

Give a star if this project helped you!

Copyright © 2020 [Pawan Paudel](https://github.com/pawanpaudel93).<br />
