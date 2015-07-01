package main

import (
	"bytes"
	"fmt"
	"github.com/apcera/termtables"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/codegangsta/cli"
	"log"
	"os"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "ebclean"
	app.Usage = "A command to lean up old Application Versions in AWS Elastic Beanstalk"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "region, r",
			Value:  "us-east-1",
			Usage:  "AWS region. Default is us-east-1",
			EnvVar: "AWS_REGION",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "inspect",
			Usage: "inspect an application",
			Action: func(c *cli.Context) {
				eb := elasticbeanstalk.New(&aws.Config{Region: c.GlobalString("region")})
				appName := c.Args().First()

				duration := c.Int("duration")

				//if flag is not set
				if duration == 0 {
					duration = 30
				}

				inspect(eb, appName, duration)
			},
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "duration, d",
					Value: 30,
					Usage: `number of days from current time that a version is considered obsolete. For example, if set to 30, 
			all versions older than 30 days from current time will be marked to delete. Default is 30`,
				},
			},
		},
		{
			Name:  "clean",
			Usage: "clean up all old versions of an application",
			Action: func(c *cli.Context) {
				eb := elasticbeanstalk.New(&aws.Config{Region: c.GlobalString("region")})
				appName := c.Args().First()

				duration := c.Int("duration")

				//if flag is not set
				if duration == 0 {
					duration = 30
				}

				disposableVersionList := inspect(eb, appName, duration)

				if len(disposableVersionList) > 0 {
					deleteSourceBundle := c.BoolT("delete-source-bundle")

					//if the flag is not set, default is true
					if !deleteSourceBundle {
						deleteSourceBundle = true
					}

					clean(eb, appName, disposableVersionList, deleteSourceBundle)
				} else {
					fmt.Println("No version can be deleted")
				}
			},
			Flags: []cli.Flag{
				cli.BoolTFlag{
					Name:  "delete-source-bundle",
					Usage: "if true, this will delete source bundle for a given application version in AWS S3. Default is true.",
				},
				cli.IntFlag{
					Name:  "duration, d",
					Value: 30,
					Usage: `number of days from current time that a version is considered obsolete. For example, if set to 30, 
			all versions older than 30 days from current time will be marked to delete. Default is 30`,
				},
			},
		},
	}

	app.Run(os.Args)
}

// inspect and print out details of disposable versions for a given application
func inspect(eb *elasticbeanstalk.ElasticBeanstalk, appName string, duration int) []*elasticbeanstalk.ApplicationVersionDescription {
	applicationVersionResp, err := eb.DescribeApplicationVersions(&elasticbeanstalk.DescribeApplicationVersionsInput{ApplicationName: &appName})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			fmt.Println(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
			if reqErr, ok := err.(awserr.RequestFailure); ok {
				fmt.Println(reqErr.Code(), reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())
			}
		} else {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}

	environmentResp, err := eb.DescribeEnvironments(&elasticbeanstalk.DescribeEnvironmentsInput{ApplicationName: &appName})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			fmt.Println(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
			if reqErr, ok := err.(awserr.RequestFailure); ok {
				fmt.Println(reqErr.Code(), reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())
			}
		} else {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}

	var disposableVersions = make([]*elasticbeanstalk.ApplicationVersionDescription, 0)

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("Application Name: %s\n", appName))
	buffer.WriteString(fmt.Sprintf("Total version: %d\n", len(applicationVersionResp.ApplicationVersions)))

	table := termtables.CreateTable()

	table.AddHeaders("Version Label", "Date Created", "Environment")

	for _, version := range applicationVersionResp.ApplicationVersions {
		var found bool = false
		for _, environment := range environmentResp.Environments {
			if *version.VersionLabel == *environment.VersionLabel {
				table.AddRow(*version.VersionLabel, version.DateCreated, *environment.EnvironmentName)
				found = true
			}
		}

		if !found {
			if time.Now().After(version.DateCreated.AddDate(0, 0, duration)) {
				disposableVersions = append(disposableVersions, version)
				table.AddRow(*version.VersionLabel, version.DateCreated, "*Disposable*")
			} else {
				table.AddRow(*version.VersionLabel, version.DateCreated, "Not disposable")
			}
		}
	}

	buffer.WriteString(table.Render())
	buffer.WriteString(fmt.Sprintf("Disposable version count: %d\n", len(disposableVersions)))

	fmt.Println(buffer.String())

	return disposableVersions
}

//clean up all Application Version inside "versions" slice
func clean(eb *elasticbeanstalk.ElasticBeanstalk, appName string, versions []*elasticbeanstalk.ApplicationVersionDescription, deleteSourceBundle bool) {
	for _, version := range versions {
		var buffer bytes.Buffer

		buffer.WriteString(fmt.Sprintf("Deleting version %s .......... ", *version.VersionLabel))

		_, err := eb.DeleteApplicationVersion(&elasticbeanstalk.DeleteApplicationVersionInput{
			ApplicationName:    &appName,
			DeleteSourceBundle: &deleteSourceBundle,
			VersionLabel:       version.VersionLabel,
		})

		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				log.Println(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
				buffer.WriteString("Error")
				if reqErr, ok := err.(awserr.RequestFailure); ok {
					log.Println(reqErr.Code(), reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())
					buffer.WriteString("Error")
				}
			} else {
				log.Println(err.Error())
				buffer.WriteString("Error")
			}
		} else {
			buffer.WriteString("OK")
		}

		fmt.Println(buffer.String())
	}
}
