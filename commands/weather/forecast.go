package weather

import (
	"fmt"
	"github.com/0x263b/porygon2"
	"github.com/0x263b/porygon2/web"
	"time"
)

func forecast(command *bot.Cmd, matches []string) (msg string, err error) {

	var location string = matches[1]
	var coords string

	if location == "" {
		location, coords = checkLocation(command.Nick)
		if coords == "" {
			return "Location not provided, nor on file. Use `-set location <location>` to save", nil
		}
	} else {
		coords = getCoords(location)
		if coords == "" {
			return fmt.Sprintf("Could not find %s", location), nil
		}
	}

	data := &Forecast{}
	err = web.GetJSON(fmt.Sprintf(DarkSkyURL, bot.Config.Weather, coords), data)
	if err != nil {
		return fmt.Sprintf("Could not get weather for: %s", location), nil
	}

	units := "°C"
	if data.Flags.Units == "us" {
		units = "°F"
	}

	output := fmt.Sprintf("Forecast | %s ", location)

	for i := range data.Daily.Data[0:4] {
		tm := time.Unix(data.Daily.Data[i].Time, 0)
		loc, _ := time.LoadLocation(data.Timezone)
		day := tm.In(loc).Weekday()
		output += fmt.Sprintf("| %s: %s %v%s/%v%s ",
			day,
			Emoji(data.Daily.Data[i].Icon),
			Round(data.Daily.Data[i].TemperatureMax),
			units,
			Round(data.Daily.Data[i].TemperatureMin),
			units,
		)
	}

	return output, nil
}

func init() {
	bot.RegisterCommand(
		"^f(?:o(?:recast)?)?(?: (.+))?$",
		forecast)
}
