package sbisec

import (
  "fmt"
  "log"
  "os"

  "github.com/go-rod/rod"
)

func PrintPageTitle(browser *rod.Browser) {
  page := browser.MustPage("https://www.sbisec.co.jp").MustWaitLoad()
  log.Printf("[DEBUG] SBI 証券のページタイトル: %s", page.MustInfo().Title)
}

func InitSBIsecPage(browser *rod.Browser) (page *rod.Page, err error) {
  defer func() {
    if r := recover(); r != nil {
      err = fmt.Errorf("InitSBIsecPage failed: %s", r)
    }
  }()

  log.Printf("[INFO] Open https://www.sbisec.co.jp and login.")

  page = browser.MustPage("https://www.sbisec.co.jp").MustWaitLoad()

  usernameDiv := page.MustElementR("h5", "ユーザーネーム").MustParent()
  usernameInput := usernameDiv.MustElement("input")
  usernameInput.MustInput(os.Getenv("SBI_SEC_USERNAME"))

  passwordDiv := page.MustElementR("h5", "パスワード").MustParent()
  passwordInput := passwordDiv.MustElement("input")
  passwordInput.MustInput(os.Getenv("SBI_SEC_LOGIN_PASSWORD"))

  page.MustScreenshot("./screenshot.png")

  page.MustElementR("input", "ログイン").MustClick()
  // page.MustWaitStable()

  return page, nil
}

func GoToScreeningPage(page *rod.Page) (err error) {
  defer func() {
    if r := recover(); r != nil {
      err = fmt.Errorf("GoToScreeningPage failed: %s", r)
    }
  }()

  log.Printf("[INFO] Go to sceening page.")

  naviDiv := page.MustElement("#navi01P")
  toDomesticStocksA := naviDiv.MustElement("img[title='国内株式']").MustParent()
  toDomesticStocksA.MustClick()
  page.MustWaitIdle()

  searchLinkBoxDiv := page.MustElement("#search_linkbox")
  toScreeningLi := searchLinkBoxDiv.MustElementR("a", "銘柄スクリーニング").MustParent()
  toScreeningLi.MustClick()
  page.MustWaitLoad()
  page.MustWaitIdle()
  page.MustWaitStable()

  return nil
}

type ScreeningRSIOption struct {
  Period string
  Lower float64
  Upper float64
}
type ScreeningMACDOption struct {
  Period string
  Signal string
}
type ScreeningOption struct {
  RSI *ScreeningRSIOption
  MACD *ScreeningMACDOption
}

func InputScreeningOptions(page *rod.Page, markets []string, screeningOption *ScreeningOption) (err error) {
  defer func() {
    if r := recover(); r != nil {
      err = fmt.Errorf("InputScreeningOptions failed: %s", r)
    }
  }()

  log.Printf("[INFO] Input screening options.")

  iframePage := page.MustElement("iframe").MustFrame()

  // 基本情報
  criteriaMenuBarDiv := iframePage.MustElement("div.CriteriaMenuBar")
  criteriaMenuBarDiv.MustElementR("li", "基本情報").MustClick()

  marketTr := iframePage.MustElement("tr.market")
  marketSelectAllCheckbox :=
    marketTr.MustElement("div.SelectionBox.selectall input[type='checkbox']")
  if marketSelectAllCheckbox.MustProperty("checked").Bool() {
    marketSelectAllCheckbox.MustParent().MustClick()
  }
  for _, market := range markets {
    marketCheckbox :=
      marketTr.MustElementR("label", market).MustElement("input[type='checkbox']")
    if !marketCheckbox.MustProperty("checked").Bool() {
      marketCheckbox.MustParent().MustClick()
    }
  }

  // テクニカル
  criteriaMenuBarDiv.MustElementR("li", "テクニカル").MustClick()

  screeningRSIOption := screeningOption.RSI
  if screeningRSIOption != nil {
    rsiTr := iframePage.MustElementR("td.title", "RSI").MustParent()
    rsiCheckbox := rsiTr.MustElement("input[type='checkbox']")
    if !rsiCheckbox.MustProperty("checked").Bool() {
      rsiCheckbox.MustParent().MustParent().MustClick()
    }

    rsiTr.MustElement("select[name='option']").MustSelect(screeningRSIOption.Period)
    rsiTr.MustElement("input[type='text'][name='lower']").MustSelectAllText().MustInput("0.1").MustSelectAllText().MustInput(fmt.Sprintf("%.2f", screeningRSIOption.Lower))
    rsiTr.MustElement("input[type='text'][name='upper']").MustSelectAllText().MustInput(fmt.Sprintf("%.2f", screeningRSIOption.Upper))
  }

  screeningMACDOption := screeningOption.MACD
  if screeningMACDOption != nil {
    rsiTr := iframePage.MustElementR("td.title", "MACD").MustParent()
    rsiCheckbox := rsiTr.MustElement("input[type='checkbox']")
    if !rsiCheckbox.MustProperty("checked").Bool() {
      rsiCheckbox.MustParent().MustParent().MustClick()
    }

    rsiTr.MustElement("select[name='option']").MustSelect(screeningMACDOption.Period)
    rsiTr.MustElement("select[name='option2']").MustSelect(screeningMACDOption.Signal)
  }

  searchBoxTopDiv := iframePage.MustElement("div.SearchBoxTop")
  searchBtn := searchBoxTopDiv.MustElement("div.searchbtn")
  searchBtn.MustClick()
  page.MustWaitStable()

  return nil
}

func DownloadScreenedCSV(browser *rod.Browser, page *rod.Page) (resultData []byte, err error) {
  defer func() {
    if r := recover(); r != nil {
      err = fmt.Errorf("DownloadScreeningCSV failed: %s", r)
    }
  }()

  log.Printf("[INFO] Download CSV of screening result.")

  iframePage := page.MustElement("iframe").MustFrame()

  wait := browser.MustWaitDownload()
  downloadDiv := iframePage.MustElement("div.download")
  downloadDiv.MustClick()
  screeningResultData := wait()

  return screeningResultData, nil
}
