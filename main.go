package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"encoding/json"

	"strconv"

	"github.com/PuerkitoBio/goquery"

	"github.com/fatih/color"
)

func main() {
	client := &http.Client{}
	reqv, err := http.NewRequest("GET", "https://steamcommunity.com/my/inventoryhistory/?app[]=730", nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	cookie := "recentlyVisitedAppHubs=1857090; timezoneOffset=10800,0; browserid=15310613048948578; sessionid=026f1d5d16c8e2e323e1a8b2; steamCountry=KZ%7C292f287301b455e899be334af20f3756; steamLoginSecure=76561198061430462%7C%7CeyAidHlwIjogIkpXVCIsICJhbGciOiAiRWREU0EiIH0.eyAiaXNzIjogInI6MDAxN18yNjRFQ0JEOF83NEZGRCIsICJzdWIiOiAiNzY1NjExOTgwNjE0MzA0NjIiLCAiYXVkIjogWyAid2ViOmNvbW11bml0eSIgXSwgImV4cCI6IDE3NDc5NDIwMTgsICJuYmYiOiAxNzM5MjE0MzQ3LCAiaWF0IjogMTc0Nzg1NDM0NywgImp0aSI6ICIwMDBBXzI2NTZBNTJDXzIwNDFGIiwgIm9hdCI6IDE3NDc3NjUyNzQsICJydF9leHAiOiAxNzY2MDYyOTE2LCAicGVyIjogMCwgImlwX3N1YmplY3QiOiAiMjEyLjk2LjY2LjE2NyIsICJpcF9jb25maXJtZXIiOiAiMjEyLjk2LjY2LjE2NyIgfQ.GrWyLIE__YIE6OJe4c8mPC3QluIfDLOQMjTCyCUQzmqr25OduYYqnVvTRMYYOZAQyLh6bd79R_c6nrVPdVv_CQ"
	reqv.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 7; WOW64; rv:51.0) Gecko/20100101 Firefox/51.0")
	reqv.Header.Add("Accept-Charset", "UTF-8")
	reqv.Header.Add("Accept-Language", "en-US")
	reqv.Header.Add("Cookie", cookie)

	resp, err := client.Do(reqv)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	item_list := CollectOpenedItems(doc)

	json_string := GetJsonString("g_rgDescriptions", doc)
	var data AppDescriptions
	err_ := json.Unmarshal([]byte(json_string), &data)
	if err_ != nil {
		log.Fatalf("Parse error %v:", err_)
	}

	cursor_string := GetJsonString("g_historyCursor", doc)

	var cursor Cursor
	err__ := json.Unmarshal([]byte(cursor_string), &cursor)
	if err__ != nil {
		log.Fatalf("Parse error %v:", err__)
	}

	println("initial cursor: ", cursor_string)

	session_id := GetJsonString("g_sessionID", doc)
	steam_id := GetJsonString("g_steamID", doc)
	println("steam_id:", steam_id)
	println("session_id: ", session_id)

	PrintItems(item_list, data)
	// // GET
	// // https://steamcommunity.com/id/fatalitiq/inventoryhistory/?ajax=1&cursor[time]=1746124298&cursor[time_frac]=0&cursor[s]=20114846436&sessionid=be1aba1cfc7b84e5517e4fbd

	// for i := 0; i < 2; i++ {
	// 	println("load try", i)
	// 	MoreLoadRequest(&cursor, session_id, cookie)
	// 	time.Sleep(4 * time.Second)
	// }

	request_count := 0
	for is_loop := true; is_loop; {
		// println("load try", request_count)
		if !MoreLoadRequest(&cursor, session_id, cookie) {
			break
		}
		request_count += 1
		time.Sleep((3 + time.Duration(rand.Float64())) * time.Second)
	}

	println("complete")

}

type Cursor struct {
	Time     int    `json:"time"`
	TimeFrac int    `json:"time_frac"`
	S        string `json:"s"`
}

type ItemDescription struct {
	// IconUrl string `json:"icon_url"`
	// IconDragUrl string `json:"icon_drag_url"`
	Name            string `json:"name"`
	MarketHashName  string `json:"market_hash_name"`
	MarketName      string `json:"market_name"`
	NameColor       string `json:"name_color"`
	BackgroundColor string `json:"background_color"`
	// Type:''
	Description []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
		Name  string `json:"name"`
	} `json:"descriptions"`
	Tags []struct {
		InternalName string `json:"internal_name"`
		Name         string `json:"name"`
		Category     string `json:"category"`
		Color        string `json:"color"`
	} `json:"tags"`
}

type AppDescriptions map[string]map[string]ItemDescription

type UpdateHistory struct {
	Success      bool            `json:"success"`
	Html         string          `json:"html"`
	Descriptions AppDescriptions `json:"descriptions"`
	NewCursor    Cursor          `json:"cursor"`
	Num          int             `json:"num"`
}

func CollectOpenedItems(doc *goquery.Document) []Item {
	item_list := []Item{}

	history_selection := doc.Find(".tradehistoryrow")
	// count_of_elements := history_selection.Length()
	// fmt.Printf("elements %d\n", count_of_elements)

	history_selection.Each(func(i int, s *goquery.Selection) {
		event_text := s.Find(".tradehistory_event_description").Text()

		if strings.Contains(event_text, "Unlocked a container") {
			s.Find(".tradehistory_items_withimages:contains('+')").Each(func(i int, s *goquery.Selection) {
				s.Find("[data-classid]").Each(func(i int, s *goquery.Selection) {
					data_classid, exist := s.Attr("data-classid")
					data_instanceid, exist := s.Attr("data-instanceid")
					_ = exist
					var added_item = Item{data_classid + "_" + data_instanceid}
					item_list = append(item_list, added_item)
				})
			})
		}
	})

	return item_list
}

func PrintItems(items []Item, app_descptions AppDescriptions) {
	for i := 0; i < len(items); i++ {
		item_id := items[i].Id

		item := app_descptions["730"][item_id]

		// print(item.Name, ": ")

		var tag_len = len(item.Tags)

		if tag_len < 5 {
			continue
		}
		color_hex := item.Tags[4].Color
		color_int, err := strconv.ParseInt(color_hex, 16, 32)
		var r, g, b int
		if err == nil {
			r = (int)((color_int >> 16) & 255)
			g = (int)((color_int >> 8) & 255)
			b = (int)((color_int) & 255)
			color.RGB(r, g, b).Println(item.Name)
		}
	}
}

func MoreLoadRequest(cursor *Cursor, session_id string, cookie string) bool {
	params := url.Values{}
	params.Add("ajax", "1")
	params.Add("cursor[time]", strconv.Itoa(cursor.Time))
	params.Add("cursor[time_frac]", strconv.Itoa(cursor.TimeFrac))
	params.Add("cursor[s]", cursor.S)
	params.Add("sessionid", session_id)
	params.Add("app[]", "730")
	fullURL := "https://steamcommunity.com/id/fatalitiq/inventoryhistory/?" + params.Encode()
	// println(fullURL)
	new_request, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	new_request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 7; WOW64; rv:51.0) Gecko/20100101 Firefox/51.0")
	new_request.Header.Add("Accept-Charset", "UTF-8")
	new_request.Header.Add("Accept-Language", "en-US")
	new_request.Header.Add("Cookie", cookie)

	new_client := &http.Client{}
	new_resp, new_resp_err := new_client.Do(new_request)
	if new_resp_err != nil {
		log.Fatal(new_resp_err)
	}

	defer new_resp.Body.Close()

	var new_data UpdateHistory

	new_body, err := io.ReadAll(new_resp.Body)
	err_2 := json.Unmarshal(new_body, &new_data)
	if err_2 != nil {
		log.Fatal(err_2)
	}

	new_doc, doc_err := goquery.NewDocumentFromReader(strings.NewReader(new_data.Html))
	if doc_err != nil {
		log.Fatal(doc_err)
	}

	new_item_list := CollectOpenedItems(new_doc)
	PrintItems(new_item_list, new_data.Descriptions)

	if new_data.NewCursor.S == "" && new_data.NewCursor.Time == 0 && new_data.NewCursor.TimeFrac == 0 {
		return false
	}
	cursor.S = new_data.NewCursor.S
	cursor.Time = new_data.NewCursor.Time
	cursor.TimeFrac = new_data.NewCursor.TimeFrac

	return true
}
