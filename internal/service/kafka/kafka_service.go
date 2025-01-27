package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	myconfig "kama_chat_server/config"
	"kama_chat_server/pkg/zlog"
	"time"
)

var ctx = context.Background()

type kafkaService struct {
	LoginWriter   *kafka.Writer
	ChatInWriter  *kafka.Writer
	ChatOutWriter *kafka.Writer
	LogoutWriter  *kafka.Writer
	KafkaConn     *kafka.Conn
}

var KafkaService = new(kafkaService)

// KafkaInit 初始化kafka
func (k *kafkaService) KafkaInit() {
	k.CreateTopic()
	k.LoginWriter = &kafka.Writer{
		Addr:                   kafka.TCP(myconfig.GetConfig().KafkaConfig.HostPort),
		Topic:                  myconfig.GetConfig().KafkaConfig.LoginTopic,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           myconfig.GetConfig().KafkaConfig.Timeout * time.Second,
		RequiredAcks:           kafka.RequireNone,
		AllowAutoTopicCreation: false,
	}
	k.ChatInWriter = &kafka.Writer{
		Addr:                   kafka.TCP(myconfig.GetConfig().KafkaConfig.HostPort),
		Topic:                  myconfig.GetConfig().KafkaConfig.ChatInTopic,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           myconfig.GetConfig().KafkaConfig.Timeout * time.Second,
		RequiredAcks:           kafka.RequireNone,
		AllowAutoTopicCreation: false,
	}
	k.ChatOutWriter = &kafka.Writer{
		Addr:                   kafka.TCP(myconfig.GetConfig().KafkaConfig.HostPort),
		Topic:                  myconfig.GetConfig().KafkaConfig.ChatOutTopic,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           myconfig.GetConfig().KafkaConfig.Timeout * time.Second,
		RequiredAcks:           kafka.RequireNone,
		AllowAutoTopicCreation: false,
	}
	k.LogoutWriter = &kafka.Writer{
		Addr:                   kafka.TCP(myconfig.GetConfig().KafkaConfig.HostPort),
		Topic:                  myconfig.GetConfig().KafkaConfig.LogoutTopic,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           myconfig.GetConfig().KafkaConfig.Timeout * time.Second,
		RequiredAcks:           kafka.RequireNone,
		AllowAutoTopicCreation: false,
	}

}

func (k *kafkaService) KafkaClose() {
	if err := k.LoginWriter.Close(); err != nil {
		zlog.Error(err.Error())
	}
	if err := k.ChatInWriter.Close(); err != nil {
		zlog.Error(err.Error())
	}
	if err := k.ChatOutWriter.Close(); err != nil {
		zlog.Error(err.Error())
	}
	if err := k.LogoutWriter.Close(); err != nil {
		zlog.Error(err.Error())
	}
	if err := k.KafkaConn.Close(); err != nil {
		zlog.Error(err.Error())
	}
}

// CreateTopic 创建topic
func (k *kafkaService) CreateTopic() {
	kafkaConfig := myconfig.GetConfig().KafkaConfig

	loginTopic := kafkaConfig.LoginTopic
	chatInTopic := kafkaConfig.ChatInTopic
	chatOutTopic := kafkaConfig.ChatOutTopic
	logoutTopic := kafkaConfig.LogoutTopic

	// 连接至任意kafka节点
	var err error
	k.KafkaConn, err = kafka.Dial("tcp", kafkaConfig.HostPort)
	if err != nil {
		zlog.Error(err.Error())
	}

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             loginTopic,
			NumPartitions:     kafkaConfig.Partition,
			ReplicationFactor: 1,
		},
		{
			Topic:             logoutTopic,
			NumPartitions:     kafkaConfig.Partition,
			ReplicationFactor: 1,
		},
		{
			Topic:             chatInTopic,
			NumPartitions:     kafkaConfig.Partition,
			ReplicationFactor: 1,
		},
		{
			Topic:             chatOutTopic,
			NumPartitions:     kafkaConfig.Partition,
			ReplicationFactor: 1,
		},
	}

	// 创建topic
	if err = k.KafkaConn.CreateTopics(topicConfigs...); err != nil {
		zlog.Error(err.Error())
	}

}
