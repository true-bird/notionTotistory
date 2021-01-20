package notion

import (
	"fmt"
	"github.com/fedesog/webdriver"
	"github.com/kjk/notionapi"
	"github.com/tebeka/selenium"
	"strings"
	"time"
)

var (
	notionLoginUrl = "https://www.notion.so/login"
	notionId       = ""
	notionPw       = ""
	notionUrl      = "https://www.notion.so/"
)

func Login(session *webdriver.Session) {
	err := session.Url(notionLoginUrl)
	if err != nil {
		panic(err)
	}
	// login
	time.Sleep(time.Second * 1)
	email, err := session.FindElement(selenium.ByCSSSelector, "#notion-email-input-1")
	email.SendKeys(notionId)
	btn, err := session.FindElement(selenium.ByCSSSelector, "#notion-app > div > div:nth-child(1) > main > section > div > div > div > div.notion-login > div:nth-child(3) > form > div:nth-child(5)")
	btn.Click()
	time.Sleep(time.Second * 2)
	password, err := session.FindElement(selenium.ByCSSSelector, "#notion-password-input-2")
	password.SendKeys(notionPw)
	btn, err = session.FindElement(selenium.ByCSSSelector, "#notion-app > div > div:nth-child(1) > main > section > div > div > div > div.notion-login > div:nth-child(3) > form > div:nth-child(6)")
	btn.Click()
	time.Sleep(time.Second * 3)
}

func ExportPages(pageList []*notionapi.Block, session *webdriver.Session) {
	for _, page := range pageList {
		url := notionUrl + strings.ReplaceAll(page.ID, "-", "")
		fmt.Println(url)
		err := session.Url(url)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 2)
		btn, err := session.FindElement(selenium.ByCSSSelector, "#notion-app > div > div.notion-cursor-listener > div.notion-frame > div:nth-child(1) > div.notion-topbar > div:nth-child(1) > div:nth-child(3) > div.notion-topbar-more-button > svg")
		btn.Click()
		time.Sleep(time.Second * 2)
		btn, err = session.FindElement(selenium.ByCSSSelector, "#notion-app > div > div.notion-overlay-container.notion-default-overlay-container > div:nth-child(2) > div > div:nth-child(2) > div:nth-child(2) > div > div > div > div > div > div:nth-child(1) > div:nth-child(6) > div:nth-child(2) > div > div:nth-child(2) > div:nth-child(1)")
		btn.Click()
		time.Sleep(time.Second * 2)
		btn, err = session.FindElement(selenium.ByCSSSelector, "#notion-app > div > div.notion-overlay-container.notion-default-overlay-container > div:nth-child(2) > div > div:nth-child(2) > div > div:nth-child(1) > div:nth-child(2)")
		btn.Click()
		time.Sleep(time.Second * 2)
		btn, err = session.FindElement(selenium.ByCSSSelector, "#notion-app > div > div.notion-overlay-container.notion-default-overlay-container > div:nth-child(3) > div > div:nth-child(2) > div:nth-child(2) > div > div > div > div > div > div > div > div:nth-child(2) > div > div > div")
		btn.Click()
		time.Sleep(time.Second * 2)
		btn, err = session.FindElement(selenium.ByCSSSelector, "#notion-app > div > div.notion-overlay-container.notion-default-overlay-container > div:nth-child(2) > div > div:nth-child(2) > div > div:nth-child(3) > div:nth-child(2)")
		btn.Click()
		time.Sleep(time.Second * 6)
	}
}
