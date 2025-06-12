package main

import (
	"fmt"
	"image/color"
	"log"

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

	cookie string

	is_on_main_screen bool = true
	is_input_cookie   bool
	is_in_process     bool
	is_watch_stat     bool

	is_show_error_msg bool
	error_message     string

	stat_result []StatResult
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic catch: %v\n%s\n", r, StackTrace(3))
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

	window := giu.NewMasterWindow("CS case open stat", 640, 480, giu.MasterWindowFlagsFloating)
	window.Run(func() {
		if is_on_main_screen {
			intro_screen()

		} else {
			program_screen()
		}
	})

}

func intro_screen() {
	giu.SingleWindow().Layout(
		giu.Align(giu.AlignCenter).To(
			giu.Label("Cookie?"),
			giu.InputText(&cookie),
			giu.Button("Begin").OnClick(func() {
				is_on_main_screen = false
				is_in_process = true
				go RequestData(cookie)
			}),
			giu.Button("Cancel").OnClick(func() {
				is_on_main_screen = false
			}),
		),
	)
}

func program_screen() {
	rows := make([]*giu.TableRowWidget, 0)
	for i := 0; i < len(ItemList); i++ {
		rows = append(rows, giu.TableRow(
			giu.Layout{
				giu.Label(ItemList[i].Date),
				giu.Custom(func() {
					giu.SameLine()
				}),
				giu.Align(giu.AlignRight).To(giu.ImageWithFile("./look.png")),
				giu.Tooltip("").Layout(giu.ImageWithURL(ItemList[i].GetIconURl())),
			},

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
