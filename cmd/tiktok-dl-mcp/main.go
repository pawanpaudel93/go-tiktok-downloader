package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/pawanpaudel93/go-tiktok-downloader/tiktok"
)

func main() {
	s := server.NewMCPServer(
		"tiktok-downloader",
		"0.2.0",
		server.WithToolCapabilities(false),
	)

	s.AddTool(mcp.NewTool("download_video",
		mcp.WithDescription("Download a TikTok video. Returns the file path where the video was saved."),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("TikTok video URL (e.g. https://www.tiktok.com/@user/video/1234567890)"),
		),
		mcp.WithString("output_dir",
			mcp.Description("Output directory (default: ./downloads)"),
		),
	), handleDownloadVideo)

	s.AddTool(mcp.NewTool("get_video_info",
		mcp.WithDescription("Get TikTok video metadata as JSON (likes, shares, comments, author info, etc.)"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("TikTok video URL"),
		),
	), handleGetVideoInfo)

	s.AddTool(mcp.NewTool("download_profile_photo",
		mcp.WithDescription("Download a TikTok user's profile picture. Returns the file path."),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("TikTok profile URL (e.g. https://www.tiktok.com/@username)"),
		),
		mcp.WithString("output_dir",
			mcp.Description("Output directory (default: ./downloads)"),
		),
		mcp.WithString("photo_size",
			mcp.Description("Photo size: thumbnail, medium, or large (default: large)"),
			mcp.Enum("thumbnail", "medium", "large"),
		),
	), handleDownloadProfilePhoto)

	s.AddTool(mcp.NewTool("get_profile_info",
		mcp.WithDescription("Get TikTok profile metadata as JSON (user info, follower counts, etc.)"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("TikTok profile URL"),
		),
	), handleGetProfileInfo)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func handleDownloadVideo(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url, err := req.RequireString("url")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	outputDir := req.GetString("output_dir", "./downloads")

	video := tiktok.Video{URL: url, BaseDIR: outputDir}
	if err := video.FetchInfo(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch video info: %v", err)), nil
	}
	path, err := video.Download()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to download video: %v", err)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("Video downloaded to: %s", path)), nil
}

func handleGetVideoInfo(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url, err := req.RequireString("url")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	video := tiktok.Video{URL: url, BaseDIR: "./downloads"}
	if err := video.FetchInfo(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch video info: %v", err)), nil
	}
	data, err := video.GetInfo()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get video info: %v", err)), nil
	}
	return mcp.NewToolResultText(data), nil
}

func handleDownloadProfilePhoto(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url, err := req.RequireString("url")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	outputDir := req.GetString("output_dir", "./downloads")
	photoSize := req.GetString("photo_size", "large")

	profile := tiktok.Profile{URL: url, BaseDIR: outputDir}
	if err := profile.FetchInfo(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch profile info: %v", err)), nil
	}
	path, err := profile.DownloadPhoto(photoSize)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to download photo: %v", err)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("Profile photo downloaded to: %s", path)), nil
}

func handleGetProfileInfo(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url, err := req.RequireString("url")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	profile := tiktok.Profile{URL: url, BaseDIR: "./downloads"}
	if err := profile.FetchInfo(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch profile info: %v", err)), nil
	}
	data, err := profile.GetProfileInfo()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get profile info: %v", err)), nil
	}
	return mcp.NewToolResultText(data), nil
}
