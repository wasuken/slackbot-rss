package nicoSearch_test

import (
	"../nicoSearch"
	"fmt"
	"github.com/BurntSushi/toml"
	"testing"
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

func TestGetNicoSearchResultTextSuccess(t *testing.T) {
	var config Config
	_, err := toml.DecodeFile("../config.tml", &config)
	if err != nil {
		panic(err)
	}
	text := nicoSearch.GetNicoSearchResultText("http://api.search.nicovideo.jp/api/v2/video/contents/search?q=VOICEROID%E5%AE%9F%E6%B3%81%E3%83%97%E3%83%AC%E3%82%A4Part1%E3%83%AA%E3%83%B3%E3%82%AF&targets=title,tags,tagsExact&_sort=-lastCommentTime&_context=apiguide&fields=title,contentId,viewCounter,mylistCounter&_limit=100")
	fmt.Println(text)
}
