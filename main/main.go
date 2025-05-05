package main

import (
	"errors"
	"strconv"
	"strings"
)

type TableRow struct {
	FailMark string
	Start    string
	Finish   string

	StartLap string
	LapTimes []struct {
		LapTime string
		AvgTime string
	}

	PenaltyStart   string
	PenaltyTime    string
	AvgPenaltyTime string

	Shots int
	Hits  int
}

type Biatlon struct {
	table map[int]TableRow
}

func (b *Biatlon) EventHandler(msg string) (string, error) {
	words := strings.Fields(msg)

	res := words[0] + " "

	id, err := strconv.Atoi(words[2])
	if err != nil {
		return "in EventHandler competitor id not conv", err
	}

	row, have := b.table[id]
	if have == false && words[1] != "1" {
		return "", errors.New("competitor is undefined")
	}

	switch words[1] {
	case "1":
		{
			b.table[id] = TableRow{}
			return res + "The competitor(" + string(id) + ") registered", nil
		}
	case "2":
		{
			row.Start = words[3]
			row.StartLap = words[3]
			return res + "The start time for the competitor(" + string(id) + ") was set by a draw to " + words[3], nil
		}
	case "3":
		{
			return res + "The competitor(" + string(id) + ") is on the start line", nil
		}
	case "4":
		{

			return res + "The competitor(" + string(id) + ") has started", nil
		}
	case "5":
		{
			row.Shots += 5
			return res + "The competitor(" + string(id) + ") is on the firing range(" + words[3] + ")", nil
		}
	case "6":
		{
			row.Hits++
			return res + "The target(" + words[3] + ") has been hit by competitor(" + string(id) + ")", nil
		}
	case "7":
		{
			return res + "The competitor(" + string(id) + ") left the firing range", nil
		}
	case "8":
		{
			row.PenaltyStart = words[0][1 : len(words[0])-1]
			return res + "The competitor(" + string(id) + ") entered the penalty laps", nil
		}
	case "9":
		{
			// count time now - row.PenaltyStart and add it to row.PenaltyTime
			return res + "The competitor(" + string(id) + ") left the penalty laps", nil
		}
	case "10":
		{
			// count new struct in row.LapTimes
			row.StartLap = words[0][1:len(words[0])]
			return res + "The competitor(" + string(id) + ") ended the main lap", nil
		}
	case "11":
		{
			row.FailMark = "NotFinished"
			tmpRes := res + "The competitor(" + string(id) + ") can`t continue:"
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

}
