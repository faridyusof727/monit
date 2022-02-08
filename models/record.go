package models

import (
	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	MonitorID uint   `json:"monitor_id,omitempty"`
	Status    string `json:"status,omitempty"`
	Code      string `json:"code,omitempty"`
}
