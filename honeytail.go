package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/a8m/kinesis-producer"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/hpcloud/tail"
)

var (
	fileName   string
	followFile bool
	reopenFile bool
	fileWhence int
	streamName string
	shardName  string
	awskey     string
	awssecret  string
)

func process(pr *producer.Producer) {
	t, err := tail.TailFile(fileName,
		tail.Config{
			Follow: followFile,
			ReOpen: reopenFile,
			Location: &tail.SeekInfo{
				Offset: 0,
				Whence: fileWhence,
			},
		})
	if err != nil {
		log.Println(err)
		return
	}
	for line := range t.Lines {
		// Process log line
		fmt.Println(line.Text)
		err := pr.Put([]byte(line.Text), shardName)
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	flag.StringVar(&awskey, "awskey", os.Getenv("AWS_KEY"), "Set environment variable AWS_KEY or specify on prompt")
	flag.StringVar(
		&awssecret, "awssecret", os.Getenv("AWS_SECRET"),
		"Set environment variable AWS_SECRET or specify on prompt",
	)
	flag.StringVar(&fileName, "file-name", "/var/log/syslog", "File to tail")
	flag.StringVar(&streamName, "stream-name", "stream-test", "Kinesis stream to write to")
	flag.StringVar(&shardName, "shard-name", "shard-test", "Kinesis shard to write to")
	flag.BoolVar(&followFile, "follow-file", true, "Follow the file, tailing new lines")
	flag.BoolVar(&reopenFile, "reopen-file", true, "Re-open the file if renamed (logrotate)")
	flag.IntVar(&fileWhence, "file-whence", 2, "0: start, 1: current, 2:end")
	flag.Parse()
	creds := credentials.NewStaticCredentials(awskey, awssecret, "")
	_, err := creds.Get()
	if err != nil {
		log.Println("Failed to create credentials", err.Error())
	}
	awsConfig := &aws.Config{
		Region:      aws.String("eu-west-1"),
		Credentials: creds,
	}
	client := kinesis.New(session.New(awsConfig))
	pr := producer.New(&producer.Config{
		StreamName:   streamName,
		BacklogCount: 2000,
		Client:       client,
	})

	pr.Start()
	process(pr)
}
