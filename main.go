package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/cobra"
)

var config struct {
	snsTopic      string
	sqsURL        string
	collectorPort string
	awsregion     string
	s3Path        string
	awsSession    *session.Session
}

var preclogfile string

func main() {

	config.collectorPort = os.Getenv("PORT")
	if config.collectorPort == "" {
		config.collectorPort = "8080"
	}

	config.snsTopic = os.Getenv("SNS_TOPIC")
	config.sqsURL = os.Getenv("SQS_URL")
	config.awsregion = os.Getenv("AWS_DEFAULT_REGION")
	config.s3Path = os.Getenv("S3_PATH")

	config.awsSession = session.Must(session.NewSession(&aws.Config{Region: aws.String(config.awsregion)}))

	var collectorCmd = &cobra.Command{
		Use:   "collect",
		Short: "Run the collector",
		Run: func(cmd *cobra.Command, args []string) {
			if config.snsTopic == "" {
				panic("SNS_TOPIC required")
			}
			startCollector()
		},
	}

	var etlCmd = &cobra.Command{
		Use:   "etl",
		Short: "Run the ETL processor",
		Run: func(cd *cobra.Command, args []string) {
			if config.sqsURL == "" {
				panic("SQS_URL required")
			}
			// ensure we have database information here
			startETL()
		},
	}

	var precipitateCmd = &cobra.Command{
		Use:   "precipitate",
		Short: "Run the cloudfront processor",
		Run: func(cd *cobra.Command, args []string) {
			if config.s3Path == "" {
				panic("S3_PATH required")
			}
			startPrecipitate()
		},
	}

	precipitateCmd.Flags().StringVarP(&preclogfile, "logfile", "l", "", "Single cloudfront log file to process")
	var rootCmd = &cobra.Command{Use: "snowblower"}
	rootCmd.AddCommand(collectorCmd)
	rootCmd.AddCommand(etlCmd)
	rootCmd.AddCommand(precipitateCmd)
	rootCmd.Execute()

}
