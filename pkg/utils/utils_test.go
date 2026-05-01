package utils

import (
	"sort"
	"testing"
)

func TestGetUsernameFromEmail(t *testing.T) {
	// 测试空字符串输入
	username := GetUsernameFromEmail("")
	if len(username) != 8 {
		t.Errorf("Expected username length of 10, got %d", len(username))
	}

	// 测试有效的电子邮件地址
	email := "test@example.com"
	expectedUsername := "test"
	username = GetUsernameFromEmail(email)
	if username != expectedUsername {
		t.Errorf("Expected username %s, got %s", expectedUsername, username)
	}

	// 测试没有 "@" 符号的字符串
	invalidEmail := "testexample.com"
	username = GetUsernameFromEmail(invalidEmail)
	if len(username) != 8 {
		t.Errorf("Expected username length of 10, got %d", len(username))
	}

}

func TestGetAgeByBirthday(t *testing.T) {
	tests := []struct {
		birthday    string
		expectedAge int
	}{
		{"1990-01-01", 34},
		{"2000-12-31", 23},
		{"2022-03-15", 2},
		{"1980-09-25", 44},
		{"1980-10-01", 43},
	}

	for _, tt := range tests {
		actualAge := GetAgeByBirthday(tt.birthday)
		if actualAge != tt.expectedAge {
			t.Errorf("Expected age %d for birthday %s, but got %d", tt.expectedAge, tt.birthday, actualAge)
		}
	}
}

// 测试 sort.slice 排序功能
func TestSortSlice(t *testing.T) {
	// 测试空切片

	slice := []struct{ Sort int }{
		{Sort: 30},
		{Sort: 11},
		{Sort: 24},
	}
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].Sort > slice[j].Sort
	})
	for i, v := range slice {
		t.Logf("slice[%d].Sort: %d", i, v.Sort)
		if i != v.Sort-1 {
			t.Errorf("Expected slice[%d].Sort to be %d, but got %d", i, v.Sort-1, v.Sort)
		}
	}

}
