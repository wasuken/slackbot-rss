package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"slackbot-rss/nicoSearch"
	"slackbot-rss/sendSlack"

	"github.com/BurntSushi/toml"
	"github.com/mmcdole/gofeed"
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

var DEEPL_URL string = "https://script.google.com/macros/s/AKfycbxpY_dGm4hJWrq3pcBdv7VokXcnoPXxtKYTz6YNtlQ9VOK59Zi9/exec"

type DeepLResult struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func deeplTranslate(text, src, target string) string {
	values := url.Values{}
	values.Add("text", text)
	values.Add("source", src)
	values.Add("target", target)
	resp, err := http.Get(DEEPL_URL + "?" + values.Encode())

	if err != nil {
		panic("get error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("body parse error")
	}
	var result DeepLResult
	if err := json.Unmarshal(body, &result); err != nil {
		panic("json parse error")
	}
	if result.Code == 400 {
		panic("api response error")
	}
	return result.Text
}

// アスキー文字率
func asciiTextRate(s string) float64 {
	s_len := utf8.RuneCountInString(s)
	cnt := 0.0
	for _, c := range s {
		if c <= 122 {
			cnt++
		}
	}
	return cnt / float64(s_len)
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
			t := item.Title
			// ここで翻訳をいれるかどうか判定する。
			if asciiTextRate(item.Title) > 0.6 {
				t = deeplTranslate(item.Title, "en", "ja")
			}

			linkList = append(linkList, "<"+item.Link+"|"+t+">")
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
