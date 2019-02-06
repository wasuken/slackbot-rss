package main

import (
	"./sendSlack"
	"github.com/BurntSushi/toml"
	"github.com/mmcdole/gofeed"
)

type Config struct {
	Hooks HooksConfig
}
type HooksConfig struct {
	Url string
}

func main() {
	postHatenaRss()
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
	slack := &sendSlack.SlackMsg{Name: "Hook君", Text: text, Channel: "random", Url: config.Hooks.Url}
	slack.PostToHookUrl()
}
