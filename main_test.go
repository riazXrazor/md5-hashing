package main

import (
	"os"
	"testing"
)

var testCases = []struct {
	ctype    string
	data     string
	expected string
}{
	// {
	// 	ctype:    "text",
	// 	data:     "ok",
	// 	expected: "c4ca4238a0b923820dcc509a6f75849b",
	// },
	{
		ctype:    "file",
		data:     "test_files/4k_file.bin",
		expected: "612b29128abdd0d13f91d40486114d59",
	},
	{
		ctype:    "file",
		data:     "test_files/empty.txt",
		expected: "d41d8cd98f00b204e9800998ecf8427e",
	},
	{
		ctype:    "file",
		data:     "test_files/hello.txt",
		expected: "b1946ac92492d2347c6235b4d2611184",
	},
}

func TestMD5(t *testing.T) {
	for _, tc := range testCases {
		switch tc.ctype {
		// case "text":
		// 	hash := md5(tc.data)
		// 	if hash != tc.expected {
		// 		t.Errorf("--- FAIL: md5sum(%s) = %s, expected=%s", tc.data, hash, tc.expected)
		// 	}
		case "file":
			fp, err := os.Open(tc.data)
			defer fp.Close()
			if err != nil {
				t.Fatalf("unable to read file %s", tc.data)
			}

			hash := md5(fp)
			if hash != tc.expected {
				t.Errorf("--- FAIL: md5sum(%s) = %s, expected=%s", tc.data, hash, tc.expected)
			}
		}
	}
}
