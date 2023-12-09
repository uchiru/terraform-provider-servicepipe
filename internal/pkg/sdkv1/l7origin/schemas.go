package l7origin

type Data struct {
	Data Result `json:"data"`
}

type Result struct {
	Result Item `json:"result"`
}

// Item represents an unmarshalled domain body from API response.
type Item struct {
	L7ResourceID int    `json:"l7ResourceId"`
	ID           int    `json:"id"`
	Weight       int    `json:"weight"`
	Mode         string `json:"mode"`
	IP           string `json:"ip"`
	CreatedAt    int    `json:"createdAt"`
	ModifiedAt   int    `json:"modifiedAt"`
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
	TotalCount int `json:"totalCount"`
	Limit      int `json:"limit"`
	Page       int `json:"page"`
}

type DataDelete struct {
	Data ResultDelete `json:"data"`
}

type ResultDelete struct {
	Result string `json:"result"`
}
