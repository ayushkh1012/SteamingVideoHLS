package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.Handle("/hls/", http.StripPrefix("/hls/", http.FileServer(http.Dir("./input"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>HLS Video Stream</title>
			<script src="https://cdn.jsdelivr.net/npm/video.js/dist/video.min.js"></script>
			<script src="https://cdn.jsdelivr.net/npm/@videojs/http-streaming/dist/videojs-http-streaming.min.js"></script>
			<link href="https://cdn.jsdelivr.net/npm/video.js/dist/video-js.min.css" rel="stylesheet">
			<style>
				.quality-display {
					position: fixed;
					top: 10px;
					right: 10px;
					background: rgba(0,0,0,0.7);
					color: white;
					padding: 10px;
					border-radius: 5px;
				}
			</style>
		</head>
		<body>
			<h1>HLS Video Stream</h1>
			<video id="hls-video" class="video-js vjs-default-skin" controls preload="auto" width="1280" height="720">
				<source src="/hls/master.m3u8" type="application/x-mpegURL">
			</video>
			<div id="quality-display" class="quality-display">Current Quality: Checking...</div>
			<script>
				// Add VideoJS HTTP Streaming (VHS) support
				var player = videojs('hls-video', {
					html5: {
						hls: {
							enableLowInitialPlaylist: true,
							smoothQualityChange: true,
							overrideNative: true
						}
					},
					controls: true,
					autoplay: false,
					preload: 'auto',
					controlBar: {
						children: [
							'playToggle',
							'progressControl',
							'volumePanel',
							'qualitySelector',
							'fullscreenToggle'
						]
					}
				});

				function updateQualityDisplay(qualityLevels) {
					const qualityDisplay = document.getElementById('quality-display');
					const selectedIndex = qualityLevels.selectedIndex;
					console.log('Updating quality display, selected index:', selectedIndex);
					
					if (selectedIndex >= 0) {
						const currentQuality = qualityLevels[selectedIndex];
						const height = currentQuality.height;
						const bitrate = Math.round(currentQuality.bitrate / 1000);
						qualityDisplay.textContent = 'Current Quality: ' + height + 'p (' + bitrate + 'kbps)';
						console.log('Updated quality display to:', height + 'p');
					} else {
						console.log('No quality level selected yet');
					}
				}

				// Add metadata logging
				player.on('loadedmetadata', function() {
					console.log('Metadata loaded');
					const qualityLevels = player.qualityLevels();
					console.log('Number of quality levels:', qualityLevels.length);
					
					if (qualityLevels.length === 0) {
						console.error('No quality levels found');
						document.getElementById('quality-display').textContent = 'Error: No quality levels found';
						return;
					}
					
					// Log initial quality levels
					for(let i = 0; i < qualityLevels.length; i++) {
						console.log('Quality level ' + i + ':', {
							width: qualityLevels[i].width,
							height: qualityLevels[i].height,
							bitrate: qualityLevels[i].bitrate,
							enabled: qualityLevels[i].enabled
						});
					}

					// Update display for initial quality
					updateQualityDisplay(qualityLevels);

					// Listen for quality changes
					qualityLevels.on('change', function() {
						console.log('Quality changed event fired');
						updateQualityDisplay(qualityLevels);
					});

					// Listen for quality level additions
					qualityLevels.on('addqualitylevel', function(event) {
						console.log('New quality level added:', {
							width: event.qualityLevel.width,
							height: event.qualityLevel.height,
							bitrate: event.qualityLevel.bitrate,
							enabled: event.qualityLevel.enabled
						});
					});
				});

				// Add loading state logging
				player.on('waiting', function() {
					console.log('Player is waiting for data');
				});

				player.on('playing', function() {
					console.log('Player is playing');
					const qualityLevels = player.qualityLevels();
					updateQualityDisplay(qualityLevels);
				});

				// Network state logging
				player.on('loadeddata', function() {
					console.log('Media data is loaded');
					const qualityLevels = player.qualityLevels();
					updateQualityDisplay(qualityLevels);
				});

				// Ready state
				player.ready(function() {
					console.log('Player is ready');
					console.log('Source URL:', player.currentSrc());
				});
			</script>
		</body>
		</html>
		`
		fmt.Fprint(w, html)
	})

	// Start the HTTP server
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
