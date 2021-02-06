

地址
----------
https://github.com/bsm/sarama-cluster



示例
==========

Consumers have two modes of operation. In the default multiplexed mode messages (and errors) of multiple topics and partitions are all passed to the single channel:	
```golang
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	cluster "github.com/bsm/sarama-cluster"
)

func main() {

	// init (custom) config, enable errors and notifications
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	// init consumer
	brokers := []string{"127.0.0.1:9092"}
	topics := []string{"my_topic", "other_topic"}
	consumer, err := cluster.NewConsumer(brokers, "my-consumer-group", topics, config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("Rebalanced: %+v\n", ntf)
		}
	}()

	// consume messages, watch signals
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				consumer.MarkOffset(msg, "")	// mark message as processed
			}
		case <-signals:
			return
		}
	}
}
```


Users who require access to individual partitions can use the partitioned mode which exposes access to partition-level consumers:
```golang
package main

import (
  "fmt"
  _ "log"
  "os"
  "os/signal"

  cluster "github.com/bsm/sarama-cluster"
)

func main() {

	// init (custom) config, set mode to ConsumerModePartitions
	config := cluster.NewConfig()
	config.Group.Mode = cluster.ConsumerModePartitions

	// init consumer
	brokers := []string{"127.0.0.1:9092"}
	topics := []string{"my_topic", "other_topic"}
	consumer, err := cluster.NewConsumer(brokers, "my-consumer-group", topics, config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume partitions
	for {
		select {
		case part, ok := <-consumer.Partitions():
			if !ok {
				return
			}

			// start a separate goroutine to consume messages
			go func(pc cluster.PartitionConsumer) {
				for msg := range pc.Messages() {
					fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
					consumer.MarkOffset(msg, "")	// mark message as processed
				}
			}(part)
		case <-signals:
			return
		}
	}
}
```


协程
```golang

package main
 
import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)
 
func main() {
	var Address = []string{"127.0.0.1:9092"}
	topic := []string{"test"}
	var wg = &sync.WaitGroup{}
	wg.Add(2)
	//广播式消费：消费者1
	go clusterConsumer(wg, Address, topic, "group-1")
	//广播式消费：消费者2
	go clusterConsumer(wg, Address, topic, "group-2")
 
	wg.Wait()
}
 
// 支持brokers cluster的消费者
func clusterConsumer(wg *sync.WaitGroup, brokers, topics []string, groupId string) {
	defer wg.Done()
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
 
	// init consumer
	consumer, err := cluster.NewConsumer(brokers, groupId, topics, config)
	if err != nil {
		log.Printf("%s: sarama-cluster.NewConsumer err, message=%s \n", groupId, err)
		return
	}
	defer consumer.Close()
 
	// trap SIGINT to trigger a shutdown
	signals := make(chan os.Signal)
	signal.Notify(signals,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		os.Interrupt,
		os.Kill,
	)
 
	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Printf("groupId=%s, Error= %s\n", groupId, err.Error())
		}
	}()
 
	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("groupId=%s, Rebalanced Info= %+v \n", groupId, ntf)
		}
	}()
 
	// consume messages, watch signals
	var successes int
Loop:
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				fmt.Fprintf(os.Stdout, "groupId=%s, topic=%s, partion=%d, offset=%d, key=%s, value=%s\n", groupId, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				consumer.MarkOffset(msg, "") // mark message as processed
				successes++
			}
		case <-signals:
			break Loop
		}
	}
	fmt.Fprintf(os.Stdout, "%s consume %d messages \n", groupId, successes)
}
```



错误
=================

panic: non-positive interval for NewTicker
-----------------
```sh
// 报错样式：
panic: non-positive interval for NewTicker
 
goroutine 59 [running]:
time.NewTicker(0x0, 0x0)
        D:/Go/src/time/tick.go:23 +0x14e
github.com/bsm/sarama-cluster.(*Consumer).cmLoop(0xc000212000, 0xc0002ba1e0)
        D:/work/mygo/pkg/mod/github.com/bsm/sarama-cluster@v2.1.15+incompatible/consumer.go:452 +0x61
github.com/bsm/sarama-cluster.(*loopTomb).Go.func1(0xc0002982a0, 0xc000288230)
        D:/work/mygo/pkg/mod/github.com/bsm/sarama-cluster@v2.1.15+incompatible/util.go:73 +0x82
created by github.com/bsm/sarama-cluster.(*loopTomb).Go
        D:/work/mygo/pkg/mod/github.com/bsm/sarama-cluster@v2.1.15+incompatible/util.go:69 +0x6d
//处理1： 找到这个consumer.go源码位置，上面的第二个报错有标注位置
github.com/bsm/sarama-cluster.(*Consumer).cmLoop(0xc000212000, 0xc0002ba1e0)
        D:/work/mygo/pkg/mod/github.com/bsm/sarama-cluster@v2.1.15+incompatible/consumer.go:452 +0x61
```


方案1：
```golang
// 修改452行，
	//	ticker := time.NewTicker(c.client.config.Consumer.Offsets.CommitInterval)
	ticker := time.NewTicker(c.client.config.Consumer.Offsets.AutoCommit.Interval)
// 保存重新build即可
```


方案2：
 
把 sarama 版本改成 从 v1.26.1 --> v1.24.1 就可以用啦 github.com/Shopify/sarama v1.24.1
 
gomod 的配置改下版本号就可以
```sh
github.com/Shopify/sarama v1.24.1
github.com/bsm/sarama-cluster v2.1.15+incompatible
```

