package Footlocker

import (
	"github.com/go-resty/resty/v2"
	"time"
	"strconv"
	"net/url"
	"strings"
	"fmt"
	struc "main/testFolder"

)

//func Confirm(URL string, ch chan bool, wg *sync.WaitGroup, taskNumber int, proxy string, returnData chan map[string]string) {

func Confirm(config interface{}, necessary interface{}){
	c1, _ := necessary.(struc.Necessary)
	defer c1.WaitGroup.Done()
	c1.CH <- true

	URL := c1.Email
	proxy := c1.Proxy
	taskNumber := c1.TaskNumber
	taskN := strconv.Itoa(taskNumber)

	fmt.Println("Starting Task " + taskN)


	decodedUrl, _ := url.QueryUnescape(URL)
	token := strings.Split(strings.Split(decodedUrl, "https://www.footlocker.com/user-activation.html?activationToken=")[1], "&ssoComplete=true&inStore=false")[0]
	client := resty.New()
	client.SetHeader("authority", "www.footlocker.com")
	client.SetHeader("accept", "application/json")
	client.SetHeader("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	client.SetHeader("content-type", "application/json")
	client.SetHeader("sec-gpc", "1")
	client.SetHeader("origin", "https://www.footlocker.com")
	client.SetHeader("sec-fetch-site", "same-origin")
	client.SetHeader("sec-fetch-mode", "cors")
	client.SetHeader("sec-fetch-dest", "empty")
	client.SetHeader("accept-language", "en-US,en;q=0.9")
	client.SetProxy(proxy)
	client.SetTimeout(30 * time.Second)


	timestamp := time.Now().UnixNano() / 1000000
	resp, err := client.R().
    SetBody(`{"activationToken":"` + token + `"}`).
    Post("https://www.footlocker.com/api/v3/activation?timestamp=" + strconv.FormatInt(int64(timestamp), 10))
    if err != nil{
    	<- c1.CH
    	c1.ReturnData <- map[string]string{"email": URL, "status": "Bad"}
    	return 
    }
    if strings.Contains(resp.String(), "Success") || strings.Contains(resp.String(), "invalid token"){
    	c1.ReturnData <- map[string]string{"email": URL, "status": "Good"}

    } else {
    	c1.ReturnData <- map[string]string{"email": URL, "status": "Bad"}

    }
    <- c1.CH
}