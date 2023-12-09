package l7origin

// CreateOpts represents requests options to create a origin.
type CreateOpts struct {
	L7ResourceID int    `json:"l7ResourceId"`
	IP           string `json:"ip"`
	Weight       int    `json:"weight"`
	Mode         string `json:"Mode"`
}

// DeleteOpts represents requests options to delete a origin.
type DeleteOpts struct {
	// L7ResourceID is the identifier of the resource.
	L7ResourceID int `json:"l7ResourceId,omitempty"`
	ID           int `json:"id,omitempty"`
}

// ListOpts represents requests options to list origins.
type ListOpts struct {
	L7ResourceID int    `url:"l7ResourceId,omitempty"`
	Sort         string `url:"sort,omitempty"`
	Page         int    `url:"page,omitempty"`
	Limit        int    `url:"limit,omitempty"`
}
