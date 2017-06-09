package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/duncan/base64x"
	"github.com/remeh/sizedwaitgroup"
)

var queue struct {
	service *sqs.SQS
	params  *sqs.ReceiveMessageInput
}

func startETL() {
	queue.service = sqs.New(config.awsSession)

	queue.params = &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(config.sqsURL),
		AttributeNames: []*string{
			aws.String("All"), // Required
		},
		MaxNumberOfMessages: aws.Int64(1),
		MessageAttributeNames: []*string{
			aws.String("All"), // Required
		},
		VisibilityTimeout: aws.Int64(3600),
		WaitTimeSeconds:   aws.Int64(10),
	}

	// while something....
	for {
		processNextBatch()
	}
}

func processNextBatch() {

	resp, err := queue.service.ReceiveMessage(queue.params)

	if err != nil {
		// A service error occurred.
		fmt.Println("Error:", err.Error())
	}

	// limit threads to 10
	swg := sizedwaitgroup.New(10)

	for _, message := range resp.Messages {
		swg.Add()
		go processSNSMessage(message, &swg)
	}
	swg.Wait()
}

func processSNSMessage(message *sqs.Message, swg *sizedwaitgroup.SizedWaitGroup) {
	defer swg.Done()
	//messageID := *message.MessageID
	// receiptHandle := *message.ReceiptHandle
	deleteParams := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(config.sqsURL),
		ReceiptHandle: aws.String(*message.ReceiptHandle),
	}
	snsMessage := SNSMessage{}

	if err := json.Unmarshal([]byte(*message.Body), &snsMessage); err != nil {
		fmt.Printf("SNS MESSAGE UNMARSHALL ERROR %s\n", err)
	} else {
		payload := CollectorPayload{}
		if err := json.Unmarshal([]byte(snsMessage.Message), &payload); err != nil {
			fmt.Printf("COLLECTOR PAYLOAD UNMARSHALL ERROR %s\n", err)
		} else {
			if !checkmode {
				// schedule for deletion
				_, delerr := queue.service.DeleteMessage(deleteParams)
				if delerr != nil {
					fmt.Println(delerr.Error())
				}
			}
			processCollectorPayload(payload)
		}
	}
}

func processCollectorPayload(cp CollectorPayload) {
	tp := TrackerPayload{}
	var e Event
	if err := json.Unmarshal([]byte(cp.Body), &tp); err != nil {
		fmt.Printf("TRACKER PAYLOAD UNMARSHALL ERROR %s\n", err)
	} else {
		for _, e = range tp.Data {
			//dsfmt.Printf("%s, %s", cp.NetworkUserID, e.AppID)
			processEvent(e, tp, cp)
		}
	}
}

func processEvent(e Event, tp TrackerPayload, cp CollectorPayload) {

	b, errue := base64x.URLEncoding.DecodeString(e.UnstructuredEventEncoded)
	if errue != nil { // Error because it's not url encoded
		b, _ = base64x.StdEncoding.DecodeString(e.UnstructuredEventEncoded)
	}

	if len(e.UnstructuredEventEncoded) > 0 {
		if err := json.Unmarshal(b, &e.UnstructuredEvent); err != nil {
			fmt.Printf("UE UNMARSHALL ERROR %s\n%s\n", err, string(b))
			return
		}

	}
	b, errco := base64x.URLEncoding.DecodeString(e.ContextsEncoded)
	if errco != nil { // Error because it's not url encoded
		b, _ = base64x.StdEncoding.DecodeString(e.ContextsEncoded)
	}
	if len(e.ContextsEncoded) > 0 {
		if err := json.Unmarshal(b, &e.Contexts); err != nil {
			fmt.Printf("CO UNMARSHALL ERROR %s\n%s\n", err, string(b))
			return
		}
	}
	// pick up details from colletor payload
	if e.UserIPAddress == "" {
		e.UserIPAddress = cp.IPAddress
	}
	e.ETLTimestamp = time.Now()
	e.ETLVersion = SBVersion
	e.CollectorTimestamp = time.Unix(cp.Timestamp, 0)
	dtm, _ := strconv.ParseInt(e.TmpDeviceTimestamp, 10, 64)
	e.DeviceTimestamp = time.Unix(dtm/1000, 0)
	e.CollectorVersion = cp.Collector
	if e.UserAgent == "" {
		e.UserAgent = cp.UserAgent
	}
	// cp.RefererURI
	e.PageURLPath = cp.Path
	e.PageURLQuery = cp.QueryString
	// cp.Headers
	e.NetworkUserID = cp.NetworkUserID

	if e.validate() {
		e.enrich()
		if checkmode == true {
			e.print()
		} else {
			e.mongosave()
		}
	}

}
