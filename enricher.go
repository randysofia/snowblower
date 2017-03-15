package main

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/oschwald/geoip2-golang"
	"github.com/ua-parser/uap-go/uaparser"
)

func (e *Event) enrich() {
	e.uaenrich()
	e.resenrich()
	e.geoenrich()
	e.urlenrich()
}

//urlenrich - URL enrichment
func (e *Event) urlenrich() {
	urldata, _ := url.Parse(e.PageURL)
	refdata, _ := url.Parse(e.PageReferrer)
	e.PageURLScheme = urldata.Scheme
	e.ReferrerURLScheme = refdata.Scheme
	e.PageURLHost = urldata.Host
	e.ReferrerURLHost = refdata.Host
	e.PageURLPath = urldata.Path
	e.ReferrerURLPath = refdata.Path
	e.PageURLQuery = urldata.RawQuery
	e.ReferrerURLQuery = refdata.RawQuery
	e.PageURLFragment = urldata.Fragment
	e.ReferrerURLFragment = refdata.Fragment
	e.PageURLPort = 80     // temporarily
	e.ReferrerURLPort = 80 // temporarily
}

//uaenrich - user agent enrichment
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
	case "iphone", "android", "ios":
		e.DeviceIsMobile = true
	}
	if strings.Contains(strings.ToLower(client.UserAgent.Family), "mobile") {
		e.DeviceIsMobile = true
	}
}

// resenrich - Add resolution data
func (e *Event) resenrich() {
	vp := strings.Split(e.Viewport, "x")
	if len(vp) > 1 {
		vpwidth, _ := strconv.Atoi(vp[0])
		vpheight, _ := strconv.Atoi(vp[1])
		e.BrViewWidth = int32(vpwidth)
		e.BrViewHeight = int32(vpheight)
	}

	res := strings.Split(e.Resolution, "x")
	if len(res) > 1 {
		reswidth, _ := strconv.Atoi(res[0])
		resheight, _ := strconv.Atoi(res[1])
		e.DeviceScreenWidth = int32(reswidth)
		e.DeviceScreenHeight = int32(resheight)
	}
}

//geoenrich Geolocation enrichment
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
		fmt.Println(err)
		return
	}
	e.GeoCity = record.City.Names["en"]
	if len(record.Subdivisions) > 0 {
		e.GeoRegion = record.Subdivisions[0].IsoCode
		e.GeoRegionName = record.Subdivisions[0].Names["en"]
	}
	e.GeoCountry = record.Country.IsoCode
	//fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)
	e.GeoZipcode = record.Postal.Code
	e.GeoTimeZone = record.Location.TimeZone
	e.GeoLatitude = float32(record.Location.Latitude)
	e.GeoLongitude = float32(record.Location.Longitude)
}

// Validate this event, returning false should prevent saving
func (e *Event) validate() bool {
	return true
}
