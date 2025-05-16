package main

import (
	"fmt"
	"net/http"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"encoding/json"

	// "io"
	// "os"

	"github.com/fatih/color"
	"strconv"
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
	reqv.Header.Add("Cookie", "timezoneOffset=10800,0; recentlyVisitedAppHubs=1138850%2C730%2C1361210%2C1643320%2C2142790; browserid=96369691529769043; strInventoryLastContext=730_2; steamCountry=RU%7Cbf50c1b444a4661436cc9f399260fbd0; sessionid=ea5097acf1061188568bde20; steamLoginSecure=76561198061430462%7C%7CeyAidHlwIjogIkpXVCIsICJhbGciOiAiRWREU0EiIH0.eyAiaXNzIjogInI6MDAwOV8yNjIwNkNEQ19EMDBGQyIsICJzdWIiOiAiNzY1NjExOTgwNjE0MzA0NjIiLCAiYXVkIjogWyAid2ViOmNvbW11bml0eSIgXSwgImV4cCI6IDE3NDc0MjQ3MTAsICJuYmYiOiAxNzM4Njk3NzUyLCAiaWF0IjogMTc0NzMzNzc1MiwgImp0aSI6ICIwMDBBXzI2NEVDQTQ3XzIyODczIiwgIm9hdCI6IDE3NDQ0Njg1NjgsICJydF9leHAiOiAxNzYyMzgwMTMwLCAicGVyIjogMCwgImlwX3N1YmplY3QiOiAiNS4xOC4yNTMuMTk4IiwgImlwX2NvbmZpcm1lciI6ICIxMDQuMjM4LjI5LjI0NyIgfQ.y3R5i9WYQSW9IhdPrT7FxTOlm1cvtsb6jny7OpoR2Bk6qnXRy0eRGr96ATId9c6pbyq5HREQZTTn5ib24-bWBg")

	resp, err := client.Do(reqv)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close() 

    // io.Copy(os.Stdout, resp.Body)
	// fmt.Print(resp.Body)
	
	// content, err := os.ReadFile("test.html")
	
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	history_selection := doc.Find(".tradehistoryrow")
	count_of_elements := history_selection.Length()
	fmt.Printf("\nelements %d\n", count_of_elements)

	item_list := []Item{}

	history_selection.Each(func(i int, s *goquery.Selection) {
		event_text := s.Find(".tradehistory_event_description").Text();
		
		if strings.Contains(event_text, "Unlocked a container") {
			println("conteiner event")
			s.Find(".tradehistory_items_withimages:contains('+')").Each(func(i int, s *goquery.Selection) {
				

				s.Find("[data-classid]").Each(func(i int, s *goquery.Selection) {
					// data_app_id, exist := s.Attr("data-appid")
					data_classid, exist := s.Attr("data-classid")
					data_instanceid, exist := s.Attr("data-instanceid")
					_ = exist
					// fmt.Println(data_classid, data_instanceid);
					
					var added_item = Item {data_classid + "_" + data_instanceid}
					item_list = append(item_list, added_item)
				})
			})
		} 
	})



	json_string := GetJsonString(doc)
	// println(len(json_string))
	// println(json_string)

	var data AppDescriptions
	err_ := json.Unmarshal([]byte(json_string), &data)
	if err_ != nil {
		log.Fatalf("Parse error %v:", err_)
	}



	for i := 0; i < len(item_list); i++ {
		item_id := item_list[i].Id;
		
		// Пример обращения к данным
		item := data["730"][item_id]
		// fmt.Println("Название предмета:", item)
		
		color_hex := item.Tags[4].Color
		color_int, err := strconv.ParseInt(color_hex, 16, 32)

		if (err == nil) {
			r := (int)((color_int >> 16) & 255);
			g := (int)((color_int >> 8) & 255);
			b := (int)((color_int) & 255);

			// fmt.Printf("%d %d %d %d", color_int, r, g, b)
			color.RGB(r, g, b).Println(item.Name);
		}
	}
}

type ItemDescription struct {
	// IconUrl string `json:"icon_url"`
	// IconDragUrl string `json:"icon_drag_url"`
	Name string `json:"name"`
	MarketHashName string `json:"market_hash_name"`
	MarketName string `json:"market_name"`
	NameColor string `json:"name_color"`
	BackgroundColor string `json:"background_color"`
	// Type:''
	Description []struct {
		Type string `json:"type"`
		Value string `json:"value"`
		Name string `json:"name"`
	} `json:"descriptions"`
	    Tags []struct {
        InternalName string `json:"internal_name"`
        Name         string `json:"name"`
        Category     string `json:"category"`
		Color		 string `json:"color"`
    } `json:"tags"`
}

type Item struct {
	Id string 
}

type AppDescriptions map[string]map[string]ItemDescription

func GetJsonString(document *goquery.Document) string {
	var scriptContent string
	document.Find("script").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.Contains(text, "g_rgDescriptions") {
			scriptContent = text
		}
	})

	jsonString := ""
	if scriptContent != "" {
		balance := 0 
		in_string := false 
		escape := false 
		
		start := strings.Index(scriptContent, "g_rgDescriptions = {")
		start_index := start + len("g_rgDescriptions = {")
		end_index := start_index 

		for end_index < len(scriptContent) {
			char := scriptContent[end_index]
			switch (char) {
				case '"':
					if !escape {
						in_string = !in_string 
					}
					escape = false
				case '\\':
					escape = !escape;
				case '}', ']':
					if (in_string) {
						balance -= 1
					}
				case '{', '[':
					if (in_string) {
						balance += 1
					}
			}
			

			if balance == 0 && !in_string {
				next_char := scriptContent[end_index + 1];
				if (next_char == ';') {
					jsonString = scriptContent[start_index-1:end_index+1];
					break;
				}
			}  

			if escape && char != '\\' {
				escape = false
			}

			end_index += 1
		}
	}

	return jsonString
}

