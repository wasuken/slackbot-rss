package nicoSearch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type NicoSearchResult struct {
	Meta Status
	Data []SearchResult
}
type Status struct {
	status     int
	totalCount int
	id         string
}
type SearchResult struct {
	Title         string
	ContentId     string
	ViewCounter   int
	MylistCounter int
}

func GetNicoSearchResultText(searchUrl string) string {
	resp, err := http.Get(searchUrl)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer resp.Body.Close()

	body, er := ioutil.ReadAll(resp.Body)
	if er != nil {
		fmt.Println(err)
		return ""
	}

	var n NicoSearchResult
	json.Unmarshal(body, &n)
	text := ""
	for _, nr := range n.Data[0:10] {
		url := "https://www.nicovideo.jp/watch/" + nr.ContentId
		text += "<" + url + "|" + nr.Title + ">\n"
	}
	return text
}
