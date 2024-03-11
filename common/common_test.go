package common

import (
	"reflect"
	"testing"
)

func TestParseIP(t *testing.T) {
	testCases := []struct {
		name     string
		ipString string
		expected []string
	}{
		{"Single IP", "192.168.1.1", []string{"192.168.1.1"}},
		{"IP Range", "192.168.1.1-192.168.1.3", []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"}},
		{"CIDR", "192.168.1.0/30", []string{"192.168.1.0", "192.168.1.1", "192.168.1.2", "192.168.1.3"}},
		// 更多测试用例...
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseIP(tc.ipString)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("ParseIP(%s) = %v, expected %v", tc.ipString, result, tc.expected)
			}
		})
	}
}

// 其他函数的测试...
