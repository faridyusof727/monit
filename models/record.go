package models

import (
	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	MonitorID     uint   `json:"monitor_id,omitempty"`
	Status        string `json:"status,omitempty"`
	Code          string `json:"code,omitempty"`
	ResponseTime  int64  `json:"response_time,omitempty"`
	SSLStatus     string `json:"ssl_status"`
	SSLIssuer     string `json:"ssl_issuer"`
	SSLCommonName string `json:"ssl_common_name"`
	SSLExpiry     string `json:"ssl_expiry"`
}
