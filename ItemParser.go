package main

import (
	"strconv"

	"github.com/fatih/color"
)

func ParseItems(items []Item, app_descptions AppDescriptions) {
	for i := 0; i < len(items); i++ {
		item_struct := items[i]
		item_id := item_struct.Id

		item := app_descptions["730"][item_id]
		var tag_len = len(item.Tags)

		if tag_len < 5 {
			continue
		}

		item_struct.Type = ItemType(item.Tags[0].InternalName)
		item_struct.Title = item.Name

		color_hex := item.Tags[4].Color
		color_int, err := strconv.ParseInt(color_hex, 16, 32)
		var r, g, b int
		if err == nil {
			r = (int)((color_int >> 16) & 255)
			g = (int)((color_int >> 8) & 255)
			b = (int)((color_int) & 255)

			color.RGB(r, g, b).Println(item.Name)
			item_struct.Color = int(color_int)
			item_struct.Rarity = item.Tags[4].Name

			DataChan <- item_struct
		}
	}

}
