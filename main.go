package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Video struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Thumbnail    string `json:"thumbnail"`
	PublishedAt  string `json:"published_at"`
	URL          string `json:"url"`
	Duration     string `json:"duration"`
	LikeCount    uint64 `json:"like_count"`
	DislikeCount uint64 `json:"dislike_count"`
	ViewCount    uint64 `json:"view_count"`
}

type ChannelStats struct {
	SubscriberCount uint64 `json:"subscriber_count"`
	ViewCount       uint64 `json:"view_count"`
	VideoCount      uint64 `json:"video_count"`
}

type Output struct {
	LastUpdated  string       `json:"last_updated"`
	Stats        ChannelStats `json:"stats"`
	LatestVideos []Video      `json:"latest_videos"`
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	durationFilter := 60 // in seconds
	err := godotenv.Load()
	if err != nil {
		logger.Warn("error in loading the .env file, relying on system variables")
	}

	apiKey := os.Getenv("YOUTUBE_API_KEY")
	channelID := os.Getenv("CHANNEL_ID")
	if apiKey == "" || channelID == "" {
		logger.Error("YOUTUBE_API_KEY and CHANNEL_ID must be set")
		os.Exit(1)
	}

	service, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		logger.Error("error in creating the Youtube service", "error", err.Error())
	}

	logger.Info("Fetching the channel stats")
	channelCall := service.Channels.List([]string{"statistics"}).Id(channelID)
	channelResponse, err := channelCall.Do()
	if err != nil {
		logger.Error("error in fetching the channel stats", "error", err.Error())
		os.Exit(1)
	}

	stats := ChannelStats{
		SubscriberCount: channelResponse.Items[0].Statistics.SubscriberCount,
		ViewCount:       channelResponse.Items[0].Statistics.ViewCount,
		VideoCount:      channelResponse.Items[0].Statistics.VideoCount,
	}

	logger.Info("fetching latest videos")
	searchCall := service.Search.List([]string{"id"}).ChannelId(channelID).Order("date").Type("video").MaxResults(50)
	searchResults, err := searchCall.Do()
	if err != nil {
		logger.Error("error in fetching the latest videos", "error", err.Error())
		os.Exit(1)
	}

	videoIDs := []string{}
	for _, item := range searchResults.Items {
		videoIDs = append(videoIDs, item.Id.VideoId)
	}

	var longFormVideos []Video

	videoCall := service.Videos.List([]string{"snippet", "contentDetails", "statistics"}).Id(videoIDs...)
	videoResults, err := videoCall.Do()
	if err != nil {
		logger.Error("error is fetching video properties", "error", err.Error())
		os.Exit(1)
	}

	for _, item := range videoResults.Items {

		duration, err := parseDuration(item.ContentDetails.Duration)
		if err != nil {
			logger.Error("error in parsing the time", "error", err)
			continue
		}
		if duration.Seconds() > float64(durationFilter) {
			v := Video{
				ID:           item.Id,
				Title:        item.Snippet.Title,
				Thumbnail:    item.Snippet.Thumbnails.High.Url,
				PublishedAt:  item.Snippet.PublishedAt,
				URL:          "https://www.youtube.com/watch?v=" + item.Id,
				LikeCount:    item.Statistics.LikeCount,
				DislikeCount: item.Statistics.DislikeCount,
				ViewCount:    item.Statistics.ViewCount,
			}
			longFormVideos = append(longFormVideos, v)
		}
	}

	finalOutput := Output{
		LastUpdated:  time.Now().Format(time.RFC3339),
		Stats:        stats,
		LatestVideos: longFormVideos,
	}

	file, err := json.MarshalIndent(finalOutput, "", " ")
	if err != nil {
		logger.Error("error in marshalling the data", "error", err.Error())
		os.Exit(1)
	}

	if _, err := os.Stat("data"); os.IsNotExist(err) {
		if err := os.Mkdir("data", 0o755); err != nil {
			logger.Error("error in creating the folder", "error", err.Error())
			os.Exit(1)
		}
	}

	err = os.WriteFile("data/youtube.json", file, 0o644)
	if err != nil {
		logger.Error("error in writing file", "error", err.Error())
		os.Exit(1)
	}

	logger.Info("Success! data is written to youtube.json file")
}

// parseDuration (Format is ISO 8601, e.g., PT1H2M10S)
func parseDuration(isoDuration string) (time.Duration, error) {
	cleanTime := strings.ToLower(isoDuration)
	cleanTime = strings.TrimPrefix(cleanTime, "pt")

	d, err := time.ParseDuration(cleanTime)
	if err != nil {
		return 0, fmt.Errorf("error in parsing the time: %w", err)
	}
	return d, nil
}
