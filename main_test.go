package main

import (
	"fmt"
	"testing"
	"time"

	speedtest "github.com/kmoz000/RouterOsSpeedTest/v1"
)

func TestConvertStructToRouterOSArrayAndBack(t *testing.T) {
	// Call the function to convert the struct to RouterOS array
	testData := speedtest.Test{
		TX:        []int64{1024},
		RX:        []int64{1024},
		PING:      95,
		TestID:    "test0",
		CreatedAt: time.Now().Unix(),
	}
	since := time.Now()
	result := speedtest.ROString(testData)
	dsince := time.Since(since)
	// Check the expected output
	expectedResult := fmt.Sprintf(`{"TestID"="test0";"TX"=(1024);"RX"=(1024);"PING"=95;"CreatedAt"=%d}`, testData.CreatedAt)
	if result != expectedResult {
		t.Errorf("ConvertStructToRouterOSArray returned %s, expected %s", result, expectedResult)
	} else {
		t.Logf("ConvertStructToRouterOSArray %dns", dsince.Nanoseconds())
		if decoded, err := speedtest.DecodeROString(result); err == nil { // decode the rosScript string into an interface
			if decodedData, ok := decoded.(*speedtest.Test); ok { // convert resulted interface into your struct
				if decodedData.TestID != testData.TestID {
					t.Errorf("DecodeROString returned %+v, expected %+v", decodedData, testData)
				} else {
					t.Logf("DecodeROString %dns", dsince.Nanoseconds())
				}
			}
		}
	}
}
func TestRouterOsArraysregex(t *testing.T) {
	// Test cases
	testCases := []struct {
		input    string
		expected bool
	}{
		{`{}`, true}, // empty array
		{`{"key1"={"key2"=1};"key2"="value2";}`, true},              // valid key-value pairs
		{`{"key1"=11111.2;"key2"="value2"}`, true},                  // valid key-value pairs without trailing ;
		{`{"key1"="value1";"key2"=11111.2;"key3"}`, true},           // incomplete key-value pair
		{`{"key1"=(1,2,3);"key2"="value2";"key3"="value3";}`, true}, // valid key-value pairs
		{`{"key1"="value1";"key2"="value2",}`, true},                // trailing comma
		{`{"key1"="value1",}`, true},                                // trailing comma with single pair
		{`{invalid}`, false},                                        // invalid syntax
		{`"key"="value";"key2"="value2";`, false},                   // missing array brackets
	}

	// Run tests
	for i, tc := range testCases {
		since := time.Now()
		result := speedtest.IsRouterOSArray(tc.input)
		dsince := time.Since(since)
		fmt.Printf("Test %d: Expected: %t, Result: %t\t speed: %dns\n", i+1, tc.expected, result, dsince.Nanoseconds())

		if result != tc.expected {
			fmt.Printf("Test %d failed!\n", i+1)
		}
	}
}

func TestConvertToMilliseconds(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expected   int64
		shouldFail bool
	}{
		{"ValidInput", "00:00:00.060258", 60, false},
		{"ValidInput", "12:34:56.789123", 45296789, false},
		{"ValidInput", "01:02:03.004567", 3723004, false},
		{"ValidInput", "23:59:59.999999", 86399999, false},
		{"InvalidInput", "invalid_time_format", 0, true},
		{"ValidInput", "12:34:56", 45296000, false},
		{"ValidInput", "25:00:00.000000", 90000000, false},
		{"ValidInput", "00:60:00.000000", 3600000, false},
		{"ValidInput", "00:00:60.000000", 60000, false},
		{"ValidInput", "00:00:00.123456789", 123, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := speedtest.ConvertToMilliseconds(test.input)

			if test.shouldFail {
				if err == nil {
					t.Errorf("Expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != test.expected {
					t.Errorf("Expected %d milliseconds, but got %d", test.expected, result)
				}
			}
		})
	}
}
