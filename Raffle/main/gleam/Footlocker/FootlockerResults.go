package Footlocker

/*
TODO
- add zipcode/address support like error handling
- improve codebase
- add time sleep
- duplicate error
Response 5:  {"errors":[{"code":"41012","message":"There was an error with this request.","subject":"Duplicate Record","type":"DuplicateRecordError"}]}
*/
import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
	"strings"
	"strconv"
	b64 "encoding/base64"
	"math/rand"
	"github.com/tidwall/gjson"
	struc "main/testFolder"


)

func Results(config interface{}, necessary interface{}){

	c1, _ := necessary.(struc.Necessary)
	c, _ := config.(struc.SiteStructure)

	defer c1.WaitGroup.Done()
	c1.CH <- true
	rand.Seed(time.Now().UnixNano())
	email := c1.Email
	proxy := c1.Proxy
	//taskNumber := c1.TaskNumber
	password := "Hoglund_6"
	shoeID := c.Shoe

	a := time.Now().UnixNano() / int64(time.Millisecond)
	requestID := strings.ToUpper(genRandom(8)) + "-" + strings.ToUpper(genRandom(4)) + "-" + strings.ToUpper(genRandom(4)) + "-" + strings.ToUpper(genRandom(4)) + "-" + strings.ToUpper(genRandom(12))
	deviceID := strings.ToUpper(genRandom(8)) + "-" + strings.ToUpper(genRandom(4)) + "-" + strings.ToUpper(genRandom(4)) + "-" + strings.ToUpper(genRandom(4)) + "-" + strings.ToUpper(genRandom(12))
	trace := genRandom(16)
	client := resty.New()
	client.SetHeader("Host", "www.footlocker.com")
	client.SetHeader("tracestate", "@nr=0-2-2684125-826625362-" + trace + "--0--" + strconv.FormatInt(int64(a), 10))
	client.SetHeader("Accept", "application/json")
	client.SetHeader("X-FL-APP-VERSION", "5.3.5")
	client.SetHeader("X-FLAPI-API-IDENTIFIER", "921B2b33cAfba5WWcb0bc32d5ix89c6b0f614")
	client.SetHeader("newrelic", generateRelic())
	client.SetHeader("X-FL-DEVICE-ID", deviceID)
	client.SetHeader("Accept-Language", "en-us")
	client.SetHeader("X-API-KEY", "m38t5V0ZmfTsRpKIiQlszub1Tx4FbnGG")
	client.SetHeader("traceparent", "00-" + genRandom(32) + "-" + trace + "-00")
	client.SetHeader("User-Agent", "FootLocker/CFNetwork/Darwin")
	client.SetHeader("X-API-COUNTRY", "US")
	client.SetHeader("X-API-LANG", "en-US")
	client.SetHeader("X-FL-REQUEST-ID", requestID)
	client.SetProxy(proxy)
	client.SetTimeout(30 * time.Second)

	resp, err := client.R().
	Get("https://www.footlocker.com/apigate/v3/session")
	fmt.Println("Response 1: ", resp.String())
	if err != nil {
		<- c1.CH
		c1.ReturnData <- map[string]string{"email": email, "status": "Bad"}
		return 
	}

    csrfToken := gjson.Get(resp.String(), "data.csrfToken").String()
    jSessionId := strings.Split(strings.Split(resp.Header()["Set-Cookie"][0], "JSESSIONID=")[1], ";")[0]

    resp2, err := client.R().
    SetBody([]byte(`{"password":"` + password + `","uid":"` + email + `"}`)).
    SetHeader("Content-Type", "application/json").
    SetHeader("X-CSRF-TOKEN", csrfToken).
    SetHeader("X-FLAPI-SESSION-ID", jSessionId).

	Post("https://www.footlocker.com/apigate/v3/auth")
	//fmt.Println("Response 2: ", resp2.String())

	time.Sleep(1 * time.Second)
	accessToken := gjson.Get(resp2.String(), "oauthToken.access_token").String()
	

	resp3, err := client.R().
	SetHeader("X-FLAPI-RESOURCE-IDENTIFIER", accessToken).
	SetHeader("X-FLAPI-TIMEOUT", "42479").
	SetHeader("X-API-COUNTRY", "US").
	SetHeader("X-FL-APP-VERSION", "5.3.5").
	SetHeader("X-FLAPI-SESSION-ID", jSessionId).
	SetHeader("Accept-Language", "en-us").
	SetHeader("Accept", "application/json").
	Get("https://www.footlocker.com/apigate/v3/users/account-info")
	//fmt.Println("Response 3: ", resp3.String())
	customerID := gjson.Get(resp3.String(), "customerID").String()
	cCoreCustomerId := gjson.Get(resp3.String(), "cCoreCustomerId").String()
	//fmt.Println(customerID, cCoreCustomerId)


	finalResp, err := client.R().
	SetHeader("X-FLAPI-RESOURCE-IDENTIFIER", accessToken).
	SetHeader("X-FLAPI-TIMEOUT", "42479").
	SetHeader("X-API-COUNTRY", "US").
	SetHeader("X-FL-APP-VERSION", "5.3.5").
	SetHeader("X-FLAPI-SESSION-ID", jSessionId).
	SetHeader("Accept-Language", "en-us").
	SetHeader("Accept", "application/json").
	SetHeader("X-CCore-Number", cCoreCustomerId).
	SetHeader("X-CUSTOMER-NUMBER", customerID).
	SetHeader("X-CSRF-TOKEN", csrfToken).
	Get("https://www.footlocker.com/apigate/release-reservation/" + shoeID)

	fmt.Println("Response 5: ", finalResp.String())
	



}





func genRandom(l int) string{
	s := ""
	alpha := "abcde123456789"
	for i := 0; i < l; i++ {
		s += string(alpha[rand.Intn(len(alpha))])
	}
	return s
}
func generateRelic() string{
	t := time.Now().UnixNano() / 10000000
	p1 := `{
"d": {
"ac": "2684125",
"ap": "826625362",
"id": "`
	p2 := genRandom(16)
	p3 := `",
"ti": `
	p4 := strconv.FormatInt(int64(t), 10)

	p5 := `,
"tr": "`
	p6 := genRandom(32)
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