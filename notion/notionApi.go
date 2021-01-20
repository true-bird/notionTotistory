package notion

import "github.com/kjk/notionapi"

var (
	tokenV2     = ""
	tablePageId = ""
)

type Client struct {
	notionapi.Client
}

func New() *Client {
	return &Client{
		notionapi.Client{
			AuthToken: tokenV2,
		},
	}
}

func (client *Client) SearchPageList() (pages []*notionapi.Block) {

	page, err := client.DownloadPage(tablePageId)
	if err != nil {
		panic(err)
	}
	tableViews := page.TableViews[0]
	for _, row := range tableViews.Rows {
		if row.Columns[1][0].Text == "발행준비 완료" {
			pages = append(pages, row.Page)
		}
	}
	return
}
