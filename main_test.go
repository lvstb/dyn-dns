package main

import (
	"context"
	"testing"
	"yourproject/mocks"

	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSetDNSRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoute53Client := mocks.NewMockRoute53API(ctrl)

	// define your test inputs
	domainName := "example.com"
	wanip := "192.0.2.0"
	ttl := int64(300)
	zoneId := "Z2K1234ABCDEFGHIJKLMNOP"

	// setup expectations
	mockRoute53Client.EXPECT().ChangeResourceRecordSets(gomock.Any(), gomock.Any()).Return(&route53.ChangeResourceRecordSetsOutput{}, nil)

	// run the function
	err := setDNSRecord(context.Background(), domainName, wanip, ttl, zoneId)

	// assert no error
	assert.NoError(t, err)
}