package main

import (
	"fmt"
	"net/http"
	// "io"
    // "os"
	"log"
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
	reqv.Header.Add("Cookie", "timezoneOffset=10800,0; recentlyVisitedAppHubs=1138850%2C730%2C1361210%2C1643320%2C2142790; browserid=96369691529769043; steamLoginSecure=76561198061430462%7C%7CeyAidHlwIjogIkpXVCIsICJhbGciOiAiRWREU0EiIH0.eyAiaXNzIjogInI6MDAwMV8yNUZCQTIzMF85MUI3OSIsICJzdWIiOiAiNzY1NjExOTgwNjE0MzA0NjIiLCAiYXVkIjogWyAid2ViOmNvbW11bml0eSIgXSwgImV4cCI6IDE3NDY5OTY0NDMsICJuYmYiOiAxNzM4MjY5NjU2LCAiaWF0IjogMTc0NjkwOTY1NiwgImp0aSI6ICIwMDBBXzI2NDg0RDdGXzI4QkNCIiwgIm9hdCI6IDE3NDIxMjQ3ODgsICJydF9leHAiOiAxNzYwMzYzMDU4LCAicGVyIjogMCwgImlwX3N1YmplY3QiOiAiNS4xOC4yNTMuMTk4IiwgImlwX2NvbmZpcm1lciI6ICIyLjU2LjE3My40MSIgfQ.P2HcfLlLDaRUonVQ-8oC-8y0TZy9FmCFATmqcjaOOI5j710VdBWEPY7ddruab_DGZEoA0RBNBY35-Qee_UzBAg; sessionid=48f57732001423a8af66502c")

	resp, err := client.Do(reqv)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close() 

    // io.Copy(os.Stdout, resp.Body)
	// fmt.Print(resp.Body)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	selection := doc.Find("[data-classid]")
	count_of_elements := selection.Length()
	fmt.Printf("\nelements %d\n", count_of_elements)


	selection.Each(func(i int, s *goquery.Selection) {

		attr, exist := s.Attr("data-appid")
		if (exist) {
			span_text := s.Find("span").Text()
			if attr == "730" {
				fmt.Printf("cs title %s\n", span_text)
			} else {
				fmt.Printf("not cs title %s\n", span_text)
			}
		} 
	})




}

// type Item struct {
// 	ItemID string
// 	ItemName string
// 	itemDescription string
// }