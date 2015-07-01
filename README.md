# aws-eb-cleanup
A tool to delete old AWS Elastic Beanstalk Application Versions (and their source bundles in S3)

# Features

1. Clean up older versions of [Elastic Beanstalk](http://aws.amazon.com/documentation/elastic-beanstalk/) applications. These versions take up your S3 storage.
2. Allow you to specify the number of day that a version is considered *old*. The script will **never** delete a running version.

# Requirements

1. You need to set your AWS ACCESS KEY ID and ACCESS KEY SECRET. Follow [this guide](https://github.com/aws/aws-sdk-go) to set them
2. Your AWS ACCESS KEY should have respective permissions. The following ElasticBeanstalk permissions must be granted: **DescribeApplicationVersions**, **DescribeEnvironments**, **DeleteApplicationVersions** (if you want to do the actual clean up)  
3. AWS_REGION environment variable must be set. Default is *us-east-1*. You can also pass this value to global flag `--region` when running `ebclean`

# Usage

1. Clone the repository
2. Install godep (https://github.com/tools/godep)
3. Run `godep restore`
4. Run `go install`
5. Execute `ebclean help` in your terminal and follow the instruction there
6. `ebclean clean <beanstalk_application_name_>` will inspect, print out results and automatically clean up all old versions of the specified application. By default, it will delete all versions that were created more than 30 days from current time.

# Known Limitations

1. No bulk clean-up. If you need to clean up old versions for multiple Elastic Beanstalk Applications, you need to do it one by one. 

# Credits

The project uses the following awesome open source libraries:

1. [AWS Official Go SDK](https://github.com/aws/aws-sdk-go)
2. [Termtables](https://github.com/apcera/termtables)
3. [Cli](github.com/codegangsta/cli)
4. [Godep](https://github.com/tools/godep)



