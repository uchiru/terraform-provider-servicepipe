package l7origin

// CreateOpts represents requests options to create a origin.
type CreateOpts struct {
	L7ResourceID int64  `json:"l7ResourceId"`
	IP           string `json:"ip"`
	Weight       int64  `json:"weight"`
	Mode         string `json:"Mode"`
}

// DeleteOpts represents requests options to delete a origin.
type DeleteOpts struct {
	// L7ResourceID is the identifier of the resource.
	L7ResourceID int64 `json:"l7ResourceId"`
	ID           int64 `json:"id"`
}

// ListOpts represents requests options to list origins.
type ListOpts struct {
	L7ResourceID int64  `url:"l7ResourceId,omitempty"`
	Sort         string `url:"sort,omitempty"`
	Page         int64  `url:"page,omitempty"`
	Limit        int64  `url:"limit,omitempty"`
}
