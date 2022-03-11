package main

import (
  "fmt"
  "github.com/go-resty/resty/v2"
  "main/gleam"
  "strings"
)

func main() {
  
  p := "http://5ox4k84S:R7QgAPUr-AnsfOZyzeq@usa.pro.stellaproxies.com:12291"

  tokURL := "https://packershoes.com/account/login"
  token := gleam.GetCaptcha("eaef486b1c8d5aa4e20d3244d207be8e", "userrecaptcha", "6LcCR2cUAAAAANS1Gpq_mDIJ2pQuJphsSQaUEuc9", tokURL, p)
  fmt.Println(token)
  client := resty.New()
  //usa.pro.stellaproxies.com:12291:5ox4k84S:R7QgAPUr-AnsfOZyzeq

  client.SetProxy(p)
  resp, err := client.R().
    SetFormData(map[string]string{
      "form_type": "create_customer",
      "utf8": "\u2713",
      "customer[first_name]": "Jenniser",
      "customer[last_name]": "Lin",
      "customer[email]": "lewisjenniferisverycool@gmail.com",
      "customer[password]": "hoglund6",
      "customer[tags][locale]": "babl#enbabl#",
      "recaptcha-v3-token": token,
    }).
	SetHeader("authority", "packershoes.com").
	SetHeader("cache-control", "max-age=0").
	SetHeader("upgrade-insecure-requests", "1").
	SetHeader("origin", "https,//packershoes.com").
	SetHeader("content-type", "application/x-www-form-urlencoded").
	SetHeader("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36").
	SetHeader("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
	SetHeader("sec-gpc", "1").
	SetHeader("sec-fetch-site", "same-origin").
	SetHeader("sec-fetch-mode", "navigate").
	SetHeader("sec-fetch-user", "?1").
	SetHeader("sec-fetch-dest", "document").
	SetHeader("referer", "https,//packershoes.com/account/register").
	SetHeader("accept-language", "en-US,en;q=0.9").
  	Post("https://packershoes.com/account")

  if err != nil {
    fmt.Println(err)
  }

  if strings.Contains(resp.String(), "To continue, let us know you&#39;re not a robot") {
    fmt.Println("Proxy Ban / Captcha Error")
  } else {
    fmt.Println("Account Created, check for email")
  }
  

  
}
//`{"shopify_domain":"jimmyjazz1.myshopify.com","shopify_product_id":"6855481655503","shopify_customer_id":"5355038900431","shopify_customer_first_name":"Wilson","shopify_customer_last_name":"Chow","shopify_customer_birthday":"01-12-2001","shopify_customer_phone":"(646) 250-1922","shopify_customer_email":"upisdownflop@gmail.com","shopify_customer_accepts_marketing":"false","variant_sku":"11574772","location_ids":["1043","1058","1035"],"has_order_history":"false"}`