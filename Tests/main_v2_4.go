package main

import (
  //"bufio"
  "fmt"
  //"log"
  "os"
  "strconv"
  "TnT_v2_3"
  //"path/filepath"
)

const (
  common_root = "../roots/"
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

func cleanup(tnt *TnT_v2_3.TnTServer) {
	tnt.Kill()
}



func setup(tag string, nserver int,nservers int) (*TnT_v2_3.TnTServer, func()) {

	var peers []string = make([]string, nservers)
	var tnt *TnT_v2_3.TnTServer = new(TnT_v2_3.TnTServer)

/*	for i:=0; i<nservers; i++ {
		peers[i] = port(tag, i)
	}
*/	peers[0]="128.30.5.191"
	peers[1]="18.9.64.13"


	tnt=TnT_v2_3.StartServer(peers, nserver, common_root+"nest"+strconv.Itoa(nserver)+"/", "WatchLog"+strconv.Itoa(nserver), common_root+"tmp"+strconv.Itoa(nserver)+"/", false)
	

	clean := func() { (cleanup(tnt)) }
	return tnt, clean
}

func main() {
	nserver:=0
	nservers:=2
	tnts, clean := setup("sync", nserver,nservers)
	defer clean()

	fmt.Println("Test: File System Syncing ...")

	fmt.Println("Enter -1 to quit the loop")
	a := 100
	b := 100
	for a >= 0 && b >= 0 {

		fmt.Printf("Sync ")
		n, err := fmt.Scanf("%d %d\n", &a)
		if err != nil {
			fmt.Println("Scanf error:", n, err)
		}
		tnts.SyncWrapper(1,"./")


		fmt.Println("-----------------------------")
	}
}
