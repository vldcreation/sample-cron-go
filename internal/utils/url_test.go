package utils_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/vldcreation/sample-cron-go/internal/utils"
)

func TestParseDate(t *testing.T) {
	test := []struct {
		name     string
		strDate  string
		expected time.Time
	}{
		{
			name:     "should return error when date is invalid",
			strDate:  "2020-01-01",
			expected: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "should return date when date is valid",
			strDate:  "2020-01-01T00:00:00Z",
			expected: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "should return date when date is valid",
			strDate:  "20230618T151930Z",
			expected: time.Date(2023, 6, 18, 15, 19, 30, 0, time.UTC),
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			_date, err := utils.ParseDate(tt.strDate)
			if err != nil {
				t.Errorf("expected nil, got %v", err)
			}

			if tt.expected != _date {
				t.Errorf("expected %v, got %v", tt.expected, tt.expected)
			}
		})
	}
}

func TestParseDateWithLocation(t *testing.T) {
	parseLocation := func(name string) *time.Location {
		loc, err := time.LoadLocation(name)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		return loc
	}

	// containsZ := func(str string) bool {
	// 	for _, c := range str {
	// 		if c == 'Z' {
	// 			return true
	// 		}
	// 	}
	// 	return false
	// }

	test := []struct {
		name     string
		strDate  string
		location *time.Location
		expected time.Time
	}{
		{
			name:     "should return error when date is invalid",
			strDate:  "2020-01-01",
			location: time.UTC,
			expected: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "should return date when date is valid",
			strDate:  "2020-01-01T00:00:00Z",
			location: parseLocation("UTC"),
			expected: time.Date(2020, 1, 1, 0, 0, 0, 0, parseLocation("UTC")),
		},
		{
			name:     "should return date when date is valid",
			strDate:  "20230618T151930Z",
			location: parseLocation("Asia/Jakarta"),
			expected: time.Date(2023, 6, 18, 15, 19, 30, 0, parseLocation("Asia/Jakarta")),
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			_date, err := utils.ParseDateWithLocation(tt.strDate, tt.location)
			if err != nil {
				t.Errorf("expected nil, got %v", err)
			}

			// if contains Z, then the location is UTC
			// if containsZ(tt.strDate) && _date.Location().String() != time.UTC.String() {
			// 	t.Errorf("error in dates %s expected location %v, got %v", tt.strDate, time.UTC, _date.Location())
			// }

			if tt.expected.Location().String() != _date.Location().String() {
				t.Errorf("error in date %s expected location %v, got %v", tt.strDate, tt.expected.Location(), _date.Location())
			}
		})
	}
}

func TestParseUrl(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected interface{}
	}{
		{
			name:     "should return expected when url is valid",
			url:      "http://127.0.0.1:9000/dci-auth/sample.jpeg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=xx&X-Amz-Date=20230618T151930Z&X-Amz-Expires=20&X-Amz-SignedHeaders=host&X-Amz-Signature=xx",
			expected: defaultTestParseUrlData(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.ParseUrl(tt.url, tt.expected)
			if err != nil {
				if err.Error() != tt.expected.(string) {
					t.Errorf("expected %v, got %v", tt.expected, err)
				}
			} else {
				if !reflect.DeepEqual(tt.expected, defaultTestParseUrlData()) {
					t.Errorf("expected %v, got %v", tt.expected, defaultTestParseUrlData())
				}
			}
		})
	}
}

func defaultTestParseUrlData() interface{} {
	type PresignUrlInfoS3 struct {
		X_AMZ_ALGORITHM     string `json:"X-Amz-Algorithm"`
		X_AMZ_CREDENTIAL    string `json:"X-Amz-Credential"`
		X_AMZ_DATE          string `json:"X-Amz-Date"`
		X_AMZ_EXPIRES       string `json:"X-Amz-Expires"`
		X_AMZ_SIGNEDHEADERS string `json:"X-Amz-SignedHeaders"`
		X_AMZ_SIGNATURE     string `json:"X-Amz-Signature"`
	}

	return PresignUrlInfoS3{
		X_AMZ_ALGORITHM:     "AWS4-HMAC-SHA256",
		X_AMZ_CREDENTIAL:    "xx",
		X_AMZ_DATE:          "20230618T151930Z",
		X_AMZ_EXPIRES:       "20",
		X_AMZ_SIGNEDHEADERS: "host",
		X_AMZ_SIGNATURE:     "xx",
	}
}
