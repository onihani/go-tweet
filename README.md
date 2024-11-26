
# Twitter Video Downloader

A simple Go application for downloading video content from Twitter. This project allows you to input a tweet URL, extract the tweet ID, and download any available video content in MP4 format.

## Features

- Extract tweet ID from a tweet URL.
- Download video from tweet URLs with video content in MP4 format.
- Save videos with bitrate information as the file name.
- Handles video resolution by selecting the highest available quality.

## Requirements

- Go 1.18+ (for the Go tools and libraries).
- `huh` library for building terminal-based forms.
- Internet connection to fetch video content.

## Installation

1. Clone this repository:
   ```bash
   git clone https://github.com/yourusername/twitter-video-downloader.git
   ```

2. Change into the project directory:
   ```bash
   cd twitter-video-downloader
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Build the project:
   ```bash
   go build
   ```

5. Run the application:
   ```bash
   ./twitter-video-downloader
   ```

## Usage

When you run the application, it will prompt you to enter a tweet URL. Once the URL is entered, the program will:

1. Extract the tweet ID from the URL.
2. Fetch tweet data from the Twitter API.
3. Display the available video resolutions.
4. Download the highest quality MP4 video.

### Example:

```bash
Enter Tweet URL: https://x.com/xghana_/status/1853065799977537588
Available videos: [320x562, 364x640]
Downloading video_632kbps.mp4 to ./downloads...
Download complete!
```

## File Structure

- `main.go`: The main application logic for fetching tweet data and handling user input.
- `tweet/tweet.go`: Defines the `Tweet` struct and methods to extract video URLs.
- `util/util.go`: Utility functions like extracting the tweet ID from the URL.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [huh](https://github.com/charmbracelet/huh): Simple library for building terminal forms.
- Twitter API (v2) for tweet data extraction.
