package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	utils "gitlab.com/stockboi/marco-polo/lib"
)

// Marco Polo born: 1254, died: 1324
var (
	portIncoming = ":1254"
	portOutgoing = ":1324"
)

func main() {
	// http server
	hashMaps := make(map[string]*utils.HashMap)
	hashMap := utils.NewHashMapDefault(utils.NewIpDb())
	hashMaps["default"] = hashMap
	initServers(hashMaps)
}

func initServers(hashMaps map[string]*utils.HashMap) {
	// Replace this with grpc endpoint
	for i := 0; i < 5; i++ {
		go seedHashMap(hashMaps["default"])
	}
	initHTTP(hashMaps) // blocking call
}

func initHTTP(hashMaps map[string]*utils.HashMap) {
	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.Dir("./static")))
	r.HandleFunc("/list", httpServeWrapperList(hashMaps))
	r.HandleFunc("/flush", httpServeWrapperAll(hashMaps))
	r.HandleFunc("/flush/{src}", httpServeWrapper(hashMaps))
	r.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("OK")) })
	http.Handle("/", r)
	log.Printf("Marco Polo - HTTP up on %v \n", portOutgoing)
	log.Fatal(http.ListenAndServe(portOutgoing, nil))
}

func httpServeWrapperList(hashMaps map[string]*utils.HashMap) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		retMap := make(map[string][]string)
		retMap["list"] = make([]string, 0)
		for k := range hashMaps {
			retMap["list"] = append(retMap["list"], k)
		}
		ret, _ := json.Marshal(retMap)
		w.Write([]byte(ret))
	}
}

func httpServeWrapper(hashMaps map[string]*utils.HashMap) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		src := mux.Vars(req)["src"]
		log.Println("Request received from:", src)
		ret := fmt.Sprintf(`{"%s":[]}`, src)
		if val, ok := hashMaps[src]; ok {
			//do something here
			ret = fmt.Sprintf(`{"%s":%s}`, src, val.FlushGeoJson())
		}
		w.Write([]byte(ret))
	}
}

func httpServeWrapperAll(hashMaps map[string]*utils.HashMap) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		retMap := make(map[string][]*utils.GeoJson)
		ret := []byte(`{"default": []}`)
		for k, v := range hashMaps {
			retMap[k] = v.GetGeoJson()
		}
		ret, _ = json.Marshal(retMap)
		w.Write(ret)
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
	getRand := func() int { return rand.Intn(max-min) + min }
	for {
		hashMap.InsertString(fmt.Sprintf(`{"ip":"2.0.%d.1", "lat":12.4, "long":13.5, "tags":{"a":"b"}}`, getRand()))
		hashMap.InsertString(fmt.Sprintf(`{"ip":"1.%d.1.1"}`, getRand()))
		hashMap.InsertString(fmt.Sprintf(`{"ip":"20.0.%d.0", "lat":12.4, "long":13.5, "tags":{"a":"b"}}`, getRand()))
		time.Sleep(1 * time.Second)
	}
}
