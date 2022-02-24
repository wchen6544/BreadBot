package Footlocker

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/x509"
    "encoding/pem"
    "time"
    jwt "github.com/golang-jwt/jwt"
    "github.com/tidwall/gjson"
    b64 "encoding/base64"
    "fmt"
    "strings"
    "github.com/go-resty/resty/v2"
    "crypto/sha256"
    "strconv"
    rand2 "math/rand"
)

type MyCustomClaims struct {
	jwt.StandardClaims
	Sub string `json:"sub"`
}

type JwtPayload struct {
	HiddenDetails  interface{} `json:"hidden_details"`
	AccountSid     string      `json:"account_sid"`
	Sid            string      `json:"sid"`
	ExpirationDate string   `json:"expiration_date"`
	DateCreated    string   `json:"date_created"`
	Status         string      `json:"status"`
	FactorSid      string      `json:"factor_sid"`
	Details        Details `json:"details"`
	ServiceSid string `json:"service_sid"`
	Identity   string `json:"identity"`
	jwt.StandardClaims

}
type Dictionary map[string]interface{}

type Details struct {
	Date    string     `json:"date"`
	Message string        `json:"message"`
	Fields  []Dictionary `json:"fields"`
}
func GenerateFactor(token string, proxy string) string{
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    publicKey := &privateKey.PublicKey
    
    _, encPub := encode(privateKey, publicKey)

    
	a := time.Now().UnixNano() / int64(time.Millisecond)
	requestID := strings.ToUpper(GenRandom(8)) + "-" + strings.ToUpper(GenRandom(4)) + "-" + strings.ToUpper(GenRandom(4)) + "-" + strings.ToUpper(GenRandom(4)) + "-" + strings.ToUpper(GenRandom(12))
	deviceID := strings.ToUpper(GenRandom(8)) + "-" + strings.ToUpper(GenRandom(4)) + "-" + strings.ToUpper(GenRandom(4)) + "-" + strings.ToUpper(GenRandom(4)) + "-" + strings.ToUpper(GenRandom(12))
	trace := GenRandom(16)

	client := resty.New()
	client.SetProxy(proxy)


	resp, err := client.R().
	SetHeader("Host", "www.footlocker.com").
	SetHeader("Accept", "application/json").
	SetHeader("X-FL-APP-VERSION", "5.4.0").
	SetHeader("X-FLAPI-API-IDENTIFIER", "921B2b33cAfba5WWcb0bc32d5ix89c6b0f614").
	SetHeader("X-FL-DEVICE-ID", deviceID).
	SetHeader("Accept-Language", "en-us").
	SetHeader("X-API-KEY", "m38t5V0ZmfTsRpKIiQlszub1Tx4FbnGG").
	SetHeader("User-Agent", "FootLocker/CFNetwork/Darwin").
	SetHeader("X-API-COUNTRY", "US").
	SetHeader("X-API-LANG", "en-US").
	SetHeader("X-FL-REQUEST-ID", requestID).
	SetHeader("User-Agent", "FootLocker; iOS; 5.4.0; 9647; iOS 14.4; iPad; TwilioVerify; 0.3.2; 1").

	Get("https://www.footlocker.com/apigate/mfa-core/v1/access-token/" + token)
	
	//fmt.Println("response1", resp.String())

	if err != nil {
        return "error"
	}
	if !strings.Contains(resp.String(), "customerId") {
		//fmt.Println("1", resp.String())
		return "error"
	}



	testProxy, err := client.R().
    SetBody([]byte(`{"customerId":"` + token + `", "factorId":"YF0361bfc4ebcbf55ad1c5593180ccb32f"}`)).
    SetHeader("Host", "www.footlocker.com").
    SetHeader("Content-Type", "application/json").
    SetHeader("Host", "www.footlocker.com").
	SetHeader("X-API-LANG", "en-US").
	SetHeader("X-NewRelic-ID", "VgAPVVdRDRAIVldUBQQEUFY=").
	SetHeader("User-Agent", "FootLocker/CFNetwork/Darwin").
	SetHeader("X-API-COUNTRY", "US").
	SetHeader("X-FL-APP-VERSION", "5.4.0").
	SetHeader("Accept-Language", "en-us").
	SetHeader("X-TIME-ZONE", "America/New_York").
	SetHeader("X-FL-DEVICE-ID", deviceID).
	SetHeader("Accept", "application/json").
	Post("https://www.footlocker.com/apigate/mfa-core/v1/send-challenge")

	if err != nil {
		return "error"
	}
	if strings.Contains(testProxy.String(), "geo") {
		fmt.Println("Geo")
		return "error"
	}


	time.Sleep(250 * time.Millisecond)
	client.SetHeader("Host", "verify.twilio.com")
	client.SetHeader("tracestate", "@nr=0-2-2684125-826625362-" + trace + "--0--" + strconv.FormatInt(int64(a), 10))
	client.SetHeader("Accept", "application/json")
	client.SetHeader("newrelic", GenerateRelic())
	client.SetHeader("Accept-Language", "en-us")
	client.SetHeader("X-NewRelic-ID", "VgAPVVdRDRAIVldUBQQEUFY=")
	client.SetHeader("User-Agent", "FootLocker; iOS; 5.4.0; 9647; iOS 14.4; iPad; TwilioVerify; 0.3.2; 1")

	client.SetHeader("traceparent", "00-" + GenRandom(32) + "-" + trace + "-00")

    accessToken := gjson.Get(resp.String(), "accessToken").String()
    encodedAccessToken := b64.StdEncoding.EncodeToString([]byte("token:" + accessToken))
    formattedPublicKey := strings.ReplaceAll(strings.Split(strings.Split(encPub, "-----BEGIN PUBLIC KEY-----")[1], "-----END PUBLIC KEY-----")[0], "\n", "")
	resp2, err := client.R().
	SetFormData(map[string]string{
	    "FriendlyName": "com.footlocker.twilio",
	    "FactorType": "push",
	    "Binding.Alg": "ES256",
	    "Binding.PublicKey": formattedPublicKey,
	    "Config.SdkVersion": "0.3.2",
	    "Config.AppId": "com.footlocker.approved",
	    "Config.NotificationPlatform": "apn",
	    "Config.NotificationToken": "c03a8aa7d9072184f63932e12bbb99c84fc872b470994dfe807c8e2cdb0d9c3f",
	}).
	SetHeader("Authorization", "Basic " + encodedAccessToken).
	Post("https://verify.twilio.com/v2/Services/VAc6913fecf8cf66f1f3cf324563d6071d/Entities/" + token + "/Factors")
	//fmt.Println("resp1.5", resp2.String())
	if err != nil {
		return "error"
	}

	if !strings.Contains(resp2.String(), "unverified") {
		
		if strings.Contains(resp2.String(), "60315") {

			fmt.Println("Factor Limit: Error")
		} else {
			fmt.Println("2", resp2.String())
		}
		return "error"
	}


    account_sid := gjson.Get(resp2.String(), "account_sid").String()
	credential_sid := gjson.Get(resp2.String(), "config.credential_sid").String()
	sid := gjson.Get(resp2.String(), "sid").String()

	authToken1 := b64.StdEncoding.EncodeToString([]byte("token:" + jwt1(credential_sid, account_sid, privateKey)))
	resp3, err := client.R().
		SetFormData(map[string]string{
			"AuthPayload": getAuthPayload(sid, privateKey),
		}).
	SetHeader("Authorization", "Basic " + authToken1).
	Post("https://verify.twilio.com/v2/Services/VAc6913fecf8cf66f1f3cf324563d6071d/Entities/" + token + "/Factors/" + sid)
	if err != nil {
		return "error"
	}

	if !strings.Contains(resp3.String(), "verified") {
		fmt.Println("3", resp3.String())
		return "error"
	}
	time.Sleep(333 * time.Millisecond)
	resp4, err := client.R().
    SetBody([]byte(`{"customerId":"` + token + `", "factorId":"` + sid + `"}`)).
    SetHeader("Host", "www.footlocker.com").
    SetHeader("Content-Type", "application/json").
    SetHeader("Host", "www.footlocker.com").
	SetHeader("X-API-LANG", "en-US").
	SetHeader("X-NewRelic-ID", "VgAPVVdRDRAIVldUBQQEUFY=").
	SetHeader("User-Agent", "FootLocker/CFNetwork/Darwin").
	SetHeader("X-API-COUNTRY", "US").
	SetHeader("X-FL-APP-VERSION", "5.4.0").
	SetHeader("Accept-Language", "en-us").
	SetHeader("X-TIME-ZONE", "America/New_York").
	SetHeader("X-FL-DEVICE-ID", deviceID).
	SetHeader("Accept", "application/json").
	Post("https://www.footlocker.com/apigate/mfa-core/v1/send-challenge")

	if err != nil {
		return "error"
	}
	if strings.Contains(resp4.String(), "geo") {
		return "error"
	}

	authToken2 := b64.StdEncoding.EncodeToString([]byte("token:" + jwt1(credential_sid, account_sid, privateKey)))
	resp5, err := client.R().
	SetHeader("Authorization", "Basic " + authToken2).
	Get("https://verify.twilio.com/v2/Services/VAc6913fecf8cf66f1f3cf324563d6071d/Entities/" + token + "/Challenges?FactorSid=" + sid + "&PageSize=20&Status=pending")
	//fmt.Println("Response 4: ", resp5.String())

	if err != nil {
		return "error"
	}
	if !strings.Contains(resp5.String(), "pending") {
		//fmt.Println("Task Number: Error")
		return "error"
	}

    challengeSid := gjson.Get(resp5.String(), "challenges.0.sid").String()
    expDate := gjson.Get(resp5.String(), "challenges.0.expiration_date").String()
    dateRn := gjson.Get(resp5.String(), "challenges.0.date_updated").String()

	time.Sleep(500 * time.Millisecond)

	authToken3 := b64.StdEncoding.EncodeToString([]byte("token:" + jwt1(credential_sid, account_sid, privateKey)))
	resp6, err := client.R().
	SetHeader("Authorization", "Basic " + authToken3).
	Get("https://verify.twilio.com/v2/Services/VAc6913fecf8cf66f1f3cf324563d6071d/Entities/" + token + "/Challenges/" + challengeSid)
	//fmt.Println("Response 5: ", resp6.String())
	
	if err != nil {
		return "error"
	}
	if !strings.Contains(resp6.String(), "pending") {
		fmt.Println("Task Number: Error")
		return "error"
	}



	authToken4 := b64.StdEncoding.EncodeToString([]byte("token:" + jwt1(credential_sid, account_sid, privateKey)))
	resp7, err := client.R().
	SetHeader("Authorization", "Basic " + authToken4).
	Get("https://verify.twilio.com/v2/Services/VAc6913fecf8cf66f1f3cf324563d6071d/Entities/" + token + "/Challenges/" + challengeSid)
	//fmt.Println("Response 6: ", resp7.String())

	if err != nil {
		return "error"
	}
	if !strings.Contains(resp7.String(), "pending") {
		fmt.Println("Task Number: Error")
		return "error"
	}


	authPayload := jwt2(account_sid, challengeSid, expDate, sid, dateRn, privateKey, token)

	authToken5 := b64.StdEncoding.EncodeToString([]byte("token:" + jwt1(credential_sid, account_sid, privateKey)))
	resp8, err := client.R().
		SetFormData(map[string]string{
			"AuthPayload": authPayload,
		}).
	SetHeader("Authorization", "Basic " + authToken5).
	Post("https://verify.twilio.com/v2/Services/VAc6913fecf8cf66f1f3cf324563d6071d/Entities/" + token + "/Challenges/" + challengeSid)
	//fmt.Println("Final Response: ", resp8.String())
	if strings.Contains(resp8.String(), "approved") {
		return sid
	} else {
		return "error"
	}

}
func EncodeSignatureJWT(sig []byte) string {
	return b64.StdEncoding.EncodeToString(sig)
}

func getAuthPayload(a string, pKey *ecdsa.PrivateKey) string{
	hash := sha256.Sum256([]byte(a))
	sig, err := ecdsa.SignASN1(rand.Reader, pKey, hash[:])
	if err != nil {
		panic(err)
	}
	return EncodeSignatureJWT(sig)
	
}


func encode(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) (string, string) {
    x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
    pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

    x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
    pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

    return string(pemEncoded), string(pemEncodedPub)
}

func decode(pemEncoded string, pemEncodedPub string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
    block, _ := pem.Decode([]byte(pemEncoded))
    x509Encoded := block.Bytes
    privateKey, _ := x509.ParseECPrivateKey(x509Encoded)

    blockPub, _ := pem.Decode([]byte(pemEncodedPub))
    x509EncodedPub := blockPub.Bytes
    genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
    publicKey := genericPublicKey.(*ecdsa.PublicKey)

    return privateKey, publicKey
}
func jwt1(credential_sid string, account_sid string, privateKey *ecdsa.PrivateKey) string{
    claims := MyCustomClaims{
		jwt.StandardClaims{
	        ExpiresAt: time.Now().Unix() + 600,
			NotBefore: time.Now().Unix(),
		},
		account_sid,
	}
    token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
    token.Header["kid"] = credential_sid
    token.Header["cty"] = "twilio-pba;v=1"
    tokenString, err := token.SignedString(privateKey)
    if err != nil {
        return "error"
    }
    return tokenString
}

func jwt2(account_sid string, challengeSid string, expDate string, sid string, dateRn string, privateKey *ecdsa.PrivateKey, token string) string{
    data := []Dictionary{}
	d := Details{
		Date: dateRn,
		Message: "FootLocker MFA verify push",
		Fields: data,
	}
	var x interface{}
	b := JwtPayload{x, account_sid, challengeSid, expDate, dateRn, "approved", sid, d, "VAc6913fecf8cf66f1f3cf324563d6071d", token, jwt.StandardClaims{}}
    token1 := jwt.NewWithClaims(jwt.SigningMethodES256, b)
    tokenString1, err := token1.SignedString(privateKey)
    if err != nil {
        return "error"
    }
    return tokenString1
}

func GenerateRelic() string{
	t := time.Now().UnixNano() / 10000000
	p1 := `{
"d": {
"ac": "2684125",
"ap": "826625362",
"id": "`
	p2 := GenRandom(16)
	p3 := `",
"ti": `
	p4 := strconv.FormatInt(int64(t), 10)

	p5 := `,
"tr": "`
	p6 := GenRandom(32)
	p7 := `",
"ty": "Mobile"
},
"v": [
0,
2
]
}`
	v := p1 + p2 + p3 + p4 + p5 + p6 + p7
	return b64.StdEncoding.EncodeToString([]byte(v))
}

func GenRandom(l int) string{
	rand2.Seed(time.Now().UnixNano())
	s := ""
	alpha := "abcde123456789"
	for i := 0; i < l; i++ {
		s += string(alpha[rand2.Intn(len(alpha))])
	}
	return s
}