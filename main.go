package main

import (
	"./sendSlack"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/mmcdole/gofeed"
	"io/ioutil"
	"net/http"
)

type Config struct {
	Hooks HooksConfig
	Nicos NicosConfig
}
type HooksConfig struct {
	Url string
}
type NicosConfig struct {
	SearchUrl string
}

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

func main() {
	postNicoSearch()
}

// TODO: パッケージを分ける。
// TODO: 検索ワードだけは別途入力分けする。
func postNicoSearch() {
	var config Config
	_, err := toml.DecodeFile("config.tml", &config)
	if err != nil {
		panic(err)
	}
	resp, err := http.Get(config.Nicos.SearchUrl)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	body, er := ioutil.ReadAll(resp.Body)
	if er != nil {
		fmt.Println(err)
		return
	}

	var n NicoSearchResult
	json.Unmarshal(body, &n)
	text := ""
	for _, nr := range n.Data[0:10] {
		url := "https://www.nicovideo.jp/watch/" + nr.ContentId
		text += nr.Title + "\n" + url + "\n"
	}
	slack := &sendSlack.SlackMsg{Name: "Hook君", Text: text, Channel: "general", Url: config.Hooks.Url}
	slack.PostToHookUrl()
}
func postHatenaRss() {
	var config Config
	_, err := toml.DecodeFile("config.tml", &config)
	if err != nil {
		panic(err)
	}
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("http://b.hatena.ne.jp/hotentry/it.rss")
	text := ""
	for _, item := range feed.Items {
		text += item.Title + "\n" + item.Link + "\n"
	}
	slack := &sendSlack.SlackMsg{Name: "Hook君", Text: text, Channel: "general", Url: config.Hooks.Url}
	slack.PostToHookUrl()
}
