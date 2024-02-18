package geoip

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path"

	"github.com/oschwald/maxminddb-golang"
	"userstyles.world/modules/log"
)

const (
	mirror = "https://dbip.mirror.framasoft.org/files/"
	dbname = "dbip-city-lite-latest" + ".mmdb" // or dbip-country-lite-latest
)

var client http.Client
var db *maxminddb.Reader

var dbpath = path.Join("data", "ipdb", dbname)
var hash_path_new = mirror + dbname + ".sha1sum"
var hash_path_cur = dbpath + ".sha1sum"

func GetNewHash() (hash string, err error) {
	hash_req, err := client.Get(hash_path_new)
	if err != nil {
		return "", fmt.Errorf("failed to get new hash")
	}
	defer hash_req.Body.Close()

	if hash_req.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get new hash (non-OK response)")
	}

	hash_new_bytes, err := io.ReadAll(hash_req.Body)
	if err != nil {
		return "", fmt.Errorf("failed to get new hash (io read)")
	}

	return string(hash_new_bytes), nil
}

func GetCurrentHash() (hash string, err error) {
	hash_current_bytes, err := os.ReadFile(hash_path_cur)
	if err != nil {
		return "", fmt.Errorf("failed to get current hash")
	}

	return string(hash_current_bytes), nil
}

func DownloadMMDB() error {
	// download file to data
	// save new hash to data
	return nil
}

func Initialize() {
	hash_new, err := GetNewHash()
	if err != nil {
		log.Warn.Print(err.Error())
	}

	hash_current, err := GetCurrentHash()
	if err != nil {
		log.Warn.Print(err.Error())
	}

	/* Logic for deciding to download MMDB: */
	/* - don't download if was not able to get latest hash */
	/* - don't download if current hash is already latest */
	/* otherwise download */
	if hash_new != "" && hash_new != hash_current {
		DownloadMMDB()
	}

	dbreader, err := maxminddb.Open(dbpath)
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
		return err
	}
	return location
}
