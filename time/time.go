package time

import "time"
import "strings"
import "strconv"
import "regexp"
import "github.com/carmo-evan/strtotime"

func Now() time.Time {
	return time.Now()
}

func Strtotime(str string) time.Time {
	reg := regexp.MustCompile(`([+-])([0-9]+)\s*(year|month|day|hour|minute|second|week)`)
	if reg == nil {
		panic("MustCompile err")
	}

	result := reg.FindAllStringSubmatch(str, -1)
	if len(result)>0 {
		base := time.Unix(time.Now().Unix(), 0)
		for _, match := range result {
			count, _ := strconv.Atoi(match[2])
			if match[1]=="-" {
				count = 0 - count
			}

			switch match[3] {
				case "year":
							base = addDate(base, count, 0, 0)
				case "month":
							base = addDate(base, 0, count, 0)
				case "day":
							base = addDate(base, 0, 0, count)
				case "week":
							base = addDate(base, 0, 0, count*7)
				case "hour":
							base = base.Add(time.Duration(count) * time.Hour)
				case "minute":
							base = base.Add(time.Duration(count) * time.Minute)
				case "second":
							base = base.Add(time.Duration(count) * time.Second)
			}
		}

		return base
	}

	u, err := strtotime.Parse(str+" +0800", time.Now().Unix())
	if err != nil {
		panic("strtotime时间格式化失败:"+str)
	}

	return time.Unix(u,0)
}

//https://libuba.com/2021/01/08/golan-adddate-%E4%B8%B4%E7%95%8C%E5%9D%91/
func addDate(t time.Time, year int, month int, day int) time.Time {
	targetDate := t.AddDate(year, month, -t.Day()+1)
	targetDay  := targetDate.AddDate(0, 1, -1).Day()

	if targetDay > t.Day() {
		targetDay = t.Day()
	}

	targetDate = targetDate.AddDate(0, 0, targetDay-1+day)
	return targetDate
}

func Date(format string, _time ...time.Time) string {
	var currentTime time.Time

	if len(_time)==0 {
		currentTime = time.Now()
	} else {
		currentTime = _time[0]
	}

	switch format {
		case "Y-m-d H:i:s":
							return currentTime.Format("2006-01-02 15:04:05")
		case "Y-m-d"      :
							return currentTime.Format("2006-01-02")
		case "Y-m"        :
							return currentTime.Format("2006-01")
		case "Y"          :
							return currentTime.Format("2006")
		case "H:i:s"      :
							return currentTime.Format("15:04:05")
		case "H:i"      :
							return currentTime.Format("15:04")
		case "H"      :
							return currentTime.Format("15")
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
					    "Monday"   : "1",
					    "Tuesday"  : "2",
					    "Wednesday": "3",
					    "Thursday" : "4",
					    "Friday"   : "5",
				    	"Saturday" : "6",
			    		"Sunday"   : "7",
					}

					var weekDay = currentTime.Format("Monday")
					format = strings.Replace(format, "N", weekDayMap[weekDay], 1)
			case "w":
					var weekDayMap = map[string]string {
					    "Monday"   : "0",
					    "Tuesday"  : "1",
					    "Wednesday": "2",
					    "Thursday" : "3",
					    "Friday"   : "4",
				    	"Saturday" : "5",
			    		"Sunday"   : "6",
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

func Timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func Readable(timestamp string) string {
	now := time.Now()

	curr := Strtotime(timestamp)

	diff := now.Sub(curr).Seconds()

	if diff>7*86400 {
		return timestamp
	}

	str := now.Sub(curr).String()

	if strings.Contains(str, "h") {
		h := str[0:strings.Index(str, "h")]
		n, _ := strconv.Atoi(h)
		if n<24 {
			return h+"小时前"
		}

		return strconv.Itoa(n/24)+"天前"
	}

	if strings.Contains(str, "m") {
		return str[0:strings.Index(str, "m")]+"分钟前"
	}

	if strings.Contains(str, "s") {
		return str[0:strings.Index(str, ".")]+"秒前"
	}

	return str
}
