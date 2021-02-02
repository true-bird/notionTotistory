package config

type Config struct {
	NotionPage NotionPage `json:"notionpage"`
	Tistory    Tistory    `json:"tistory"`
	Login      Login      `json:"login"`
	Etc        Etc        `json:"etc"`
}

type NotionPage struct {
	TokenV2        string `json:"token_v2"`
	TablePageId    string `json:"table_pageid"`
}

type Tistory struct {
	BlogName    string `json:"blogname"`
	AccessToken string `json:"accesstoken"`
}

type Login struct {
	NotionId string `json:"notion_id"`
	NotionPw string `json:"notion_pw"`
}

type Etc struct {
	ChromeDriver string `json:"chromedriver"`
	HtmlDir      string `json:"html_dir"`
	TargetDir    string `json:"target_dir"`
}
