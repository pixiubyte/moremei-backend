package encrypt

import (
	"fmt"
	"testing"
)

func TestSha256(t *testing.T) {
	str := fmt.Sprintf("%s%d", "12345abcdef", 1704855303)
	hashStr := Sha256(str)
	if hashStr == "" {
		t.Errorf("err")
	}
	fmt.Printf(hashStr)
}

func TestIsSha256(t *testing.T) {
	str := fmt.Sprintf("%s%d", "12345abcdef", 1704855303)
	result := "28e5b9923ba82737a5d2e07381ebc3e592b6982491099edb18d3e08588999b90"
	if !IsSha256(str, result) {
		t.Errorf("err:%v", str)
	} else {
		fmt.Printf("success")
	}
}
