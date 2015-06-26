package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "ebclean"
	app.Usage = "A command to lean up old Application Versions in AWS Elastic Beanstalk"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "region, r",
			Value:  "us-east-1",
			Usage:  "AWS region",
			EnvVar: "AWS_REGION",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "list application, environment and how many versions, which is running",
			Action: func(c *cli.Context) {
				eb := elasticbeanstalk.New(&aws.Config{Region: c.String("region")})
				resp, err := eb.DescribeApplications(nil)

				if err != nil {
					if awsErr, ok := err.(awserr.Error); ok {
						fmt.Println(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())

						if reqErr, ok := err.(awserr.RequestFailure); ok {
							fmt.Println(reqErr.Code(), reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())
						}
					} else {
						fmt.Println(err.Error())
					}
				}

				fmt.Println(awsutil.StringValue(resp))
			},
		},
	}

	app.Run(os.Args)
}
