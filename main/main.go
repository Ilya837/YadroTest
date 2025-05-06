package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Laps        int    `json:"laps"`
	LapLen      int    `json:"lapLen"`
	PenaltyLen  int    `json:"penaltyLen"`
	FiringLines int    `json:"firingLines"`
	Start       string `json:"start"`
	StartDelta  string `json:"startDelta"`
}

type TableRow struct {
	FailMark string
	Start    time.Time
	Finish   time.Time

	StartLap time.Time
	LapTimes []time.Time

	PenaltyStart     time.Time
	TotalPenaltyTime time.Time
	PenaltyCount     uint

	Shots uint
	Hits  uint
}

type Biatlon struct {
	table map[int]TableRow
}

func (b *Biatlon) EventHandler(msg string) (string, error) {
	words := strings.Fields(msg)
	timeFormat := "15:04:05.000"

	res := words[0] + " "

	id, err := strconv.Atoi(words[2])
	if err != nil {
		return "in EventHandler competitor id not conv", err
	}

	row, have := b.table[id]
	defer func() { b.table[id] = row }()

	if have == false && words[1] != "1" {
		return "", errors.New("competitor is undefined")
	}

	var nowTime time.Time
	var timeDiff time.Duration

	switch words[1] {
	case "1":
		{
			b.table[id] = TableRow{}
			return res + "The competitor(" + strconv.Itoa(id) + ") registered", nil
		}
	case "2":
		{
			row.Start, err = time.Parse(timeFormat, words[3])

			if err != nil {
				return "in EventHandler competitor time not conv", err
			}

			row.StartLap = row.Start
			return res + "The start time for the competitor(" + strconv.Itoa(id) + ") was set by a draw to " + words[3], nil
		}
	case "3":
		{
			return res + "The competitor(" + strconv.Itoa(id) + ") is on the start line", nil
		}
	case "4":
		{

			return res + "The competitor(" + strconv.Itoa(id) + ") has started", nil
		}
	case "5":
		{
			row.Shots += 5
			return res + "The competitor(" + strconv.Itoa(id) + ") is on the firing range(" + words[3] + ")", nil
		}
	case "6":
		{
			row.Hits++
			return res + "The target(" + words[3] + ") has been hit by competitor(" + strconv.Itoa(id) + ")", nil
		}
	case "7":
		{
			return res + "The competitor(" + strconv.Itoa(id) + ") left the firing range", nil
		}
	case "8":
		{
			row.PenaltyStart, err = time.Parse(timeFormat, words[0][1:len(words[0])-1])

			if err != nil {
				return "in EventHandler competitor time not conv", err
			}

			row.PenaltyCount++

			return res + "The competitor(" + strconv.Itoa(id) + ") entered the penalty laps", nil
		}
	case "9":
		{
			nowTime, err = time.Parse(timeFormat, words[0][1:len(words[0])-1])

			if err != nil {
				return "in EventHandler competitor time not conv", err
			}

			timeDiff = nowTime.Sub(row.PenaltyStart)
			row.TotalPenaltyTime = row.TotalPenaltyTime.Add(timeDiff)
			return res + "The competitor(" + strconv.Itoa(id) + ") left the penalty laps", nil
		}
	case "10":
		{
			nowTime, err = time.Parse(timeFormat, words[0][1:len(words[0])-1])

			if err != nil {
				return "in EventHandler competitor time not conv", err
			}

			timeDiff = nowTime.Sub(row.StartLap)

			row.LapTimes = append(row.LapTimes, time.Time{}.Add(timeDiff))

			row.StartLap = nowTime
			return res + "The competitor(" + strconv.Itoa(id) + ") ended the main lap", nil
		}
	case "11":
		{
			row.FailMark = "NotFinished"
			tmpRes := res + "The competitor(" + strconv.Itoa(id) + ") can`t continue:"
			for i := 3; i < len(words); i++ {
				tmpRes += " " + words[i]
			}
			return tmpRes, nil
		}
	default:
		return "", errors.New("unknown event")
	}

}

func (b *Biatlon) GetResult(cfg Config) string {

	timeFormat := "15:04:05.000"

	Totalres := ""

	for k, v := range b.table {
		res := ""
		if v.FailMark != "" {
			res += "[" + v.FailMark + "] "
		} else {
			res += "[" + time.Time{}.Add(v.Finish.Sub(v.Start)).Format(timeFormat) + "] "
		}

		res += strconv.Itoa(k)
		res += " ["

		for i := 0; i < len(v.LapTimes); i++ {
			res += "{" + v.LapTimes[i].Format(timeFormat) + ", "

			tmp := strconv.FormatFloat((float64(cfg.LapLen) * 1000 / float64(v.LapTimes[i].Sub(time.Time{}).Milliseconds())), 'f', 4, 64)
			res += tmp[0 : len(tmp)-1]
			res += "}, "
		}

		for i := len(v.LapTimes); i < cfg.Laps; i++ {
			res += "{,}, "
		}

		res = res[0 : len(res)-2]

		if v.PenaltyCount == 0 {
			res += "] {" + time.Time{}.Format(timeFormat) + ", 0.000} "
		} else {

			penaltyAvgTime := time.Time{}.Add(time.Duration(int(float64(v.TotalPenaltyTime.Sub(time.Time{}).Nanoseconds()) / float64(v.PenaltyCount))))

			res = res[0 : len(res)-2]
			res += "] {" + penaltyAvgTime.Format(timeFormat) + ", "

			tmp := strconv.FormatFloat((float64(cfg.PenaltyLen) * 1000 / float64(penaltyAvgTime.Sub(time.Time{}).Milliseconds())), 'f', 4, 64)
			res += tmp[0 : len(tmp)-1]

			res += "} "

		}

		res += strconv.Itoa(int(v.Hits)) + "/" + strconv.Itoa(int(v.Shots))

		Totalres += res + "\n"
	}

	return Totalres
}

func main() {
	// input := "09:59:05.321"
	// input2 := "09:59:06.323"
	// t, _ := time.Parse("15:04:05.000", input)
	// t2, _ := time.Parse("15:04:05.000", input2)
	// fmt.Print(time.Time{}.Add(t2.Sub(t)).Format("15:04:05.000"))

	// configFile, err := os.Open("./sunny_5_skiers/config.json")
	// defer configFile.Close()
	// if err != nil {
	// 	fmt.Print(err)
	// }

	// config := Config{}

	// jsonParser := json.NewDecoder(configFile)
	// jsonParser.Decode(&config)

	// fmt.Print(config)
	fmt.Print(strconv.FormatFloat((float64(3651*1000000) / float64(29*60*1000000+3*1000000+872000)), 'f', 4, 64))
}
