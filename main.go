package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

var domainName string = "home.wingu.dev"
var zoneId string = ""
var ttl int64 = 60

func getWanIP() (wanip string, err error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	wanipBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	wanip = string(wanipBytes)
	return
}

func setDNSRecord(ctx context.Context, domainName string, wanip string, ttl int64, zoneId string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("unable to load SDK config, %v", err)
	}

	client := route53.NewFromConfig(cfg)
	// a struct defining all config needed for the changeresourcerecordset method
	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeActionCreate,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: &domainName,
						Type: types.RRTypeA,
						TTL:  &ttl,
						ResourceRecords: []types.ResourceRecord{
							{
								Value: &wanip,
							},
						},
					},
				},
			},
			Comment: aws.String("Creating an A record"),
		},
		HostedZoneId: &zoneId,
	}

	_, err = client.ChangeResourceRecordSets(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create DNS record, %v", err)
	}

	fmt.Println("DNS record successfully created")
	return nil
}

func main() {
	//
	ctx := context.Background()
	wanip, err := getWanIP()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	err = setDNSRecord(ctx, domainName, wanip, ttl, zoneId)
	if err != nil {
		fmt.Println("Error creating DNS record:", err)
	}
	// fmt.Println("WAN IP:", wanip)
}