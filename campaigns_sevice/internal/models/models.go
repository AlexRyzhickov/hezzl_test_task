package models

import "time"

type Campaign struct {
	Id   int
	Name string
}

type Item struct {
	Id          int       `json:"id" gorm:"primaryKey"`
	CampaignId  int       `json:"campaign_id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	CreatedAt   time.Time `json:"created_at"`
}

type Error struct {
	Error string `json:"error"`
}

type Msg struct {
	Msg string `json:"msg"`
}
