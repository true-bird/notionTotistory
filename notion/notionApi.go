package notion

import (
	"github.com/kjk/notionapi"
	"github.com/true-bird/notionTotistory/config"
)

type Client struct {
	notionPage *config.NotionPage
	notionapi.Client
}

func New(page *config.NotionPage) *Client {
	return &Client{
		page,
		notionapi.Client{
			AuthToken: page.TokenV2,
		},
	}
}

func (client *Client) SearchPageList() (pages []*notionapi.Block) {

	page, err := client.DownloadPage(client.notionPage.TablePageId)
	if err != nil {
		panic(err)
	}
	tableViews := page.TableViews[0]
	for _, row := range tableViews.Rows {
		if len(row.Columns[1]) == 0 {
			continue
		}
		if row.Columns[1][0].Text == "발행준비 완료" {
			pages = append(pages, row.Page)
		}
	}
	return
}
