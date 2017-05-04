package main

import "time"

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

// Event represents an individual atomic event. This is a mixture of what is
// needed to decode JSON messages from clients as well as what columns in the
// database these are serialized to
type Event struct {

	// Application
	AppID    string `json:"aid,omitempty" bson:"app_id"`
	Platform string `json:"p,omitempty" bson:"platform"`

	// Date/Time
	ETLTimestamp       time.Time `json:"ETLTimestamp,string" bson:"etl_tstamp"`
	CollectorTimestamp time.Time `json:"CollectorTimestamp,string" bson:"collector_tstamp"`
	DeviceTimestamp    time.Time `json:"dvce_tstamp,string,omitempty" bson:"dvce_tstamp"`
	TmpDeviceTimestamp string    `json:"dtm,omitempty" bson:"-"`

	// Event
	Event         string `json:"e,omitempty" bson:"event"`
	EventID       string `json:"eid,omitempty" bson:"event_id,omitempty"`
	TransactionID string `json:"tid,omitempty" bson:"txn_id,omitempty"`

	// Namespaceing and versioning
	TrackerName      string `json:"tna,omitempty" bson:"name_tracker,omitempty"`
	TrackerVersion   string `json:"tv,omitempty" bson:"v_tracker"`
	CollectorVersion string `bson:"v_collector"`
	ETLVersion       string `bson:"v_etl"`

	// User and visit
	UserID           string `json:"uid,omitempty" bson:"user_id,omitempty"`
	UserIPAddress    string `json:"ip,omitempty" bson:"user_ipaddress"`
	UserFingerprint  string `bson:"user_fingerprint,omitempty"`
	DomainUserID     string `json:"duid,omitempty" bson:"domain_userid,omitempty"`
	DomainSessionIDX int16  `bson:"domain_sessionidx,omitempty"`
	NetworkUserID    string `json:"tnuid,omitempty" bson:"network_userid,omitempty"`

	// Location
	GeoCountry    string  `bson:"geo_country,omitempty"`
	GeoRegion     string  `bson:"geo_region,omitempty"`
	GeoCity       string  `bson:"geo_city,omitempty"`
	GeoZipcode    string  `bson:"geo_zipcode,omitempty"`
	GeoLatitude   float32 `bson:"geo_latitude,omitempty"`
	GeoLongitude  float32 `bson:"geo_longitude,omitempty"`
	GeoRegionName string  `bson:"geo_region_name,omitempty"`

	// IP Lookups
	IPISP          string `bson:"ip_isp,omitempty"`
	IPOrganization string `bson:"ip_organization,omitempty"`
	IPDomain       string `bson:"ip_domain,omitempty"`
	IPNetspeed     string `bson:"ip_netspeed,omitempty"`

	// Page
	PageURL      string `json:"url,omitempty" bson:"page_url,omitempty"`
	PageTitle    string `json:"page,omitempty" bson:"page_title,omitempty"`
	PageReferrer string `json:"refr,omitempty" bson:"page_referrer,omitempty"`

	// Page URL Components
	PageURLScheme   string `bson:"page_urlscheme,omitempty"`
	PageURLHost     string `bson:"page_urlhost,omitempty"`
	PageURLPort     int32  `bson:"page_urlport,omitempty"`
	PageURLPath     string `bson:"page_urlpath,omitempty"`
	PageURLQuery    string `bson:"page_urlquery,omitempty"`
	PageURLFragment string `bson:"page_urlfragment,omitempty"`

	// Referrer URL Components
	ReferrerURLScheme   string `bson:"refr_urlscheme,omitempty"`
	ReferrerURLHost     string `bson:"refr_urlhost,omitempty"`
	ReferrerURLPort     int32  `bson:"refr_urlport,omitempty"`
	ReferrerURLPath     string `bson:"refr_urlpath,omitempty"`
	ReferrerURLQuery    string `bson:"refr_urlquery,omitempty"`
	ReferrerURLFragment string `bson:"refr_urlfragment,omitempty"`

	// Referrer Details
	ReferrerMedium string `bson:"refr_medium,omitempty"`
	ReferrerSource string `bson:"refr_source,omitempty"`
	ReferrerTerm   string `bson:"refr_term,omitempty"`

	// Marketing
	MarketingMedium   string `bson:"mkt_medium,omitempty"`
	MarketingSource   string `bson:"mkt_source,omitempty"`
	MarketingTerm     string `bson:"mkt_term,omitempty"`
	MarketingContent  string `bson:"mkt_content,omitempty"`
	MarketingCampaign string `bson:"mkt_campaign,omitempty"`

	// Custom Contexts
	Contexts        map[string]interface{} `json:"co,string,omitempty" bson:"contexts,omitempty"`
	ContextsEncoded string                 `json:"cx,omitempty" bson:"contexts_enc,omitempty"`

	// Custom Structured Event
	StructuredEventCategory string `json:"se_ca,omitempty" bson:"se_category,omitempty"`
	StructuredEventAction   string `json:"se_ac,omitempty" bson:"se_action,omitempty"`
	StructuredEventLabel    string `json:"se_la,omitempty" bson:"se_label,omitempty"`
	StructuredEventProperty string `json:"se_pr,omitempty" bson:"se_property,omitempty"`
	StructuredEventValue    string `json:"se_va,omitempty" bson:"se_value,omitempty"`

	// Unstructured Event
	UnstructuredEvent map[string]interface{} `json:"ue_pr,string,omitempty" bson:"unstruct_event,omitempty"`

	UnstructuredEventEncoded string `json:"ue_px,omitempty,omitempty" bson:"unstruct_event_enc,omitempty"`

	// Ecommerce
	TransactionOrderID      string `json:"tr_id,omitempty" bson:"tr_orderid,omitempty"`
	TransactionAffiliation  string `json:"tr_af,omitempty" bson:"tr_affiliation,omitempty"`
	TransactionTotal        string `json:"tr_tt,omitempty" bson:"tr_total,omitempty"`
	TransactionTax          string `json:"tr_tx,omitempty" bson:"tr_tax,omitempty"`
	TransactionShipping     string `json:"tr_sh,omitempty" bson:"tr_shipping,omitempty"`
	TransactionCity         string `json:"tr_ci,omitempty" bson:"tr_city,omitempty"`
	TransactionState        string `json:"tr_st,omitempty" bson:"tr_state,omitempty"`
	TransactionCountry      string `json:"tr_co,omitempty" bson:"tr_country,omitempty"`
	TransactionItemID       string `json:"ti_id,omitempty" bson:"ti_orderid,omitempty"`
	TransactionItemSKU      string `json:"ti_sk,omitempty" bson:"ti_sku,omitempty"`
	TransactionItemName     string `json:"ti_nm,omitempty" bson:"ti_name,omitempty"`
	TransactionItemCategory string `json:"ti_ca,omitempty" bson:"ti_category,omitempty"`
	TransactionItemPrice    string `json:"ti_pr,omitempty" bson:"ti_price,omitempty"`
	TransactionItemQuantity string `json:"ti_qu,omitempty" bson:"ti_quantity,omitempty"`

	// Page Ping
	PPXOffsetMin int32 `json:"pp_mix,omitempty" bson:"pp_xoffset_min,omitempty"`
	PPXOffsetMax int32 `json:"pp_max,omitempty" bson:"pp_xoffset_max,omitempty"`
	PPYOffsetMin int32 `json:"pp_miy,omitempty" bson:"pp_yoffset_min,omitempty"`
	PPYOffsetMax int32 `json:"pp_may,omitempty"bson:"pp_yoffset_max,omitempty"`

	// Temporary fields to handle string conversion from GET -> JSON requests
	TmpPPXOffsetMin int32 `json:"string_pp_mix,string,omitempty" bson:"-"`
	TmpPPXOffsetMax int32 `json:"string_pp_max,string,omitempty" bson:"-"`
	TmpPPYOffsetMin int32 `json:"string_pp_miy,string,omitempty" bson:"-"`
	TmpPPYOffsetMax int32 `json:"string_pp_may,string,omitempty" bson:"-"`

	// User Agent
	UserAgent string `json:"ua,omitempty" bson:"useragent"`

	// Browser
	BrName           string `bson:"br_name,omitempty"`
	BrFamily         string `bson:"br_family,omitempty"`
	BrVersion        string `bson:"br_version,omitempty"`
	BrType           string `bson:"br_type,omitempty"`
	BrRenderer       string `bson:"br_renderengine,omitempty"`
	BrLangauge       string `json:"lang,omitempty" bson:"br_lang,omitempty"`
	BrFeatPDF        bool   `json:"f_pdf,omitempty"  bson:"br_features_pdf,omitempty"`
	BrFeatFl         bool   `json:"f_fla,omitempty" bson:"br_features_flash,omitempty"`
	BrFeatJava       bool   `json:"f_java,omitempty" bson:"br_features_java,omitempty"`
	BrFeatDir        bool   `json:"f_dir,omitempty" bson:"br_features_director,omitempty"`
	BrFeatQT         bool   `json:"f_qt,omitempty" bson:"br_features_quicktime,omitempty"`
	BrFeatRealPlayer bool   `json:"f_realp,omitempty" bson:"br_features_realplayer,omitempty"`
	BrFeatWinMedia   bool   `json:"f_wma,omitempty" bson:"br_features_windowsmedia,omitempty"`
	BrFeatGears      bool   `json:"f_gears,omitempty" bson:"br_features_gears,omitempty"`
	BrCookies        bool   `json:"cookie,omitempty" bson:"br_cookies,omitempty"`

	// Temporary fields to handle string conversion from GET -> JSON requests
	TmpBrFeatPDF        bool `json:"string_f_pdf,string,omitempty" bson:"-"`
	TmpBrFeatFl         bool `json:"string_f_fla,string,omitempty" bson:"-"`
	TmpBrFeatJava       bool `json:"string_f_java,string,omitempty" bson:"-"`
	TmpBrFeatDir        bool `json:"string_f_dir,string,omitempty" bson:"-"`
	TmpBrFeatQT         bool `json:"string_f_qt,string,omitempty" bson:"-"`
	TmpBrFeatRealPlayer bool `json:"string_f_realp,string,omitempty" bson:"-"`
	TmpBrFeatWinMedia   bool `json:"string_f_wma,string,omitempty" bson:"-"`
	TmpBrFeatGears      bool `json:"string_f_gears,string,omitempty" bson:"-"`
	TmpBrCookies        bool `json:"string_cookie,string,omitempty" bson:"-"`

	BrFeatSilver bool   `bson:"br_features_silverlight,omitempty"`
	BrColorDepth string `bson:"br_colordepth,omitempty"`
	BrViewWidth  int32  `bson:"br_viewwidth,omitempty"`
	BrViewHeight int32  `bson:"br_viewheight,omitempty"`

	// Operating System
	OSName         string `bson:"os_name,omitempty"`
	OSFamily       string `bson:"os_family,omitempty"`
	OSManufacturer string `bson:"os_manufacturer,omitempty"`
	OSTimeZone     string `bson:"os_timezone,omitempty,omitempty"`

	// Device/Hardware
	DeviceType         string `bson:"dvce_type,omitempty"`
	DeviceIsMobile     bool   `bson:"dvce_ismobile,omitempty"`
	DeviceScreenWidth  int32  `bson:"dvce_screenwidth,omitempty"`
	DeviceScreenHeight int32  `bson:"dvce_screenheight,omitempty"`

	// Document
	DocCharset string `bson:"doc_charset,omitempty"`
	DocWidth   int32  `bson:"doc_width,omitempty"`
	DocHeight  int32  `bson:"doc_height,omitempty"`

	// Currency
	TransactionCurrency     string `json:"tr_cu,omitempty" bson:"tr_currency,omitempty"`
	TransactionTotalBase    string `bson:"tr_total_base,omitempty"`
	TransactionTaxBase      string `bson:"tr_tax_base,omitempty"`
	TransactionShippingBase string `bson:"tr_shipping_base,omitempty"`
	TransactionItemCurrency string `json:"ti_cu,omitempty" bson:"ti_currency,omitempty"`
	TransactionItemBase     string `bson:"ti_price_base,omitempty"`
	BaseCurrency            string `bson:"base_currency,omitempty"`

	// Geolocation
	GeoTimeZone string `bson:"geo_timezone,omitempty"`

	// Click ID
	MarketClickID string `bson:"mkt_clickid,omitempty"`
	MarketNetwork string `bson:"mkt_network,omitempty"`

	// ETL Tags
	ETLTags string `bson:"etl_tags,omitempty"`

	// Time event was sent
	DeviceSentTimestamp string `bson:"dvce_sent_tstamp,omitempty"`

	// Referer
	ReferrerDomainUserID    string `bson:"refr_domain_userid,omitempty"`
	ReferrerDeviceTimestamp string `bson:"refr_dvce_timestamp,omitempty"`

	// Contexts
	DerivedContexts string `bson:"derived_contexts,omitempty"`

	// SessionID
	DomainSessionID string `bson:"domain_sessionid,omitempty"`

	// Derived Timestamp
	DerivedTimestamp string `bson:"derived_tstamp,omitempty"`

	// JSON ONLY PROPERTIES ON INCOMING EVENT

	//Namespace  string `json:"tna,omitempty"`
	Timezone   string `json:"tz,omitempty" bson:"timezone,omitempty"`
	Resolution string `json:"res,omitempty" bson:"resolution,omitempty"`
	// Language   string `json:"lang,omitempty"`
	ColorDepth string `json:"cd,omitempty" bson:"colordepth,omitempty"`
	Viewport   string `json:"vp,omitempty" bson:"viewport,omitempty"`
}
