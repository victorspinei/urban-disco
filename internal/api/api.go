package api

import (
	"fmt"
	"regexp"
	"errors"
	"github.com/gocolly/colly/v2"
)

func GetSongsFromAlbum(albumName string) ([]string, error) {
	var trackListing []string

	c := colly.NewCollector(
		colly.AllowedDomains("www.wikipedia.org", "en.wikipedia.org"),
	)

	// Handle errors
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	// Response logging
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})

	re := regexp.MustCompile(`"([^"]+)"`)

	// Form submission logic - This is where we handle the search
	c.OnHTML("table.tracklist", func(e *colly.HTMLElement) {
		//fmt.Println("Found tracklist table. Extracting song names...")

		// Iterate over each row in the table
		e.ForEach("tbody > tr", func(rowIndex int, rowElement *colly.HTMLElement) {
			if rowIndex == 0 { // Skip header row
				return
			}

			// Extract song title from the second column (index 1)
			rawTitle := rowElement.ChildText("td:nth-of-type(1)")

			// Find song titles inside quotes using regex
			matches := re.FindStringSubmatch(rawTitle)
			if len(matches) > 1 { // The first match is the entire match, second is the text inside quotes
				songTitle := matches[1]
				//fmt.Println("Song Title:", songTitle)
				trackListing = append(trackListing, songTitle)
			}
		})
	})

	// Log once scraping is complete
	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, "scraped!")
	})

	// Start by visiting the Wikipedia homepage
	searchURL := fmt.Sprintf("https://en.wikipedia.org/w/index.php?title=Special:Search&search=%s", albumName)
	c.Visit(searchURL)

	if len(trackListing) != 0 {
		return trackListing, nil
	} else {
		return []string{}, errors.New("no songs found in the album")
	}
}
