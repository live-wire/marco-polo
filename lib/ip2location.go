package lib

import (
	"sync"
	"log"
	"github.com/ip2location/ip2location-go"
)

var dataPath string = "db/IP2LOCATION-LITE-DB5.BIN"
type IpDb struct {
	Db *ip2location.DB
	L sync.Mutex
}

func NewIpDb() *IpDb {
	db, err := ip2location.OpenDB(dataPath)
	if err != nil {
		log.Fatal("IP2Location couldn't be set up", err)
	} else {
		log.Println("IpDb setup complete.")
	}
	ipdb := IpDb{}
	ipdb.Db = db
	return &ipdb
}

func (x *IpDb) Lookup(ip string) (ip2location.IP2Locationrecord,error) {
	// x.L.Lock()
	// defer x.L.Unlock()
	return x.Db.Get_all(ip)
}
