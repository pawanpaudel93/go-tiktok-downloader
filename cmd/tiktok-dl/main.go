package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pawanpaudel93/go-tiktok-downloader/tiktok"
)

const usage = `Usage:
  tiktok-dl video [flags] <url>
  tiktok-dl profile [flags] <url>

Commands:
  video      Download a TikTok video or get its info
  profile    Download a profile photo or get profile info

Video flags:
  --info           Print video info as JSON instead of downloading
  -o, --output     Output directory (default: ./downloads)
  --proxy          Proxy URL (host:port or user:pass@host:port)

Profile flags:
  --info           Print profile info as JSON instead of downloading
  --photo          Photo size: thumbnail, medium, large (default: large)
  -o, --output     Output directory (default: ./downloads)
  --proxy          Proxy URL (host:port or user:pass@host:port)
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "video":
		runVideo(os.Args[2:])
	case "profile":
		runProfile(os.Args[2:])
	case "help", "--help", "-h":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}
}

func runVideo(args []string) {
	fs := flag.NewFlagSet("video", flag.ExitOnError)
	info := fs.Bool("info", false, "Print video info as JSON")
	output := fs.String("output", "./downloads", "Output directory")
	fs.StringVar(output, "o", "./downloads", "Output directory")
	proxy := fs.String("proxy", "", "Proxy URL")
	fs.Parse(args)

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Error: video URL required")
		os.Exit(1)
	}
	url := fs.Arg(0)

	video := tiktok.Video{URL: url, BaseDIR: *output, Proxy: *proxy}
	if err := video.FetchInfo(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if *info {
		data, err := video.GetInfo()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(data)
		return
	}

	path, err := video.Download()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(path)
}

func runProfile(args []string) {
	fs := flag.NewFlagSet("profile", flag.ExitOnError)
	info := fs.Bool("info", false, "Print profile info as JSON")
	output := fs.String("output", "./downloads", "Output directory")
	fs.StringVar(output, "o", "./downloads", "Output directory")
	proxy := fs.String("proxy", "", "Proxy URL")
	photo := fs.String("photo", "large", "Photo size: thumbnail, medium, large")
	fs.Parse(args)

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Error: profile URL required")
		os.Exit(1)
	}
	url := fs.Arg(0)

	profile := tiktok.Profile{URL: url, BaseDIR: *output, Proxy: *proxy}
	if err := profile.FetchInfo(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if *info {
		data, err := profile.GetProfileInfo()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(data)
		return
	}

	path, err := profile.DownloadPhoto(*photo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(path)
}
