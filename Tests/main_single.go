package main

import (
  //"bufio"
  //"log"
  "fmt"
  "io/ioutil"
  "os"
  "strconv"
  "TnT_single_v2_1"
)

const (
  common_root = "../roots/root"
)

func port(tag string, host int) string {
  s := "/var/tmp/824-"
  s += strconv.Itoa(os.Getuid()) + "/"
  os.Mkdir(s, 0777)
  s += "tnt-"
  s += strconv.Itoa(os.Getpid()) + "-"
  s += tag + "-"
  s += strconv.Itoa(host)
  return s
}

func cleanup(tnts []*TnT_single_v2_1.TnTServer) {
  for i:=0; i < len(tnts); i++ {
    tnts[i].Kill()
  }
}

func printfiles(nservers int, fname string) {
  for i:=0; i<nservers; i++ {
    path := common_root + strconv.Itoa(i) + "/" + fname
    data, err := ioutil.ReadFile(path)
    if err != nil {
        fmt.Println(path, ": <!>\n")
    } else {
        fmt.Println(path, ":", string(data))
    }
  }
}

func setup(tag string, nservers int, fname string) ([]*TnT_single_v2_1.TnTServer, func()) {

  var peers []string = make([]string, nservers)
  var tnts []*TnT_single_v2_1.TnTServer = make([]*TnT_single_v2_1.TnTServer, nservers)

  for i:=0; i<nservers; i++ {
    peers[i] = port(tag, i)
  }

  for i:=0; i<nservers; i++ {
    tnts[i] = TnT_single_v2_1.StartServer(peers, i, common_root+strconv.Itoa(i)+"/", fname)
  }

  clean := func() { (cleanup(tnts)) }
  return tnts, clean
}

func main() {

  const nservers = 3
  const fname = "foo.txt"

  printfiles(nservers, fname)

  tnts, clean := setup("sync", nservers, fname)
  defer clean()

  fmt.Println("Test: Single File Syncing ...")

  fmt.Println("Enter -1 to quit the loop")
  a := 100
  b := 100
  for a >= 0 && b >= 0 {

      fmt.Printf("Sync? Enter (who) and (from): ")
      n, err := fmt.Scanf("%d %d\n", &a, &b)
      if err != nil {
          fmt.Println("Scanf error:", n, err)
      }

      if 0 <= a && a < nservers && 0 <= b && b < nservers && a != b {
          tnts[a].SyncNow(b)
          printfiles(nservers, fname)
      }

      fmt.Println("-----------------------------")
  }
}
