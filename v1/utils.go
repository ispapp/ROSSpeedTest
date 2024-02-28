package speedtest

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ROString(t interface{}) string {
	if t == nil {
		return ""
	}
	val := reflect.ValueOf(t)
	typ := reflect.TypeOf(t)
	switch val.Kind() {
	case reflect.Struct:
		result := "{"
		for i := 0; i < val.NumField(); i++ {
			key := typ.Field(i).Name
			valueType := typ.Field(i).Type.Kind()
			field := val.Field(i)

			if i > 0 {
				result += ";"
			}

			switch {
			case valueType == reflect.String:
				result += fmt.Sprintf("\"%s\"=\"%s\"", key, field.String())
			case valueType == reflect.Bool:
				result += fmt.Sprintf("\"%s\"=%t", key, field.Bool())
			case valueType == reflect.Array || valueType == reflect.Slice:
				result += fmt.Sprintf("\"%s\"=(%s)", key, arrayToString(field))
			case field.CanFloat():
				result += fmt.Sprintf("\"%s\"=%f", key, field.Float())
			case field.CanInt():
				result += fmt.Sprintf("\"%s\"=%d", key, field.Int())
			case valueType == reflect.Struct:
				result += fmt.Sprintf("\"%s\"=%s", key, ROString(field.Interface()))
			}
		}
		result += "}"
		return result
	default:
		return "{}"
	}
}
func arrayToString(arr reflect.Value) string {
	result := ""
	for i := 0; i < arr.Len(); i++ {
		if i > 0 {
			result += ","
		}
		elem := arr.Index(i)
		switch {
		case elem.Kind() == reflect.String:
			result += fmt.Sprintf("\"%s\"", elem.String())
		case elem.Kind() == reflect.Bool:
			result += fmt.Sprintf("%t", elem.Bool())
		case elem.Kind() == reflect.Struct:
			result += ROString(elem.Interface())
		case elem.CanFloat():
			result += fmt.Sprintf("%f", elem.Float())
		case elem.CanInt():
			result += fmt.Sprintf("%d", elem.Int())
		}
	}
	return result
}

func IsRouterOSArray(s string) bool {
	// a surface check that avoid us some lines of code
	re := regexp.MustCompile(`^\{("([0-9a-zA-Z]*)"(=?)(("?{?\(?)([0-9a-zA-Z."]*(,?))*("?}?\)?,?;?))(;?))*\}$`)
	return re.MatchString(s)
}
func getPacketSize(req *http.Request) int {
	bytesSize := 0
	bytesSize += len("Host:") + len(req.Host) + 3 // space (0x32) + \r\n
	for k, v := range req.Header {
		bytesSize += len(k) + 3
		for _, _v := range v {
			bytesSize += len(_v) + 1 // space (0x32)
		}
	}
	bytesSize += len(req.Method) + 1
	bytesSize += len(req.RequestURI) + 1
	bytesSize += 8  // Protocol Version: HTTP/1.1 (8 bytes)
	bytesSize += 32 // Transmission Control Protocol (headers)
	bytesSize += 20 // Internet Protocol Version 4 (headers)
	bytesSize += 14 // Ethernet II (headers)
	return bytesSize
}

func Avg(numbers []int64) float64 {
	if len(numbers) == 0 {
		return 0.0
	}
	var sum float64
	for _, num := range numbers {
		sum += float64(num)
	}
	average := sum / float64(len(numbers))
	return average
}
func DecodeROString(s string) (interface{}, error) {
	var t interface{}
	if !IsRouterOSArray(s) {
		return nil, fmt.Errorf("decodeROString: the second argument must be a non-nil pointer to a struct")
	}
	val := reflect.ValueOf(t)
	// typ := reflect.TypeOf(t)

	if val.Kind() != reflect.Ptr || val.IsNil() {
		return nil, fmt.Errorf("decodeROString: the second argument must be a non-nil pointer to a struct")
	}

	if val.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("decodeROString: the second argument must be a pointer to a struct")
	}

	val = val.Elem()
	// typ = typ.Elem()

	// Removing curly braces from the string
	s = strings.Trim(s, "{}")

	// Splitting the string into key-value pairs
	pairs := strings.Split(s, ";")

	for _, pair := range pairs {
		// Splitting each key-value pair into key and value
		kv := strings.Split(pair, "=")

		if len(kv) != 2 {
			return nil, fmt.Errorf("decodeROString: invalid key-value pair: %s", pair)
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		// Finding the field in the struct
		field := val.FieldByName(key)
		if field.IsNil() {
			return nil, fmt.Errorf("decodeROString: field not found in struct: %s", key)
		}

		// Setting the value of the field based on its type
		switch field.Kind() {
		case reflect.String:
			field.SetString(value)
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(value)
			if err != nil {
				return nil, fmt.Errorf("decodeROString: error parsing boolean value for field %s: %v", key, err)
			}
			field.SetBool(boolValue)
		case reflect.Float64:
			floatValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("decodeROString: error parsing float value for field %s: %v", key, err)
			}
			field.SetFloat(floatValue)
		case reflect.Int:
			intValue, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("decodeROString: error parsing integer value for field %s: %v", key, err)
			}
			field.SetInt(int64(intValue))
		default:
			return nil, fmt.Errorf("decodeROString: unsupported field type: %s", field.Kind())
		}
	}

	return t, nil
}
func ConvertToMilliseconds(timeStr string) (int64, error) {
	// Parse the time string
	var MytimeStr []string
	fragments := strings.Split(timeStr, ":")
	if len(fragments) != 3 {
		return 0, fmt.Errorf("invalid format")
	}
	MytimeStr = append(MytimeStr, fmt.Sprintf("%sh", fragments[0]))
	MytimeStr = append(MytimeStr, fmt.Sprintf("%sm", fragments[1]))
	MytimeStr = append(MytimeStr, fmt.Sprintf("%ss", fragments[2]))
	formated := strings.Join(MytimeStr, "")
	parsedTime, err := time.ParseDuration(formated)
	if err != nil {
		return 0, err
	}
	// Convert to milliseconds
	milliseconds := parsedTime.Milliseconds()

	return milliseconds, nil
}

// func getResponceSize(req *http.Response) int {
// 	bytesSize := 0
// 	for k, v := range req.Header {
// 		bytesSize += len(k) + len(v[0]) // Assuming single value per header
// 	}
// 	bytesSize += int(req.ContentLength)
// 	bytesSize += len(req.Status)
// 	bytesSize += len(req.Proto)
// 	return bytesSize
// }

// \{ // definitly need
// 		(
//			"([0-9a-zA-Z]*)"= // key
//				("?)([0-9a-zA-Z().,"]*)("?)(;?) // value (we have problem with the "" signs)
//		)* // repeaded 0 or more
// \} // definitly need
