package kafka

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"time"

// 	"bitbucket.org/inseinc/report-service/kafka/protocol"
// 	"github.com/google/uuid"
// 	"github.com/segmentio/kafka-go"
// )

// var bootstrapServer = []string{"localhost:9092"}

// const (
// 	Topic    = "transaction"
// 	GroupKey = "group1"
// )

// var transactionTypes = []string{
// 	"digital_sell",
// 	"retail_sell",
// 	"void",
// 	"cash_ticket",
// 	"top_up",
// 	"pin",
// }

// var entityTypes = []string{"user_wallet", "pos"}

// type Transaction struct {
// 	EntityType     string      `json:"entityType"`
// 	EntityId       string      `json:"entityId"`
// 	Performer      string      `json:"performer"`
// 	Time           string      `json:"time"`
// 	OperationId    string      `json:"operationId"`
// 	TransactionId  string      `json:"transactionId"`
// 	Amount         float64     `json:"amount"`
// 	Currency       string      `json:"currency"`
// 	OperationType  string      `json:"operationType"`
// 	PaymentDetails interface{} `json:"paymentDetails"`
// 	Metadata       interface{} `json:"metadata"`
// }

// var listProduct = []string{
// 	"17d1074e-f443-4b97-b82e-3164d68e36d5",
// 	"ae85647f-4960-458e-8972-5266b2e1d179",
// 	"5dfde466-b9ad-4130-8dcb-c849585f68ac",
// 	// "765e5b3a-bc07-4f98-a198-36353a6d4a75",
// 	// "4dab199d-28c9-4a76-a466-f3e836a1f68f",
// }

// var listPos = []string{
// 	"71d5428c-ce67-40f1-9858-f8bfeccf9490",
// 	"c2f1d588-5bca-4138-ab7d-78fdec4f3874",
// 	"d7f0ee5a-b79e-4368-9d59-1753d9f46b78",
// 	"828a588f-2051-4914-bb7d-f1f42b5e35b2",
// 	"73852a3a-bccf-4107-9844-c07cad4dc68e",
// 	// "89074851-314b-4cd9-9725-63852830291b",
// 	// "ca436a43-73ce-4cea-90bb-810ec13abb80",
// }

// // func (t *TestStruct) String() string {
// // 	return fmt.Sprintf("(id: %s, name: %s, value: %d)", t.Id, t.Name, t.Value)
// // }

// func TestKafka() {
// 	producer := setupKafkaProducer()
// 	producer.Initialize()
// 	// consumers := setupConsumers(1)
// 	// components := make([]common.Initiable, len(consumers)+1)
// 	// components[0] = producer
// 	// for i, consumer := range consumers {
// 	// 	components[i+1] = consumer
// 	// }

// 	// initialise components
// 	// for _, component := range components {
// 	// 	err := component.Initialize()
// 	// 	if err != nil {
// 	// 		panic(fmt.Sprintf("Error while iniatializing component: [%s]; Error: %s", component.Name(), err))
// 	// 	}
// 	// }

// 	// starting producers and consumers
// 	// for i, consumer := range consumers {
// 	// 	go startConsumer(consumer, fmt.Sprintf(GroupKey, i))
// 	// }

// 	channel := make(chan *protocol.Message)
// 	go messageProducer(5, channel)
// 	for message := range channel {
// 		publishMessage(producer, message)
// 	}

// 	// destroy components
// 	// for _, component := range components {
// 	// 	err := component.Destroy()
// 	// 	if err != nil {
// 	// 		log.Printf("Error while destrying component: [%s]; Error: %s", component.Name(), err)
// 	// 	}
// 	// }
// }

// func setupConsumers(count int) []*protocol.KafkaMessageConsumer {
// 	consumers := make([]*protocol.KafkaMessageConsumer, count)
// 	for i := 0; i < count; i++ {
// 		group := fmt.Sprintf(GroupKey, i)
// 		consumers[i] = setupKafkaConsumer(group)
// 	}
// 	return consumers
// }

// func messageProducer(count int, msgChannel chan *protocol.Message) {
// 	for i := 0; i < count; i++ {
// 		randInt := rand.Int31n(10_000)
// 		time.Sleep(time.Millisecond * time.Duration(randInt))
// 		msgChannel <- buildMessage(i)
// 	}
// 	close(msgChannel)
// }

// func publishMessage(producer *protocol.KafkaMessageProducer, msg *protocol.Message) {
// 	fmt.Println("Publish message to kafka")
// 	err := producer.PublishMessage(msg)
// 	if err != nil {
// 		panic(fmt.Sprintf("Error while producing message: %s", err))
// 	}
// }

// func setupKafkaProducer() *protocol.KafkaMessageProducer {
// 	producerConfig := &protocol.KafkaProducerConfiguration{
// 		BootstrapServers:  bootstrapServer,
// 		Topic:             Topic,
// 		MaxAttempts:       10,
// 		Balancer:          &kafka.Hash{},
// 		TopicAutoCreation: true,
// 	}
// 	producer := protocol.NewKafkaMessageProducer(producerConfig)
// 	return producer
// }

// func setupKafkaConsumer(group string) *protocol.KafkaMessageConsumer {
// 	consumerConfig := &protocol.KafkaConsumerConfiguration{
// 		BootstrapServers: bootstrapServer,
// 		Topic:            Topic,
// 		GroupID:          group,
// 	}
// 	consumer := protocol.NewKafkaMessageConsumer(consumerConfig)
// 	return consumer
// }

// func startConsumer(consumer *protocol.KafkaMessageConsumer, group string) {
// 	consumer.StartConsumer(func(message *protocol.Message) error {
// 		ts := &Transaction{}
// 		err := json.Unmarshal(message.Data, ts)
// 		if err != nil {
// 			log.Printf("Error while unmarshalling message data: %s", err)
// 		}
// 		log.Printf("Recieved message on consumer group: [%s]; Message: %s", group, ts)
// 		return nil
// 	})
// }

// func buildMessage(seq int) *protocol.Message {
// 	key := fmt.Sprintf("key-%d", seq)

// 	ts := &Transaction{
// 		EntityType:    entityTypes[random(0, 2)],
// 		Performer:     fmt.Sprintf("Performer %v", seq),
// 		Time:          time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
// 		OperationId:   uuid.New().String(),
// 		TransactionId: uuid.New().String(),
// 		Amount:        randomAmount(10, 20),
// 		Currency:      "USD",
// 	}

// 	if ts.EntityType == "pos" {
// 		ts.EntityId = listPos[random(0, 5)]
// 		ts.OperationType = []string{"retail_sell", "cash_ticket"}[random(0, 2)]
// 	} else if ts.EntityType == "user_wallet" {
// 		ts.EntityId = uuid.New().String() // user id
// 		ts.OperationType = []string{"digital_sell", "top_up"}[random(0, 2)]
// 	}

// 	if ts.OperationType == "digital_sell" || ts.OperationType == "retail_sell" {
// 		ts.Metadata = map[string]string{
// 			"productId": listProduct[random(0, 3)],
// 		}
// 	}

// 	if ts.OperationType == "digital_sell" || ts.OperationType == "cash_ticket" {
// 		ts.Amount = randomAmount(1, 5) * -1
// 	}

// 	if ts.OperationType == "cash_ticket" {
// 		ts.Metadata = map[string]string{
// 			"ticketId": uuid.New().String(),
// 		}
// 	}

// 	if ts.OperationType == "top_up" {
// 		ts.PaymentDetails = map[string]string{
// 			"paymentIdProvider": uuid.New().String(),
// 			"paymentId":         uuid.New().String(),
// 			"paymentMethodId":   uuid.New().String(),
// 		}
// 	}

// 	fmt.Println("Message: ", ts)
// 	data, err := json.Marshal(ts)
// 	if err != nil {
// 		panic("Marshalling error")
// 	}

// 	return &protocol.Message{
// 		Topic:     Topic,
// 		Timestamp: time.Now(),
// 		Key:       key,
// 		Data:      data,
// 	}
// }

// func random(min, max int) int {
// 	r := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	return r.Intn(max-min) + min
// }

// func randomAmount(min, max int) float64 {
// 	r := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	amount := r.Intn(max-min) + min
// 	decimal := float64(r.Intn(1/0.05)) * 0.05
// 	return float64(amount)*10 + decimal
// }
