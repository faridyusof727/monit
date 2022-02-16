package models

import (
	"database/sql"
	"errors"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
)

const (
	TypeHttps = "https"
	TypeHttp  = "http"

	ReqPost = "post"
	ReqGet  = "get"

	DefaultInterval = 5
	DefaultTimeout  = 30
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
	Records          []Record       `json:"records"`
}

func (a Monitor) Validate() error {

	u, err := url.Parse(a.Url)
	if err == nil {
		a.Url = u.Host + u.Path + u.RawQuery + u.Fragment
	}

	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required),
		validation.Field(&a.Url, validation.Required),

		validation.Field(&a.Type, validation.Required, validation.In(TypeHttp, TypeHttps)),
		validation.Field(&a.RequestMethod, validation.Required, validation.In(ReqGet, ReqPost)),
		validation.Field(&a.IntervalInMinute, validation.By(a.isDefault)),
		validation.Field(&a.TimeoutInSecond, validation.By(a.isDefault)),
	)
}

func (a Monitor) isDefault(value interface{}) error {

	if value == DefaultInterval || value == 0 {
		return nil
	}

	if value == DefaultTimeout || value == 0 {
		return nil
	}

	return errors.New("field is not default")
}
