package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Songmu/prompter"
	"github.com/antchfx/htmlquery"
	"gopkg.in/toast.v1"
)

func main() {
	// cli code
	answers := struct {
		URL        string
		TimePeriod string
		Xpath      string
	}{}

	// perform questions with prompter
	answers.URL = prompter.Prompt("❔ Enter URL of the site", "")
	answers.TimePeriod = prompter.Prompt("❔ Enter time period in seconds (optional)", "5")
	answers.Xpath = prompter.Prompt("❔ Enter Xpath of target element (optional)", "//body")

	// extracting input values from survey answers and assigning default values
	url := answers.URL
	xpath := answers.Xpath
	timePeriod, err := strconv.Atoi(answers.TimePeriod)
	if err != nil {
		timePeriod = 5
	}

	// initial call to getContent
	Ogcontent := getContent(url, xpath)
	arr := [4]rune{'|', '/', '-', '\\'}
	var index int
	index = 0
	for {
		time.Sleep(time.Second * time.Duration(timePeriod))
		newContent := getContent(url, xpath)
		if newContent != Ogcontent {
			pushNotification(url)
			break
		}
		fmt.Printf("\b%c", arr[index])
		index = index + 1
		index = index % 4
	}
	fmt.Printf("\bexiting..")
}

func getContent(url string, xpath string) string {
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		panic(err)
	}
	s := htmlquery.FindOne(doc, xpath)
	return htmlquery.InnerText(s)
}

func pushNotification(url string) {

	notification := toast.Notification{
		AppID:   "pageAlert",
		Title:   "Content Changed",
		Message: "the content in " + url + " has been updated",
		// Icon: "go.png", // This file must exist (remove this line if it doesn't)
		Actions: []toast.Action{
			{Type: "protocol", Label: "go to site", Arguments: url},
			{Type: "protocol", Label: "Dismiss", Arguments: ""},
		},
	}
	err := notification.Push()
	if err != nil {
		log.Fatalln(err)
	}
}
