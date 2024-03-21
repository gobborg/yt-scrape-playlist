package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/api/youtube/v3"
)

// MAKE A NEW PLAYLIST
func makePlaylist(ytClient *youtube.Service) (*youtube.Playlist, error) {
	newPlaylist := &youtube.Playlist{ // y.P struct
		Snippet: &youtube.PlaylistSnippet{ // y.PS struct
			Title: "THE DUMPING PARADE 100",
		},
	}
	createPlaylist, err := ytClient.Playlists.Insert([]string{"snippet"}, newPlaylist).Do()
	if err != nil {
		log.Fatalf("Error creating playlist: %v", err)
	}
	fmt.Printf("PLAYLIST CREATED. PLAYLIST ID: %v\n", createPlaylist.Id)
	//playlistId := createPlaylist.Id
	return createPlaylist, nil
}

// POPULATE THAT NEW PLAYLIST
func addToPlaylist(ytClient *youtube.Service, playlist *youtube.Playlist, slugs []string, i int) error {
	if i >= len(slugs) {
		return nil
	}
	var j int64 //type conv for Snippet.Position
	j = int64(i)
	pid := playlist.Id
	videoSlugs := slugs[i]
	tracksForPlaylist := &youtube.PlaylistItem{
		Snippet: &youtube.PlaylistItemSnippet{
			PlaylistId: pid,
			Position:   j,
			ResourceId: &youtube.ResourceId{
				Kind:    "youtube#video",
				VideoId: videoSlugs,
			},
		},
	}
	delay := 300 * time.Millisecond
	retries := 5
	sendTracks, err := ytClient.PlaylistItems.Insert([]string{"snippet"}, tracksForPlaylist).Do()
	switch {
	case sendTracks.HTTPStatusCode == 409 && retries > 0:
		retries--
		time.Sleep(delay * 2)
		return addToPlaylist(ytClient, playlist, slugs, i+1)
	case sendTracks.HTTPStatusCode == 409 && retries == 0:
		return fmt.Errorf("Error 409 after 5 tries. Quitting.")
	case sendTracks.HTTPStatusCode == 200:
		retries = 5
		time.Sleep(delay)
		return addToPlaylist(ytClient, playlist, slugs, i+1)
	default:
		if err != nil {
			return fmt.Errorf("Error sending tracks: %v", err)
		}
	}
	return nil
}

func main() {
	log.SetOutput(os.Stdout)
	slugs, err := Scrape()
	if err != nil {
		log.Fatalf("Error getting slugs for main(). Info: %v", err)
	}
	i := 0
	auth, err := OAuth()
	playlist, err := makePlaylist(auth)
	if err != nil {
		log.Fatalf("Could not make playlist due to: %v", err)
	}
	addToPlaylist(auth, playlist, slugs, i)

}
