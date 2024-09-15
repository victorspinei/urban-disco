package main

import (
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
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
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

var apiKey string
var token string

func main() {
	// Load environment variables from .env file
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	apiKey = os.Getenv("API_KEY")
	token = os.Getenv("TOKEN")

	app := fiber.New()

	if os.Getenv("ENV") == "development" {
		app.Use(cors.New(cors.Config{
			AllowOrigins: "http://localhost:5173/",
			AllowHeaders: "Origin,Content-Type,Accept",
		}))
	} else if os.Getenv("ENV") == "production" {
		app.Static("/", "./client/dist")
	}

	// Rate limiting configuration
	app.Use("/api/song", limiter.New(limiter.Config{
		Max:        5,               // Maximum number of requests
		Expiration: 1 * time.Minute, // Rate limit duration
	}))

	app.Get("/api/tracklist/:q", getTrackList)
	app.Get("/api/song", getSongFile)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "5000"
	}
	log.Fatal(app.Listen("0.0.0.0:" + PORT))

}

func getTrackList(c *fiber.Ctx) error {
	// Get the query parameter and format it
	query := c.Params("q")
	query = strings.ReplaceAll(query, "%20", " ")

	// Fetch the tracklist and artist from the album
	tracklist, artist, err := GetSongsFromAlbum(token, query)
	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("Error fetching tracklist: %v", err))
	}

	// Create a response with body, artist, and downloaded=false for each track
	var response []map[string]interface{}

	for _, track := range tracklist {
		response = append(response, map[string]interface{}{
			"body":       track,
			"artist":     artist,
			"downloaded": false,
		})
	}

	// Return the JSON response
	return c.Status(200).JSON(response)
}

func getSongFile(c *fiber.Ctx) error {
	// Read the query parameters
	songName := c.Query("name")
	artist := c.Query("artist")

	if songName == "" || artist == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing 'name' or 'artist' query parameters")
	}

	// Create a directory for temporary files if it does not exist
	os.MkdirAll("tmp", os.ModePerm)

	// Format the query for YouTube search
	ytQuery := fmt.Sprintf("%s %s (oficial)", songName, artist)

	// Search for the track on YouTube
	videoID, err := SearchYouTube(apiKey, ytQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error searching YouTube: %v", err))
	}

	if videoID == "" {
		return c.Status(fiber.StatusNotFound).SendString("No video found for the query")
	}

	// Download the audio from YouTube as MP3
	title := strings.ReplaceAll(songName, " ", "_") // Make file name safe
	mp3FileName := fmt.Sprintf("tmp/%s.mp3", title)
	if err := DownloadYouTubeAudioAsMP3(videoID, mp3FileName); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error downloading audio: %v", err))
	}

	// Serve the MP3 file for download
	defer os.Remove(mp3FileName) // Clean up after download
	return c.Download(mp3FileName, title+".mp3")
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
func SearchYouTube(apiKey, query string) (string, error) {
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", err
	}

	// Set up the search call
	call := service.Search.List([]string{"id", "snippet"}).
		Q(query).      // The search query
		Type("video"). // Filter out playlists and channels
		MaxResults(1)  // Limit results to 1

	response, err := call.Do()
	if err != nil {
		return "", err
	}

	// Print the entire response for debugging
	fmt.Printf("API Response: %+v\n", response)

	// Check if we got any results
	if len(response.Items) > 0 {
		videoID := response.Items[0].Id.VideoId
		fmt.Printf("Found Video ID: %s\n", videoID)
		return videoID, nil
	}

	return "", fmt.Errorf("no results found for query: %s", query)
}

// downloadYouTubeAudioAsMP3 downloads the audio of a YouTube video and converts it to MP3
func DownloadYouTubeAudioAsMP3(videoID, title string) error {
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

	audioFileName := "tmp/" + videoID + ".m4a"
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
	if err := ConvertToMP3(audioFileName, title); err != nil {
		return fmt.Errorf("error converting to mp3: %v", err)
	}

	log.Printf("Converted audio to MP3: %s", title)

	// Remove temporary audio file
	os.Remove(audioFileName)

	return nil
}

// convertToMP3 converts an audio file to MP3 using ffmpeg
func ConvertToMP3(inputFile, outputFile string) error {
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-vn", "-ab", audioBitrate, "-ar", sampleRate, "-y", outputFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("Converting %s to %s", filepath.Base(inputFile), filepath.Base(outputFile))
	return cmd.Run()
}
