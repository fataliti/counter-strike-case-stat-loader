package main

import (
	"fmt"
	"net/http"
	"io"
    "os"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	client := &http.Client{}
	reqv, err := http.NewRequest("GET", "https://steamcommunity.com/my/inventoryhistory", nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	reqv.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 7; WOW64; rv:51.0) Gecko/20100101 Firefox/51.0")
	reqv.Header.Add("Accept-Charset", "UTF-8") 
	reqv.Header.Add("Accept-Language", "en-US")
	reqv.Header.Add("Cookie", "timezoneOffset=10800,0; recentlyVisitedAppHubs=1138850%2C730%2C1361210%2C1643320%2C2142790; browserid=96369691529769043; steamLoginSecure=76561198061430462%7C%7CeyAidHlwIjogIkpXVCIsICJhbGciOiAiRWREU0EiIH0.eyAiaXNzIjogInI6MDAwOV8yNjIwNkNEQ19EMDBGQyIsICJzdWIiOiAiNzY1NjExOTgwNjE0MzA0NjIiLCAiYXVkIjogWyAid2ViOmNvbW11bml0eSIgXSwgImV4cCI6IDE3NDY5MTAxNTAsICJuYmYiOiAxNzM4MTgzNTI5LCAiaWF0IjogMTc0NjgyMzUyOSwgImp0aSI6ICIwMDBBXzI2NDQyRTcwXzZBMUI0IiwgIm9hdCI6IDE3NDQ0Njg1NjgsICJydF9leHAiOiAxNzYyMzgwMTMwLCAicGVyIjogMCwgImlwX3N1YmplY3QiOiAiNS4xOC4yNTMuMTk4IiwgImlwX2NvbmZpcm1lciI6ICIxMDQuMjM4LjI5LjI0NyIgfQ.2p65wfla1AmAq-23l8frbh8XX_cz9XDr2F4ipJs_5zhx8HC0GKN8g81s7T4Hp3-Tvzmpw8IxEl9GKvIRLFT5Cw; sessionid=48f57732001423a8af66502c")

	resp, err := client.Do(reqv)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close() 
    io.Copy(os.Stdout, resp.Body)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {

	}
	title := doc.Find("h1").Text()
	fmt.Println("Заголовок:", title)


}

type Item struct {
	ItemID string
	ItemName string
	itemDescription string
}