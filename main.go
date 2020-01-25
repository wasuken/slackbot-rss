package main

import (
	"./nicoSearch"
	"./sendSlack"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/mmcdole/gofeed"
	"os"
	"strings"
	"time"
)

type Config struct {
	Hooks HooksConfig
	Nicos NicosConfig
	Rss   RssConfig
}
type HooksConfig struct {
	Url string
}
type NicosConfig struct {
	SearchUrl string
}

type RssConfig struct {
	UrlList []string
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "nico":
			postNicoSearch()
		case "rss":
			postRss()
		default:
			fmt.Println("not found cmd")
		}
	} else {
		fmt.Println("please input cmd")
	}
}

// 存在確認ができたパス文字列を返す。
// 全部ダメならnil
func loadFiles(filepaths []string) string {
	_, err := os.Stat(filepaths[0])
	if err != nil {
		return loadFiles(filepaths[1:])
	} else {
		return filepaths[0]
	}
}

var DEFAULT_LOAD_FILES []string = []string{
	"config.tml",
	"~/.config/slackbot-rss/config.tml",
	"/etc/slackbot/config.tml"}

func postNicoSearch() {

	var config Config
	_, err := toml.DecodeFile(loadFiles(DEFAULT_LOAD_FILES), &config)
	if err != nil {
		panic(err)
	}
	text := nicoSearch.GetNicoSearchResultText(config.Nicos.SearchUrl)

	slack := &sendSlack.SlackMsg{Name: "Hook君", Text: text, Channel: "voiceroid", Url: config.Hooks.Url}
	slack.PostToHookUrl()
}
func postRss() {
	var config Config
	_, err := toml.DecodeFile(loadFiles(DEFAULT_LOAD_FILES), &config)
	if err != nil {
		panic(err)
	}
	sep := 0
	for _, url := range config.Rss.UrlList {
		linkList := []string{}
		text := ""
		fp := gofeed.NewParser()
		feed, _ := fp.ParseURL(url)
		text += "## " + url + " ##\n"
		for _, item := range feed.Items {
			linkList = append(linkList, "<"+item.Link+"|"+item.Title+">")
		}
		sep = int(len(linkList) / 2)
		text = strings.Join(linkList[0:sep], "\n")
		slack := &sendSlack.SlackMsg{Name: "Hook君", Text: text, Channel: "news", Url: config.Hooks.Url}
		slack.PostToHookUrl()

		text = strings.Join(linkList[sep:], "\n")
		slack = &sendSlack.SlackMsg{Name: "Hook君", Text: text, Channel: "news", Url: config.Hooks.Url}
		slack.PostToHookUrl()
		time.Sleep(time.Second * 5)
	}
}
