package external

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
)

type ExplorerBitcoingoldOrgAnswer struct {
	Blocks []struct {
		Height   int    `json:"height"`
		Size     int    `json:"size"`
		Hash     string `json:"hash"`
		Time     int    `json:"time"`
		Txlength int    `json:"txlength"`
		PoolInfo struct {
		} `json:"poolInfo"`
	} `json:"blocks"`
	Length     int `json:"length"`
	Pagination struct {
		Next      string `json:"next"`
		Prev      string `json:"prev"`
		CurrentTs int    `json:"currentTs"`
		Current   string `json:"current"`
		IsToday   bool   `json:"isToday"`
		More      bool   `json:"more"`
	} `json:"pagination"`
}


func GetHeightFromExplorerBitcoingoldOrg() (int, string) {
  resp, err := http.Get("https://explorer.bitcoingold.org/insight-api/blocks")
  if err != nil {
    return 0, fmt.Sprintf(`error: GetHeightFromBtgRpc: Failed to run get request %s`, err)
  }
  
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetHeightFromBtgRpc: Failed to read body, err: '%s'`, err)
  }
  
  var v ExplorerBitcoingoldOrgAnswer
  err = json.Unmarshal(body, &v)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetHeightFromBtgRpc: json.Unmarshal, err: '%s'`, err)
  }
  if len(v.Blocks) == 0 {
    return 0, fmt.Sprintf(`error: len(v.Blocks) == 0`)
  }
  
  return v.Blocks[0].Height, ""
  
}
