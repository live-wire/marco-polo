package lib

import (
	"log"
	"sync"

	"github.com/ip2location/ip2location-go"
)

var dataPath string = "db/IP2LOCATION-LITE-DB5.BIN"

// IpDb is a wrapper for an IP2location Database to make it thread safe
type IpDb struct {
	Db *ip2location.DB
	L  sync.Mutex
}

// NewIpDb creates a new instance of an IpDb
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

// Lookup is a wrapper for Ip2Location Get_all
func (x *IpDb) Lookup(ip string) (ip2location.IP2Locationrecord, error) {
	// x.L.Lock()
	// defer x.L.Unlock()
	return x.Db.Get_all(ip)
}
