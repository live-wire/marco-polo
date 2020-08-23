package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	utils "github.com/live-wire/marco-polo/lib"
	pb "github.com/live-wire/marco-polo/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type marcoPoloService struct {
	IPDb           *utils.IpDb
	HashMaps       map[string]*utils.HashMap
	SeedRandomData *bool
}

func (s *marcoPoloService) Consume(ctx context.Context, in *pb.Message) (*pb.Burp, error) {
	valid, message := Validate(in)
	if !valid {
		return &pb.Burp{Status: 400, Message: message}, errors.New(message)
	}
	if _, ok := s.HashMaps[in.Src]; !ok {
		s.HashMaps[in.Src] = utils.NewHashMapDefault(s.IPDb)
	}
	s.HashMaps[in.Src].Insert(getMarcoPoloMessage(in))
	return &pb.Burp{Status: 200, Message: "ok"}, nil
}

// Validate validates input to MarcoPolo
func Validate(in *pb.Message) (valid bool, message string) {
	if len(in.Ip) == 0 {
		return false, "Ip field is compulsary"
	}
	if len(in.Src) == 0 {
		in.Src = "default"
	}
	return true, ""
}

func getMarcoPoloMessage(in *pb.Message) *utils.MarcoPoloMessage {
	return &utils.MarcoPoloMessage{
		Ip:   in.Ip,
		Src:  in.Src,
		Lat:  in.Lat,
		Long: in.Long,
		Tags: in.Tags,
	}
}

// Marco Polo born: 1254, died: 1324
var (
	portIncoming = ":1254"
	portOutgoing = ":1324"
)

func main() {
	// http server
	seedRandomData := flag.Bool("dummy", false, "Seed random data ?")
	flag.Parse()
	hashMaps := make(map[string]*utils.HashMap)
	ipDb := utils.NewIpDb()
	hashMap := utils.NewHashMapDefault(ipDb)
	hashMaps["default"] = hashMap
	marcoPoloService := &marcoPoloService{IPDb: ipDb, HashMaps: hashMaps, SeedRandomData: seedRandomData}
	initServers(marcoPoloService)
}

func initServers(marcoPoloService *marcoPoloService) {
	// Seeds random data to default's HashMap
	if *marcoPoloService.SeedRandomData {
		for i := 0; i < 5; i++ {
			go seedHashMap(marcoPoloService.HashMaps["default"])
		}
	}
	go initGrpc(marcoPoloService)
	initHTTP(marcoPoloService) // blocking call
}

func initGrpc(marcoPoloService *marcoPoloService) {
	lis, err := net.Listen("tcp", portIncoming)
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}
	server := grpc.NewServer()
	pb.RegisterMarcoPoloServiceServer(server, marcoPoloService)
	log.Printf("Marco Polo - gRPC up on %v \n", portIncoming)
	if err := server.Serve(lis); err != nil {
		log.Fatal(errors.Wrap(err, "Failed to start server!"))
	}
}

func initHTTP(marcoPoloService *marcoPoloService) {
	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.Dir("./static")))
	r.HandleFunc("/list", marcoPoloService.httpServeWrapperList)
	r.HandleFunc("/flush", marcoPoloService.httpServeWrapperAll)
	r.HandleFunc("/flush/{src}", marcoPoloService.httpServeWrapper)
	r.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("OK")) })
	http.Handle("/", r)
	log.Printf("Marco Polo - HTTP up on %v \n", portOutgoing)
	log.Fatal(http.ListenAndServe(portOutgoing, nil))
}

func (s *marcoPoloService) httpServeWrapperList(w http.ResponseWriter, req *http.Request) {
	retMap := make(map[string][]string)
	retMap["list"] = make([]string, 0)
	for k := range s.HashMaps {
		retMap["list"] = append(retMap["list"], k)
	}
	ret, _ := json.Marshal(retMap)
	w.Write([]byte(ret))
}

func (s *marcoPoloService) httpServeWrapper(w http.ResponseWriter, req *http.Request) {
	src := mux.Vars(req)["src"]
	log.Println("Request received from:", src)
	ret := fmt.Sprintf(`{"%s":[]}`, src)
	if val, ok := s.HashMaps[src]; ok {
		//do something here
		ret = fmt.Sprintf(`{"%s":%s}`, src, val.FlushGeoJson())
	}
	w.Write([]byte(ret))
}

func (s *marcoPoloService) httpServeWrapperAll(w http.ResponseWriter, req *http.Request) {
	retMap := make(map[string][]*utils.GeoJson)
	ret := []byte(`{"default": []}`)
	for k, v := range s.HashMaps {
		retMap[k] = v.GetGeoJson()
	}
	ret, _ = json.Marshal(retMap)
	w.Write(ret)
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
