package main

import (
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/oschwald/geoip2-golang"
	"github.com/ua-parser/uap-go/uaparser"
)

func (e *Event) enrich() {
	e.uaenrich()
	e.resenrich()
	e.geoenrich()
}

func (e *Event) uaenrich() {
	parser, err := uaparser.New("./regexes.yaml")
	if err != nil {
		log.Fatal(err)
	}

	client := parser.Parse(e.UserAgent)
	e.BrFamily = client.UserAgent.Family
	e.BrVersion = client.UserAgent.ToVersionString()
	e.OSFamily = client.Os.Family
	e.DeviceType = client.Device.Family
	e.OSName = client.Os.ToString()
	e.BrName = client.UserAgent.ToString()
	e.OSManufacturer = client.Device.Brand
	switch strings.ToLower(client.Device.Family) {
	case "iphone", "android":
		e.DeviceIsMobile = true
	}
}

func (e *Event) resenrich() {
	vp := strings.Split(e.Viewport, "x")
	vpwidth, _ := strconv.Atoi(vp[0])
	vpheight, _ := strconv.Atoi(vp[1])
	e.BrViewWidth = int32(vpwidth)
	e.BrViewHeight = int32(vpheight)

	res := strings.Split(e.Resolution, "x")
	reswidth, _ := strconv.Atoi(res[0])
	resheight, _ := strconv.Atoi(res[1])
	e.DeviceScreenWidth = int32(reswidth)
	e.DeviceScreenHeight = int32(resheight)
}

func (e *Event) geoenrich() {
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// If you are using strings that may be invalid, check that ip is not nil
	ip := net.ParseIP(e.UserIPAddress)
	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}
	e.GeoCity = record.City.Names["en"]
	e.GeoRegion = record.Subdivisions[0].Names["en"]
	e.GeoCountry = record.Country.Names["en"]
	//fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)
	e.GeoZipcode = record.Postal.Code
	e.GeoTimeZone = record.Location.TimeZone
	e.GeoLatitude = float32(record.Location.Latitude)
	e.GeoLongitude = float32(record.Location.Longitude)
}
