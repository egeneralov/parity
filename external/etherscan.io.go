package external

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "strings"
  "net/url"
  "strconv"
  "encoding/json"
  "github.com/gorilla/websocket"
  "github.com/antchfx/htmlquery"
)


type ethWsMessage struct {
  Dashb struct {
    Marketcap       string `json:"marketcap"`
    Price           string `json:"price"`
    Lastblock       string `json:"lastblock"`
    DecOpenPrice    string `json:"decOpenPrice"`
    DecCurrentPrice string `json:"decCurrentPrice"`
  } `json:"dashb"`
  Blocks []struct {
    BNo       string `json:"b_no"`
    BTime     string `json:"b_time"`
    BMiner    string `json:"b_miner"`
    BMinerTag string `json:"b_miner_tag"`
    BTxns     string `json:"b_txns"`
    BMtime    string `json:"b_mtime"`
    BReward   string `json:"b_reward"`
  } `json:"blocks"`
  Txns []struct {
    THash            string `json:"t_hash"`
    TFrom            string `json:"t_from"`
    TTo              string `json:"t_to"`
    TContractAddress string `json:"t_contractAddress"`
    TAmt             string `json:"t_amt"`
    TTime            string `json:"t_time"`
  } `json:"txns"`
}


func GetEthHeightFromEtherscanIoWebsocket () (int, string) {
  
  u := url.URL{Scheme: "wss", Host: "etherscan.io:443", Path: "/wshandler"}
  // log.Printf("connecting to %s", u.String())
  c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
  defer c.Close()
  if err != nil {
    return 0, fmt.Sprintf(`error: GetEthHeightFromEtherscanIoWebsocket: Failed to get url: '%s', err: '%s'`, u.String(), err)
  }
  
  // welcome msg
  _, srcMessage, err := c.ReadMessage()
  if err != nil {
    return 0, fmt.Sprintf(`error: GetEthHeightFromEtherscanIoWebsocket: Failed read msg from ws, err: '%s'`, err)
  }
  // fmt.Println(string(srcMessage))
  
  // request block subscription
  err = c.WriteMessage(websocket.TextMessage, []byte(`{"event": "gs"}`))
  if err != nil {
    return 0, fmt.Sprintf("error: GetEthHeightFromEtherscanIoWebsocket: Failed to write the message to ws, err: '%s'", err)
  }

  
  // message with good payload
  _, srcMessage, err = c.ReadMessage()
  if err != nil {
    return 0, fmt.Sprintf(`error: GetEthHeightFromEtherscanIoWebsocket: Failed read msg from ws, err: '%s'`, err)
  }
  // fmt.Println(string(srcMessage))
  
  var v ethWsMessage
  err = json.Unmarshal(srcMessage, &v)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetEthHeightFromEtherscanIoWebsocket: json.Unmarshal err: '%s'`, err)
  }
  
  LastEtherBlock, err := strconv.Atoi(v.Dashb.Lastblock)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetEthHeightFromEtherscanIoWebsocket: strconv.Atoi err: '%s'`, err)
  }
  
  return LastEtherBlock, ""
}








func GetEthHeightFromEtherscanIo () (int, string) {
  // Generated by curl-to-Go: https://mholt.github.io/curl-to-go
  req, err := http.NewRequest("GET", "https://etherscan.io", nil)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetEthHeightFromEtherscanIo: GET https://etherscan.io err: '%s'`, err)
  }
  
  req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:71.0) Gecko/20100101 Firefox/71.0")
  req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
  req.Header.Set("Accept-Language", "ru,en-US;q=0.7,en;q=0.3")
  req.Header.Set("Dnt", "1")
  req.Header.Set("Connection", "keep-alive")
  req.Header.Set("Upgrade-Insecure-Requests", "1")
  req.Header.Set("Pragma", "no-cache")
  req.Header.Set("Cache-Control", "no-cache")
  
  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetEthHeightFromEtherscanIo: http.DefaultClient.Do err: '%s'`, err)
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetEthHeightFromEtherscanIo: Failed to read body, err: '%s'`, err)
  }
  
  doc, err := htmlquery.Parse(strings.NewReader(string(body)))
  if err != nil {
    return 0, fmt.Sprintf(`error: GetEthHeightFromEtherscanIo: htmlquery.Parse err: '%s'`, err)
  }

  list, err := htmlquery.QueryAll(doc, `//span[@id="lastblock"]`)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetEthHeightFromEtherscanIo: htmlquery.QueryAll err: '%s'`, err)
  }
  
  var str string
  for _, n := range list { str = htmlquery.InnerText(n) }
  
  LastEtherBlock, err := strconv.Atoi(str)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetEthHeightFromEtherscanIo: strconv.Atoi err: '%s'`, err)
  }
  
  return LastEtherBlock, ""
}



















