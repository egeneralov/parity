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

  "gitlab.com/egeneralov/parity/external"
)


var (
  HttpBindTo = "127.0.0.1:8090"
  WorkingMode = "parity-eth"
  LocalNodeRpcUrl = "http://127.0.0.1:7345"
  AllowedBlockLag = 5
  
  // logic
  RemoteLastBlock = 0
  LocalLastBlock = 0
  isReady = false
  
  // __system__
  errorString = ""
  PossibleLocalLastBlock = 0
  PossibleRemoteLastBlock = 0
)


func main() {
  flag.StringVar(&WorkingMode, "mode", "parity-eth", "working mode [parity-eth, parity-etc]")
  flag.StringVar(&HttpBindTo, "bind", "0.0.0.0:8090", "golang web server bind to")
  flag.StringVar(&LocalNodeRpcUrl, "rpcurl", "http://127.0.0.1:7345", "url to rpc server (default is parity rpc url)")
  flag.IntVar(&AllowedBlockLag, "lag", 5, "allowed lag between explorer and local node")
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
      
/*
      message := fmt.Sprintf("# HELP LocalLastBlock LocalLastBlock\n# TYPE LocalLastBlock gauge\nLocalLastBlock %d\n", LocalLastBlock)
      message = fmt.Sprintf("%s# HELP RemoteLastBlock RemoteLastBlock\n# TYPE RemoteLastBlock gauge\nRemoteLastBlock{WorkingMode=\"%s\"} %d\n", message, WorkingMode, RemoteLastBlock)
      message = fmt.Sprintf("%s# HELP AllowedBlockLag AllowedBlockLag\n# TYPE AllowedBlockLag gauge\nAllowedBlockLag %d\n", message, AllowedBlockLag)
      message = fmt.Sprintf("%s# HELP diff diff\n# TYPE diff gauge\ndiff %d\n", message, RemoteLastBlock - LocalLastBlock)
      switch WorkingMode {
        default:
          log.Printf(`Invalid WorkingMode: %s`, WorkingMode)
          os.Exit(1)
        
        case "parity-eth":
          message = fmt.Sprintf("%s# HELP GetEthHeightFromEtherscanIo GetEthHeightFromEtherscanIo\n# TYPE GetEthHeightFromEtherscanIo gauge\nGetEthHeightFromEtherscanIo %d\n", message, RemoteLastBlock)
        case "parity-etc":
          message = fmt.Sprintf("%s# HELP GetEtcHeightFromGastrackerIo GetEtcHeightFromGastrackerIo\n# TYPE GetEtcHeightFromGastrackerIo gauge\nGetEtcHeightFromGastrackerIo %d\n", message, RemoteLastBlock)
      }
*/
      
      message = fmt.Sprintf(
        "%s# HELP CurrentHeight CurrentHeight\n# TYPE CurrentHeight gauge\nCurrentHeight{type=\"remote\", daemon=\"%s\"} %d\n",
        message,
        WorkingMode,
        RemoteLastBlock,
      )
      
      message = fmt.Sprintf(
        "%s# HELP CurrentHeight CurrentHeight\n# TYPE CurrentHeight gauge\nCurrentHeight{type=\"local\", daemon=\"%s\"} %d\n",
        message,
        WorkingMode,
        LocalLastBlock,
      )
      
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
      
      PossibleLocalLastBlock, errorString = external.GetHeightFromParityRpc(LocalNodeRpcUrl)
      
      if errorString != "" {
        log.Printf(`RemoteHeight request error: %s`, errorString)
      } else {
        LocalLastBlock = PossibleLocalLastBlock
        log.Println(`LocalLastBlock:`, LocalLastBlock)
      }
      
      time.Sleep(time.Second)
    }
  }()
  
  
  // decision flow (isReady)
  go func() {
    for {
      time.Sleep(time.Second)
      
      if (RemoteLastBlock == 0) {
        isReady = false
        log.Println("isReady: false because: RemoteLastBlock = 0")
        continue
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
      
    } // switch
  } // for
} // func main

