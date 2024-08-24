package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func getWanIP(url string) (string, error) {
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	wanipBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(wanipBytes), nil
}

type R53set interface {
	ChangeResourceRecordSets(context.Context, *route53.ChangeResourceRecordSetsInput, ...func(*route53.Options)) (*route53.ChangeResourceRecordSetsOutput, error)
}

func setDNSRecord(ctx context.Context, domainName string, wanip string, zoneId string, client R53set) error {
	ttl := int64(300)
	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeActionUpsert,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: &domainName,
						Type: types.RRTypeA,
						TTL:  aws.Int64(ttl),
						ResourceRecords: []types.ResourceRecord{
							{
								Value: &wanip,
							},
						},
					},
				},
			},
			Comment: aws.String("Updating the record"),
		},
		HostedZoneId: &zoneId,
	}

	_, err := client.ChangeResourceRecordSets(ctx, input)
	if err != nil {
		fmt.Printf("Error details: %+v\n", err)
		return fmt.Errorf("failed to create DNS record, %v", err)
	}

	fmt.Println("DNS record successfully created")
	return nil
}

func main() {
	domainName := flag.String("domain", "", "Domain name to update")
	zoneId := flag.String("zone", "", "Route53 Hosted Zone ID")

	flag.Parse()

	if *domainName == "" || *zoneId == "" {
		fmt.Println("Error: domain and zone flags are required")
		flag.Usage()
		return
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		fmt.Println("Error loading AWS config:", err)
		return
	}

	r53Client := route53.NewFromConfig(cfg)
	wanip, err := getWanIP("https://api.ipify.org")
	if err != nil {
		fmt.Println("Error getting WAN IP:", err)
		return
	}

	err = setDNSRecord(ctx, *domainName, wanip, *zoneId, r53Client)
	if err != nil {
		fmt.Println("Error creating DNS record:", err)
		return
	}

	fmt.Println("DNS record updated successfully. New IP:", wanip)
}