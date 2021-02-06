package main
 
import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)
 
func main() {
	var address = []string{"127.0.0.1:9092"}
	var topic string = "test01"

	asyncProducer(address, topic)
	
	time.Sleep(3*time.Second)
}
 
// 异步生产者
func asyncProducer(address []string, topic string) {
	config := sarama.NewConfig()
	// 等待服务器所有副本都保存成功后，再返回响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 随机向partition发送消息
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 是否等待成功和失败后的响应，只有上面的RequireAcks设置不是NoResponse，这里才有用。
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	// 设置读写超时时间为2秒，默认为10秒
	config.Producer.Timeout = 2 * time.Second
	// 尝试发送消息最大次数
	config.Producer.Retry.Max = 3
	// 设置使用的kafka版本，如果低于V0_10_0_0版本，消息中的timestamp没有作用。需要消费者和生产者同时配置。
	// 注意，如果设置了版本，但设置不对，则kafka会返回很奇怪的错误，并且无法成功发送消息。
	//config.Version = sarama.V0_10_0_1
 
	fmt.Println("start to make a producer")
	// 使用配置，新建一个异步生产者
	producer, e := sarama.NewAsyncProducer(address, config)
	if e != nil {
		fmt.Println("fail to make a producer, error=", e)
		return
	}
	defer producer.AsyncClose()
 
	// 循环判断哪个通道发送过来数据。
	fmt.Println("start goroutine to get response")
	go func(p sarama.AsyncProducer) {
		for {
			select {
			case suc := <-p.Successes():
				if suc != nil {
					fmt.Printf("succeed, offset=%d, timestamp=%s, partitions=%d\n", suc.Offset, suc.Timestamp.String(), suc.Partition)
					//fmt.Println("offset: ", suc.Offset, "timestamp: ", suc.Timestamp.String(), "partitions: ", suc.Partition)
				}
			case fail := <-p.Errors():
				if fail != nil {
					fmt.Printf("error= %v\n", fail.Err)
				}
			}
		}
	}(producer)
 
	// 发送消息
	strKey := "key"
	srcValue := "async: this is a message, index=%d "
	for i := 0; i < 5; i++ {
		time.Sleep(500 * time.Millisecond)
		value := fmt.Sprintf(srcValue, i)
		// 发送的消息对应的主题。
		// 注意：这里的msg必须是新构建的变量。不然，发送过去的消息内容都是一样的，因为批次发送消息的关系。
		msg := &sarama.ProducerMessage{
			Topic: topic,
		}
 
		// 设置消息的key
		msg.Key = sarama.StringEncoder(strKey)
		// 设置消息的value，将字符串转化为字节数组
		msg.Value = sarama.ByteEncoder(value)
		//fmt.Println(value)
 
		// 使用通道发送
		producer.Input() <- msg
	}
}