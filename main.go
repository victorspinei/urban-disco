package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	kkdaiYoutube "github.com/kkdai/youtube/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// Structs matching Discogs API response
type DiscogsSearchResponse struct {
	Results []ResultSearch `json:"results"`
}

type ResultSearch struct {
	MasterURL string `json:"master_url"`
}

type DiscogsAlbumResponse struct {
	Tracklist []Track   `json:"tracklist"`
	Artists   []Artists `json:"artists"`
}

type Track struct {
	Title string `json:"title"`
}

type Artists struct {
	Name string `json:"name"`
}

const discogsBaseURL = "https://api.discogs.com/"
const audioBitrate = "192k"
const sampleRate = "44100"

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")
	token := os.Getenv("TOKEN")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter album name: ")
	query, _ := reader.ReadString('\n')
	query = strings.TrimSpace(query)

	tracklist, artist, err := GetSongsFromAlbum(token, query)
	if err != nil {
		log.Fatal(err)
	}

	for _, track := range tracklist {
		query := track + " " + artist
		videoID, err := searchYouTube(apiKey, query)
		if err != nil {
			log.Fatalf("Error searching YouTube: %v", err)
		}

		if err := downloadYouTubeAudioAsMP3(videoID, track); err != nil {
			log.Fatalf("Error downloading audio: %v", err)
		}
	}
}

// GetSongsFromAlbum queries Discogs for the album and returns the tracklist and artist name
func GetSongsFromAlbum(token, albumName string) ([]string, string, error) {
	var trackListing []string
	query := strings.ReplaceAll(albumName, " ", "+")

	// Get album data from Discogs
	resp, err := http.Get(discogsBaseURL + "database/search?token=" + token + "&q=" + query)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	var discogsResponse DiscogsSearchResponse
	if err := json.Unmarshal(body, &discogsResponse); err != nil {
		return nil, "", fmt.Errorf("error parsing search response: %w", err)
	}

	// Ensure there is a result
	if len(discogsResponse.Results) == 0 {
		return nil, "", errors.New("no results found")
	}

	masterURL := discogsResponse.Results[0].MasterURL

	// Get tracklist from album master URL
	resp, err = http.Get(masterURL)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	var albumResponse DiscogsAlbumResponse
	if err := json.Unmarshal(body, &albumResponse); err != nil {
		return nil, "", fmt.Errorf("error parsing album response: %w", err)
	}

	// Check for tracklist
	if len(albumResponse.Tracklist) == 0 {
		return nil, "", errors.New("no tracks found for this album")
	}

	for _, track := range albumResponse.Tracklist {
		trackListing = append(trackListing, track.Title)
	}

	return trackListing, albumResponse.Artists[0].Name, nil
}

// searchYouTube searches for a video on YouTube using the given API key and query, and returns the video ID
func searchYouTube(apiKey, query string) (string, error) {
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", err
	}

	call := service.Search.List([]string{"id", "snippet"}).Q(query).MaxResults(1)
	response, err := call.Do()
	if err != nil {
		return "", err
	}

	if len(response.Items) > 0 {
		return response.Items[0].Id.VideoId, nil
	}

	return "", fmt.Errorf("no results found for query: %s", query)
}

// downloadYouTubeAudioAsMP3 downloads the audio of a YouTube video and converts it to MP3
func downloadYouTubeAudioAsMP3(videoID, title string) error {
	client := kkdaiYoutube.Client{}

	// Fetch the video info
	video, err := client.GetVideo(videoID)
	if err != nil {
		return fmt.Errorf("error getting video info: %v", err)
	}

	// Get audio format
	format := video.Formats.Type("audio")
	if len(format) == 0 {
		return fmt.Errorf("no audio format found for video: %s", videoID)
	}

	selectedFormat := &format[0]

	// Download the audio stream
	stream, _, err := client.GetStream(video, selectedFormat)
	if err != nil {
		return fmt.Errorf("error downloading audio stream: %v", err)
	}

	audioFileName := videoID + ".m4a"
	file, err := os.Create(audioFileName)
	if err != nil {
		return fmt.Errorf("error creating audio file: %v", err)
	}
	defer file.Close()

	if _, err = file.ReadFrom(stream); err != nil {
		return fmt.Errorf("error saving audio file: %v", err)
	}

	log.Printf("Downloaded audio to: %s", audioFileName)

	// Convert to MP3
	mp3FileName := title + ".mp3"
	if err := convertToMP3(audioFileName, mp3FileName); err != nil {
		return fmt.Errorf("error converting to mp3: %v", err)
	}

	log.Printf("Converted audio to MP3: %s", mp3FileName)

	// Remove temporary audio file
	os.Remove(audioFileName)

	return nil
}

// convertToMP3 converts an audio file to MP3 using ffmpeg
func convertToMP3(inputFile, outputFile string) error {
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-vn", "-ab", audioBitrate, "-ar", sampleRate, "-y", outputFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("Converting %s to %s", filepath.Base(inputFile), filepath.Base(outputFile))
	return cmd.Run()
}
