package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/a8m/kinesis-producer"
	"github.com/aws/aws-sdk-go/aws"
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
	flag.StringVar(&fileName, "file-name", "/var/log/syslog", "File to tail")
	flag.StringVar(&streamName, "stream-name", "test-stream", "Kinesis stream to write to")
	flag.StringVar(&shardName, "shard-name", "test-shard", "Kinesis shard to write to")
	flag.BoolVar(&followFile, "follow-file", true, "Follow the file, tailing new lines")
	flag.BoolVar(&reopenFile, "reopen-file", true, "Re-open the file if renamed (logrotate)")
	flag.IntVar(&fileWhence, "file-whence", 2, "0: start, 1: current, 2:end")
	flag.Parse()
	client := kinesis.New(session.New(aws.NewConfig()))
	pr := producer.New(&producer.Config{
		StreamName:   streamName,
		BacklogCount: 2000,
		Client:       client,
	})

	pr.Start()
	process(pr)
}
