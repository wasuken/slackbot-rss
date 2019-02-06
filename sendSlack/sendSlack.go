package sendSlack

import (
	"bytes"
	"fmt"
	"net/http"
)

type SlackMsg struct {
	Name    string
	Text    string
	Channel string
	Url     string
}

func (self *SlackMsg) PostToHookUrl() {
	jsonStr := `{"channel":"` + self.Channel + `","username":"` + self.Name + `","text":"` + self.Text + `"}`
	req, err := http.NewRequest(
		"POST",
		self.Url,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	if err != nil {
		fmt.Print(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err)
	}

	fmt.Print(resp)
	defer resp.Body.Close()
}
