package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
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
			if b.table == nil {
				b.table = make(map[int]TableRow, 8)
			}
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
	case "32":
		{

			row.FailMark = "NotStarted"
			return res + "The competitor(" + strconv.Itoa(id) + ") is disqualified", nil

		}
	case "33":
		{
			row.Finish, err = time.Parse(timeFormat, words[0][1:len(words[0])-1])

			if err != nil {
				return "in EventHandler competitor time not conv", err
			}

			return res + "The competitor(" + strconv.Itoa(id) + ") has finished", nil
		}
	default:
		return "", errors.New("unknown event")
	}

}

func (b *Biatlon) GetResult(cfg Config) []string {

	timeFormat := "15:04:05.000"

	Totalres := make([]string, 0)

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

			res += "] {" + penaltyAvgTime.Format(timeFormat) + ", "

			tmp := strconv.FormatFloat((float64(cfg.PenaltyLen) * 1000 / float64(penaltyAvgTime.Sub(time.Time{}).Milliseconds())), 'f', 4, 64)
			res += tmp[0 : len(tmp)-1]

			res += "} "

		}

		res += strconv.Itoa(int(v.Hits)) + "/" + strconv.Itoa(int(v.Shots))

		Totalres = append(Totalres, res)
	}

	return Totalres
}

func remove(slice []int, value int) []int {
	newSlice := []int{}
	for _, v := range slice {
		if v != value {
			newSlice = append(newSlice, v)
		}
	}

	return newSlice
}

func main() {

	configFile, err := os.Open("./sunny_5_skiers/config.json")
	defer configFile.Close()
	if err != nil {
		fmt.Print(err)
	}

	config := Config{}

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	tmpDelta, _ := time.Parse("15:04:05", config.StartDelta)
	deltaTime := tmpDelta.Sub(time.Time{})

	inputFile, err := os.Open("./sunny_5_skiers/events")
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)

	b := Biatlon{}
	disqualif := make([]int, 0)

	for scanner.Scan() {

		line := scanner.Text()
		words := strings.Fields(line)

		nowTime, _ := time.Parse("15:04:05.000", words[0][1:len(words[0])-1])
		nowTime = nowTime.AddDate(-1, 0, -1)
		disqualifDel := make([]int, 0)
		for i := 0; i < len(disqualif); i++ {

			if nowTime.Add(-deltaTime).After(b.table[disqualif[i]].Start) {
				sum := b.table[disqualif[i]].Start.Add(deltaTime)
				strId := strconv.Itoa(disqualif[i])
				res, err := b.EventHandler("[" + sum.Format("15:04:05.000") + "] 32 " + strId)
				if err != nil {
					fmt.Print(err)
					panic(err)
				}

				fmt.Println(res)

				disqualifDel = append(disqualifDel, disqualif[i])
			}

		}

		for _, v := range disqualifDel {
			disqualif = remove(disqualif, v)
		}

		res, err := b.EventHandler(line)

		if err != nil {
			fmt.Print(err)
			panic(err)
		}

		fmt.Println(res)

		if words[1] == "10" {
			if id, _ := strconv.Atoi(words[2]); len(b.table[id].LapTimes) == config.Laps {
				res, err := b.EventHandler(words[0] + " 33 " + words[2])
				if err != nil {
					fmt.Print(err)
					panic(err)
				}
				fmt.Println(res)
			}
		}

		if words[1] == "2" {
			id, _ := strconv.Atoi(words[2])
			disqualif = append(disqualif, id)
		}

		if words[1] == "4" {
			id, _ := strconv.Atoi(words[2])
			disqualif = remove(disqualif, id)
		}

	}

	result := b.GetResult(config)

	sort.Slice(result, func(i, j int) bool {
		wordsI := strings.Fields(result[i])
		wordsJ := strings.Fields(result[j])
		timeI, _ := time.Parse("15:04:05.000", wordsI[0][1:len(wordsI[0])-1])
		timeJ, _ := time.Parse("15:04:05.000", wordsJ[0][1:len(wordsI[0])-1])

		return timeI.Before(timeJ)
	})

	fmt.Println()
	for _, v := range result {
		fmt.Println(v)
	}

}
