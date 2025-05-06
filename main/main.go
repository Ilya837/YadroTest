package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type TableRow struct {
	FailMark string
	Start    time.Time
	Finish   time.Time

	StartLap time.Time
	LapTimes []struct {
		LapTime  time.Time
		AvgSpeed float32
	}

	PenaltyStart     time.Time
	PenaltyTime      time.Time
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
			// count new struct in row.LapTimes
			nowTime, err = time.Parse(timeFormat, words[0][1:len(words[0])-1])

			if err != nil {
				return "in EventHandler competitor time not conv", err
			}

			timeDiff = nowTime.Sub(row.StartLap)

			row.LapTimes = append(row.LapTimes,
				struct {
					LapTime  time.Time
					AvgSpeed float32
				}{time.Time{}.Add(timeDiff), 0})

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

func main() {
	input := "09:59:05.321"
	input2 := "09:59:06.323"
	t, _ := time.Parse("15:04:05.000", input)
	t2, _ := time.Parse("15:04:05.000", input2)
	fmt.Print(time.Time{}.Add(t2.Sub(t)).Format("15:04:05.000"))
}
