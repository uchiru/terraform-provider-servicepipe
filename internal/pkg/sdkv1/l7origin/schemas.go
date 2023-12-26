package l7origin

type Data struct {
	Data Result `json:"data"`
}

type Result struct {
	Result Item `json:"result"`
}

// Item represents an unmarshalled domain body from API response.
type Item struct {
	L7ResourceID int64  `json:"l7ResourceId"`
	ID           int64  `json:"id"`
	Weight       int64  `json:"weight"`
	Mode         string `json:"mode"`
	IP           string `json:"ip"`
	CreatedAt    int64  `json:"createdAt"`
	ModifiedAt   int64  `json:"modifiedAt"`
}

type DataItems struct {
	DataItems ResultItems `json:"data"`
}

type ResultItems struct {
	ResultItems Items `json:"result"`
}

type Items struct {
	Items []Item `json:"items"`
	Info  Info   `json:"info"`
}

type Info struct {
	TotalCount int64 `json:"totalCount"`
	Limit      int64 `json:"limit"`
	Page       int64 `json:"page"`
}

type DataDelete struct {
	Data ResultDelete `json:"data"`
}

type ResultDelete struct {
	Result string `json:"result"`
}
