package external

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
)

type moneroRpcAnswer struct {
	Height    int    `json:"height"`
	Status    string `json:"status"`
	Untrusted bool   `json:"untrusted"`
}

func GetHeightFromMoneroRpc(url string, suffix string) (int, string) {
  resp, err := http.Get(url + suffix)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetHeightFromMoneroRpc: Failed to run get request %s`, err)
  }
  
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetHeightFromMoneroRpc: Failed to read body, err: '%s'`, err)
  }
  
//   fmt.Println(string(body))



//   var v interface{}
  var v moneroRpcAnswer
  err = json.Unmarshal(body, &v)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetHeightFromMoneroRpc: json.Unmarshal, err: '%s'`, err)
  }
  return v.Height, ""
/*
//   fmt.Println(string(body))
  
  mm := v.(map[string]interface{})
//   mm := v.(map[string]int)
  vvv := mm["Height"].(map[int]int)
  
  return int(mm["Height"]), ""
*/
  
}
