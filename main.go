package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

func main() {

	const cacheFilename = "dynroute53.cache"
	now := time.Now().Format("2006-01-02 15:04:05")

	domain := flag.String("domain", "", "fully qualified domain name to update, eg. blog.mydomain.com")
	hostedZoneID := flag.String("hostedZoneID", "", "AWS hosted zone ID (user must have permissions")
	flag.Parse()

	if *domain == "" || *hostedZoneID == "" {
		fmt.Println("Domain and hostedZoneId must be set")
		return
	}

	ip, err := getCurrentExternalIP()
	if err != nil {
		fmt.Println("Error while retrieving external IP", err)
		return
	}

	bytes, err := ioutil.ReadFile(cacheFilename)
	if err != nil {
		fmt.Println("Error while reading cache for external IP", err)
	}

	cachedIP := string(bytes)

	if cachedIP == ip {
		fmt.Printf("[%s] No need to update:\t%s -> %s (cached %s)\n", now, *domain, ip, cachedIP)
		return
	}

	if err := updateAWSRoute53(*domain, *hostedZoneID, ip); err != nil {
		fmt.Println("Error while updating route53 record", err)
		return
	}

	if err := ioutil.WriteFile(cacheFilename, []byte(ip), os.ModePerm); err != nil {
		fmt.Println("Error while updating local cache", err)
	}

	fmt.Printf("[%s] Update complete:\t%s -> %s\n", now, *domain, ip)
}

func getCurrentExternalIP() (string, error) {
	r, err := http.Get("http://ipv4.myexternalip.com/raw")
	if err != nil || r.StatusCode != 200 {
		return "", err
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	return strings.TrimSpace(string(bodyBytes)), nil
}

func updateAWSRoute53(domain string, hostedZoneID string, ip string) error {
	params := route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(domain),
						Type: aws.String("A"),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(ip),
							},
						},
						TTL: aws.Int64(300),
					},
				},
			},
		},
		HostedZoneId: aws.String(hostedZoneID),
	}
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	}))

	svc := route53.New(sess)
	_, err := svc.ChangeResourceRecordSets(&params)

	return err
}
