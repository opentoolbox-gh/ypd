package main

import (
	"fmt"
	"log"
	"ypd/utils"
)

func main() {
	songs := utils.ListSongs("https://www.youtube.com/playlist?list=PLyL6BB2aK1Xvu4EG_Z5v8H-fbJz9bwAgf")
	for i := 0; i < len(songs); i++ {
		fmt.Print(songs[i].Url)
		downLoadUrl := utils.GetSongDownloadUrl(songs[i].Url)
		err := utils.DownloadFile(downLoadUrl, songs[i].Info.Title)
		if err != nil {
			log.Print(err.Error())
		}
	}

}
