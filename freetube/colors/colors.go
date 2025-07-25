// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

// https://github.com/FreeTubeApp/FreeTube/blob/1f36288ae1548c44f2be590f98484b1c1b7aba5d/src/renderer/helpers/colors.js
package colors

import (
	"math/rand/v2"
)

type Color struct {
	Name, Value string
}

func (c Color) ColorLuminance() string {
	cutHex := c.Value[1:7]
	colorValueR := fromHex(cutHex[0], cutHex[1])
	colorValueG := fromHex(cutHex[2], cutHex[3])
	colorValueB := fromHex(cutHex[4], cutHex[5])
	luminance := (0.299*float64(colorValueR) + 0.587*float64(colorValueG) + 0.114*float64(colorValueB)) / 255
	if luminance > 0.5 {
		return "#000000"
	} else {
		return "#FFFFFF"
	}
}

func RandomColor() Color {
	return colors[rand.IntN(len(colors))] // #nosec G404
}

var colors = []Color{
	{Name: "Red", Value: "#d50000"},
	{Name: "Pink", Value: "#C51162"},
	{Name: "Purple", Value: "#AA00FF"},
	{Name: "DeepPurple", Value: "#6200EA"},
	{Name: "Indigo", Value: "#304FFE"},
	{Name: "Blue", Value: "#2962FF"},
	{Name: "LightBlue", Value: "#0091EA"},
	{Name: "Cyan", Value: "#00B8D4"},
	{Name: "Teal", Value: "#00BFA5"},
	{Name: "Green", Value: "#00C853"},
	{Name: "LightGreen", Value: "#64DD17"},
	{Name: "Lime", Value: "#AEEA00"},
	{Name: "Yellow", Value: "#FFD600"},
	{Name: "Amber", Value: "#FFAB00"},
	{Name: "Orange", Value: "#FF6D00"},
	{Name: "DeepOrange", Value: "#DD2C00"},
	{Name: "CatppuccinFrappeRosewater", Value: "#f2d5cf"},
	{Name: "CatppuccinFrappeFlamingo", Value: "#eebebe"},
	{Name: "CatppuccinFrappePink", Value: "#f4b8e4"},
	{Name: "CatppuccinFrappeMauve", Value: "#ca9ee6"},
	{Name: "CatppuccinFrappeRed", Value: "#e78284"},
	{Name: "CatppuccinFrappeMaroon", Value: "#ea999c"},
	{Name: "CatppuccinFrappePeach", Value: "#ef9f76"},
	{Name: "CatppuccinFrappeYellow", Value: "#e5c890"},
	{Name: "CatppuccinFrappeGreen", Value: "#a6d189"},
	{Name: "CatppuccinFrappeTeal", Value: "#81c8be"},
	{Name: "CatppuccinFrappeSky", Value: "#99d1b"},
	{Name: "CatppuccinFrappeSapphire", Value: "#85c1dc"},
	{Name: "CatppuccinFrappeBlue", Value: "#8caaee"},
	{Name: "CatppuccinFrappeLavender", Value: "#babbf1"},
	{Name: "CatppuccinMochaRosewater", Value: "#F5E0DC"},
	{Name: "CatppuccinMochaFlamingo", Value: "#F2CDCD"},
	{Name: "CatppuccinMochaPink", Value: "#F5C2E7"},
	{Name: "CatppuccinMochaMauve", Value: "#CBA6F7"},
	{Name: "CatppuccinMochaRed", Value: "#F38BA8"},
	{Name: "CatppuccinMochaMaroon", Value: "#EBA0AC"},
	{Name: "CatppuccinMochaPeach", Value: "#FAB387"},
	{Name: "CatppuccinMochaYellow", Value: "#F9E2AF"},
	{Name: "CatppuccinMochaGreen", Value: "#A6E3A1"},
	{Name: "CatppuccinMochaTeal", Value: "#94E2D5"},
	{Name: "CatppuccinMochaSky", Value: "#89DCEB"},
	{Name: "CatppuccinMochaSapphire", Value: "#74C7EC"},
	{Name: "CatppuccinMochaBlue", Value: "#89B4FA"},
	{Name: "CatppuccinMochaLavender", Value: "#B4BEFE"},
	{Name: "DraculaCyan", Value: "#8BE9FD"},
	{Name: "DraculaGreen", Value: "#50FA7B"},
	{Name: "DraculaOrange", Value: "#FFB86C"},
	{Name: "DraculaPink", Value: "#FF79C6"},
	{Name: "DraculaPurple", Value: "#BD93F9"},
	{Name: "DraculaRed", Value: "#FF5555"},
	{Name: "DraculaYellow", Value: "#F1FA8C"},
	{Name: "EverforestarkRed", Value: "#E67E80"},
	{Name: "EverforestarkOrange", Value: "#E69875"},
	{Name: "EverforestarkYellow", Value: "#bBC7F"},
	{Name: "EverforestarkGreen", Value: "#A7C080"},
	{Name: "EverforestarkAqua", Value: "#83C092"},
	{Name: "EverforestarkBlue", Value: "#7FBBB3"},
	{Name: "EverforestarkPurple", Value: "#D699B6"},
	{Name: "EverforestLightRed", Value: "#D83532"},
	{Name: "EverforestLightOrange", Value: "#D55D0F"},
	{Name: "EverforestLightYellow", Value: "#A96E00"},
	{Name: "EverforestLightGreen", Value: "#6D8100"},
	{Name: "EverforestLightAqua", Value: "#25976C"},
	{Name: "EverforestLightBlue", Value: "#2a84b5"},
	{Name: "EverforestLightPurple", Value: "#CF59aa"},
	{Name: "GruvboxarkGreen", Value: "#b8bb26"},
	{Name: "GruvboxarkYellow", Value: "#fabd2f"},
	{Name: "GruvboxarkBlue", Value: "#83a593"},
	{Name: "GruvboxarkPurple", Value: "#d3869b"},
	{Name: "GruvboxarkAqua", Value: "#8ec07c"},
	{Name: "GruvboxarkOrange", Value: "#fe8019"},
	{Name: "GruvboxLightRed", Value: "#9d0006"},
	{Name: "GruvboxLightBlue", Value: "#076678"},
	{Name: "GruvboxLightPurple", Value: "#8f3f71"},
	{Name: "GruvboxLightOrange", Value: "#af3a03"},
	{Name: "SolarizedYellow", Value: "#b58900"},
	{Name: "SolarizedOrange", Value: "#cb4b16"},
	{Name: "SolarizedRed", Value: "#dc322f"},
	{Name: "SolarizedMagenta", Value: "#d33682"},
	{Name: "SolarizedViolet", Value: "#6c71c4"},
	{Name: "Solarizeblue", Value: "#268bd2"},
	{Name: "SolarizedCyan", Value: "#2aa198"},
	{Name: "SolarizedGreen", Value: "#859900"},
}
