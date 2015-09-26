// this is used to give the ability to NULL times in the database
// source: https://github.com/jinzhu/gorm/issues/10

package nullables

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

const (
	HTMLFormTime         = "15:04"
	HTMLFormDate         = "2006-01-02"
	HTMLFormDateTime     = "2006-01-02 3:04 PM"
	NZHTMLFormDate       = "02-01-2006"
	NZHTMLFormDateTime   = "02-01-2006 3:04 PM"
	DBTime               = "15:04:05:000"
	HTMLFormDateTime24Hr = "2006-01-02 15:04"
	DBDateTime           = "2006-01-02 15:04:05:000"
	TimeWithSeconds      = "15:04:05"
)

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
// Can scan from both time.Time and NullTime interfaces
func (nt *NullTime) Scan(value interface{}) error {
	if value == nil {
		nt.Valid = false
		return nil
	}

	switch iface := value.(type) {
	case time.Time:
		nt.Time, nt.Valid = iface, true
		return nil
	case NullTime:
		nt.Time, nt.Valid = iface.Time, iface.Valid
		return nil
	default:
		return nil
	}
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

//this function is used when JSON tries to marshal this struct.
//Without it it will be marshalled into a JSON object with fields Time and Valid.
//This is far more elegant.
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return nt.Time.MarshalJSON()
	} else {
		return json.Marshal(nil)
	}
}

func (nt *NullTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" || string(b) == "" {
		nt.Valid = false
		return nil
	}
	v, err := time.Parse(time.RFC3339, string(b[1:len(b)-1]))
	if err != nil {
		return err
	}
	nt.Time = v
	nt.Valid = true
	return nil
}

func (nt NullTime) GetHTMLDateTime() string {
	if nt.Valid {
		return nt.Time.Format(HTMLFormDateTime)
	}
	return "N/A"
}

//This function is used to convert an HTML form value (returned from a create or edit, for instance) to a NullTime.
//It will first try parse it as a Time, if that does not work it will parse it as a Date, if that doesn't work it will parse it as a datetime
func NullTimeConverter(b string) reflect.Value {
	decodedTime := NullTime{}
	//first check if we have been given a time
	v, err := time.Parse(HTMLFormTime, b)
	if err == nil {
		v := v.AddDate(1, 0, 0) //Nasty hack so that the year is not 0000, which is valid in Go but not MsSQL
		decodedTime.Time = v
		decodedTime.Valid = true
		return reflect.ValueOf(decodedTime)
	}
	//now check if it was a date
	v, err = time.Parse(NZHTMLFormDate, b)
	if err == nil {
		decodedTime.Time = v
		decodedTime.Valid = true
		return reflect.ValueOf(decodedTime)
	}

	v, err = time.Parse(HTMLFormDate, b)
	if err == nil {
		decodedTime.Time = v
		decodedTime.Valid = true
		return reflect.ValueOf(decodedTime)
	}
	//now check if it was a datetime
	v, err = time.Parse(HTMLFormDateTime, b)
	if err == nil {
		decodedTime.Time = v
		decodedTime.Valid = true
		return reflect.ValueOf(decodedTime)
	}
	//now check if it was a date
	v, err = time.Parse(NZHTMLFormDateTime, b)
	if err == nil {
		decodedTime.Time = v
		decodedTime.Valid = true
		return reflect.ValueOf(decodedTime)
	}
	v, err = time.Parse(HTMLFormDateTime24Hr, b)
	if err == nil {
		decodedTime.Time = v
		decodedTime.Valid = true
		return reflect.ValueOf(decodedTime)
	}
	return reflect.ValueOf(decodedTime)
}

type NullString struct {
	String string
	Valid  bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullString) Scan(value interface{}) error {
	if value == nil {
		nt.Valid = false
		return nil
	}
	nt.String, nt.Valid = fmt.Sprintf("%s", value), true
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullString) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.String, nil
}

//this function is used when JSON tries to marshal this struct.
//Without it it will be marshalled into a JSON object with fields Time and Valid.
//This is far more elegant.
func (nt NullString) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return []byte(fmt.Sprintf("\"%v\"", nt.String)), nil
	} else {
		return json.Marshal(nil)
	}
}

func (nt *NullString) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		nt.Valid = false
		return nil
	}
	nt.String = string(b[1 : len(b)-1])
	nt.Valid = true
	return nil
}

func ConvertBool(value string) reflect.Value {
	if value == "on" {
		return reflect.ValueOf(true)
	} else if v, err := strconv.ParseBool(value); err == nil {
		return reflect.ValueOf(v)
	}

	return reflect.ValueOf(false)
}

//This function is used to convert an HTML form value (returned from a create or edit, for instance) to a NullTime.
//It will first try parse it as a Time, if that does not work it will parse it as a Date.
func NullStringConverter(s string) reflect.Value {
	decodedString := NullString{}
	//If string is empty we'll make it null
	if s == "" {
		decodedString.Valid = false
	} else {
		decodedString.String = s
		decodedString.Valid = true
	}
	return reflect.ValueOf(decodedString)
}

type NullInt64 struct {
	Int64 int64
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullInt64) Scan(value interface{}) error {
	if value == nil {
		nt.Valid = false
		return nil
	}
	nt.Int64, nt.Valid = value.(int64), true
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullInt64) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Int64, nil
}

//this function is used when JSON tries to marshal this struct.
//Without it it will be marshalled into a JSON object with fields Time and Valid.
//This is far more elegant.
func (nt NullInt64) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return []byte(fmt.Sprintf("%v", nt.Int64)), nil
	} else {
		return json.Marshal(nil)
	}
}

func (nt *NullInt64) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		nt.Valid = false
		return nil
	}
	var err error
	nt.Int64, err = strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	nt.Valid = true
	return nil
}

func NullInt64Converter(i string) reflect.Value {
	decodedInt64 := NullInt64{}
	//If string is empty we'll make it null
	if i == "" {
		decodedInt64.Valid = false
	} else {
		var err error
		decodedInt64.Int64, err = strconv.ParseInt(i, 10, 64)
		if err != nil {
			decodedInt64.Valid = false
		} else {
			decodedInt64.Valid = true
		}
	}
	return reflect.ValueOf(decodedInt64)
}

type NullFloat64 struct {
	Float64 float64
	Valid   bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullFloat64) Scan(value interface{}) error {
	if value == nil {
		nt.Valid = false
		return nil
	}
	nt.Float64, nt.Valid = value.(float64), true
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullFloat64) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Float64, nil
}

//this function is used when JSON tries to marshal this struct.
//Without it it will be marshalled into a JSON object with fields Time and Valid.
//This is far more elegant.
func (nt NullFloat64) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return []byte(fmt.Sprintf("%v", nt.Float64)), nil
	} else {
		return json.Marshal(nil)
	}
}

func (nt *NullFloat64) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		nt.Valid = false
		return nil
	}
	var err error
	nt.Float64, err = strconv.ParseFloat(string(b), 64)
	if err != nil {
		return err
	}
	nt.Valid = true
	return nil
}

func NullFloat64Converter(i string) reflect.Value {
	decodedFloat64 := NullFloat64{}
	//If string is empty we'll make it null
	if i == "" {
		decodedFloat64.Valid = false
	} else {
		var err error
		decodedFloat64.Float64, err = strconv.ParseFloat(i, 64)
		if err != nil {
			decodedFloat64.Valid = false
		} else {
			decodedFloat64.Valid = true
		}
	}
	return reflect.ValueOf(decodedFloat64)
}

type NullBool struct {
	Bool  bool
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullBool) Scan(value interface{}) error {
	if value == nil {
		nt.Valid = false
		return nil
	}
	nt.Bool, nt.Valid = value.(bool), true
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullBool) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Bool, nil
}

//this function is used when JSON tries to marshal this struct.
//Without it it will be marshalled into a JSON object with fields Time and Valid.
//This is far more elegant.
func (nt NullBool) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		if nt.Bool {
			return []byte("true"), nil
		}
		return []byte("false"), nil
	} else {
		return json.Marshal(nil)
	}
}

func (nt *NullBool) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		nt.Valid = false
		return nil
	}
	if string(b) == "true" {
		nt.Valid = true
		nt.Bool = true
		return nil
	}
	if string(b) == "false" {
		nt.Valid = true
		nt.Bool = false
	}
	return errors.New(string(b) + " is not a valid JSON bool")
}

func NullBoolConverter(i string) reflect.Value {
	decodedBool := NullBool{}
	//If string is empty we'll make it null
	if i == "" {
		decodedBool.Valid = false
	} else {
		var err error
		decodedBool.Bool, err = strconv.ParseBool(i)
		if err != nil {
			decodedBool.Valid = false
		} else {
			decodedBool.Valid = true
		}
	}
	return reflect.ValueOf(decodedBool)
}
