package main

import (
  "os"
  "time"
  "os/signal"
  "syscall"
  "fmt"
  "log"
  "net/http"
  "flag"

  "github.com/egeneralov/parity/external"
)


var (
  HttpBindTo = "127.0.0.1:8090"
  WorkingMode = "parity-eth"
  LocalNodeRpcUrl = "http://127.0.0.1:7345"
  AllowedBlockLag = 5
  RefreshLocalState = 2
  
  // logic
  RemoteLastBlock = 0
  LocalLastBlock = 0
  isReady = false
  notReadyWithoutExternal = true
  
  // __system__
  errorString = ""
  PossibleLocalLastBlock = 0
  PossibleRemoteLastBlock = 0
  rpcUser string
  rpcPassword string
)


func main() {
  flag.StringVar(&WorkingMode, "mode", "parity-eth", "working mode [parity-eth, parity-etc, xmc, xmr, btg]")
  flag.StringVar(&HttpBindTo, "bind", "0.0.0.0:8090", "golang web server bind to")
  flag.StringVar(&LocalNodeRpcUrl, "rpcurl", "http://127.0.0.1:7345", "url to rpc server (default is parity rpc url)")
  flag.IntVar(&AllowedBlockLag, "lag", 5, "allowed lag between explorer and local node")
  flag.IntVar(&RefreshLocalState, "refresh", 5, "refresh local state every X seconds")
  flag.BoolVar(&notReadyWithoutExternal, "notReadyWithoutExternal", true, "not ready without external explorer answer")

  flag.StringVar(&rpcUser, "rpcUser", "usr", "valid only for zec")
  flag.StringVar(&rpcPassword, "rpcPassword", "pwd", "valid only for zec")
  flag.Parse()
  
  log.Printf(`WorkingMode: %s`, WorkingMode)
  
  // thread for webserver
  go func() {
    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
      
      message := fmt.Sprintf(
        `{"LocalLastBlock": "%d", "RemoteLastBlock": "%d", "AllowedBlockLag": "%d", "diff": "%d"}`,
        LocalLastBlock, RemoteLastBlock, AllowedBlockLag, RemoteLastBlock - LocalLastBlock,
      )
      
      // do not spam if browser is comming
      if (r.URL.String() == "/favicon.ico") { return }
      
      // log http request
      log.Printf("handle http request at url '%s', answer: '%s'", r.URL, message)
      
      // http code for k8s readiness check
      if isReady {
        w.WriteHeader(200)
      } else {
        w.WriteHeader(500)
      }
      // write response body
      fmt.Fprintf(w, message)
    })

    http.HandleFunc("/metrics", func (w http.ResponseWriter, r *http.Request) {
      log.Printf("handle http request at url '%s'", r.URL)
      message := fmt.Sprintf("# HELP AllowedBlockLag AllowedBlockLag\n# TYPE AllowedBlockLag gauge\nAllowedBlockLag %d\n", AllowedBlockLag)
      
      if RemoteLastBlock != 0 {
        message = fmt.Sprintf(
          "%s# HELP CurrentHeight CurrentHeight\n# TYPE CurrentHeight gauge\nCurrentHeight{type=\"remote\", daemon=\"%s\"} %d\n",
          message,
          WorkingMode,
          RemoteLastBlock,
        )
      }
      
      if LocalLastBlock != 0 {
        message = fmt.Sprintf(
          "%s# HELP CurrentHeight CurrentHeight\n# TYPE CurrentHeight gauge\nCurrentHeight{type=\"local\", daemon=\"%s\"} %d\n",
          message,
          WorkingMode,
          LocalLastBlock,
        )
      }
      
      // write response body
      fmt.Fprintf(w, message)
    })
    
    log.Printf(`Starting web server: %s`, HttpBindTo)
    log.Println(http.ListenAndServe(HttpBindTo, nil))
  }()
  
  
  //init terminating signal channel
  sigchan := make(chan os.Signal, 1)
  signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
  // check for exit signal
  go func() {
    select {
      case sig := <-sigchan:
        fmt.Printf("\nCaught signal %v: terminating\n", sig)
        os.Exit(0)
    }
  }()
  
  
  // thread for catch last block from local node
  go func() {
    for {
      
      
      switch WorkingMode {
        default:
          log.Printf(`Invalid WorkingMode: %s`, WorkingMode)
          os.Exit(1)
        case "btg":
          PossibleLocalLastBlock, errorString = external.GetHeightFromBtgRpc(LocalNodeRpcUrl)
        case "parity-eth":
          PossibleLocalLastBlock, errorString = external.GetHeightFromParityRpc(LocalNodeRpcUrl)
        case "parity-etc":
          PossibleLocalLastBlock, errorString = external.GetHeightFromParityRpc(LocalNodeRpcUrl)
        case "xmc":
          PossibleLocalLastBlock, errorString = external.GetHeightFromMoneroRpc(LocalNodeRpcUrl, "/getheight")
        case "xmr":
          PossibleLocalLastBlock, errorString = external.GetHeightFromMoneroRpc(LocalNodeRpcUrl, "/get_height")
        case "zec":
          PossibleLocalLastBlock, errorString = external.GetHeightFromZecRpc(LocalNodeRpcUrl, rpcUser, rpcPassword)
      }
      
      if errorString != "" {
        log.Printf(`RemoteHeight request error: %s`, errorString)
      } else {
        LocalLastBlock = PossibleLocalLastBlock
        log.Println(`LocalLastBlock:`, LocalLastBlock)
      }
      
      time.Sleep(time.Second * time.Duration(RefreshLocalState))
    }
  }()
  
  
  // decision flow (isReady)
  go func() {
    for {
      time.Sleep(time.Second)
      if notReadyWithoutExternal {
        if (RemoteLastBlock == 0) {
          isReady = false
          log.Println("isReady: false because: RemoteLastBlock = 0")
          continue
        }
      }
      
      if (LocalLastBlock == 0) {
        isReady = false
        log.Println("isReady: false because: LocalLastBlock = 0")
        continue
      }
      
      diff := RemoteLastBlock - LocalLastBlock
      
      log.Printf("RemoteLastBlock - LocalLastBlock = %d", diff)
      if (diff > AllowedBlockLag) {
        isReady = false
        log.Printf("isReady: false because: diff > AllowedBlockLag = %d", diff)
      } else {
        isReady = true
        log.Println("isReady: true")
      }
    } // for
  }() // decision flow (isReady)
  
  
  for {
    switch WorkingMode {
      default:
        log.Printf(`Invalid WorkingMode: %s`, WorkingMode)
        os.Exit(1)
      
      case "btg":
        time.Sleep(time.Second)
        
        PossibleRemoteLastBlock, errorString = external.GetHeightFromExplorerBitcoingoldOrg()
        if errorString != "" {
          log.Printf(`RemoteHeight request error: %s`, errorString)
        } else {
          RemoteLastBlock = PossibleRemoteLastBlock
          log.Println(`GetHeightFromExplorerBitcoingoldOrg:`, RemoteLastBlock)
        }
            
      case "parity-eth":
        time.Sleep(time.Second)
        
        PossibleRemoteLastBlock, errorString = external.GetEthHeightFromEtherscanIo()
        if errorString != "" {
          log.Printf(`RemoteHeight request error: %s`, errorString)
        } else {
          RemoteLastBlock = PossibleRemoteLastBlock
          log.Println(`GetEthHeightFromEtherscanIo:`, RemoteLastBlock)
        }
      
      case "parity-etc":
        time.Sleep(time.Second)
        
        PossibleRemoteLastBlock, errorString = external.GetEtcHeightFromGastrackerIo()
        if errorString != "" {
          log.Printf(`RemoteHeight request error: %s`, errorString)
        } else {
          RemoteLastBlock = PossibleRemoteLastBlock
          log.Println(`GetEtcHeightFromGastrackerIo:`, RemoteLastBlock)
        }
      
      case "xmc":
        time.Sleep(time.Second)
        
//         PossibleRemoteLastBlock, errorString = external.GetXmcHeightFromMoneroClassicOrg()
        PossibleRemoteLastBlock, errorString = external.GetHeightFromXmcFairHashAnswer()
        if errorString != "" {
          log.Printf(`RemoteHeight request error: %s`, errorString)
        } else {
          RemoteLastBlock = PossibleRemoteLastBlock
          log.Println(`GetXmcHeightFromMoneroClassicOrg:`, RemoteLastBlock)
        }
      
      case "xmr":
        time.Sleep(time.Second)
        
        PossibleRemoteLastBlock, errorString = external.GetXmrHeightFromMoneroBlocsInfo()
        if errorString != "" {
          log.Printf(`RemoteHeight request error: %s`, errorString)
        } else {
          RemoteLastBlock = PossibleRemoteLastBlock
          log.Println(`GetXmrHeightFromMoneroBlocsInfo:`, RemoteLastBlock)
        }
      
      case "zec":
        time.Sleep(time.Second)
        
        PossibleRemoteLastBlock, errorString = external.GetHeightFromExplorerZchaInAnswer()
        if errorString != "" {
          log.Printf(`RemoteHeight request error: %s`, errorString)
        } else {
          RemoteLastBlock = PossibleRemoteLastBlock
          log.Println(`GetXmrHeightFromMoneroBlocsInfo:`, RemoteLastBlock)
        }
    } // switch
  } // for
} // func main

