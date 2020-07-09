package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/google/uuid"
)

type collector struct {
	publisher Publisher
}

func (c *collector) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	var networkID string
	spCookie, err := request.Cookie("sp")
	if err != nil {
		newuuid, _ := uuid.NewRandom()
		networkID = newuuid.String()
	} else {
		networkID = spCookie.Value
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "sp",
		Value:   networkID,
		Expires: time.Now().AddDate(1, 0, 0),
		Domain:  os.Getenv("COOKIE_DOMAIN"),
	})

	// TODO set P3P header

	switch request.Method {
	case "POST":
		c.servePost(w, request, networkID)
	default:
		//fmt.Println(request.URL.Query())
		c.serveGet(w, request, networkID)
	}
}

// Handle Get Request
func (c *collector) serveGet(w http.ResponseWriter, request *http.Request,
	networkID string) {

	if request.URL.Query().Get("aid") == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Missing GET param")
		return
	}

	bodyBytes, err := urlValuesToBodyBytes(request.URL.Query())
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	SNSerr := c.jsonInputToSNS(bodyBytes, realRemoteAddr(request), request.UserAgent(), requestHeadersAsArray(request), networkID)

	if SNSerr == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Handle post request
func (c *collector) servePost(
	w http.ResponseWriter,
	request *http.Request,
	networkID string,
) {
	bodyBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write(nil) // TODO make nice message
		return
	}

	// this deserialization is used to check if we have any real data. It’s
	// slightly wasteful, but it’s certainly not a deal breaker right now and
	// the savings we get from not shipping empty events is huge

	if checkmode == true {
		fmt.Println("Checkmode POST request: " + string(bodyBytes))
	}

	trackerPayload := TrackerPayload{}
	if err := json.Unmarshal(bodyBytes, &trackerPayload); err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Bad schema")
		return
	}

	// decode trackerPayload.Data[] into proper datatypes

	if len(trackerPayload.Data) > 0 {
		SNSerr := c.jsonInputToSNS(bodyBytes, realRemoteAddr(request), request.UserAgent(), requestHeadersAsArray(request), networkID)
		if SNSerr == nil {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func (c *collector) jsonInputToSNS(bodyBytes []byte, IPAddress string, UserAgent string, Headers []string, networkID string) error {

	collectorPayload := CollectorPayload{
		Schema:        CollectorPayloadSchema,
		IPAddress:     IPAddress,
		Timestamp:     time.Now().Unix(),
		Collector:     SBVersion,
		UserAgent:     UserAgent,
		Body:          string(bodyBytes),
		Headers:       Headers,
		NetworkUserID: networkID,
		//Encoding:      "",
		//RefererURI:    "",
		//Path:          "",
		//QueryString:   "",
		//ContentType:   "",
		//Hostname:      "",
	}
	messageBytes, err := json.Marshal(collectorPayload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return err
	}
	message := string(messageBytes)
	if checkmode == true {
		fmt.Println("Checkmode SNS payload: " + message)
		if enrichcheck == true {
			processCollectorPayload(collectorPayload)
		}
	} else {
		c.publisher.publish(message)
	}

	return nil
}

/* The following works towards converting the query string to a typed
json object to later be unmarshalled into the appropiate data structures */
func urlValuesToBodyBytes(requestdata map[string][]string) ([]byte, error) {
	var tmpreturn []byte
	fixedrequestdata := make(map[string]string)

	// Remove values from arrays in Map and prepare for string -> (bool|int)
	// conversions
	for k, v := range requestdata {
		if v[0] == "0" {
			v[0] = "false"
		} else if v[0] == "1" {
			v[0] = "true"
		}
		switch k {
		case "f_pdf", "f_fla", "f_java", "f_dir", "f_qt", "f_realp", "f_wma",
			"f_gears", "pp_mix", "pp_max", "pp_miy", "pp_may", "cookie":
			k = "string_" + k
		default:

		}
		fixedrequestdata[k] = v[0]
	}

	jsonrequest, _ := json.MarshalIndent(fixedrequestdata, "", "\t")
	if checkmode == true {
		fmt.Println("Checkmode GET vars: " + string(jsonrequest))
	}
	var eventPayload []Event
	singleEvent := Event{}

	if err := json.Unmarshal(jsonrequest, &singleEvent); err != nil {
		return tmpreturn, err
	}

	// Handle string -> bool conversion in json built from GET params
	singleEvent.BrFeatPDF = singleEvent.TmpBrFeatPDF
	singleEvent.TmpBrFeatPDF = false

	singleEvent.BrFeatFl = singleEvent.TmpBrFeatFl
	singleEvent.TmpBrFeatFl = false

	singleEvent.BrFeatJava = singleEvent.TmpBrFeatJava
	singleEvent.TmpBrFeatJava = false

	singleEvent.BrFeatDir = singleEvent.TmpBrFeatDir
	singleEvent.TmpBrFeatDir = false

	singleEvent.BrFeatQT = singleEvent.TmpBrFeatQT
	singleEvent.TmpBrFeatQT = false

	singleEvent.BrFeatRealPlayer = singleEvent.TmpBrFeatRealPlayer
	singleEvent.TmpBrFeatRealPlayer = false

	singleEvent.BrFeatWinMedia = singleEvent.TmpBrFeatWinMedia
	singleEvent.TmpBrFeatWinMedia = false

	singleEvent.BrFeatGears = singleEvent.TmpBrFeatGears
	singleEvent.TmpBrFeatGears = false

	singleEvent.BrCookies = singleEvent.TmpBrCookies
	singleEvent.TmpBrCookies = false

	// Handle string - > int32 conversion in json built from GET params
	singleEvent.PPXOffsetMin = singleEvent.TmpPPXOffsetMin
	singleEvent.TmpPPXOffsetMin = 0

	singleEvent.PPXOffsetMax = singleEvent.TmpPPXOffsetMax
	singleEvent.TmpPPXOffsetMax = 0

	singleEvent.PPYOffsetMin = singleEvent.TmpPPYOffsetMin
	singleEvent.TmpPPYOffsetMin = 0

	singleEvent.PPYOffsetMax = singleEvent.TmpPPYOffsetMax
	singleEvent.TmpPPYOffsetMax = 0

	eventPayload = append(eventPayload, singleEvent)

	trackerPayload := TrackerPayload{
		Schema: CollectorPayloadSchema,
		Data:   eventPayload,
	}

	bodyBytes, err := json.Marshal(trackerPayload)
	if err != nil {
		return tmpreturn, err
	}
	return bodyBytes, nil
}

func startCollector() {

	healthHandler := &health{}
	http.Handle("/", healthHandler)
	http.Handle("/api/health", healthHandler)
	http.Handle("/health", healthHandler)

	collectorHandler := &collector{
		publisher: &SNSPublisher{
			service: sns.New(config.awsSession),
			topic:   config.snsTopic,
		},
	}

	http.Handle("/com.snowplowanalytics.snowplow/tp2", collectorHandler)
	http.Handle("/i", collectorHandler)

	portString := fmt.Sprintf(":%s", config.collectorPort)
	log.Printf("Starting server on %s", portString)
	http.ListenAndServe(portString, nil)
}
