package Footlocker

/*
TODO
- add zipcode/address support like error handling
- improve codebase
- add time sleep
*/
import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
	"strings"
	"strconv"
	"github.com/tidwall/gjson"
	rand "math/rand"
	struc "main/testFolder"
)
type Size struct {
	SelectedRegion string `json:"selectedRegion"`
	Us             string `json:"us"`
	Eu             string `json:"eu"`
	Uk             string `json:"uk"`
}
type Attributes struct {
	OnlineParticipating bool     `json:"onlineParticipating"`
	StoreIds            []string `json:"storeIds"`
	LaunchEventID       string   `json:"launchEventId"`
	Size                Size `json:"size"`
	Platform           string `json:"platform"`
	StoreParticipating bool   `json:"storeParticipating"`
}
type MoreData struct {
	Type       string `json:"type"`
	Attributes Attributes `json:"attributes"`
}
type FinalData struct {
	CustomerID string `json:"customerId"`
	FactorID   string `json:"factorId"`
	Data      MoreData `json:"data"`
}

type TaskReturn struct {
	Email string
	Status string
}

func Start(config interface{}, necessary interface{}){
	c1, _ := necessary.(struc.Necessary)
	c, _ := config.(struc.SiteStructure)
	defer c1.WaitGroup.Done()
	c1.CH <- true
	rand.Seed(time.Now().UnixNano())
	email := c1.Email
	proxy := c1.Proxy
	taskNumber := c1.TaskNumber
	password := "Hoglund_6"
	zipcode := c.Store
	shoeID := c.Shoe
	randSize := strings.Split(c.Size, ",")
	size := randSize[rand.Intn(len(randSize))]
	taskN := strconv.Itoa(taskNumber)
	fmt.Println("Starting Task " + taskN)
	storeIDs := []string{}
	a := time.Now().UnixNano() / int64(time.Millisecond)
	requestID := strings.ToUpper(GenRandom(8)) + "-" + strings.ToUpper(GenRandom(4)) + "-" + strings.ToUpper(GenRandom(4)) + "-" + strings.ToUpper(GenRandom(4)) + "-" + strings.ToUpper(GenRandom(12))
	deviceID := strings.ToUpper(GenRandom(8)) + "-" + strings.ToUpper(GenRandom(4)) + "-" + strings.ToUpper(GenRandom(4)) + "-" + strings.ToUpper(GenRandom(4)) + "-" + strings.ToUpper(GenRandom(12))
	trace := GenRandom(16)
	client := resty.New()
	client.SetHeader("Host", "www.footlocker.com")
	client.SetHeader("tracestate", "@nr=0-2-2684125-826625362-" + trace + "--0--" + strconv.FormatInt(int64(a), 10))
	client.SetHeader("Accept", "application/json")
	client.SetHeader("X-FL-APP-VERSION", "5.4.0")
	client.SetHeader("X-FLAPI-API-IDENTIFIER", "921B2b33cAfba5WWcb0bc32d5ix89c6b0f614")
	client.SetHeader("newrelic", GenerateRelic())
	client.SetHeader("X-FL-DEVICE-ID", deviceID)
	client.SetHeader("Accept-Language", "en-us")
	client.SetHeader("X-API-KEY", "m38t5V0ZmfTsRpKIiQlszub1Tx4FbnGG")
	client.SetHeader("traceparent", "00-" + GenRandom(32) + "-" + trace + "-00")
	client.SetHeader("User-Agent", "FootLocker/CFNetwork/Darwin")
	client.SetHeader("X-API-COUNTRY", "US")
	client.SetHeader("X-API-LANG", "en-US")
	client.SetHeader("X-FL-REQUEST-ID", requestID)
	client.SetProxy(proxy)
	resp, err := client.R().
	Get("https://www.footlocker.com/apigate/v3/session")
	if err != nil {
		fmt.Println("Task Number " + taskN + " : ", err)
		<- c1.CH
		c1.ReturnData <- map[string]string{"email": email, "status": "Bad"}
		return 
	}
	//fmt.Println("Response 1: ", resp.String())
	if !strings.Contains(resp.String(), "csrfToken") {
		fmt.Println("Task Number " + taskN + " : Cannot Generate Session")
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
	if err != nil {
		fmt.Println("Task Number " + taskN + " : ", err)
		<- c1.CH
		c1.ReturnData <- map[string]string{"email": email, "status": "Bad"}
		return 
	}

	if strings.Contains(resp2.String(), "geo") {
		fmt.Println("Task Number " + taskN + " : Proxy Ban")
		<- c1.CH
		c1.ReturnData <- map[string]string{"email": email, "status": "Bad"}
		return
	}

	time.Sleep(1 * time.Second)
	accessToken := gjson.Get(resp2.String(), "oauthToken.access_token").String()
	//fmt.Println(accessToken)

	resp3, err := client.R().
	SetHeader("X-FLAPI-RESOURCE-IDENTIFIER", accessToken).
	SetHeader("X-FLAPI-TIMEOUT", "42479").
	SetHeader("X-API-COUNTRY", "US").
	SetHeader("X-FL-APP-VERSION", "5.4.0").
	SetHeader("X-FLAPI-SESSION-ID", jSessionId).
	SetHeader("Accept-Language", "en-us").
	SetHeader("Accept", "application/json").
	Get("https://www.footlocker.com/apigate/v3/users/account-info")

	if err != nil {
		fmt.Println("Task Number " + taskN + " : ", err)
		<- c1.CH
		c1.ReturnData <- map[string]string{"email": email, "status": "Bad"}
		return 
	}

	//fmt.Println("Response 3: ", resp3.String())

	if strings.Contains(resp3.String(), "geo") {
		fmt.Println("Task Number " + taskN + " : Proxy Ban", )
		<- c1.CH
		c1.ReturnData <- map[string]string{"email": email, "status": "Bad"}
		return
	}

	customerID := gjson.Get(resp3.String(), "customerID").String()
	cCoreCustomerId := gjson.Get(resp3.String(), "cCoreCustomerId").String()
	//fmt.Println(customerID, cCoreCustomerId)

	resp4, err := client.R().
	SetQueryParams(map[string]string{
		"procedure": "2",
		"sku": shoeID,
		"address": zipcode,
    }).
	SetHeader("X-FLAPI-RESOURCE-IDENTIFIER", accessToken).
	SetHeader("X-FLAPI-SESSION-ID", jSessionId).
	Get("https://www.footlocker.com/apigate/launch-stores")

	if err != nil {
		fmt.Println("Task Number " + taskN + " : ", err)
		<- c1.CH
		c1.ReturnData <- map[string]string{"email": email, "status": "Bad"}
		return 
	}

	//fmt.Println("Response 4: ", resp4.String())

	time.Sleep(500 * time.Millisecond)
	if strings.Count(resp4.String(), "displayName") >= 1 {
		for i := 0; i < 3; i++ {
			if gjson.Get(resp4.String(), "stores." + strconv.Itoa(i) + ".id").String() != "" {
				storeIDs = append(storeIDs, gjson.Get(resp4.String(), "stores." + strconv.Itoa(i) + ".id").String())
			}

		}
	}
	//fmt.Println(storeIDs, len(storeIDs))
	if len(storeIDs) == 0 {
		fmt.Println("Task Number " + taskN + " : No Available Stores")
		<- c1.CH
		c1.ReturnData <- map[string]string{"email": email, "status": "Bad"}
		return 
	}
	//fmt.Println(cCoreCustomerId)
	factor := GenerateFactor(cCoreCustomerId, proxy)
	if factor == "error" {
		fmt.Println("Task Number " + taskN + " : Factor Error")
		<- c1.CH
		c1.ReturnData <- map[string]string{"email": email, "status": "Bad"}
		return 
	}
	
	sizeData := Size{
		SelectedRegion: "us",
		Us: size,
		Eu: "",
		Uk: "",
	}
	attributeData := Attributes{
		OnlineParticipating: false,
		StoreIds: storeIDs,
		LaunchEventID: shoeID,
		Size: sizeData,
		Platform: "ios",
		StoreParticipating: true,
	}
	someData := MoreData{
		Type: "reservation",
		Attributes: attributeData,
	}
	combinedData := FinalData{
		CustomerID: cCoreCustomerId,
		FactorID: factor,
		Data: someData,

	}
	time.Sleep(1500 * time.Millisecond)
	client.SetCookies(resp2.Cookies())
	coo := resp2.Cookies()
	value := ""
	for _, e := range coo {
		if e.Name == "JSESSIONID" {
			value = e.Value
		}
	}
	finalResp, err := client.R().
	SetBody(combinedData).
	SetHeader("Host", "www.footlocker.com").
	SetHeader("X-API-KEY", "m38t5V0ZmfTsRpKIiQlszub1Tx4FbnGG").
	SetHeader("X-FLAPI-RESOURCE-IDENTIFIER", accessToken).
	SetHeader("X-API-LANG", "en-US").
	SetHeader("User-Agent", "FootLocker/CFNetwork/Darwin").
	SetHeader("X-API-COUNTRY", "US").
	SetHeader("X-CCore-Number", cCoreCustomerId).
	SetHeader("X-FLAPI-TIMEOUT", "41581").
	SetHeader("X-CUSTOMER-NUMBER", customerID).
	SetHeader("X-FL-APP-VERSION", "5.4.0").
	SetHeader("X-FLAPI-SESSION-ID", value).
	SetHeader("X-CSRF-TOKEN", csrfToken).
	SetHeader("Accept-Language", "en-us").
	SetHeader("X-FLAPI-API-IDENTIFIER", "921B2b33cAfba5WWcb0bc32d5ix89c6b0f614").
	SetHeader("X-TIME-ZONE", "America/New_York").
	SetHeader("Content-Type", "application/json").
	SetHeader("Accept", "application/json").
	SetHeader("X-NewRelic-ID", "VgAPVVdRDRAIVldUBQQEUFY=").

	Post("https://www.footlocker.com/apigate/reservations/")
	//fmt.Println(accessToken, jSessionId, cCoreCustomerId, customerID, csrfToken, resp2.Cookies())
	//fmt.Println(strings.Split(strings.Split(coo, "JSESSIONID=")[1], ";")[0])
	if err != nil {
		fmt.Println("Task Number " + taskN + " : ", err)
		<- c1.CH
		c1.ReturnData <- map[string]string{"email": email, "status": "Bad"}
		return 
	}

	//fmt.Println("Response 5: ", finalResp.String())
	if strings.Contains(finalResp.String(), `"reservations":true`) {
		c1.ReturnData <- map[string]string{"email": email, "status": "Success"}
	} else if strings.Contains(finalResp.String(), "Duplicate") {
		fmt.Println(finalResp.String())
		c1.ReturnData <- map[string]string{"email": email, "status": "Bad"}
	} else {
		fmt.Println(finalResp.String())
		c1.ReturnData <- map[string]string{"email": email, "status": "Bad"}
	}

	<- c1.CH

}

