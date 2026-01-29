package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Video struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Thumbnail   string `json:"thumbnail"`
	PublishedAt string `json:"published_at"`
	URL         string `json:"url"`
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
	err := godotenv.Load()
	if err != nil {
		slog.Warn("error in loading the .env file, relying on system variables")
	}

	apiKey := os.Getenv("YOUTUBE_API_KEY")
	channelID := os.Getenv("CHANNEL_ID")
	if apiKey == "" || channelID == "" {
		slog.Error("YOUTUBE_API_KEY and CHANNEL_ID must be set")
	}

	service, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("error in creating the Youtube service: %w", err)
	}

	slog.Info("Fetching the channel stats")
	channelCall := service.Channels.List([]string{"statistics"}).Id(channelID)
	channelResponse, err := channelCall.Do()
	if err != nil {
		log.Fatalf("error in fecthing the changel stats: %w", err)
	}

	stats := ChannelStats{
		SubscriberCount: channelResponse.Items[0].Statistics.SubscriberCount,
		ViewCount:       channelResponse.Items[0].Statistics.ViewCount,
		VideoCount:      channelResponse.Items[0].Statistics.VideoCount,
	}

	slog.Info("fetching latest videos")
	searchCall := service.Search.List([]string{"snippet"}).ChannelId(channelID).Order("date").Type("video").MaxResults(6)
	searchResults, err := searchCall.Do()
	if err != nil {
		log.Fatalf("error in fetching the latest videos: %v", err.Error())
	}

	videos := make([]Video, 0, 6)
	for _, item := range searchResults.Items {
		v := Video{
			ID:          item.Id.VideoId,
			Title:       item.Snippet.Title,
			Thumbnail:   item.Snippet.Thumbnails.High.Url,
			PublishedAt: item.Snippet.PublishedAt,
			URL:         "https://www.youtube.com/" + item.Id.VideoId,
		}
		videos = append(videos, v)
	}

	finalOutput := Output{
		LastUpdated:  time.Now().Format(time.RFC3339),
		Stats:        stats,
		LatestVideos: videos,
	}

	file, err := json.MarshalIndent(finalOutput, "", " ")
	if err != nil {
		log.Fatalf("error in marhalling the data: %v", err.Error())
	}

	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0o755)
	}

	err = os.WriteFile("data/youtube.json", file, 0o644)
	if err != nil {
		log.Fatalf("error in writing file: %v", err.Error)
	}

	slog.Info("Success! data is writen to youtube.json file")
}
