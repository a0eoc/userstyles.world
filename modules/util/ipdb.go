package util

import (
	"net"

	"github.com/oschwald/maxminddb-golang"
	"userstyles.world/modules/log"
)

var db *maxminddb.Reader

func InitIPDB() {
	dbreader, err := maxminddb.Open("data/ipdb/dbip-city-lite-latest.mmdb") // pr dbip-country-lite-latest
	if err != nil {
		log.Warn.Fatal("Failed to read mmdb: ", err.Error())
	}
	db = dbreader
}

func GetLocation(addr net.IP) any {
	var location any
	err := db.Lookup(addr, &location)
	if err != nil {
		log.Warn.Print("Failed to lookup address: ", err.Error())
		// need to return something indicating this?
	}
	return location
}
