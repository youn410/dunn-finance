package sbisec

import (
  "log"
  "github.com/go-rod/rod"
)

func PrintPageTitle(browser *rod.Browser) {
  page := browser.MustPage("https://www.sbisec.co.jp").MustWaitLoad()
  log.Printf("[DEBUG] SBI 証券のページタイトル: %s", page.MustInfo().Title)
}
