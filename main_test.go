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
		TX:        1024,
		RX:        1024,
		PING:      95,
		Size:      0,
		TestID:    "test0",
		CreatedAt: time.Now().Unix(),
	}
	since := time.Now()
	result := speedtest.ROString(testData)
	dsince := time.Since(since)
	// Check the expected output
	expectedResult := fmt.Sprintf(`{"TestID"="test0";"TX"=1024;"RX"=1024;"PING"=95;"Size"=0;"CreatedAt"=%d}`, testData.CreatedAt)
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
