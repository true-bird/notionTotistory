package tistory

import (
	"encoding/json"
	"github.com/true-bird/notionTotistory/config"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	values url.Values
}

func New(tistory *config.Tistory) *Client {
	return &Client{
		url.Values{
			"access_token": {tistory.AccessToken},
			"blogName":     {tistory.BlogName},
			"output":       {"json"},
		},
	}
}

func (client *Client) GetValues() url.Values {
	return client.values
}

func Post(values url.Values) {
	res, err := http.PostForm("https://www.tistory.com/apis/post/write", values)
	defer res.Body.Close()
	if err != nil {
		panic(err)
	}
}

func (client *Client) GetCategoryId(target string) string {
	url := "https://www.tistory.com/apis/category/list?" + client.values.Encode()
	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var result map[string]map[string]map[string][]map[string]string
	json.Unmarshal(data, &result)
	categories := result["tistory"]["item"]["categories"]
	for _, category := range categories {
		if category["label"] == target {
			return category["id"]
		}
	}
	return ""
}
