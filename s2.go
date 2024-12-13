package main

// importing all the necessary modules
import (
	"log"
	"net/http"
	"path/filepath"
	"os"
)

func main() {

	// Defining the directory where HLS files are stored (shared with Server 1)
	hlsPath := "C:/nginx-rtmp-win64-master/html/hls"

	// Serve static files from the HLS directory
	fs := http.FileServer(http.Dir(hlsPath))

	// Handle requests to /hls/ and serve the files accordingly
	http.Handle("/hls/", http.StripPrefix("/hls/", fs))

	// Handle playlist.m3u8 specifically for real-time streaming
	http.HandleFunc("/hls/playlist.m3u8", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving playlist.m3u8")
	
		// Prevent caching - prevents browser or any client form caching previous playlist - Real Time only
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
	
		// Chech if the playlist exists or not
		playlistPath := filepath.Join(hlsPath, "playlist.m3u8")

		// Throw an error if there is no playlist
		if _, err := os.Stat(playlistPath); os.IsNotExist(err) {
			http.Error(w, "Playlist not found", http.StatusNotFound)
			return
		}
	
		// If the playlist exists, serve it to the client
		http.ServeFile(w, r, playlistPath)
	})

	// Define the port
	port := ":8080"
	
	// Start and listen to port 
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v\n", err)
	}

	log.Printf("Serving HLS content on http://localhost%s/hls/\n", port)
}
