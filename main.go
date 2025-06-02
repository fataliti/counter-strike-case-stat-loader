package main

import (
	"fmt"
	"image/color"
	"log"
	"runtime/debug"

	"github.com/AllenDang/giu"
)

type StatResult struct {
	Color  color.RGBA
	Title  string
	Amount int
}

const (
	ParseComplete    int = 0
	UploadStartError int = 1
)

var (
	DataChan   chan Item   = make(chan Item)
	EventsChan chan int    = make(chan int)
	ErrorChan  chan string = make(chan string)
	ItemList   []Item

	is_input_cookie bool
	cookie          string

	is_in_process bool
	is_watch_stat bool

	is_show_error_msg bool
	error_message     string

	stat_result []StatResult
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic catch: %v\n%s\n", r, debug.Stack())
		}
	}()

	go func() {
		for {
			data, ok := <-DataChan
			if ok {
				ItemList = append(ItemList, data)
				sum_index := -1

				for i := 0; i < len(stat_result); i++ {
					if stat_result[i].Title == data.Rarity {
						sum_index = i
						break
					}
				}

				if sum_index >= 0 {
					stat_result[sum_index].Amount += 1
				} else {
					stat_result = append(stat_result, StatResult{
						Color:  data.GetColorStruct(),
						Title:  data.Rarity,
						Amount: 1,
					})
				}
			}
		}
	}()

	go func() {
		for {
			event, ok := <-EventsChan
			if ok {
				println("event: ", event)

				switch event {
				case ParseComplete:
					is_in_process = false
					is_watch_stat = true
					fmt.Println("parse complete")

				case UploadStartError:
					is_in_process = false
					fmt.Println("error")
				}
			}
		}
	}()

	go func() {
		for {
			err, ok := <-ErrorChan
			if ok {
				is_show_error_msg = true
				error_message = err
			}
		}
	}()

	// reader := bufio.NewReader(os.Stdin)
	// input_cookie, input_err := reader.ReadString('\n')
	// if input_err != nil {
	// 	log.Fatal(input_err)
	// }
	// input_cookie = input_cookie[:len(input_cookie)-1]
	// input_cookie = strings.TrimSpace(input_cookie)
	// RequestData("recentlyVisitedAppHubs=1857090; timezoneOffset=10800,0; browserid=15310613048948578; sessionid=2ba0912b213d5f32bcfdcfa3; steamDidLoginRefresh=1748819715; steamCountry=RU%7Cbf50c1b444a4661436cc9f399260fbd0; steamLoginSecure=76561198061430462%7C%7CeyAidHlwIjogIkpXVCIsICJhbGciOiAiRWREU0EiIH0.eyAiaXNzIjogInI6MDAwQl8yNjU2QTUyQl8zMjUzNSIsICJzdWIiOiAiNzY1NjExOTgwNjE0MzA0NjIiLCAiYXVkIjogWyAid2ViOmNvbW11bml0eSIgXSwgImV4cCI6IDE3NDg5MDc3MDgsICJuYmYiOiAxNzQwMTc5NzE1LCAiaWF0IjogMTc0ODgxOTcxNSwgImp0aSI6ICIwMDBBXzI2NjM4MTkwXzBGRUVBIiwgIm9hdCI6IDE3NDc4NTQzMTAsICJydF9leHAiOiAxNzY1OTU4NTkwLCAicGVyIjogMCwgImlwX3N1YmplY3QiOiAiMjEyLjE1NC4yMTIuNDciLCAiaXBfY29uZmlybWVyIjogIjIxMi45Ni43NS4yMDEiIH0.X59TVBChbxdWyvCZy-dYh1SE3hwyx1Q_-NkcboLKMkJ19KW3AH0Ir8Zyq5fKz0zVmmVeVPJjoAdbIt-mk7-qDg")
	window := giu.NewMasterWindow("CS case open stat", 640, 480, giu.MasterWindowFlagsFloating)
	window.Run(func() {
		rows := make([]*giu.TableRowWidget, 0)
		for i := 0; i < len(ItemList); i++ {
			rows = append(rows, giu.TableRow(
				giu.Label(ItemList[i].Date),
				giu.Layout{
					giu.Custom(func() {
						giu.PushStyleColor(giu.StyleColorText, ItemList[i].GetColorStruct())
					}),
					giu.Label(ItemList[i].Title),
					giu.Custom(func() {
						giu.PopStyleColor()
					}),
				},
			))
		}

		giu.SingleWindowWithMenuBar().Layout(
			giu.MenuBar().Layout(
				giu.Menu("Menu").Layout(
					giu.MenuItem("Start").OnClick(func() {
						if !is_in_process {
							go RequestData(cookie)
							is_in_process = true
						}
					}),

					giu.MenuItem("Input cookie").OnClick(func() {
						is_input_cookie = true
					}),

					giu.MenuItem("watch stat").OnClick(func() {
						is_watch_stat = true
					}),
				),
			),

			giu.Table().Columns(
				giu.TableColumn("Date").Flags(giu.TableColumnFlagsWidthFixed).InnerWidthOrWeight(120),
				giu.TableColumn("Name"),
			).Rows(rows...),

			giu.PrepareMsgbox(),
			giu.Custom(func() {
				if is_show_error_msg {
					giu.Msgbox("Error", error_message)
					is_show_error_msg = false
				}
			}),
		)

		if is_input_cookie {
			giu.Window("Cookie?").IsOpen(&is_input_cookie).Flags(giu.WindowFlagsNoResize|giu.WindowFlagsNoDocking).Size(420, 120).Layout(
				giu.Align(giu.AlignCenter).To(
					giu.InputText(&cookie).Label("Cookie"),
					giu.Button("ok").OnClick(func() {
						is_input_cookie = false
					}),
				),
			)
		}

		if is_watch_stat {
			giu.Window("Stat").IsOpen(&is_watch_stat).Flags(giu.WindowFlagsNoDocking).Size(240, 400).Layout(
				get_labels()...,
			)
		}
	})

}

func get_labels() []giu.Widget {
	labels := make([]giu.Widget, 0)
	for i := 0; i < len(stat_result); i++ {
		stat := stat_result[i]

		labels = append(labels, giu.Layout{
			giu.Custom(func() {
				giu.PushStyleColor(giu.StyleColorText, stat.Color)
			}),
			giu.Label(fmt.Sprintf("%s: %d", stat.Title, stat.Amount)),
			giu.Custom(func() {
				giu.PopStyleColor()
			}),
		})
	}

	return labels
}
