package l7origin

// CreateOpts represents requests options to create a domain.
type CreateOpts struct {
	L7ResourceID int    `json:"l7ResourceId"`
	IP           string `json:"ip"`
	Weight       int    `json:"weight"`
	Mode         string `json:"Mode"`
}

// CreateOpts represents requests options to create a domain.
type DeleteOpts struct {
	// L7ResourceID is the identifier of the resource.
	L7ResourceID int `json:"l7ResourceId,omitempty"`
	ID           int `json:"id,omitempty"`
}
