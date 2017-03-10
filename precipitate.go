package main

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/google/uuid"
)

// GetS3 Get the specific file from S3 and return the local path
func GetS3(filepath string) (string, error) {
	var ofilename string
	urldata, _ := url.Parse(filepath)

	params := &s3.GetObjectInput{
		Bucket: aws.String(urldata.Host),
		Key:    aws.String(urldata.Path),
	}

	newuuid, _ := uuid.NewRandom()
	ofilename = "/tmp/" + newuuid.String() + ".gz"

	output, err := os.Create(ofilename)
	if err != nil {
		return "", err
	}
	defer output.Close()

	downloader := s3manager.NewDownloader(config.awsSession)
	numBytes, err := downloader.Download(output, params)
	fmt.Printf("Downloaded %s to %s, %d bytes\n", filepath, ofilename, numBytes)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return ofilename, nil
}

func decompressgz(gzFilePath string, dstFilePath string) (int64, error) {
	gzFile, err := os.Open(gzFilePath)
	if err != nil {
		return 0, fmt.Errorf("Failed to open file %s for unpack: %s", gzFilePath, err)
	}
	dstFile, err := os.OpenFile(dstFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		return 0, fmt.Errorf("Failed to create destination file %s for unpack: %s", dstFilePath, err)
	}

	ioReader, ioWriter := io.Pipe()

	go func() { // goroutine leak is possible here
		gzReader, _ := gzip.NewReader(gzFile)
		// it is important to close the writer or reading from the other end of the
		// pipe or io.copy() will never finish
		defer func() {
			gzFile.Close()
			gzReader.Close()
			ioWriter.Close()
		}()

		io.Copy(ioWriter, gzReader)
	}()

	written, err := io.Copy(dstFile, ioReader)
	if err != nil {
		return 0, err // goroutine leak is possible here
	}
	ioReader.Close()
	dstFile.Close()

	os.Remove(gzFilePath)

	return written, nil
}

func processCSV(filename string) error {
	fs, err := os.Open(filename)
	defer fs.Close()
	if err != nil {
		return err
	}
	csvreader := csv.NewReader(bufio.NewReader(fs))
	csvreader.Comma = '\t'
	csvreader.FieldsPerRecord = 24
	for {
		record, err := csvreader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			continue
		}
		query, _ := url.QueryUnescape(record[11])
		u, err := url.ParseQuery(query)
		bodyBytes, err := urlValuesToBodyBytes(u)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		//fmt.Println(string(bodyBytes))
		collectorHandler := &collector{
			publisher: &SNSPublisher{
				service: sns.New(config.awsSession),
				topic:   config.snsTopic,
			},
		}
		newuuid, _ := uuid.NewRandom()
		networkID := newuuid.String()
		useragent, _ := url.QueryUnescape(record[10])
		var header []string
		header = append(header, record[15])
		collectorHandler.jsonInputToSNS(bodyBytes, record[4], strings.Replace(useragent, "%20", " ", -1), header, networkID)

	}
	return nil
}

func orchestratePrecipitation(filename string) {
	gzfile, _ := GetS3(config.s3Path + "/" + filename)
	csvfile := strings.Replace(gzfile, "gz", "csv", 1)
	decompressgz(gzfile, csvfile)
	processCSV(csvfile)
	os.Remove(csvfile)
}

func startPrecipitate() {
	if preclogfile != "" {
		orchestratePrecipitation(preclogfile)
	} else {
		urldata, _ := url.Parse(config.s3Path)

		svc := s3.New(config.awsSession)
		params := &s3.ListObjectsInput{
			Bucket: aws.String(urldata.Host),
		}
		resp, _ := svc.ListObjects(params)
		for _, key := range resp.Contents {
			if filepath.Ext(*key.Key) == ".gz" {
				filedata := strings.Split(*key.Key, "/")
				orchestratePrecipitation(filedata[len(filedata)-1])
			}
		}
	}
}
