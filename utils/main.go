package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

type VideoInfo struct {
	Title       string `json:"title"`
	Author_name string `json:"author_name"`
}

type Video struct {
	Url  string    `json:"url"`
	Info VideoInfo `json:"info"`
}

type ConvertVideoResp struct {
	Status   string `json:"status"`
	Message  string `json:"mess"`
	DLink    string `json:"dlink"`
	Title    string `json:"title"`
	C_Status string `json:"c_status"`
}

func makePlayList(playListUrl string) string {
	var url string = fmt.Sprintf("https://loader.to/api/ajax/playlistJSON?format=1080&api=dfcb6d76f2f6a9894gjkege8a4ab232222&limit=100&url=%s", playListUrl)
	fmt.Println("Resolving playlist songs.....")
	return url
}

func ListSongs(url string) []Video {
	resp, err := http.Get(makePlayList(url))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var videos []Video

	err = json.NewDecoder(resp.Body).Decode(&videos)
	if err != nil {
		fmt.Println("Maarshaal error:", err)
	}
	return videos

}

func ConvertVideo(vidId string, k string) ConvertVideoResp {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("vid", vidId)
	writer.WriteField("k", k)

	req, _ := http.NewRequest(
		"POST",
		"https://www.y2mate.com/mates/convertV2/index",
		bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var response ConvertVideoResp

	client := &http.Client{}
	resp, _ := client.Do(req)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Request failed with response code: %d\n", resp.StatusCode)
	}

	error := json.NewDecoder(resp.Body).Decode(&response)
	if error != nil {
		panic(error)
	}
	return response

}

func GetSongDownloadUrl(videoUrl string) string {
	var dlink string

	u, err := url.Parse(videoUrl)
	if err != nil {
		panic(err)
	}

	queryObj, _ := url.ParseQuery(u.RawQuery)
	vidId := queryObj["v"][0]

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("k_query", videoUrl)
	writer.WriteField("k_page", "home")
	writer.WriteField("ht", "en")
	writer.WriteField("q_auto", "0")

	req, _ := http.NewRequest(
		"POST",
		"https://www.y2mate.com/mates/analyzeV2/ajax",
		bytes.NewReader(body.Bytes()))

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Origin", "https://tomp3.cc")
	req.Header.Set("Referer", "https://tomp3.cc/youtube-to-mp3/6zr6HXG5WoI")

	var response struct {
		Status string `json:"status"`
		Links  map[string]map[string]struct {
			Size string `json:"size"`
			F    string `json:"f"`
			K    string `json:"k"`
		}
	}

	client := &http.Client{}
	resp, _ := client.Do(req)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Request failed with response code: %d", resp.StatusCode)
	}

	error := json.NewDecoder(resp.Body).Decode(&response)
	if error != nil {
		fmt.Println(error)
	} else {

	SearchSong:
		for i, value := range response.Links {
			if i == "mp3" {
				for _, v := range value {
					convertVidResp := ConvertVideo(vidId, v.K)
					if convertVidResp.C_Status == "CONVERTED" {
						dlink = convertVidResp.DLink
						break SearchSong
					}

				}
			}
		}
	}
	return dlink

}

func DownloadFile(fileUrl, fileName string) error {
	resp, err := http.Get(fileUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s\n", resp.Status)
	}

	out, err := os.Create(fileName + ".mp3")
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("File downloaded successfully: %s\n", fileName)
	return nil
}

func ProcessInputPlaylist(urlArg string) string {
	u, err := url.Parse(urlArg)
	if err != nil {
		log.Print(err)
		return ""
	}

	queryObj, _ := url.ParseQuery(u.RawQuery)
	playListId := queryObj["list"][0]

	return fmt.Sprintf("https://youtube.com/playlist?list=%s", playListId)

}
