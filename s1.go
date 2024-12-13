package main

// importing all the necessary modules
import (
	"log"
	"os"
	"os/exec"
)

func main() {

	// Defining path for rtmp input - RTMP strema URL
	rtmpInput := "rtmp://localhost/live/my-stream" 

	// Defining path for segemented output - hls output directory
	// Shared storage
	hlsOutputDir := "C:/nginx-rtmp-win64-master/html/hls" 
	
	// Defining the playlist file (.m3u8)
	hlsPlaylist := hlsOutputDir + "/playlist.m3u8" 


	// Create the HLS output directory if it doesn't exist
	err := os.MkdirAll(hlsOutputDir, 0755)
	// 0755 - Owner - RWX, Group - RX, Other - Read only

	if err != nil {
		log.Fatalf("Failed to create HLS output directory: %v\n", err)
	}

	log.Println("Starting FFmpeg to transcode RTMP to HLS...")

	// FFmpeg command to transcode RTMP to HLS
	// Using the ffmpeg command with all the necessary flags for desired transcoding
	cmd := exec.Command(
		"ffmpeg",
		"-i", rtmpInput,                // RTMP input stream URL
		"-c:v", "libx264",              // Video codec
		"-preset", "veryfast",          // Transcoding speed
		"-g", "60",                     // Group of Pictures (keyframe interval)
		"-sc_threshold", "0",           // Scene change threshold
		"-force_key_frames", "expr:gte(t,n_forced*2)", // Force keyframes at 2-second intervals
		"-hls_time", "2",               // Segment duration in seconds
		"-hls_playlist_type", "event",  // Event-based playlist (non-looping)
		"-hls_list_size", "5",          // Keep the last 5 segments in the playlist
		"-hls_flags", "delete_segments", // Delete old segments once they are no longer needed
		"-f", "hls",                    // Output format: HLS
		hlsPlaylist,                   // Output directory for the HLS playlist and segments
	)

	// log everything for debugging
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the FFmpeg command
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to run FFmpeg: %v\n", err)
	}

	log.Println("Transcoding completed. HLS files are ready.")
}
