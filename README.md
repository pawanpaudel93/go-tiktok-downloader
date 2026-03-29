<h1 align="center">Welcome to go-tiktok-downloader</h1>
<p>
  <img alt="Version" src="https://img.shields.io/badge/version-0.3.0-blue.svg?cacheSeconds=2592000" />
</p>

> A Go package, CLI tool, and MCP server to download TikTok videos and profile pictures.

Also, Checkout the TikTok downloader webapp [govue-tiktok-downloader](https://github.com/pawanpaudel93/govue-tiktok-downloader)

### [Homepage](https://github.com/pawanpaudel93/go-tiktok-downloader)

## Install

**As a Go package:**
```bash
go get github.com/pawanpaudel93/go-tiktok-downloader
```

**As a CLI tool:**
```bash
go install github.com/pawanpaudel93/go-tiktok-downloader/cmd/tiktok-dl@latest
```

**As an MCP server (for Claude Code, Codex, etc.):**
```bash
go install github.com/pawanpaudel93/go-tiktok-downloader/cmd/tiktok-dl-mcp@latest
```

## CLI Usage

```
tiktok-dl video [flags] <url>
tiktok-dl profile [flags] <url>
```

### Download a video
```bash
tiktok-dl video https://www.tiktok.com/@user/video/1234567890
```

### Get video info as JSON
```bash
tiktok-dl video --info https://www.tiktok.com/@user/video/1234567890
```

### Download a profile photo
```bash
tiktok-dl profile https://www.tiktok.com/@user
tiktok-dl profile --photo thumbnail https://www.tiktok.com/@user
```

### Get profile info as JSON
```bash
tiktok-dl profile --info https://www.tiktok.com/@user
```

### Flags
| Flag | Description | Default |
|------|-------------|---------|
| `--info` | Print info as JSON instead of downloading | `false` |
| `-o`, `--output` | Output directory | `./downloads` |
| `--proxy` | Proxy URL (`host:port` or `user:pass@host:port`) | |
| `--photo` | Photo size: `thumbnail`, `medium`, `large` (profile only) | `large` |

## MCP Server

The MCP server exposes 4 tools for AI agents (Claude Code, Codex CLI, Cursor, etc.):

| Tool | Description |
|------|-------------|
| `download_video` | Download a TikTok video, returns file path |
| `get_video_info` | Get video metadata (likes, shares, author, etc.) |
| `download_profile_photo` | Download a profile picture, returns file path |
| `get_profile_info` | Get profile metadata (followers, bio, etc.) |

### Claude Code

Add to `~/.claude/settings.json`:
```json
{
  "mcpServers": {
    "tiktok": {
      "command": "tiktok-dl-mcp",
      "args": []
    }
  }
}
```

Or if running from source:
```json
{
  "mcpServers": {
    "tiktok": {
      "command": "go",
      "args": ["run", "./cmd/tiktok-dl-mcp"]
    }
  }
}
```

### Codex CLI

Add to your `codex` MCP config:
```json
{
  "mcpServers": {
    "tiktok": {
      "command": "tiktok-dl-mcp"
    }
  }
}
```

## Go Package Usage

- [Video](#video)
     - [FetchInfo](#fetchinfo)
     - [Download](#download)
     - [GetInfo](#getinfo)
- [Profile](#profile)
     - [FetchInfo](#fetchinfo-1)
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

### ``FetchInfo``
```go
err := profile.FetchInfo()
```
Fetches profile metadata from the TikTok page. Must be called before `DownloadPhoto`, `GetProfileInfo`, or `GetPPInfo`.

### ``DownloadPhoto``
```go
photoPath, err := profile.DownloadPhoto("medium")
```
Downloads the profile picture. PhotoType accepts `thumbnail`, `medium`, or any other string for the large version.

### ``GetProfileInfo``
```go
profileJSON, err := profile.GetProfileInfo()
```
Returns profile information as JSON including: user details (id, uniqueId, nickname, avatars, signature, verified), user stats (followingCount, followerCount, heartCount, videoCount), and meta info (title, description).

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

func main() {
	baseDIR := "./downloads"
	video := tiktok.Video{URL: "https://www.tiktok.com/@berywambeatbox/video/6897238157025086721", BaseDIR: baseDIR}
	if err := video.FetchInfo(); err != nil {
		fmt.Println("Error:", err)
		return
	}
	if path, err := video.Download(); err == nil {
		fmt.Println("Video saved to:", path)
	}

	profile := tiktok.Profile{URL: "https://www.tiktok.com/@berywambeatbox", BaseDIR: baseDIR}
	if err := profile.FetchInfo(); err != nil {
		fmt.Println("Error:", err)
		return
	}
	if path, err := profile.DownloadPhoto("large"); err == nil {
		fmt.Println("Photo saved to:", path)
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
