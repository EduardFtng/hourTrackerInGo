package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	layout = "15:04:05"
)

type ContentSlice struct {
	Content []WorkDay `json:"data"`
}

type WorkDay struct {
	Date    string `json:"date"`
	Perc100 bool   `json:"100_perc"`
	Start   string `json:"start"`
	End     string `json:"end"`
}

func openJSONfile(path string) ContentSlice {
	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var data ContentSlice

	json.Unmarshal(byteValue, &data)

	return data

}

func printJSONData(data ContentSlice) {
	fmt.Print("\n")
	for i := 0; i < len(data.Content); i++ {
		fmt.Println("Date: " + data.Content[i].Date)
		fmt.Println("Sunday/Holiday: " + strconv.FormatBool(data.Content[i].Perc100))
		fmt.Println("Start: " + data.Content[i].Start)
		fmt.Println("End: " + data.Content[i].End + "\n")
	}
}

func calcAllHours(data ContentSlice) time.Duration {
	var sum time.Duration

	for i := 0; i < len(data.Content); i++ {
		start, _ := time.Parse(layout, data.Content[i].Start)
		end, _ := time.Parse(layout, data.Content[i].End)
		diff := end.Sub(start)
		sum += diff
	}
	return sum
}

func calcHoursTill20(data ContentSlice) time.Duration {
	var sum time.Duration
	limit, _ := time.Parse(layout, "20:00:00")

	for i := 0; i < len(data.Content); i++ {
		start, _ := time.Parse(layout, data.Content[i].Start)
		end, _ := time.Parse(layout, data.Content[i].End)
		extra := data.Content[i].Perc100

		if extra == false {
			if end.Before(limit) {
				diff := end.Sub(start)
				sum += diff
			} else {
				diff := limit.Sub(start)
				sum += diff
			}
		}
	}
	return sum
}

func calcHoursPast20(data ContentSlice) time.Duration {
	var sum time.Duration
	limit, _ := time.Parse(layout, "20:00:00")

	for i := 0; i < len(data.Content); i++ {
		end, _ := time.Parse(layout, data.Content[i].End)
		extra := data.Content[i].Perc100

		if extra == false {
			if end.After(limit) {
				diff := end.Sub(limit)
				sum += diff
			}
		}
	}
	return sum
}

func calcExtraHours(data ContentSlice) time.Duration {

	var sum time.Duration

	for i := 0; i < len(data.Content); i++ {
		start, _ := time.Parse(layout, data.Content[i].Start)
		end, _ := time.Parse(layout, data.Content[i].End)
		extra := data.Content[i].Perc100

		if extra == true {
			diff := end.Sub(start)
			sum += diff
		}
	}
	return sum
}

func printAllCalcHours(fp string) {
	fmt.Print("\n")
	fmt.Println("All hours till 8PM: ", calcHoursTill20(openJSONfile(fp)))
	fmt.Println("All hours after 8PM: ", calcHoursPast20(openJSONfile(fp)))
	fmt.Println("All +100% paid hours: ", calcExtraHours(openJSONfile(fp)))
	fmt.Println("Sum of all worked hours: ", calcAllHours(openJSONfile(fp)))
	fmt.Println("*******************************")
	fmt.Print("\n")
}

func main() {
	//filePath := "info.json"

	fmt.Println("****Welcome To Time Tracker****")
	reader := bufio.NewReader(os.Stdin)

	var filePath string
outer: // Label to break out from the outer loop
	for {
		if filePath == "" {
			fmt.Print("Enter the name of your JSON file: ")
			filePath, _ = reader.ReadString('\n')
			filePath = strings.Replace(filePath, "\n", "", -1) // Removing \n from the string to get a clean string
		} else {
			for {
				fmt.Print("Enter 'd' to see the data or enter 'c' to calculate the hours or 'q' to quit: ")
				text, _ := reader.ReadString('\n')
				text = strings.Replace(text, "\n", "", -1)
				if text == "d" {
					printJSONData(openJSONfile(filePath))
				} else if text == "c" {
					printAllCalcHours(filePath)
				} else if text == "q" {
					fmt.Println("Exiting.....")
					break outer
				} else {
					fmt.Println("Enter a valid option.")
				}
			}
		}
	}
}
