package speedtest

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
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
func getRequestSize(req *http.Request) int {
	bytesSize := 0
	for k, v := range req.Header {
		bytesSize += len(k) + len(v[0]) // Assuming single value per header
	}
	bytesSize += len(req.Method)
	bytesSize += len(req.RequestURI)
	bytesSize += len(req.UserAgent())
	bytesSize += len(req.RemoteAddr)
	bytesSize += len(req.Referer())
	bytesSize += int(req.ContentLength)
	bytesSize += len(req.Proto)
	return bytesSize
}
func getResponceSize(req *http.Response) int {
	bytesSize := 0
	for k, v := range req.Header {
		bytesSize += len(k) + len(v[0]) // Assuming single value per header
	}
	bytesSize += int(req.ContentLength)
	bytesSize += len(req.Status)
	bytesSize += len(req.Proto)
	return bytesSize
}

// \{ // definitly need
// 		(
//			"([0-9a-zA-Z]*)"= // key
//				("?)([0-9a-zA-Z().,"]*)("?)(;?) // value (we have problem with the "" signs)
//		)* // repeaded 0 or more
// \} // definitly need
