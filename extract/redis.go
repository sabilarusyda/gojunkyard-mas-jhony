package extract

import (
	"encoding/json"
	"reflect"
	"strconv"
	"time"
)

var layout = "2006-01-02T15:04:05Z"

// HMGet ..
func HMGet(result interface{}, v reflect.Value, fields []string, data []string) error {
	var a = make(map[string]interface{})

	for i := 0; i < v.NumField(); i++ {
		types := v.Field(i).Type().String()

		if data[i] == "" {
			a[fields[i]] = nil
			continue
		}
		switch types {
		case "int":
			num, err := strconv.ParseInt(data[i], 10, 10)
			if err != nil {
				return err
			}
			a[fields[i]] = num

		case "*int":
			num, err := strconv.ParseInt(data[i], 10, 10)
			if err != nil {
				return err
			}
			a[fields[i]] = &num

		case "bool":
			num, err := strconv.ParseBool(data[i])
			if err != nil {
				return err
			}
			a[fields[i]] = num

		case "*bool":
			num, err := strconv.ParseBool(data[i])
			if err != nil {
				return err
			}
			a[fields[i]] = &num

		case "float64":
			num, err := strconv.ParseFloat(data[i], 64)
			if err != nil {
				return err
			}
			a[fields[i]] = num

		case "*float64":
			num, err := strconv.ParseFloat(data[i], 64)
			if err != nil {
				return err
			}
			a[fields[i]] = &num

		case "int64":
			num, err := strconv.ParseInt(data[i], 10, 64)
			if err != nil {
				return err
			}
			a[fields[i]] = num

		case "*int64":
			num, err := strconv.ParseInt(data[i], 10, 64)
			if err != nil {
				return err
			}
			a[fields[i]] = &num

		case "string":
			a[fields[i]] = data[i]

		case "*string":
			a[fields[i]] = &data[i]

		case "time.Time":
			times, _ := time.Parse(layout, data[i])
			a[fields[i]] = times

		case "*time.Time":
			times, _ := time.Parse(layout, data[i])
			a[fields[i]] = &times

		default:
			a[fields[i]] = nil
		}
	}

	b, err := json.Marshal(a)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &result)
	if err != nil {
		return err
	}
	return nil
}
