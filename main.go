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
var checkmode bool
var enrichcheck bool

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
			if config.snsTopic == "" && !checkmode {
				panic("SNS_TOPIC required")
			}
			if enrichcheck && !checkmode {
				panic("Cannot do an enrichcheck unless checkmode is also set.")
			}
			startCollector()
		},
	}

	var etlCmd = &cobra.Command{
		Use:   "etl",
		Short: "Run the ETL processor",
		Run: func(cd *cobra.Command, args []string) {
			if config.sqsURL == "" && !checkmode {
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

	collectorCmd.Flags().BoolVarP(&checkmode, "check", "c", false, "Checkmode, verbose output; does not write to SQS or SNS. Use for debugging.")
	collectorCmd.Flags().BoolVarP(&enrichcheck, "echeck", "e", false, "When checkmode is defined, pass events onto enricher for debugging as well.")
	precipitateCmd.Flags().BoolVarP(&checkmode, "check", "c", false, "Checkmode, verbose output; does not write to SQS, SNS and will not move S3 logs to completed. Use for debugging.")
	precipitateCmd.Flags().StringVarP(&preclogfile, "logfile", "l", "", "Single cloudfront log file to process")
	etlCmd.Flags().BoolVarP(&checkmode, "check", "c", false, "Checkmode, verbose output; does not write to DB and will not delete SQS items. Use for debugging.")
	var rootCmd = &cobra.Command{Use: "snowblower"}
	rootCmd.AddCommand(collectorCmd)
	rootCmd.AddCommand(etlCmd)
	rootCmd.AddCommand(precipitateCmd)
	rootCmd.Execute()

}
