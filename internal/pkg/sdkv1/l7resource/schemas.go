package l7resource

type Data struct {
	Data Result `json:"data"`
}

type Result struct {
	Result Item `json:"result"`
}

// Item represents an unmarshalled domain body from API response.
type Item struct {
	// PartnerClientAccountId is the identifier of partner client.
	PartnerClientAccountId int `json:"partnerClientAccountId"`

	// L7ResourceID is the identifier of the resource.
	L7ResourceID        int    `json:"l7ResourceId"`
	L7ResourceName      string `json:"l7ResourceName"`
	L7ResourceIsActive  int    `json:"l7ResourceIsActive"`
	L7ProtectionDisable int    `json:"l7ProtectionDisable"`
	UseCustomSsl        int    `json:"useCustomSsl"`
	UseLetsencryptSsl   int    `json:"useLetsencryptSsl"`
	CustomSslKey        string `json:"customSslKey"`
	CustomSslCrt        string `json:"customSslCrt"`

	// CreateDate represents Unix timestamp when resource has been created.
	DeletedAt string `json:"deletedAt"`
	Sort      string `json:"sort"`
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`

	// CreateDate represents Unix timestamp when domain has been created.
	СreatedAt             int    `json:"сreatedAt"`
	Forcessl              int    `json:"forcessl"`
	ServiceHttp2          int    `json:"serviceHttp2"`
	GeoipMode             int    `json:"geoipMode"`
	GeoipList             string `json:"geoipList"`
	GlobalWhitelistActive int    `json:"globalWhitelistActive"`

	Http2https int `json:"http2https"`
	Https2http int `json:"https2http"`

	ProtectedIp   string `json:"protectedIp"`
	ModifiedAt    int    `json:"modifiedAt"`
	SslExpireDate int    `json:"SslExpireDate"`
	Wwwredir      int    `json:"wwwredir"`
	Cdn           int    `json:"cdn"`
	CdnHost       string `json:"cdnHost"`
	CdnProxyHost  string `json:"cdnProxyHost"`
}

type DataItems struct {
	DataItems ResultItems `json:"data"`
}

type ResultItems struct {
	ResultItems Items `json:"result"`
}

type Items struct {
	Items []Item `json:"items"`
}

type DataDelete struct {
	Data ResultDelete `json:"data"`
}

type ResultDelete struct {
	ResultItems string `json:"result"`
}
