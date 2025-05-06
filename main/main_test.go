package main

import (
	"testing"
)

func TestEvents(t *testing.T) {

	b := Biatlon{make(map[int]TableRow, 4)}

	var tests = []struct {
		name   string
		input  string
		output string
	}{
		{"Event 1", "[09:05:59.867] 1 1", "[09:05:59.867] The competitor(1) registered"},
		{"Event 2", "[09:15:00.841] 2 1 09:30:00.000", "[09:15:00.841] The start time for the competitor(1) was set by a draw to 09:30:00.000"},
		{"Event 3", "[09:29:45.734] 3 1", "[09:29:45.734] The competitor(1) is on the start line"},
		{"Event 4", "[09:30:01.005] 4 1", "[09:30:01.005] The competitor(1) has started"},
		{"Event 5", "[09:49:31.659] 5 1 2", "[09:49:31.659] The competitor(1) is on the firing range(2)"},
		{"Event 6", "[09:49:33.123] 6 1 2", "[09:49:33.123] The target(2) has been hit by competitor(1)"},
		{"Event 7", "[09:49:38.339] 7 1", "[09:49:38.339] The competitor(1) left the firing range"},
		{"Event 8", "[09:49:55.915] 8 1", "[09:49:55.915] The competitor(1) entered the penalty laps"},
		{"Event 9", "[09:51:48.391] 9 1", "[09:51:48.391] The competitor(1) left the penalty laps"},
		{"Event 10", "[09:59:03.872] 10 1", "[09:59:03.872] The competitor(1) ended the main lap"},
		{"Event 11", "[09:59:03.872] 11 1 Lost in the forest", "[09:59:03.872] The competitor(1) can`t continue: Lost in the forest"},
		{"Event 32", "[09:59:03.873] 32 1", "[09:59:03.873] The competitor(1) is disqualified"},
		{"Event 33", "[09:59:03.874] 33 1", "[09:59:03.874] The competitor(1) has finished"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans, err := b.EventHandler(tt.input)
			if err != nil || ans != tt.output {
				t.Errorf(`EventHandler("`+tt.input+`") = %q, %v, want match for %#q, nil`, ans, err, tt.output)
			}
		})
	}

}

func TestResultFunc1(t *testing.T) {
	b := Biatlon{make(map[int]TableRow, 4)}

	b.EventHandler("[09:05:59.867] 1 1")
	b.EventHandler("[09:15:00.841] 2 1 09:30:00.000")
	b.EventHandler("[09:29:45.734] 3 1")
	b.EventHandler("[09:30:01.005] 4 1")
	b.EventHandler("[09:49:31.659] 5 1 2")
	b.EventHandler("[09:49:33.123] 6 1 2")
	b.EventHandler("[09:49:38.339] 7 1")
	b.EventHandler("[09:49:55.915] 8 1")
	b.EventHandler("[09:51:48.391] 9 1")
	b.EventHandler("[09:59:03.872] 10 1")
	b.EventHandler("[09:59:03.872] 11 1 Lost in the forest")

	config := Config{2, 3651, 50, 1, "09:30:00", "00:00:30"}
	res := b.GetResult(config)
	if res[0] != "[NotFinished] 1 [{00:29:03.872, 2.093}, {,}] {00:01:52.476, 0.444} 1/5" {
		t.Errorf(`Result = %v, want match for [NotFinished] 1 [{00:29:03.872, 2.093}, {,}] {00:01:52.476, 0.444} 1/5`, res[0])
	}
}

func TestResultFunc2(t *testing.T) {
	b := Biatlon{make(map[int]TableRow, 4)}

	b.EventHandler("[09:05:59.867] 1 1")
	b.EventHandler("[09:15:00.841] 2 1 09:30:00.000")
	b.EventHandler("[09:29:45.734] 3 1")
	b.EventHandler("[09:30:01.005] 4 1")
	b.EventHandler("[09:49:31.659] 5 1 1")
	b.EventHandler("[09:49:33.123] 6 1 1")
	b.EventHandler("[09:49:33.123] 6 1 2")
	b.EventHandler("[09:49:33.123] 6 1 3")
	b.EventHandler("[09:49:33.123] 6 1 4")
	b.EventHandler("[09:49:33.123] 6 1 5")
	b.EventHandler("[09:49:38.339] 7 1")
	b.EventHandler("[10:30:00.000] 10 1")
	b.EventHandler("[10:30:00.000] 33 1")

	config := Config{1, 3651, 50, 1, "09:30:00", "00:00:30"}
	res := b.GetResult(config)
	if res[0] != "[01:00:00.000] 1 [{01:00:00.000, 1.014}] {00:00:00.000, 0.000} 5/5" {
		t.Errorf(`Result = %v, want match for [01:00:00.000] 1 [{01:00:00.000, 1.014}] {00:00:00.000, 0.000} 5/5`, res[0])
	}
}
