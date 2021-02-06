package main
 
 
import (
	"github.com/Shopify/sarama"
	"time"
	"log"
	"fmt"
	"os"
)
 
func main()  {
	var address = []string{"127.0.0.1:9092"}
	topic := "test01"
	version := "2.1.1"
	syncProducer(address, version, topic)

	time.Sleep(2*time.Second)
}
 
// 同步生产消息模式
func syncProducer(address []string, version, topic string) error {
	kafkaVersion, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		return err
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = kafkaVersion
	config.Producer.Timeout = 3 * time.Second

	p, err := sarama.NewSyncProducer(address, config)
	if err != nil {
		log.Printf("sarama.NewSyncProducer err, message=%s \n", err)
		return err
	}
	defer p.Close()
 
	strKey := "key"
	srcValue := "sync: this is a message, index=%d "
	
	for i:=0; i<5; i++ {
		value := fmt.Sprintf(srcValue, i)
		msg := &sarama.ProducerMessage{
			Key:sarama.StringEncoder(strKey),
			Topic:topic,
			Value:sarama.ByteEncoder(value),
		}

		part, offset, err := p.SendMessage(msg)
		if err != nil {
			log.Printf("send message(%s) err=%v \n", value, err)
		} else {
			fmt.Fprintf(os.Stdout, value + "发送成功, partition=%d, offset=%d \n", part, offset)
		}
	}

	return nil
}