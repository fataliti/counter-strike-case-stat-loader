package main

import (
	"fmt"
	"image/color"
)

type ItemType string

const (
	CSGO_Type_Pistol      ItemType = "CSGO_Type_Pistol"
	CSGO_Type_Rifle       ItemType = "CSGO_Type_Rifle"
	CSGO_Type_SMG         ItemType = "CSGO_Type_SMG"
	CSGO_Type_SniperRifle ItemType = "CSGO_Type_SniperRifle"
	CSGO_Type_Shotgun     ItemType = "CSGO_Type_Shotgun"
	CSGO_Type_Machinegun  ItemType = "CSGO_Type_Machinegun"
	CSGO_Type_Knife       ItemType = "CSGO_Type_Knife"
	CSGO_Type_C4          ItemType = "CSGO_Type_C4"
	CSGO_Type_Grenade     ItemType = "CSGO_Type_Grenade"
	CSGO_Type_Equipment   ItemType = "CSGO_Type_Equipment"
)

type Item struct {
	Id      string
	Color   int
	Type    ItemType
	Title   string
	Date    string
	Rarity  string
	IconUrl string
}

func (item Item) GetColorStruct() color.RGBA {
	r := (int)((item.Color >> 16) & 255)
	g := (int)((item.Color >> 8) & 255)
	b := (int)((item.Color) & 255)

	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
}

func (item Item) GetIconURl() string {
	return fmt.Sprintf("https://community.fastly.steamstatic.com/economy/image/%s/480x160", item.IconUrl)
}
