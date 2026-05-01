package auth

import (
	"testing"
)

func TestParseProviderName(t *testing.T) {
	tests := []struct {
		name             string
		provider         string
		expectedTenantId string
		expectedApp      string
		expectedProvider string
	}{
		{
			name:             "正常情况",
			provider:         "tenant1:app1:provider1",
			expectedTenantId: "tenant1",
			expectedApp:      "app1",
			expectedProvider: "provider1",
		},
		{
			name:             "分隔符数量不对",
			provider:         "tenant1:app1",
			expectedTenantId: "",
			expectedApp:      "",
			expectedProvider: "",
		},
		{
			name:             "有空值情况",
			provider:         "tenant1::provider1",
			expectedTenantId: "tenant1",
			expectedApp:      "",
			expectedProvider: "provider1",
		},
		{
			name:             "只有 provider 的情况",
			provider:         "::provider1",
			expectedTenantId: "",
			expectedApp:      "",
			expectedProvider: "provider1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tenantId, app, provider := ParseProviderName(tt.provider)
			if tenantId != tt.expectedTenantId || app != tt.expectedApp || provider != tt.expectedProvider {
				t.Errorf("期望的结果 (%s, %s, %s)，实际得到的结果 (%s, %s, %s)", tt.expectedTenantId, tt.expectedApp, tt.expectedProvider, tenantId, app, provider)
			}
		})
	}
}
