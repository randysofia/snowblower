package main

import (
	"github.com/spark451/snowman/snowplow"
)

// Event is an alias for snowplow.Event
type Event snowplow.Event

// CollectorPayloadSchema ...
const CollectorPayloadSchema = "iglu:com.snowplowanalytics.snowplow/CollectorPayload/thrift/1-0-0"

// SBVersion Version of collector / enricher
const SBVersion = "Snowblower/0.0.1"

// CollectorPayload defines the structure of data posted from Snowplow collectors
type CollectorPayload struct {
	Schema        string   `json:"schema"`
	IPAddress     string   `json:"ipAddress"`
	Timestamp     int64    `json:"timestamp"`
	Encoding      string   `json:"encoding"`
	Collector     string   `json:"collector"`
	UserAgent     string   `json:"userAgent,omitempty"`
	RefererURI    string   `json:"refererUri,omitempty"`
	Path          string   `json:"path,omitempty"`
	QueryString   string   `json:"querystring,omitempty"`
	Body          string   `json:"body,omitempty"`
	Headers       []string `json:"headers,omitempty"`
	ContentType   string   `json:"contentType,omitEmpty"`
	Hostname      string   `json:"hostname,omitEmpty"`
	NetworkUserID string   `json:"networkUserId,omitEmpty"`
}

// TrackerPayload defines the structure of data posted from Snowplow trackers
type TrackerPayload struct {
	Schema string  `json:"schema"`
	Data   []Event `json:"data"`
}
