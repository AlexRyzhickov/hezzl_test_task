package models

type Campaign struct {
	Id   int
	Name string
}

type Item struct {
	Id          int
	Campaign_id int    `json:"campaign_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"priority,omitempty"`
	Priority    int    `json:"order_uid,omitempty"`
	Removed     bool   `json:"removed,omitempty"`
	Created_at  string `json:"created_at,omitempty"`
}
