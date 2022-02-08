package models

import (
	"database/sql"
	"github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
)

const (
	TypeHttps = "https"
	TypeHttp  = "http"

	ReqPost = "post"
	ReqGet  = "get"
)

type Monitor struct {
	gorm.Model
	Owner            string         `json:"owner,omitempty"`
	Type             string         `json:"type,omitempty"`
	Name             string         `json:"name,omitempty"`
	Url              string         `json:"url,omitempty"`
	IntervalInMinute int            `json:"interval_in_minute,omitempty" gorm:"default:5"`
	TimeoutInSecond  int            `json:"timeout_in_second,omitempty" gorm:"default:30"`
	RequestMethod    string         `json:"request_method,omitempty"`
	RequestHeader    sql.NullString `json:"request_header,omitempty"`
	SuccessCodes     sql.NullString `json:"success_codes,omitempty"`
	FailCodes        sql.NullString `json:"fail_codes,omitempty"`
}

func (a Monitor) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required),
		validation.Field(&a.Url, validation.Required),

		validation.Field(&a.Type, validation.Required, validation.In(TypeHttp, TypeHttps)),
		validation.Field(&a.RequestMethod, validation.Required, validation.In(ReqGet, ReqPost)),
	)
}
