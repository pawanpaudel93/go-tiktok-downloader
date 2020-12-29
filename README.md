<h1 align="center">Welcome to go-tiktok-downloader</h1>
<p>
  <img alt="Version" src="https://img.shields.io/badge/version-0.1.0-blue.svg?cacheSeconds=2592000" />
</p>

> A go package to download tiktok videos and profile pictures. 

Also, Checkout the TikTok downloader webapp [govue-tiktok-downloader](https://github.com/pawanpaudel93/govue-tiktok-downloader)

### 🏠 [Homepage](https://github.com/pawanpaudel93/go-tiktok-downloader)

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
It fetches the required info about video from tiktok website and returns error if any else nil.

### ``Download``
```go
videoPath, err := video.Download()
```
It downloads the video and returns the path where its saved and err.

### ``GetInfo``
```go
videoJSON, err := video.GetInfo()
```
It returns the videoJSON and err, VideoJSON contains information about the video like ID, URL, likes, shares, comments, played, createdTime, description, cover and about author like uniqueID, nickname, url, followers, following, hearts, and videos.

### Profile
```go
profile := tiktok.Profile{URL: "https://www.tiktok.com/@berywambeatbox", BaseDIR: "./downloads", Proxy: ""}
```
You can also provide proxy with authentication (username:password@host:port) or simply (host:port) or leave it if you don't want to use proxy.

### ``FetchInfo``
```go
err := profile.FetchInfo()
```
It fetches the required info about profile from tiktok website and returns error if any else nil.

### ``DownloadPhoto``
```go
photoPath, err := profile.DownloadPhoto("medium")
```
DownloadPhtot takes argument for parameter PhotoType which takes either medium, thumbnail or any string for large photo.
It downloads the profile picture and returns the path where its saved and err.

### ``GetProfileInfo``
```go
profileJSON, err := profile.GetProfileInfo()
```
It returns the profileJSON and err. profileJSON contains information about tiktok user(id, uniqueId, nickname, avararThumb, avatarMedium, avatarLarger, signature, verified), user stats(followingCount, followerCount, heartCount, videoCount) and user meta params(title, keywords, description, canonicalHref).

### ``GetPPInfo``
```go
photoJSON, err := profile.GetPPInfo()
```
It returns tiktok user profile picture urls (Larger, Medium, Thumbnail) and err.

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

👤 **Pawan Paudel**

* Github: [@pawanpaudel93](https://github.com/pawanpaudel93)

## 🤝 Contributing

Contributions, issues and feature requests are welcome!<br />Feel free to check [issues page](https://github.com/pawanpaudel93/go-tiktok-downloader/issues). 

## Show your support

Give a ⭐️ if this project helped you!

Copyright © 2020 [Pawan Paudel](https://github.com/pawanpaudel93).<br />
