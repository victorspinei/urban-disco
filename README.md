# Urban Disco

Urban Disco is a Go-based application that facilitates finding, downloading, and managing songs from a specific album by a given artist. The application fetches a list of tracks from the album, downloads them as MP3 files from YouTube, and retrieves their lyrics.

## Features

- **Album and Artist Search:** Enter the name of an album and the artist to retrieve all tracks.
- **MP3 Download:** Download each track as an MP3 file.
- **Lyrics Retrieval:** Fetch lyrics for each track.
- **Rate Limiting:** Protect the API from excessive requests.
- **Environment Configuration:** Manage API tokens and environment settings through a `.env` file.

## Getting Started

### Prerequisites

- Go 1.18 or higher
- ffmpeg (for MP3 conversion)
- A Discogs API token
- A YouTube Data API key

### Clone the Repository

```bash
git clone https://github.com/victorspinei/urban-disco.git
cd urban-disco
```

### Install Dependencies

Install Go dependencies by running:

```bash
go mod tidy
```

### Environment Configuration

Create a `.env` file in the root directory of the project with the following content:

```plaintext
TOKEN=your_discogs_token_here
API_KEY=your_youtube_api_key_here
ENV=production
PORT=5000
```

Replace `your_discogs_token_here` and `your_youtube_api_key_here` with your actual API tokens.

### Run the Server

To start the server, use:

```bash
go run main.go
```

The server will be available at `http://localhost:5000` by default.

## API Endpoints

### 1. Get Track List

- **Endpoint:** `GET /api/tracklist/:q`
- **Parameters:**
  - `q` (required): The name of the album.
- **Description:** Retrieves all tracks from the specified album and returns them with `body`, `artist`, and `downloaded=false`.

**Example Request:**

```plaintext
http://localhost:5000/api/tracklist/Powerslave
```

**Example Response:**

```json
[
  {
    "body": "Aces High",
    "artist": "Iron Maiden",
    "downloaded": false
  },
  {
    "body": "2 Minutes to Midnight",
    "artist": "Iron Maiden",
    "downloaded": false
  }
]
```

### 2. Download Song

- **Endpoint:** `GET /api/song`
- **Parameters:**
  - `name` (required): The title of the song.
  - `artist` (required): The name of the artist.
- **Description:** Searches for the specified song on YouTube, downloads it as an MP3 file, and returns the file.

**Example Request:**

```plaintext
http://localhost:5000/api/song?name=Aces High&artist=Iron%20Maiden
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes.

---

Feel free to adjust the content according to your project's specifics or personal preferences!