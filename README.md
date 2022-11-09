# Continuous Release Radar

A go service that moves tracks from user's Release Radar to a long term playlist (named Continuous Release Radar by default).

I have provided a cronjob script that automates the process and runs the script every week on Tuesday at noon. 

## Set up
1. Start by registering your application at the following page:

    https://developer.spotify.com/my-applications/.

    You'll get a client ID and secret key for your application. 
2. Create a .env file and set `SPOTIFY_ID` to the client ID and `SPOTIFY_SECRET` to the secret key. 
3. If you haven't installed Golang previously, go ahead and install it. 
4.  If you would like to name the the playlist something other than "Continuous Release Radar", edit the `PLAYLIST_NAME` in the `move_tracks.go` file.
6. Run `go run .` in the terminal in the project folder. 
7. The app will print a url. Go to that url and sign in and allow the permissions the app requests. The app will then create the playlist and copy over this week's Release Radar. 
8. Edit the cronjob.sh file by setting the variables. `CRR_HOME` should be wherever you have this project located on your machine. `GO_PATH` should be wherever Golang is installed on your machine. A common location for `GO_PATH` is `/usr/local/go/bin/go`.
9. To prevent duplicates, if current time is between Tuesday at noon and Thursday night, continue to the steps below. If the current time is after Friday and before Tuesday at noon, then delete the Continuous Release Radar playlist from the Spotify client and continue with the steps below. Don't worry, it will get recreated on Tuesday at noon.
10. Run the cronjob.sh file. 
11. That's it! As long as the computer is awake every Tuesday at noon, the tracks in that week's Release Radar will be copied over into the playlist. 

## Notes
The cronjob works on Linux and Mac. I don't know if Windows supports cron, I presume not though. 

## Questions
Feel free to throw any questions into the issues page. I'd be happy to help out. 

PRs are also welcome.