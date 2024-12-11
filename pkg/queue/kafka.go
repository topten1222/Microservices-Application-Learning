package queue

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"

	"github.com/IBM/sarama"
	"github.com/go-playground/validator/v10"
)

func ConnectProducer(brokerUrls []string, apiKey, secret string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	if apiKey != "" && secret != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = apiKey
		config.Net.SASL.Password = secret
		config.Net.SASL.Mechanism = "PLAIN"
		config.Net.SASL.Handshake = true
		config.Net.SASL.Version = sarama.SASLHandshakeV1
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = &tls.Config{
			InsecureSkipVerify: true,
			ClientAuth:         tls.NoClientCert,
		}
	}
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3

	producer, err := sarama.NewSyncProducer(brokerUrls, config)
	if err != nil {
		return nil, errors.New("error: faild to connect to producer")
	}
	return producer, nil
}

func PushMessageWithKeyToQueue(brokerUrls []string, apiKey, secret, topic, key string, message []byte) error {
	producer, err := ConnectProducer(brokerUrls, apiKey, secret)
	if err != nil {
		return errors.New("error to connect producer: " + err.Error())
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
		Key:   sarama.StringEncoder(key),
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return errors.New("error: faild to send message")
	}
	log.Printf("Message is stored in topic (%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}

func ConnectionConsumer(brokerUrls []string, apiKey, secret string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	if apiKey != "" || secret != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = apiKey
		config.Net.SASL.Password = secret
		config.Net.SASL.Mechanism = "PLAIN"
		config.Net.SASL.Handshake = true
		config.Net.SASL.Version = sarama.SASLHandshakeV1
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = &tls.Config{
			InsecureSkipVerify: true,
			ClientAuth:         tls.NoClientCert,
		}
	}
	config.Consumer.Return.Errors = true
	config.Consumer.Fetch.Max = 3
	consumer, err := sarama.NewConsumer(brokerUrls, config)
	if err != nil {
		return nil, errors.New("error to connect consumer: " + err.Error())
	}

	return consumer, nil
}

func DecodeMessage(obj any, value []byte) error {
	if err := json.Unmarshal(value, obj); err != nil {
		return errors.New("error to decode message: " + err.Error())
	}
	validate := validator.New()
	if err := validate.Struct(obj); err != nil {
		return errors.New("error to validate message: " + err.Error())
	}

	return nil
}
