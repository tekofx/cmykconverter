package models

type GithubRelease struct {
	Url       string        `json:"url"`
	TagName   string        `json:"tag_name"`
	Extension string        `json:"extension"`
	Assets    []GithubAsset `json:"assets"`
}

type GithubAsset struct {
	BrowserDownloadUrl string `json:"browser_download_url"`
	Name               string `json:"name"`
}
