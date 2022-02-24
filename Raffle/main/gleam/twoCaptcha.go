package gleam


import (
  "io/ioutil"
  "github.com/go-resty/resty/v2"
  "time"
  "strings"
  "fmt"
  "sync"
  "math/rand"
  //"net/url"
  "net/http"

)

func GetCaptcha(apiKey string, method string, siteKey string, url string, p string) string{
			client := resty.New()
	    captchaType := ""
	    data := ""
	    if method == "hcaptcha" {
	    	captchaType = "sitekey"
	    } else {
	    	captchaType = "googlekey"
	    }
	    if p == "" {
	    	data = `{"key": "` + apiKey + `", "method":"` + method + `", "` + captchaType + `":"` + siteKey + `", ` + `"pageurl":` + `"` + url + `"}`

	    } else {
	    	data = `{"key": "` + apiKey + `", "method":"` + method + `", "` + captchaType + `":"` + siteKey + `", ` + `"pageurl":` + `"` + url + `", ` + `"userAgent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36", "min_score": 0.9, "proxy": "` + p + `", "proxytype": "HTTP"}`
	    	//data = `{"key": "` + apiKey + `", "method":"` + method + `", "` + captchaType + `":"` + siteKey + `", ` + `"pageurl":` + `"` + url + `", ` + `"proxy": "` + proxy + `"}`
	    }
	    resp, err := client.R().
	    SetBody(data).
	    Post("http://2captcha.com/in.php")
	    if err != nil {
	    	fmt.Println(err)
	    }
    	captchaID := strings.Split(resp.String(), "|")[1]
    	var notCaptchaRecieved bool = true
    	for notCaptchaRecieved == true {
    		    time.Sleep(8 * time.Second)
    		    url := "http://2captcha.com/res.php?key=" + apiKey + "&action=get&id=" + captchaID
    		    resp2, err := client.R().
    		    Get(url)
    		    if err != nil {
    		    	fmt.Println(err)
    		    }
    		    if resp2.String() != "CAPCHA_NOT_READY" {
    		    	notCaptchaRecieved = false
    		    	return strings.Split(resp2.String(), "|")[1]
    		    }
    		    

    	}	
    	return "ERROR"

}

func Testing(job string, ch chan bool, wg *sync.WaitGroup, proxy string, website string, email string){
			//proxySplitted := strings.Split(proxy, ":")
			//user, pass, host, port := proxySplitted[2], proxySplitted[3], proxySplitted[0], proxySplitted[1]
			//proxyUrl, _ := url.Parse("http://" + user + ":" + pass + "@" + host + ":" + port)
			// As soon as the current goroutine finishes (job done!), notify back WaitGroup.
			defer wg.Done()
			ch <- true
	    client := http.Client{Transport: &http.Transport{
    		//Proxy: http.ProxyURL(proxyUrl),
  		}}
	    rand.Seed(time.Now().UnixNano())

	   // waiters := []int{5, 6, 7, 8, 9}
	   // time.Sleep(time.Duration(rand.Intn(len(waiters))) * time.Second)
	   	time.Sleep(5 * time.Second)
	    req, _ := http.NewRequest("GET", "https://httpbin.org/get", nil)
	    resp, _ := client.Do(req)
	    body4, _ := ioutil.ReadAll(resp.Body)
	    fmt.Println(job, string(body4), proxy, website, email)
			<- ch
}

func Testing2(job string, ch chan bool, wg *sync.WaitGroup, returnData chan map[string]string) {
			defer wg.Done()
			ch <- true
			fmt.Println(job, "in")
			a := make(map[string]string)
			a["email"] = job
			a["Status"] = "Good"
	    waiters := []int{5, 6, 7, 8, 9}
	    time.Sleep(time.Duration(waiters[rand.Intn(len(waiters))]) * time.Second)
	    fmt.Println(job, "out")
	   	returnData <- a
			<- ch

}