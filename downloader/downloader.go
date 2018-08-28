package downloader

import (
	"os"
	"encoding/json"
	"io/ioutil"
	"log"
	"fmt"
	"net/http"
	"io"
)

type downloader struct {
	config struct {
		GoogleApi string `json:"google_api"`
		GoogleCx string `json:"google_cx"`
	}
}

type apiResult struct {
	Items []struct{
		Link string
	}
}

var (
	d downloader
)

func Init() error {
	d = downloader{}

	configFile, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer configFile.Close()

	bytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &d.config)

	log.Println(d.config.GoogleApi, d.config.GoogleCx)

	return err
}

func SaveImage(i int, cat string, link string) error {
	log.Println(link)
	resp, err := http.Get(link)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	filename := fmt.Sprintf("./img/%s/%d.jpg", cat, i)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	io.Copy(file, resp.Body)
	file.Close()
	log.Println(filename + " saved")
	return nil
}


func Download(text string) error {
	var (
		start = 1
		result = apiResult{}
	)
	url := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&searchType=image&num=10&q=%s&cx=%s&start=%d",
		d.config.GoogleApi, text, d.config.GoogleCx, start)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println(string(bytes))
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return err
	}

	json.Unmarshal(bytes, &result)

	for i, item := range result.Items {
		if err = SaveImage(i, text, item.Link); err != nil {
			return err
		}
	}

	return nil
}