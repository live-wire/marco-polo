package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"math/rand"
	utils "gitlab.com/stockboi/marco-polo/lib"
)

// Marco Polo born: 1254, died: 1324
var (
	portIncoming = ":1254"
	portOutgoing = ":1324"
)

func main() {
	// http server
	hashMap := utils.NewHashMapDefault(utils.NewIpDb())
	initServers(hashMap)
}

func initServers(hashMap *utils.HashMap) {
	// Replace this with grpc endpoint
	for i:=0; i<5; i++ {
		go seedHashMap(hashMap)
	}
	initHTTP(hashMap) // blocking call
}

func initHTTP(hashMap *utils.HashMap) {
  	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/flush", httpServeWrapper(hashMap))
	http.HandleFunc("/healthcheck", func (w http.ResponseWriter, r *http.Request) {w.Write([]byte("OK"))})
	log.Printf("Marco Polo - HTTP up on %v \n", portOutgoing)
	log.Fatal(http.ListenAndServe(portOutgoing, nil))
}

func httpServeWrapper(hashMap *utils.HashMap) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(hashMap.FlushGeoJson()))
	}
}

func playWithIps() {
	ipdb := utils.NewIpDb()
	val, err := ipdb.Lookup("1.1.1.1")
	if err != nil {
		fmt.Println("IP lookup error", err)
	} else {
		fmt.Printf("latitude: %f\n", val.Latitude)
		fmt.Printf("longitude: %f\n", val.Longitude)
	}
}

func seedHashMap(hashMap *utils.HashMap) {
	max := 255
	min := 0
	getRand := func() int {return rand.Intn(max - min) + min}
	for {
		hashMap.Insert(fmt.Sprintf(`{"ip":"2.0.%d.1", "lat":12.4, "long":13.5, "tags":{"a":"b"}}`, getRand()))
		hashMap.Insert(fmt.Sprintf(`{"ip":"1.%d.1.1"}`, getRand()))
		hashMap.Insert(fmt.Sprintf(`{"ip":"20.0.%d.0", "lat":12.4, "long":13.5, "tags":{"a":"b"}}`, getRand()))
		time.Sleep(1 * time.Second)
	}
}
