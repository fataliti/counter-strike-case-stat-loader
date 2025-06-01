package main

import (
	"fmt"

	"github.com/AllenDang/giu"
)

var DataChan chan Item

func main() {
	DataChan = make(chan Item)

	go func() {
		for {
			data, ok := <-DataChan
			if ok {
				// color_int := data.Color
				// r := (int)((color_int >> 16) & 255)
				// g := (int)((color_int >> 8) & 255)
				// b := (int)((color_int) & 255)

				println(data.Title)
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

		rows := make([]*giu.TableRowWidget, 10)
		for i := 0; i < 10; i++ {
			rows[i] = giu.TableRow(
				giu.Label(fmt.Sprintf("%d", i)), // № (индекс)
				giu.Label("galil"),              // Значение
			)
		}

		giu.SingleWindowWithMenuBar().Layout(
			giu.MenuBar().Layout(
				giu.Menu("Menu").Layout(
					giu.MenuItem("Start").OnClick(func() {
						// println("clicked")

						go RequestData("recentlyVisitedAppHubs=1857090; timezoneOffset=10800,0; browserid=15310613048948578; sessionid=2ba0912b213d5f32bcfdcfa3; steamDidLoginRefresh=1748819715; steamCountry=RU%7Cbf50c1b444a4661436cc9f399260fbd0; steamLoginSecure=76561198061430462%7C%7CeyAidHlwIjogIkpXVCIsICJhbGciOiAiRWREU0EiIH0.eyAiaXNzIjogInI6MDAwQl8yNjU2QTUyQl8zMjUzNSIsICJzdWIiOiAiNzY1NjExOTgwNjE0MzA0NjIiLCAiYXVkIjogWyAid2ViOmNvbW11bml0eSIgXSwgImV4cCI6IDE3NDg5MDc3MDgsICJuYmYiOiAxNzQwMTc5NzE1LCAiaWF0IjogMTc0ODgxOTcxNSwgImp0aSI6ICIwMDBBXzI2NjM4MTkwXzBGRUVBIiwgIm9hdCI6IDE3NDc4NTQzMTAsICJydF9leHAiOiAxNzY1OTU4NTkwLCAicGVyIjogMCwgImlwX3N1YmplY3QiOiAiMjEyLjE1NC4yMTIuNDciLCAiaXBfY29uZmlybWVyIjogIjIxMi45Ni43NS4yMDEiIH0.X59TVBChbxdWyvCZy-dYh1SE3hwyx1Q_-NkcboLKMkJ19KW3AH0Ir8Zyq5fKz0zVmmVeVPJjoAdbIt-mk7-qDg")
					}),
				),
			),
			giu.Table().Columns(
				giu.TableColumn("Date"),
				giu.TableColumn("Name"),
			).Rows(rows...),
		)

	})
}
