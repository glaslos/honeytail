package kinesis

import (
	"github.com/a8m/kinesis-producer"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

// Producer publishes aggregated messages
type Producer struct {
	pr        *producer.Producer
	shardName string
}

func newProducer(awskey string, awssecret string, streamName string, shardName string) (Producer, error) {
	creds := credentials.NewStaticCredentials(awskey, awssecret, "")
	_, err := creds.Get()
	if err != nil {
		return Producer{}, err
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
	k := Producer{
		pr:        pr,
		shardName: shardName,
	}
	return k, nil
}

func (k Producer) publish(payload string) error {
	err := k.pr.Put([]byte(payload), k.shardName)
	if err != nil {
		return err
	}
	return nil
}
