package feishu

import (
	"fmt"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("cli_a792280fdcf6900b", "Wx7U9EqNP3unsNnrQsYCIcUWuSJiaduY")
	err := client.WriteSheet("CIovsx8m0hJvfAtk1R8cIvM3nwb", "c65ee4", [][]interface{}{
		{"11", "33", "555"},
		{"22", "33", "555"},
		{"33", "33", "555"},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", err)

}
