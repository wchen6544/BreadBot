package main

import (
    "errors"
    "fmt"
    "reflect"
    footlocker "main/gleam/Footlocker"
    struc "main/testFolder"
    "os"
    "encoding/json"
    "io"
    "sync"
    "strings"
    rand "math/rand"
    "time"
      "github.com/denisbrodbeck/machineid"
  "github.com/tidwall/gjson"
  "github.com/go-resty/resty/v2"
  "log"
  "strconv"

)

type stubMapping map[string]interface{}

var (
    wg sync.WaitGroup
    StubStorage = stubMapping{}
)



func main() {
    colorRed := "\033[31m"
    prox := getProxies()
    emai := getEmails()
    if len(prox) < len(emai) {
        log.Fatalln(string(colorRed), strconv.Itoa(len(emai) - len(prox)) + " additional proxies required to initialize")
    } else {
        log.Println("Beginning")
        begin()
    }

}


func begin() {
    breadBot := `
____                     _ _           _   
|  _ \                   | | |         | |  
| |_) |_ __ ___  __ _  __| | |__   ___ | |_ 
|  _ <| '__/ _ \/ _' |/ _' | '_ \ / _ \| __|
| |_) | | |  __/ (_| | (_| | |_) | (_) | |_ 
|____/|_|  \___|\__,_|\__,_|_.__/ \___/ \__|


[1] Footlocker
[2] Footlocker Acc Gen
[3] Footlocker Acc Confirm
    `
    colorReset := "\033[0m"
    colorRed := "\033[31m"
    colorGreen := "\033[32m"
    //colorYellow := "\033[33m"
    colorCyan := "\033[36m"

    if validLicense() {

    // --------------------

        rand.Seed(time.Now().UnixNano())
        Fails := []string{}
        Success := []string{}
        fmt.Println(string(colorCyan), breadBot + string(colorReset))
        fmt.Print("Enter your site selection: ")
        var userInput string
        fmt.Scanln(&userInput)
        jsonFile, err := os.Open("config.json")
        if err != nil {
            fmt.Println(err)
        }
        defer jsonFile.Close()

        byteValue, _ := io.ReadAll(jsonFile)

        var result map[string]interface{}
        json.Unmarshal([]byte(byteValue), &result)
        site := getSiteData(getSite(userInput), result)
        workerPool := make(chan bool, site.Tasks)
        resultData := make(chan map[string]string)
        fmt.Println("BEGIN")
        var emailList []string = getEmails()  // gets the list of emails
        var proxyList []string = getProxies() // gets the list of proxies
        
        wg.Add(len(emailList))

        StubStorage = map[string]interface{}{
            "Footlocker": footlocker.Start,
            "Footlocker Acc Gen": footlocker.Gen,
            "Footlocker Acc Confirm": footlocker.Confirm,
        }

        for taskNumber, email := range emailList {
            mm := struc.Necessary{email, proxyList[taskNumber], resultData, &wg, workerPool, taskNumber}
            go call(getSite(userInput), site, mm)
        }


        for i := 0; i < len(emailList); i++ {
            returnData := <-resultData
            if returnData["status"] == "Bad" {
                fmt.Println(string(colorRed) + "Failed to Enter: ", returnData["email"] + string(colorReset))
                Fails = append(Fails, returnData["email"])
            } else {
                fmt.Println(string(colorGreen) + "Successfully Entered: ", returnData["email"] + string(colorReset))
                Success = append(Success, returnData["email"])
            }
        }
        wg.Wait()
        fmt.Println("END")
        fmt.Print(strings.Join(Fails[:], "\n"))
        fmt.Println("\n__________________________")
        fmt.Println(strings.Join(Success[:], "\n"))
        begin()
    }
}

func call(funcName string, params ... interface{}) (err error) {
    f := reflect.ValueOf(StubStorage[funcName])
    if len(params) != f.Type().NumIn() {
        err = errors.New("Error with Task Manager")
        return
    }
    in := make([]reflect.Value, len(params))
    for k, param := range params {
        in[k] = reflect.ValueOf(param)
    }
    f.Call(in)
    return
}

func getSiteData(site string, result map[string]interface{}) struc.SiteStructure {
    var siteData struc.SiteStructure
    defaultJsonBody, _ := json.Marshal(result["Sites"].(map[string]interface{})[site])
    err3 := json.Unmarshal(defaultJsonBody, &siteData)
    if err3 != nil { fmt.Println(err3) }
    return siteData
}

func getEmails() []string{
    var emailList []string
    emailFile, _ := os.Open("emails.txt")
    emails, _ := io.ReadAll(emailFile)
    for _, e := range strings.Split(string(emails),"\n") {
        emailList = append(emailList, e)
    }
    return emailList
}

func getProxies() []string{
    var proxyList []string
    proxiesFile, _ := os.Open("proxies.txt")
    proxy, _ := io.ReadAll(proxiesFile)
    for _, e := range strings.Split(string(proxy),"\n") {
        listProxy := strings.Split(e, ":")
        finalE := "http://" + listProxy[2] + ":" + listProxy[3] + "@" + listProxy[0] + ":" + listProxy[1]
        proxyList = append(proxyList, finalE)
    }
    return proxyList
}

func getGeneralData(t string, configData string) string{
    switch t {
    case "license": return gjson.Get(configData, "General.licenseKey").String()
    case "2CaptchaKey": return gjson.Get(configData, "General.2CaptchaKey").String()
    case "discordWebhook": return gjson.Get(configData, "General.discordWebhook").String()
    }
    return ""
}

func getSite(userI string) string{

    array := [3]string{"Footlocker", "Footlocker Acc Gen", "Footlocker Acc Confirm"}

    selected, _ := strconv.Atoi(userI)


    return array[selected - 1]





}

func validLicense() bool{
    jsonFile, err := os.Open("config.json")
    if err != nil {
        fmt.Println(err)
    }
    defer jsonFile.Close()

    byteValue, _ := io.ReadAll(jsonFile)

    id, err := machineid.ID()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(id)
    client := resty.New()
    resp, err := client.R().
    SetFormData(map[string]string{
        "key": getGeneralData("license", string(byteValue[:])),
        "uuid": "6BE50D76-0B2C-5950-B3D4-C98578D93A9B",
    }).
    Post("https://randomthhinglicense.herokuapp.com/license")
    if resp.String() == "Success" {
        return true
    }
    return false
}