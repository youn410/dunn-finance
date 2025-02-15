package main

import (
  "log"
  "dunn-finance/pkg/browser"
  "dunn-finance/pkg/browser/sites/sbisec"
)

func main() {
  log.Println("[INFO] rsi-fetcher starts.")

  browserInstance := browser.NewBrowser()
  defer browserInstance.Close()

  sbisec.PrintPageTitle(browserInstance)

  log.Println("[INFO] rsi-fetcher ends.")
}
