package encrypt

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestSignPKCS1v15(t *testing.T) {

	stamp := fmt.Sprintf("%d", time.Now().Unix())
	nonceStr := "4r2t43gg1"
	str := stamp + "&" + nonceStr
	key := `-----BEGIN PRIVATE KEY-----
MIICeAIBADANBgkqhkiG9w0BAQEFAASCAmIwggJeAgEAAoGBAJfVRaVN+RwyQTWjwAzciJy1SbOybeTBOBaCu+lvrQChbSqHGNfFSr56d/RCgvRLRhooOPXbwIoaifk1WYIvHZDyyVspB2vGFFF7AFa7GybIpSFgKzj5b4Jm2brBwoHSASBtDDNWPAfL4ctRybOv7srXkYDWrvwUWyvamYBsIfSXAgMBAAECgYAUD9qdfoGJb7x9DW99Y5BrgGrGIy/xF3TFSK3yqM5ntGy7v1ERhaCcLYt3C2jJRb70CteH6Or2dI5BjwKOOJKD1e6i00FIbtlFiBwxg7FC9eOx4DMgxu6j93COAvbsGl8a+bc0+CZpz12xS64IM5YN0loGmbkCOTDAbT1zOuZtSQJBAMZD0KTeoWEd84cVnnglHJ72Y0T3E5uHKAOM2FeGGy/6BmcenEjDIyT17dAOHD+Zi7G9UPcA2NihhnLCIBhAcJ8CQQDEDBMC7ADy5xJdK0B9pmwZq3BBZhuji6WUKszCdB//mlvozaCIIdiruamkdaJq/UUqpvFAjf1T6ItRDUzGtqEJAkEAlmu5Dn0CPyZ0LxbN5iVx84DHi/lQ3PzL9PWU5cKPOfUdinsE44d5UH9tcB5kfDRIcg9KMDxqSOEzmjmCFCQ/zQJBAJpqeD8A7O5mGwzPmIhfoR3G7zBT4Mk8oTrHS2iOVvXY+zOvYxZWsnbwUjJ7hWaH/wbNX5DdRf/lVnaM50BNcSECQQCiwXSBy55lnYbhS20fO95rsVv3D4cmRvfMOinQXpHw29Yxhf5m1S7PgMPLYJwtzuwY2/nJOT+uSpkCUs8YIcoP
-----END PRIVATE KEY-----`
	var s Sha256WithRsa

	sign, err := s.SignPKCS1v15(key, str)
	if err != nil {
		panic(err)
	}
	signBase64 := base64.StdEncoding.EncodeToString([]byte(sign))
	req, err := http.NewRequest(http.MethodGet, "https://557n539z00.vicp.fun/api/v1/tpos/customer/customer-search", nil)
	req.Header.Set("tpos-timestamp", stamp)
	req.Header.Set("tpos-nonce-str", nonceStr)
	req.Header.Set("tops-account", "moremei")
	req.Header.Set("tpos-sign", signBase64)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
}
