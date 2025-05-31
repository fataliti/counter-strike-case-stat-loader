package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var DataChan chan Item

const RowHeight float32 = 5

func main() {
	DataChan = make(chan Item)

	a := app.New()
	w := a.NewWindow("CS case open list")

	control_menu := fyne.NewMenu("Menu",
		fyne.NewMenuItem("Start  parsing", func() {
			entry := widget.NewEntry()
			entry.SetPlaceHolder("Input cookie")
			entry.Resize(fyne.NewSize(300, 100))

			dialog.ShowCustomConfirm(
				"Cookie?",
				"ok",
				"no",
				entry,
				func(b bool) {
					if b {
						go RequestData(entry.Text)
					}
				},
				w,
			)
		}),
	)

	mainMenu := fyne.NewMainMenu(
		control_menu,
	)
	w.SetMainMenu(mainMenu)

	grid := container.NewVBox()

	scroll := container.NewScroll(grid)
	scroll.SetMinSize(fyne.NewSize(640, 480))

	go func() {
		for {
			data, ok := <-DataChan
			if ok {
				color_int := data.Color
				r := (int)((color_int >> 16) & 255)
				g := (int)((color_int >> 8) & 255)
				b := (int)((color_int) & 255)

				fyne.Do(func() {
					date_label := canvas.NewText(data.Date, color.White) //
					label := canvas.NewText(data.Title, color.RGBA{uint8(r), uint8(g), uint8(b), 255})

					left_side := container.NewStack(date_label)
					left_side.Resize(fyne.NewSize(200, RowHeight))
					right_side := container.NewStack(label)

					row := container.NewHBox(left_side, right_side)
					row.Resize(fyne.NewSize(640, RowHeight))

					grid.Add(row)
					grid.Refresh()
				})
			}
		}
	}()
	w.SetContent(scroll)
	w.Resize(fyne.NewSize(640, 480))
	w.ShowAndRun()

	// reader := bufio.NewReader(os.Stdin)
	// input_cookie, input_err := reader.ReadString('\n')
	// if input_err != nil {
	// 	log.Fatal(input_err)
	// }
	// input_cookie = input_cookie[:len(input_cookie)-1]
	// input_cookie = strings.TrimSpace(input_cookie)
	// RequestData("timezoneOffset=10800,0; recentlyVisitedAppHubs=1138850%2C730%2C1361210%2C1643320%2C2142790%2C1088710; browserid=96369691529769043; strInventoryLastContext=730_2; sessionid=69f27bab2c8c8f02bdc43788; steamDidLoginRefresh=1748616596; steamCountry=RU%7Cbf50c1b444a4661436cc9f399260fbd0; steamLoginSecure=76561198061430462%7C%7CeyAidHlwIjogIkpXVCIsICJhbGciOiAiRWREU0EiIH0.eyAiaXNzIjogInI6MDAwOV8yNjIwNkNEQ19EMDBGQyIsICJzdWIiOiAiNzY1NjExOTgwNjE0MzA0NjIiLCAiYXVkIjogWyAid2ViOmNvbW11bml0eSIgXSwgImV4cCI6IDE3NDg3MDM5NjAsICJuYmYiOiAxNzM5OTc2NTk2LCAiaWF0IjogMTc0ODYxNjU5NiwgImp0aSI6ICIwMDBBXzI2NUZFMTU3X0QyRkVGIiwgIm9hdCI6IDE3NDQ0Njg1NjgsICJydF9leHAiOiAxNzYyMzgwMTMwLCAicGVyIjogMCwgImlwX3N1YmplY3QiOiAiNS4xOC4yNTMuMTk4IiwgImlwX2NvbmZpcm1lciI6ICIxMDQuMjM4LjI5LjI0NyIgfQ.nDWRwajwl31LKMQ7lzsCYF3PmVfcY38O5m4Pox0Au4BIbeB1v-wyGeJ9huXHyYnlcR9oLg1wVFbZlGBOnfaPAA")
}
