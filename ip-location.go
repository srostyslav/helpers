package helpers

import (
	"errors"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type IPLocation struct {
	Country  string  `json:"country"`
	IP       string  `json:"ip"`
	Postal   string  `json:"postal"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Timezone string  `json:"timezone"`
	Org      string  `json:"org"`
	City     string  `json:"city"`
	Hostname string  `json:"hostname"`
	Region   string  `json:"region"`
	Loc      string  `json:"loc"`
}

func GetIPLocation(ip string) (*IPLocation, error) {
	var (
		loc IPLocation
		err error
	)

	tokens := strings.Split(os.Getenv("ipinfo_tokens"), ",")
	if len(tokens) == 0 {
		return &loc, errors.New("'ipinfo_tokens' is not found")
	}

	_, _, day := time.Now().Date()
	interval := math.Ceil(float64(31) / float64(len(tokens)))

	indx := int(math.Floor(float64(day) / interval))
	if indx >= len(tokens) {
		indx = len(tokens) - 1
	}
	token := tokens[indx]

	req := &Request{Url: "https://ipinfo.io/" + ip, Params: map[string]string{"token": token}, Json: true}
	if err = req.Get(); err != nil {
		return &loc, err
	} else if err = req.Decode(&loc); err == nil {
		if l := strings.Split(loc.Loc, ","); len(l) == 2 {
			loc.Lat, _ = strconv.ParseFloat(l[0], 64)
			loc.Lon, _ = strconv.ParseFloat(l[1], 64)
		}
	}

	return &loc, err
}
