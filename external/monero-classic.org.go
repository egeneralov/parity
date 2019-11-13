package external

import (
  "fmt"
  "strconv"
  "strings"
  "io/ioutil"
  "net/http"
  "github.com/antchfx/htmlquery"
)


func GetXmcHeightFromMoneroClassicOrg () (int, string) {
  resp, err := http.Get("http://explorer.monero-classic.org")
  if err != nil {
    return 0, fmt.Sprintf(`error: GetXmcHeightFromMoneroClassicOrg: Failed to execute request, err: '%s'`, err)
  }
  defer resp.Body.Close()
  
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetXmcHeightFromMoneroClassicOrg: Failed to read body, err: '%s'`, err)
  }
  
  doc, err := htmlquery.Parse(strings.NewReader(string(body)))
  if err != nil {
    return 0, fmt.Sprintf(`error: GetXmcHeightFromMoneroClassicOrg: htmlquery.Parse err: '%s'`, err)
  }


  list, err := htmlquery.QueryAll(doc, `//a[contains(@href,'block')]`)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetXmcHeightFromMoneroClassicOrg: htmlquery.QueryAll err: '%s'`, err)
  }
  
  if len(list) == 0 {
    return 0, fmt.Sprintf(`error: GetXmcHeightFromMoneroClassicOrg: list is empty`)
  }
  
  str := htmlquery.InnerText(list[0])
  
  LastBlock, err := strconv.Atoi(str)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetXmcHeightFromMoneroClassicOrg: strconv.Atoi err: '%s'`, err)
  }
  
  return LastBlock, ""
}




