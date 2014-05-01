package main

import (
  "bufio"
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "strconv"
  "TnT"
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

func cleanup(tnts []*TnT.TnTServer) {
  for i:=0; i < len(tnts); i++ {
    tnts[i].Kill()
  }
}

func readLines(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  return lines, scanner.Err()
}

func printfiles(nservers int, fname string) {
  for i:=0; i<nservers; i++ {
    path := common_root + strconv.Itoa(i) + "/" + fname
    //lines, err := readLines(path)
    data, err := ioutil.ReadFile(path)
    if err != nil {
        log.Println("printfiles: Error opening file:", err)
    } else {
        //fmt.Println(path, ":", lines, "EOF")
        fmt.Println(path, ":", string(data))
    }
  }
}

func setup(tag string, nservers int, fname string) ([]*TnT.TnTServer, func()) {

  var peers []string = make([]string, nservers)
  var tnts []*TnT.TnTServer = make([]*TnT.TnTServer, nservers)

  for i:=0; i<nservers; i++ {
    peers[i] = port(tag, i)
  }

  for i:=0; i<nservers; i++ {
    tnts[i] = TnT.StartServer(peers, i, common_root+strconv.Itoa(i)+"/", fname)
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