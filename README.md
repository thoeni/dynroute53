# Dynroute53

### Abstract
This small program updates an Amazon AWS Route53 `A` Record for the specified `-domain` attached to a specified `-hostedZoneID` with the external IP address exposed to the internet.

### Use case
I built this to run in on a **Raspberry 🍓 PI** and I scheduled it with a `cron` job to run every minute: the program creates a local cache on the filesystem and if the external IP is equal to the one stored in the `dynroute.cache` file it won't update **Route53**, otherwise it will call the Amazon service to update the `A` record.

This helps me to reach my Raspberry PI by using a domain name even if the connection gets reset as within at most one minute the domain will be reattached to the correct (new) IP address.

### How to
To run this, you need:

- your AWS credentials stored in the (usual) `~/.aws/credentials` file [[more info]](http://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html)
- the user linked to those credentials needs to have access rights to the `route53:ChangeResourceRecordSets` action (see *AWS policy example* below)

<details><summary>AWS policy example</summary><p>

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "route53:ChangeResourceRecordSets"
            ],
            "Resource": [
                "arn:aws:route53:::hostedzone/%YOUR_HOSTED_ZONE%"
            ]
        }
    ]
}
```
</p></details>

#### The program requires two parameters:

- `-domain`: the entire subdomain/domain (excluding the protocol), for example `myblog.mydomain.com`
- `-hostedZoneID`: the AWS hosted zone ID as `arn` (in the policy example, replace `%YOUR_HOSTED_ZONE%` with it as well to allow the script to operate only on this specific zone/domain


This is what a call looks like:

```
./dynroute53 -domain=myblog.mydomain.com -hostedZoneID=Z123ABCDE45F67
```

### Build

If you want to build it for a **Raspberry PI** bear in mind I had an issue with the AWS Go SDK and I had to cross compile the program on my Mac with these env vars:

```
env GOOS=linux GOARCH=arm GOARM=7 go build main.go
```
