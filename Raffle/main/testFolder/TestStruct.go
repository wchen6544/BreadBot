package testFolder

import (
    "sync"
)

type SiteStructure struct {
    Url  string `json:"URL"`
    Tasks int    `json:"Tasks"`
    Size string `json:"Size"`
    Name  string `json:"Name"`
    Store string `json:"Store"`
    Release string `json:"Release"`
    Shoe string `json:"Shoe"`
}

type Necessary struct {
    Email string `json:"Email"`
    Proxy string `json:"Proxy"`
    ReturnData chan map[string]string `json:"ReturnData"`
    WaitGroup *sync.WaitGroup `json:"WaitGroup"`
    CH chan bool `json:"CH"`
    TaskNumber int `json"TaskNumber"`
}