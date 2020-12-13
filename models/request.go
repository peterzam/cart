package models

import "time"

//Request model for saving in db
type Request struct {
	URL       string        `json:"url" form:"url" query:"url"`
	Data      string        `json:"data" form:"data" query:"data"`
	CreatedAt time.Duration `gorm:"autoCreateTime:milli"`
}
