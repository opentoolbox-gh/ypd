package main

import (
	"fmt"
	"log"
	"os"
	"ypd/utils"
)

func main() {

	playList := utils.ProcessInputPlaylist(os.Args[1])
	songs := utils.ListSongs(playList)

	fmt.Printf("This playlist has %d songs\n", len(songs))
	for i := 0; i < len(songs); i++ {
		fmt.Print(songs[i].Url)
		fmt.Printf("Attempting to download %s\n", songs[i].Info.Title)
		downLoadUrl := utils.GetSongDownloadUrl(songs[i].Url)
		err := utils.DownloadFile(downLoadUrl, songs[i].Info.Title)
		if err != nil {
			log.Print(err.Error())
		} else {
			fmt.Printf("Downloaded %s into %s.mp3\n", songs[i].Info.Title, songs[i].Info.Title)
		}
	}

}
