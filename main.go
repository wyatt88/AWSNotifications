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

type Finding struct {

}

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
	//Trusted Advisor
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

	//AWS Health
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

	//Spot Instances Interruption Notice
	InstanceID     string `json:"instance-id"`
	InstanceAction string `json:"instance-action"`

	// Guard Duty
	AccountID   string `json:"accountId"`
	Arn         string `json:"arn"`
	Confidence  int    `json:"confidence"`
	CreatedAt   string `json:"createdAt"`
	Description string `json:"description"`
	ID          string `json:"id"`
	Partition   string `json:"partition"`
	Region      string `json:"region"`
	Resource    struct {
		AccessKeyDetails struct {
			AccessKeyID string `json:"accessKeyId"`
			PrincipalID string `json:"principalId"`
			UserName    string `json:"userName"`
			UserType    string `json:"userType"`
		} `json:"accessKeyDetails"`
		InstanceDetails struct {
			AvailabilityZone   string `json:"availabilityZone"`
			IamInstanceProfile struct {
				Arn string `json:"arn"`
				ID  string `json:"id"`
			} `json:"iamInstanceProfile"`
			ImageDescription  string `json:"imageDescription"`
			ImageID           string `json:"imageId"`
			InstanceID        string `json:"instanceId"`
			InstanceState     string `json:"instanceState"`
			InstanceType      string `json:"instanceType"`
			LaunchTime        string `json:"launchTime"`
			NetworkInterfaces []struct {
				Ipv6Addresses      []string `json:"ipv6Addresses"`
				NetworkInterfaceID string   `json:"networkInterfaceId"`
				PrivateDNSName     string   `json:"privateDnsName"`
				PrivateIPAddress   string   `json:"privateIpAddress"`
				PrivateIPAddresses []struct {
					PrivateDNSName   string `json:"privateDnsName"`
					PrivateIPAddress string `json:"privateIpAddress"`
				} `json:"privateIpAddresses"`
				PublicDNSName  string `json:"publicDnsName"`
				PublicIP       string `json:"publicIp"`
				SecurityGroups []struct {
					GroupID   string `json:"groupId"`
					GroupName string `json:"groupName"`
				} `json:"securityGroups"`
				SubnetID string `json:"subnetId"`
				VpcID    string `json:"vpcId"`
			} `json:"networkInterfaces"`
			OutpostArn   string `json:"outpostArn"`
			Platform     string `json:"platform"`
			ProductCodes []struct {
				Code        string `json:"code"`
				ProductType string `json:"productType"`
			} `json:"productCodes"`
			Tags []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"tags"`
		} `json:"instanceDetails"`
		ResourceType string `json:"resourceType"`
	} `json:"resource"`
	SchemaVersion string `json:"schemaVersion"`
	GuardDutyService       struct {
		Action struct {
			ActionType       string `json:"actionType"`
			AwsAPICallAction struct {
				API           string `json:"api"`
				CallerType    string `json:"callerType"`
				DomainDetails struct {
					Domain string `json:"domain"`
				} `json:"domainDetails"`
				RemoteIPDetails struct {
					City struct {
						CityName string `json:"cityName"`
					} `json:"city"`
					Country struct {
						CountryCode string `json:"countryCode"`
						CountryName string `json:"countryName"`
					} `json:"country"`
					GeoLocation struct {
						Lat int `json:"lat"`
						Lon int `json:"lon"`
					} `json:"geoLocation"`
					IPAddressV4  string `json:"ipAddressV4"`
					Organization struct {
						Asn    string `json:"asn"`
						AsnOrg string `json:"asnOrg"`
						Isp    string `json:"isp"`
						Org    string `json:"org"`
					} `json:"organization"`
				} `json:"remoteIpDetails"`
				ServiceName string `json:"serviceName"`
			} `json:"awsApiCallAction"`
			DNSRequestAction struct {
				Domain string `json:"domain"`
			} `json:"dnsRequestAction"`
			NetworkConnectionAction struct {
				Blocked             bool   `json:"blocked"`
				ConnectionDirection string `json:"connectionDirection"`
				LocalIPDetails      struct {
					IPAddressV4 string `json:"ipAddressV4"`
				} `json:"localIpDetails"`
				LocalPortDetails struct {
					Port     int    `json:"port"`
					PortName string `json:"portName"`
				} `json:"localPortDetails"`
				Protocol        string `json:"protocol"`
				RemoteIPDetails struct {
					City struct {
						CityName string `json:"cityName"`
					} `json:"city"`
					Country struct {
						CountryCode string `json:"countryCode"`
						CountryName string `json:"countryName"`
					} `json:"country"`
					GeoLocation struct {
						Lat int `json:"lat"`
						Lon int `json:"lon"`
					} `json:"geoLocation"`
					IPAddressV4  string `json:"ipAddressV4"`
					Organization struct {
						Asn    string `json:"asn"`
						AsnOrg string `json:"asnOrg"`
						Isp    string `json:"isp"`
						Org    string `json:"org"`
					} `json:"organization"`
				} `json:"remoteIpDetails"`
				RemotePortDetails struct {
					Port     int    `json:"port"`
					PortName string `json:"portName"`
				} `json:"remotePortDetails"`
			} `json:"networkConnectionAction"`
			PortProbeAction struct {
				Blocked          bool `json:"blocked"`
				PortProbeDetails []struct {
					LocalIPDetails struct {
						IPAddressV4 string `json:"ipAddressV4"`
					} `json:"localIpDetails"`
					LocalPortDetails struct {
						Port     int    `json:"port"`
						PortName string `json:"portName"`
					} `json:"localPortDetails"`
					RemoteIPDetails struct {
						City struct {
							CityName string `json:"cityName"`
						} `json:"city"`
						Country struct {
							CountryCode string `json:"countryCode"`
							CountryName string `json:"countryName"`
						} `json:"country"`
						GeoLocation struct {
							Lat int `json:"lat"`
							Lon int `json:"lon"`
						} `json:"geoLocation"`
						IPAddressV4  string `json:"ipAddressV4"`
						Organization struct {
							Asn    string `json:"asn"`
							AsnOrg string `json:"asnOrg"`
							Isp    string `json:"isp"`
							Org    string `json:"org"`
						} `json:"organization"`
					} `json:"remoteIpDetails"`
				} `json:"portProbeDetails"`
			} `json:"portProbeAction"`
		} `json:"action"`
		Archived       bool   `json:"archived"`
		Count          int    `json:"count"`
		DetectorID     string `json:"detectorId"`
		EventFirstSeen string `json:"eventFirstSeen"`
		EventLastSeen  string `json:"eventLastSeen"`
		Evidence       struct {
			ThreatIntelligenceDetails []struct {
				ThreatListName string   `json:"threatListName"`
				ThreatNames    []string `json:"threatNames"`
			} `json:"threatIntelligenceDetails"`
		} `json:"evidence"`
		ResourceRole string `json:"resourceRole"`
		ServiceName  string `json:"serviceName"`
		UserFeedback string `json:"userFeedback"`
	} `json:"service"`
	Severity  int    `json:"severity"`
	Title     string `json:"title"`
	Type      string `json:"type"`
	UpdatedAt string `json:"updatedAt"`
	
	// Inspector
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
	Detail     Detail        `json:"detail"`
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
			switch snsMessage.Source {
			case "aws.health":
				message = TeamsMessage{
					Type:       "MessageCard",
					Context:    "http://schema.org/extensions",
					ThemeColor: "fff30b",
					Summary:    snsMessage.DetailType,
					Sections: []Section{
						{
							ActivityTitle:    snsMessage.Detail.EventTypeCode,
							ActivitySubtitle: snsMessage.Detail.EventTypeCategory,
							ActivityImage:    "",
							Facts: []Fact{
								{
									Name:  "Description",
									Value: snsMessage.Detail.EventDescription[0].LatestDescription,
								},
							},
							Markdown: true,
						},
					},
					PotentialActions: nil,
				}
			case "aws.trustedadvisor":
				taUri := fmt.Sprintf("https://console.aws.amazon.com/trustedadvisor/home?region=%s#/category/service-limits", snsMessage.Detail.CheckItemDetail.Region)
				colorMap := map[string]string{
					"Red":    "d7000b",
					"Yellow": "fff30b",
				}
				message = TeamsMessage{
					Type:       "MessageCard",
					Context:    "http://schema.org/extensions",
					ThemeColor: colorMap[snsMessage.Detail.CheckItemDetail.Status],
					Summary:    snsMessage.Detail.CheckItemDetail.LimitName,
					Sections: []Section{
						{
							ActivityTitle:    "AWS service limit monitor",
							ActivitySubtitle: snsMessage.Detail.CheckName,
							ActivityImage:    "",
							Facts: []Fact{
								{
									Name:  "Service",
									Value: snsMessage.Detail.CheckItemDetail.Service,
								},
								{
									Name:  "Region",
									Value: snsMessage.Detail.CheckItemDetail.Region,
								},
								{
									Name:  "CurrentValue",
									Value: snsMessage.Detail.CheckItemDetail.CurrentUsage,
								},
								{
									Name:  "LimitAmount",
									Value: snsMessage.Detail.CheckItemDetail.LimitAmount,
								},
							},
							Markdown: true,
						},
					},
					PotentialActions: []PotentialAction{
						{
							Type:   "OpenUri",
							Name:   "More details",
							Inputs: nil,
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
			case "aws.ec2":
				colorMap := map[string]string{
					"terminate":    "d7000b",
					"stop": "fff30b",
					"hibernate": "fff30b",
				}
				summary := fmt.Sprintf("Spot instance %s will %s",snsMessage.Detail.InstanceID,snsMessage.Detail.InstanceAction)
				message = TeamsMessage{
					Type:       "MessageCard",
					Context:    "http://schema.org/extensions",
					ThemeColor: colorMap[snsMessage.Detail.InstanceAction],
					Summary: summary,
					Sections: []Section{
						{
							ActivityTitle:    "AWS Spot Instance will " + snsMessage.Detail.InstanceAction,
							ActivitySubtitle: "Interruption Notice",
							ActivityImage:    "",
							Facts: []Fact{
								{
									Name:  "Instance ID",
									Value: snsMessage.Detail.InstanceID,
								},
								{
									Name:  "Action",
									Value: snsMessage.Detail.InstanceAction,
								},
							},
							Markdown: true,
						},
					},
					PotentialActions: []PotentialAction{
					},
				}
			//case "aws.guardduty":
			//	if len(snsMessage.Detail.Findings) >0 {
			//		for _, finding := range snsMessage.Detail.Findings {
			//			fmt.Print(finding.ID)
			//		}
			//	}
			default:
				//b, _ := json.Marshal(snsMessage.Detail)
				fmt.Print(snsMessage)
				message = TeamsMessage{
					Type:       "MessageCard",
					Context:    "http://schema.org/extensions",
					ThemeColor: "fff30b",
					Summary:    "Coming soon",
					Sections: []Section{
						{
							ActivityTitle:    snsMessage.Detail.EventTypeCode,
							ActivitySubtitle: snsMessage.Detail.EventTypeCategory,
							ActivityImage:    "",
							Facts: []Fact{
								{
									Name:  "Source Type",
									Value: snsMessage.Source,
								},
								{
									Name:  "message",
									Value: record.SNS.Message,
								},
							},
							Markdown: true,
						},
					},
					PotentialActions: nil,
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
