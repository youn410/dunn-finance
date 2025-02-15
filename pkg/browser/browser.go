package browser

import (
  "github.com/go-rod/rod"
  "github.com/go-rod/rod/lib/launcher"
)

func NewBrowser() *rod.Browser {
  url := launcher.New().MustLaunch()
  browser := rod.New().ControlURL(url).MustConnect()

  return browser
}
