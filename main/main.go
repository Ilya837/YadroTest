package main

import (
	"errors"
	"strings"
)

func EventHandler(msg string) (string, error) {
	words := strings.Fields(msg)

	res := words[0] + " "

	switch words[1] {
	case "1":
		{
			return res + "The competitor(" + words[2] + ") registered", nil
		}
	case "2":
		{
			return res + "The start time for the competitor(" + words[2] + ") was set by a draw to " + words[3], nil
		}
	case "3":
		{
			return res + "The competitor(" + words[2] + ") is on the start line", nil
		}
	case "4":
		{
			return res + "The competitor(" + words[2] + ") has started", nil
		}
	case "5":
		{
			return res + "The competitor(" + words[2] + ") is on the firing range(" + words[3] + ")", nil
		}
	case "6":
		{
			return res + "The target(" + words[3] + ") has been hit by competitor(" + words[2] + ")", nil
		}
	case "7":
		{
			return res + "The competitor(" + words[2] + ") left the firing range", nil
		}
	case "8":
		{
			return res + "The competitor(" + words[2] + ") entered the penalty laps", nil
		}
	case "9":
		{
			return res + "The competitor(" + words[2] + ") left the penalty laps", nil
		}
	case "10":
		{
			return res + "The competitor(" + words[2] + ") ended the main lap", nil
		}
	case "11":
		{
			tmpRes := res + "The competitor(" + words[2] + ") can`t continue:"
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
