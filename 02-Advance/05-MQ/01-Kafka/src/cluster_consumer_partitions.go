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
	topics := []string{"test01",}
	consumer, err := cluster.NewConsumer(brokers, "group-1", topics, config)
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