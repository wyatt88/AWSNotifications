package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	//"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/aws/awserr"
	//"github.com/aws/aws-sdk-go/aws/session"
)

type Section struct {
	ActivityTitle    string `json:"activityTitle"`
	ActivitySubtitle string `json:"activitySubtitle"`
	ActivityImage    string `json:"activityImage"`
	Facts            []Fact `json:"facts"`
	Markdown         bool   `json:"markdown"`
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
	Type    string                  `json:"@type"`
	Name    string                  `json:"name"`
	Inputs  []PotentialActionInput  `json:"inputs"`
	Targets []PotentialActionTarget `json:"targets"`
	Actions []PotentialActionAction `json:"actions"`
}

type PotentialActionTarget struct {
	OS  string `json:"os"`
	URI string `json:"uri"`
}

type TeamsMessage struct {
	Type             string            `json:"@type"`
	Context          string            `json:"@context"`
	ThemeColor       string            `json:"themeColor"`
	Summary          string            `json:"summary"`
	Sections         []Section         `json:"sections"`
	PotentialActions []PotentialAction `json:"potentialAction"`
}

type Detail struct {
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
	EventArn          string `json:"eventArn"`
	Service           string `json:"service"`
	EventTypeCode     string `json:"eventTypeCode"`
	EventTypeCategory string `json:"eventTypeCategory"`
	StartTime         string `json:"startTime"`
	EventDescription  []struct {
		Language          string `json:"language"`
		LatestDescription string `json:"latestDescription"`
	} `json:"eventDescription"`
	AffectedEntities []struct {
		EntityValue string `json:"entityValue"`
	} `json:"affectedEntities"`
}

type SNSMessage struct {
	Version    string        `json:"version"`
	ID         string        `json:"id"`
	DetailType string        `json:"detail-type"`
	Source     string        `json:"source"`
	Account    string        `json:"account"`
	Time       time.Time     `json:"time"`
	Region     string        `json:"region"`
	Resources  []interface{} `json:"resources"`
	Detail     Detail         `json:"detail"`
}

func sendToWebhook(chatApplication string, webhook string, message TeamsMessage) error {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(message)
	res, err := http.Post(webhook, "application/json; charset=utf-8", b)
	if err != nil {
		return err
	}
	fmt.Print(res.Status)
	return nil
}

//func setMessage(interface{}) TeamsMessage {
//	var message TeamsMessage
//
//	return TeamsMessage{}
//}

func handler(ctx context.Context, event events.SNSEvent) {
	var chatApplication = os.Getenv("Chat_Application")
	var webhookURL = os.Getenv("Chat_Webhook")
	var snsMessage SNSMessage
	if len(event.Records) > 0 {
		for _, record := range event.Records {
			var message TeamsMessage
			err := json.Unmarshal([]byte(record.SNS.Message), &snsMessage)
			if err != nil {
				fmt.Print(err)
			}
			if snsMessage.Source == "aws.health" {
				message = TeamsMessage{
					Type:             "MessageCard",
					Context:          "http://schema.org/extensions",
					ThemeColor:       "fff30b",
					Summary:          snsMessage.DetailType,
					Sections:         []Section{
						{
							ActivityTitle:    snsMessage.Detail.EventTypeCode,
							ActivitySubtitle: snsMessage.Detail.EventTypeCategory,
							ActivityImage:    "",
							Facts:            []Fact{
								{
									Name:  "Description",
									Value: snsMessage.Detail.EventDescription[0].LatestDescription,
								},
							},
							Markdown:         true,
						},
					},
					PotentialActions: nil,
				}
			} else if snsMessage.Source == "aws.trustedadvisor" {
				taUri := fmt.Sprintf("https://console.aws.amazon.com/trustedadvisor/home?region=%s#/category/service-limits",snsMessage.Detail.CheckItemDetail.Region)
				colorMap := map[string]string{
					"Red": "d7000b",
					"Yellow": "fff30b",
				}
				message = TeamsMessage{
					Type:             "MessageCard",
					Context:          "http://schema.org/extensions",
					ThemeColor:       colorMap[snsMessage.Detail.CheckItemDetail.Status],
					Summary:          snsMessage.Detail.CheckItemDetail.LimitName,
					Sections:         []Section{
						{
							ActivityTitle:    "AWS service limit monitor",
							ActivitySubtitle: snsMessage.Detail.CheckName,
							ActivityImage:    "",
							Facts:            []Fact{
								{
									Name:  "Service",
									Value: snsMessage.Detail.CheckItemDetail.Service,
								},
								{
									Name:  "Region",
									Value: snsMessage.Detail.CheckItemDetail.Region,
								},
								{
									Name: "CurrentValue",
									Value: snsMessage.Detail.CheckItemDetail.CurrentUsage,
								},
								{
									Name: "LimitAmount",
									Value: snsMessage.Detail.CheckItemDetail.LimitAmount,
								},
							},
							Markdown:         true,
						},
					},
					PotentialActions: []PotentialAction{
						{
							Type:    "OpenUri",
							Name:    "More details",
							Inputs:  nil,
							Targets: []PotentialActionTarget{
								{
									OS:  "default",
									URI: taUri,
								},
							},
							Actions: nil,
						},
					},
				}
			}
			err = sendToWebhook(chatApplication, webhookURL, message)
			if err != nil {
				fmt.Print(err)
			}
		}
	}
}

func main() {
	lambda.Start(handler)
}
