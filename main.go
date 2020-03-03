package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"os"
	"time"

	//"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/aws/awserr"
	//"github.com/aws/aws-sdk-go/aws/session"
)

type Section struct {
	ActivityTitle    string `json:"activityTitle"`
	ActivitySubtitle string `json:"activitySubtitle"`
	ActivityImage    string `json:"activityImage"`
	Facts            []Fact `json:"facts"`
	Markdown bool `json:"markdown"`
}

type Fact struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type PotentialActionInput struct {
	Type        string `json:"@type"`
	ID          string `json:"id"`
	IsMultiline bool   `json:"isMultiline"`
	Title       string `json:"title"`
}
type PotentialActionAction struct {
	Type   string `json:"@type"`
	Name   string `json:"name"`
	Target string `json:"target"`
}

type PotentialAction struct {
	Type   string `json:"@type"`
	Name   string `json:"name"`
	Inputs []PotentialActionInput `json:"inputs"`
	Actions []PotentialActionAction `json:"actions"`
}

type Message struct {
	Type       string `json:"@type"`
	Context    string `json:"@context"`
	ThemeColor string `json:"themeColor"`
	Summary    string `json:"summary"`
	Sections   []Section `json:"sections"`
	PotentialActions []PotentialAction `json:"potentialAction"`
}

type LimitMessage struct {
	Version    string        `json:"version"`
	ID         string        `json:"id"`
	DetailType string        `json:"detail-type"`
	Source     string        `json:"source"`
	Account    string        `json:"account"`
	Time       time.Time     `json:"time"`
	Region     string        `json:"region"`
	Resources  []interface{} `json:"resources"`
	Detail     struct {
		CheckName       string `json:"check-name"`
		CheckItemDetail struct {
			Status       string `json:"Status"`
			CurrentUsage string `json:"Current Usage"`
			LimitName    string `json:"Limit Name"`
			Region       string `json:"Region"`
			Service      string `json:"Service"`
			LimitAmount  string `json:"Limit Amount"`
		} `json:"check-item-detail"`
		Status     string `json:"status"`
		ResourceID string `json:"resource_id"`
		UUID       string `json:"uuid"`
	} `json:"detail"`
}

func sendToWebhook(chatApplication string,webhook string,message Message) error {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(message)
	res, err := http.Post(webhook, "application/json; charset=utf-8", b)
	if err != nil {
		return err
	}
	fmt.Print(res.Status)
	return nil
}

func handler(ctx context.Context, event events.SNSEvent) {
	var chatApplication = os.Getenv("Chat_Application")
	var webhookURL = os.Getenv("Chat_Webhook")
	var limitMessage LimitMessage
	if len(event.Records) > 0 {
		for _, record := range event.Records {
			fmt.Print(record.SNS.Message)
			err := json.Unmarshal([]byte(record.SNS.Message),&limitMessage)
			if err != nil {
				fmt.Print(err)
			}
			var message Message
			message.Type = "MessageCard"
			message.Context = "http://schema.org/extensions"
			message.ThemeColor = "0076D7"
			var section Section
			message.Summary = limitMessage.Detail.CheckItemDetail.LimitName
			section.Markdown = true
			section.ActivityTitle = "Limit Monitor"
			var fact Fact
			fact.Name = "Service"
			fact.Value = limitMessage.Detail.CheckItemDetail.Service
			section.Facts = append(section.Facts,fact)
			fact.Name = "CurrentValue"
			fact.Value = limitMessage.Detail.CheckItemDetail.CurrentUsage
			section.Facts = append(section.Facts,fact)
			fact.Name = "LimitAmount"
			fact.Value = limitMessage.Detail.CheckItemDetail.LimitAmount
			section.Facts = append(section.Facts,fact)
			var potentialAction PotentialAction
			var potentialActionInput PotentialActionInput
			var potentialActionAction PotentialActionAction
			potentialAction.Type = "ActionCard"
			potentialAction.Name = "More details"
			potentialActionAction.Type = "OpenUri"
			potentialActionAction.Target = "https://console.aws.amazon.com/trustedadvisor/home?region=us-east-1#/category/service-limits"
			potentialAction.Inputs = append(potentialAction.Inputs,potentialActionInput)
			potentialAction.Actions = append(potentialAction.Actions,potentialActionAction)
			message.PotentialActions = append(message.PotentialActions,potentialAction)
			message.Sections = append(message.Sections,section)
			err = sendToWebhook(chatApplication,webhookURL,message)
			if err != nil {
				fmt.Print(err)
			}
		}
	}
}

func main() {
	lambda.Start(handler)
}