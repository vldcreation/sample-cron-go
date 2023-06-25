package utils

import (
	"encoding/json"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func ParseUrl(_url string, sc interface{}) error {
	urlParsed, err := url.Parse(_url)

	if err != nil {
		log.Printf("error parse url: %v\n", err)
		return err
	}

	switch sc.(type) {
	case *url.URL:
		*sc.(*url.URL) = *urlParsed
		break
	case *url.Values:
		*sc.(*url.Values) = urlParsed.Query()
		break
	default:
		queryParams := urlParsed.Query()
		attributes := make(map[string]string)
		for key, values := range queryParams {
			attributes[key] = strings.Join(values, ",")
		}
		jsonStr, _ := json.Marshal(attributes)
		json.Unmarshal(jsonStr, &sc)
	}

	return nil
}

func ParseDate(strDate string) (time.Time, error) {
	// accepted layouts
	// since golang 1.15, RFC3339 is supported
	// https://golang.org/doc/go1.15#time
	// so make sure to use golang 1.15 or above
	var date time.Time
	var err error

	layouts := []string{
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05Z",
		"20060102T150405Z",
		"2006-01-02",
	}

	for _, layout := range layouts {
		date, err = time.Parse(layout, strDate)
		if err == nil {
			break
		}
	}

	return date, err
}

func ParseUnix(unixStr string) (time.Time, error) {
	unixConv, err := strconv.Atoi(unixStr)
	if err != nil {
		return time.Time{}, err
	}

	expiresInt, _ := time.Parse("2006-01-02 15:04:05", time.Unix(0, 0).String())
	return expiresInt.Add(time.Second * time.Duration(unixConv)), nil
}

// @strDate: date string to be parsed
// @location: location to be used
//
//	(e.g. time.UTC, time.Local, time.LoadLocation("Asia/Jakarta"))
//	by default, time.Local is used
//
// @return: parsed date with location
func ParseDateWithLocation(strDate string, location *time.Location) (time.Time, error) {
	// accepted layouts
	// since golang 1.15, RFC3339 is supported
	// https://golang.org/doc/go1.15#time
	// so make sure to use golang 1.15 or above
	var date time.Time
	var err error

	layouts := []string{
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05Z",
		"20060102T150405Z",
		"2006-01-02",
	}

	for _, layout := range layouts {
		date, err = time.ParseInLocation(layout, strDate, location)
		if err == nil {
			break
		}
	}

	return date, err
}
