package main

import (
	"./nicoSearch"
	"./sendSlack"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/mmcdole/gofeed"
	"os"
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

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "nico":
			postNicoSearch()
		case "hatena":
			postHatenaRss()
		default:
			fmt.Println("not found cmd")
		}
	} else {
		fmt.Println("please input cmd")
	}
}

func postNicoSearch() {
	var config Config
	_, err := toml.DecodeFile("config.tml", &config)
	if err != nil {
		_, err := toml.DecodeFile("~/.config/slackbot-rss/config.tml", &config)
		if err != nil {
			panic(err)
		}
	}
	text := nicoSearch.GetNicoSearchResultText(config.Nicos.SearchUrl)

	slack := &sendSlack.SlackMsg{Name: "Hook君", Text: text, Channel: "general", Url: config.Hooks.Url}
	slack.PostToHookUrl()
}
func postHatenaRss() {
	var config Config
	_, err := toml.DecodeFile("config.tml", &config)
	if err != nil {
		_, err := toml.DecodeFile("~/.config/slackbot-rss/config.tml", &config)
		if err != nil {
			panic(err)
		}
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
