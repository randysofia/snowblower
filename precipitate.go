package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

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

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	// Pretty-print the response data.
	fmt.Println(numBytes)

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

func startPrecipitate() {
	gzfile, _ := GetS3(config.s3Path + "/" + preclogfile)
	csvfile := strings.Replace(gzfile, "gz", "csv", 1)
	decompressgz(gzfile, csvfile)

}
