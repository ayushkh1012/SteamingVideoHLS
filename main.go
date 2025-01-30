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
				.video-selector {
					margin: 20px 0;
				}
				.video-selector button {
					padding: 10px;
					margin: 0 10px;
					cursor: pointer;
				}
				.active {
					background: #2196F3;
					color: white;
					border: none;
					border-radius: 4px;
				}
				.play-button {
					margin-left: 20px;
					padding: 10px 20px;
					background: #2196F3;
					color: white;
					border: none;
					border-radius: 4px;
					cursor: pointer;
				}
				.play-button:hover {
					background: #1976D2;
				}
			</style>
		</head>
		<body>
			<h1>HLS Video Stream</h1>
			<div class="video-selector">
				<span>Now Playing: <span id="current-video">Big Bunny</span></span>
				<button id="playButton" class="play-button">Start Playlist</button>
			</div>
			<video id="hls-video" class="video-js vjs-default-skin" controls preload="auto" width="1280" height="720">
				<source src="/hls/bigbunny/1080p/playlist.m3u8" type="application/x-mpegURL">
			</video>
			<div id="quality-display" class="quality-display">Current Quality: Checking...</div>
			<script>
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

				// Define the playlist
				const videoPlaylist = [
					{
						name: 'Big Bunny',
						src: '/hls/bigbunny/1080p/playlist.m3u8'
					},
					{
						name: 'Jelly',
						src: '/hls/jelly/1080p/playlist.m3u8'
					},
					{
						name: 'Sintel',
						src: '/hls/sintel/1080p/playlist.m3u8'
					}
				];

				let currentVideoIndex = 0;
				let playlistStarted = false;

				function playNextVideo() {
					currentVideoIndex = (currentVideoIndex + 1) % videoPlaylist.length;
					const nextVideo = videoPlaylist[currentVideoIndex];
					document.getElementById('current-video').textContent = nextVideo.name;
					player.src({
						src: nextVideo.src,
						type: 'application/x-mpegURL'
					});
					if (playlistStarted) {
						player.play();
					}
				}

				// Play next video when current one ends
				player.on('ended', playNextVideo);

				// Add click handler for the play button
				document.getElementById('playButton').addEventListener('click', function() {
					playlistStarted = true;
					player.play().then(() => {
						console.log('Playback started successfully');
					}).catch(error => {
						console.log('Playback failed:', error);
					});
					this.style.display = 'none';  // Hide the button after starting
				});

				function updateQualityDisplay(qualityLevels) {
					const qualityDisplay = document.getElementById('quality-display');
					const selectedIndex = qualityLevels.selectedIndex;
					
					if (selectedIndex >= 0) {
						const currentQuality = qualityLevels[selectedIndex];
						const height = currentQuality.height;
						const bitrate = Math.round(currentQuality.bitrate / 1000);
						qualityDisplay.textContent = 'Current Quality: ' + height + 'p (' + bitrate + 'kbps)';
					}
				}

				player.on('loadedmetadata', function() {
					const qualityLevels = player.qualityLevels();
					
					if (qualityLevels.length === 0) {
						document.getElementById('quality-display').textContent = 'Error: No quality levels found';
						return;
					}
					
					updateQualityDisplay(qualityLevels);

					qualityLevels.on('change', function() {
						updateQualityDisplay(qualityLevels);
					});
				});

				player.on('error', function(e) {
					console.error('Video error:', e);
					if (playlistStarted) {
						playNextVideo();
					}
				});

				player.ready(function() {
					console.log('Player is ready');
					document.getElementById('current-video').textContent = videoPlaylist[currentVideoIndex].name;
				});
			</script>
		</body>
		</html>
		`
		fmt.Fprint(w, html)
	})

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
