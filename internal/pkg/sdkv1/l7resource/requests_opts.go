package l7resource

// CreateOpts represents requests options to create a domain.
type CreateOpts struct {
	// L7ResourceName represents valid domain name (without www.).
	L7ResourceName string `json:"l7ResourceName,omitempty"`

	// OriginData represents valid IP v4 address.
	OriginData string `json:"originData,omitempty"`
	Wwwredir   int    `json:"wwwredir,omitempty"`
}

// CreateOpts represents requests options to create a domain.
type DeleteOpts struct {
	// L7ResourceID is the identifier of the resource.
	L7ResourceID int `json:"l7ResourceId,omitempty"`
}
