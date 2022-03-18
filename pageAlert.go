package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/antchfx/htmlquery"
	"gopkg.in/toast.v1"
)

// survey config
// the questions to ask
var qs = []*survey.Question{
	{
		Name:     "url",
		Prompt:   &survey.Input{Message: "Enter url of site:"},
		Validate: survey.Required,
	},
	{
		Name:   "timePeriod",
		Prompt: &survey.Input{Message: "enter time period in seconds(optional):"},
	},
	{
		Name:   "xpath",
		Prompt: &survey.Input{Message: "enter Xpath of target element (optional):"},
	},
}

func main() {
	// cli code
	answers := struct {
		URL        string
		TimePeriod string
		Xpath      string
	}{}

	// perform the questions
	err := survey.Ask(qs, &answers, survey.WithShowCursor(true))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// extracting input values from survey answers and assigning default values
	url := answers.URL
	xpath := answers.Xpath
	if xpath == "" {
		xpath = "//body"
	}
	timePeriod, err := strconv.Atoi(answers.TimePeriod)
	if err != nil {
		timePeriod = 10
	}

	// initial call to getContent
	Ogcontent := getContent(url, xpath)
	for {
		time.Sleep(time.Second * time.Duration(timePeriod))
		newContent := getContent(url, xpath)
		if newContent != Ogcontent {
			pushNotification(url)
		}
		println("..")
	}

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
			{"protocol", "go to site", ""},
			{"protocol", "Dismiss", ""},
		},
	}
	err := notification.Push()
	if err != nil {
		log.Fatalln(err)
	}
}
