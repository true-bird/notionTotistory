package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/fedesog/webdriver"
	"github.com/true-bird/notionTotistory/notion"

	//"github.com/fedesog/webdriver"
	"github.com/true-bird/notionTotistory/config"
	//"github.com/true-bird/notionTotistory/notion"
	"github.com/true-bird/notionTotistory/tistory"
	"github.com/true-bird/notionTotistory/util"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	cfg config.Config
)

func main() {
	_cfg, err := LoadConfig()
	cfg = _cfg
	nowTime := time.Now()
	c := notion.New(&cfg.NotionPage)
	pageList := c.SearchPageList()

	chromedriver := webdriver.NewChromeDriver("./chromedriver")
	defer chromedriver.Stop()
	err = chromedriver.Start()
	if err != nil {
		panic(err)
	}
	desired := webdriver.Capabilities{"Platform": "Mac"}
	required := webdriver.Capabilities{}
	session, err := chromedriver.NewSession(desired, required)
	defer session.Delete()
	if err != nil {
		panic(err)
	}

	notion.Login(&cfg.Login,session)
	notion.ExportPages(pageList, session)

	downloadHtml(nowTime)

}

func LoadConfig() (config.Config, error) {
	var c config.Config
	file, err := os.Open("./config/sample.json")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		panic(err)
	}
	return c, err
}

func downloadHtml(nowTime time.Time) {
	files, err := ioutil.ReadDir(cfg.Etc.TargetDir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "Export") {
			if nowTime.Before(file.ModTime()) {
				filePath := fmt.Sprintf("%v/%v", cfg.Etc.TargetDir, file.Name())
				unzip, err := util.Unzip(filePath, cfg.Etc.HtmlDir)
				if err != nil {
					panic(err)
				}
				post(unzip[0])
			}
		}
	}
}

func post(filePath string) {
	fmt.Println("filePath" + filePath)
	var status, category, categoryId, tagsStr string

	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		panic(err)
	}
	client := tistory.New(&cfg.Tistory)

	doc.Find("meta").Remove()
	doc.Find("title").Remove()
	doc.Find("style").Remove()

	article := doc.Find("article")
	article.AddClass("Notion_P")
	article.Find("tr").Each(func(i int, tr *goquery.Selection) {
		colName := tr.Find("th").Text()
		switch colName {
		case "상태":
			status = tr.Find("td").Text()
		case "카테고리":
			category = tr.Find("td").Text()
			categoryId = client.GetCategoryId(category)
		case "태그":
			tr.Find("td").Find("span").Each(func(i int, tag *goquery.Selection) {
				tagsStr += tag.Text() + ", "
			})
			tagsStr = tagsStr[:len(tagsStr)-2]
		}
	})

	article.Find("table").Each(func(i int, table *goquery.Selection) {
		table.Remove()
	})
	title := article.Find("header").Find("h1")
	titleText := title.Text()
	title.Remove()

	article.Find("code").Each(func(i int, code *goquery.Selection) {
		switch {
		case category == "자바" || category == "알고리즘":
			code.AddClass("java")
		case category == "go":
			code.AddClass("go")
		}
	})

	article.Find("img").Each(func(i int, img *goquery.Selection) {
		imagePath, exists := img.Attr("src")
		if exists == false {
			return
		}
		decodedValue, err := url.QueryUnescape(imagePath)
		if err != nil {
			return
		}
		imagePath = fmt.Sprintf("%v/%v", cfg.Etc.HtmlDir, decodedValue)
		imgUrl, err := UploadImage(imagePath)
		if err != nil {
			return
		}

		img.SetAttr("src", imgUrl)
		img.ParentFiltered("a").SetAttr("href", imgUrl)

	})

	html, err := doc.Html()
	if err != nil {
		panic(err)
	}
	values := client.GetValues()
	values.Set("title", titleText)
	values.Set("content", html)
	values.Set("visibility", "0")
	values.Set("category", categoryId)
	values.Set("tag", tagsStr)
	tistory.Post(values)
}

func UploadImage(filePath string) (imgUrl string, err error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("uploadedfile", filepath.Base(file.Name()))
	if err != nil {
		return
	}
	io.Copy(part, file)
	writer.Close()

	params := url.Values{
		"access_token": {cfg.Tistory.AccessToken},
		"blogName":     {cfg.Tistory.BlogName},
		"output":       {"json"},
	}
	url := "https://www.tistory.com/apis/post/attach?" + params.Encode()

	r, err := http.NewRequest("POST", url, body)
	if err != nil {
		return
	}
	r.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(r)
	defer res.Body.Close()
	if err != nil {
		return
	}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	var b map[string]map[string]string
	json.Unmarshal(resBody, &b)
	fileId := b["tistory"]["url"]
	fileId = fileId[strings.LastIndex(fileId, "/")+1 : len(fileId)-4]
	fmt.Println(fileId)
	imgUrl = fmt.Sprintf("https://t1.daumcdn.net/cfile/tistory/%s?original", fileId)
	fmt.Println(imgUrl)
	return
}
