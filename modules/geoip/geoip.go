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

var dbpath = path.Join("data", "geoip", dbname)
var hash_path_latest = mirror + dbname + ".sha1sum"
var hash_path_current = dbpath + ".sha1sum"

func GetLatestHash() (hash string, err error) {
	hash_req, err := client.Get(hash_path_latest)
	if err != nil {
		return "", fmt.Errorf("failed to get latest hash")
	}
	defer hash_req.Body.Close()

	if hash_req.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get latest hash (non-OK response)")
	}

	hash_latest_bytes, err := io.ReadAll(hash_req.Body)
	if err != nil {
		return "", fmt.Errorf("failed to get latest hash (io read)")
	}

	return string(hash_latest_bytes), nil
}

func GetCurrentHash() (hash string, err error) {
	hash_current_bytes, err := os.ReadFile(hash_path_current)
	if err != nil {
		return "", fmt.Errorf("failed to get current hash")
	}

	return string(hash_current_bytes), nil
}

func SaveCurrentHash(hash string) error {
	err := os.WriteFile(hash_path_current, []byte(hash), 0o755)
	if err != nil {
		return err
	}
	return nil
}

func DownloadMMDB() error {
	// can download .gz and extract. 2 times less bandwidth and time spent downloading.

	log.Warn.Print("Downloading latest MMDB...")
	file_req, err := http.Get(mirror + dbname)
	if err != nil {
		return err
	}

	if file_req.StatusCode != http.StatusOK {
		return fmt.Errorf("non-OK response")
	}

	file_bytes, err := io.ReadAll(file_req.Body)
	if err != nil {
		return err
	}

	err = os.WriteFile(dbpath, file_bytes, 0o755)
	if err != nil {
		return err
	}
	return nil
}

func Initialize() {
	// can verify hash of actual file instead of reading it from .sha1sum file

	hash_latest, err := GetLatestHash()
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
	if hash_latest != "" && hash_latest != hash_current {
		err = DownloadMMDB()
		if err != nil {
			log.Warn.Fatal("Failed to download latest mmdb: ", err.Error())
		} else {
			/* only save latest hash if the download was successful */
			SaveCurrentHash(hash_latest)
		}
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
