# defog-sync-go
A lightweight Go utility that interfaces with the YouTube Data API v3 to generate static JSON datasets for Astro builds.

## Setup

### Prerequisites
- Go 1.16 or higher
- YouTube Data API v3 credentials

### Installation

1. Clone the repository:
```bash
git clone git@github.com:Defog-online/content-sync-go.git
cd defog-content-sync
```

2. Create a `.env` file from the template:
```bash
cp .env.Template .env
```

3. Add your credentials to `.env`:
```
YOUTUBE_API_KEY=your_api_key_here
CHANNEL_ID=your_channel_id_here
```

4. Run the utility:
```bash
go run main.go
```

## Output

The utility generates a `data/youtube.json` file containing:
- Channel statistics (subscriber count, view count, video count)
- Latest videos (ID, title, thumbnail, publish date, URL, duration, engagement metrics)

## Configuration

Set the following environment variables:
- `YOUTUBE_API_KEY`: Your YouTube Data API v3 key
- `CHANNEL_ID`: Your YouTube channel ID
