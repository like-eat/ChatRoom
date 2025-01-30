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
	LoginWriter  *kafka.Writer
	ChatWriter   *kafka.Writer
	LogoutWriter *kafka.Writer
	LoginReader  *kafka.Reader
	ChatReader   *kafka.Reader
	LogoutReader *kafka.Reader
	KafkaConn    *kafka.Conn
}

var KafkaService = new(kafkaService)

// KafkaInit 初始化kafka
func (k *kafkaService) KafkaInit() {
	//k.CreateTopic()
	kafkaConfig := myconfig.GetConfig().KafkaConfig
	k.LoginWriter = &kafka.Writer{
		Addr:                   kafka.TCP(kafkaConfig.HostPort),
		Topic:                  kafkaConfig.LoginTopic,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           kafkaConfig.Timeout * time.Second,
		RequiredAcks:           kafka.RequireNone,
		AllowAutoTopicCreation: false,
	}
	k.ChatWriter = &kafka.Writer{
		Addr:                   kafka.TCP(kafkaConfig.HostPort),
		Topic:                  kafkaConfig.ChatTopic,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           kafkaConfig.Timeout * time.Second,
		RequiredAcks:           kafka.RequireNone,
		AllowAutoTopicCreation: false,
	}
	k.LogoutWriter = &kafka.Writer{
		Addr:                   kafka.TCP(kafkaConfig.HostPort),
		Topic:                  kafkaConfig.LogoutTopic,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           kafkaConfig.Timeout * time.Second,
		RequiredAcks:           kafka.RequireNone,
		AllowAutoTopicCreation: false,
	}
	k.LoginReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaConfig.HostPort},
		Topic:          kafkaConfig.LoginTopic,
		CommitInterval: kafkaConfig.Timeout * time.Second,
		GroupID:        "login",
		StartOffset:    kafka.LastOffset,
	})
	k.ChatReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaConfig.HostPort},
		Topic:          kafkaConfig.ChatTopic,
		CommitInterval: kafkaConfig.Timeout * time.Second,
		GroupID:        "chat",
		StartOffset:    kafka.LastOffset,
	})
	k.LogoutReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaConfig.HostPort},
		Topic:          kafkaConfig.LogoutTopic,
		CommitInterval: kafkaConfig.Timeout * time.Second,
		GroupID:        "logout",
		StartOffset:    kafka.LastOffset,
	})
}

func (k *kafkaService) KafkaClose() {
	if err := k.LoginWriter.Close(); err != nil {
		zlog.Error(err.Error())
	}
	if err := k.ChatWriter.Close(); err != nil {
		zlog.Error(err.Error())
	}
	if err := k.LogoutWriter.Close(); err != nil {
		zlog.Error(err.Error())
	}
	if err := k.KafkaConn.Close(); err != nil {
		zlog.Error(err.Error())
	}
	if err := k.LoginReader.Close(); err != nil {
		zlog.Error(err.Error())
	}
	if err := k.ChatReader.Close(); err != nil {
		zlog.Error(err.Error())
	}
	if err := k.LogoutReader.Close(); err != nil {
		zlog.Error(err.Error())
	}
}

// CreateTopic 创建topic
func (k *kafkaService) CreateTopic() {
	// 如果已经有topic了，就不创建了
	kafkaConfig := myconfig.GetConfig().KafkaConfig

	loginTopic := kafkaConfig.LoginTopic
	chatTopic := kafkaConfig.ChatTopic
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
			Topic:             chatTopic,
			NumPartitions:     kafkaConfig.Partition,
			ReplicationFactor: 1,
		},
	}

	// 创建topic
	if err = k.KafkaConn.CreateTopics(topicConfigs...); err != nil {
		zlog.Error(err.Error())
	}

}
