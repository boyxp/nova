package main

import "time"
import "fmt"
import "strings"

func main() {
	fmt.Println(Date("H:i:s", time.Now().Add(time.Hour * 3)))
}

func Date(format string, _time ...time.Time) string {
	var currentTime time.Time

	if len(_time)==0 {
		currentTime = time.Now()
	} else {
		currentTime = _time[0]
	}

	if format=="Y-m-d H:i:s" {
		return currentTime.Format("2006.01.02 15:04:05")
	}

	if format=="Y-m-d" {
		return currentTime.Format("2006.01.02")
	}

	if format=="H:i:s" {
		return currentTime.Format("15:04:05")
	}

	for i:=0;i<len(format);i++ {
		switch string(format[i]) {
			case "Y":
					format = strings.Replace(format, "Y", currentTime.Format("2006"), 1)
			case "m":
					format = strings.Replace(format, "m", currentTime.Format("01"), 1)
			case "d":
					format = strings.Replace(format, "d", currentTime.Format("02"), 1)
			case "H":
					format = strings.Replace(format, "H", currentTime.Format("15"), 1)
			case "i":
					format = strings.Replace(format, "i", currentTime.Format("04"), 1)
			case "s":
					format = strings.Replace(format, "s", currentTime.Format("05"), 1)
			case "y":
					format = strings.Replace(format, "y", currentTime.Format("06"), 1)
			case "D":
					format = strings.Replace(format, "D", currentTime.Format("Mon"), 1)
			case "j":
					format = strings.Replace(format, "j", currentTime.Format("2"), 1)
			case "l":
					format = strings.Replace(format, "l", currentTime.Format("Monday"), 1)
			case "N":
					var weekDayMap = map[string]string {
					    "Monday":    "1",
					    "Tuesday":   "2",
					    "Wednesday": "3",
					    "Thursday":  "4",
					    "Friday":    "5",
				    	"Saturday":  "6",
			    		"Sunday":    "7",
					}

					var weekDay = currentTime.Format("Monday")
					format = strings.Replace(format, "N", weekDayMap[weekDay], 1)
			case "w":
					var weekDayMap = map[string]string {
					    "Monday":    "0",
					    "Tuesday":   "1",
					    "Wednesday": "2",
					    "Thursday":  "3",
					    "Friday":    "4",
				    	"Saturday":  "5",
			    		"Sunday":    "6",
					}

					var weekDay = currentTime.Format("Monday")
					format = strings.Replace(format, "w", weekDayMap[weekDay], 1)
			case "F":
					format = strings.Replace(format, "F", currentTime.Format("June"), 1)
			case "M":
					format = strings.Replace(format, "M", currentTime.Format("Jan"), 1)
			case "n":
					format = strings.Replace(format, "n", currentTime.Format("1"), 1)
			case "a":
					format = strings.Replace(format, "a", currentTime.Format("pm"), 1)
			case "A":
					format = strings.Replace(format, "A", currentTime.Format("PM"), 1)
			case "h":
					format = strings.Replace(format, "h", currentTime.Format("3"), 1)
			case "c":
					format = strings.Replace(format, "c", currentTime.Format(time.RFC3339), 1)
		}
	}

	return format
}
