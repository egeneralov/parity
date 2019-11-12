package external

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
)


type moneroBlocsInfoAnswer struct {
	Difficulty    int64   `json:"difficulty"`
	Height        int     `json:"height"`
	Hashrate      float64 `json:"hashrate"`
	TotalEmission string  `json:"total_emission"`
	LastReward    int64   `json:"last_reward"`
	LastTimestamp int     `json:"last_timestamp"`
}

func GetXmrHeightFromMoneroBlocsInfo () (int, string) {
  resp, err := http.Get("https://moneroblocks.info/api/get_stats")
  if err != nil {
    return 0, fmt.Sprintf(`error: GetXmrHeightFromMoneroBlocsInfo: Failed to run get request '%s'`, err)
  }
  
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetXmrHeightFromMoneroBlocsInfo: Failed to read body '%s'`, err)
  }
  
  var v moneroBlocsInfoAnswer
  err = json.Unmarshal(body, &v)
  if err != nil {
    return 0, fmt.Sprintf(`error: GetXmrHeightFromMoneroBlocsInfo: Failed json.Unmarshal '%s'`, err)
  }
  
  return v.Height, ""
}

