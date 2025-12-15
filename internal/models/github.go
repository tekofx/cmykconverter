package models

type GithubRelease struct {
	Url       string        `json:"url"`
	TagName   string        `json:"tag_name"`
	Extension string        `json:"extension"`
	Assets    []GithubAsset `json:"assets"`
}

type GithubAsset struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}
