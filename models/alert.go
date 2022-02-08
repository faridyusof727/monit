package models

import (
	"gorm.io/gorm"
)

type Alert struct {
	gorm.Model
	MonitorID uint   `json:"monitor_id,omitempty"`
	Type      string `json:"status,omitempty" gorm:"default:telegram"`
	Key       string `json:"code,omitempty"`
}
