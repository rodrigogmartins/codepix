package kafka

import (
	"fmt"
	"os"

	"github.com/codeedu/imersao/codepix-go/application/factory"
	appmodel "github.com/codeedu/imersao/codepix-go/application/model"
	"github.com/codeedu/imersao/codepix-go/application/usecase"
	"github.com/codeedu/imersao/codepix-go/domain/model"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jinzhu/gorm"
)

type KafkaProcessor struct {
	Database        *gorm.DB
	Producer        *ckafka.Producer
	DeliveryChannel chan ckafka.Event
}

func NewKafkaProcessor(database *gorm.DB, producer *ckafka.Producer, deliveryChannel chan ckafka.Event) *KafkaProcessor {
	return &KafkaProcessor{
		Database:        database,
		Producer:        producer,
		DeliveryChannel: deliveryChannel,
	}
}

func (kafkaProcessor *KafkaProcessor) Consume() {
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": os.Getenv("kafkaBootstrapServers"),
		"group.id":          os.Getenv("kafkaConsumerGroupId"),
		"auto.offset.reset": "earliest",
	}

	consumer, err := ckafka.NewConsumer(configMap)

	if err != nil {
		panic(err)
	}

	topics := []string{
		os.Getenv("kafkaTransactionTopic"),
		os.Getenv("kafkaTransactionConfirmationTopic"),
	}
	consumer.SubscribeTopics(topics, nil)

	fmt.Println("kafka consumer has been started")

	for {
		message, err := consumer.ReadMessage(-1)

		if err == nil {
			kafkaProcessor.processMessage(message)
		}
	}
}

func (kafkaProcessor *KafkaProcessor) processMessage(message *ckafka.Message) {
	transactionsTopic := "transactions"
	transactionConfirmationTopic := "transaction_confirmation"

	switch topic := message.TopicPartition.Topic; topic {
	case &transactionsTopic:
		kafkaProcessor.processTransaction(message)
	case &transactionConfirmationTopic:
		kafkaProcessor.processTransactionConfirmation(message)
	default:
		fmt.Println("not a valid topic", string(message.Value))
	}
}

func (kafkaProcessor *KafkaProcessor) processTransaction(message *ckafka.Message) error {
	transaction := appmodel.NewTransaction()
	err := transaction.ParseJson(message.Value)

	if err != nil {
		return err
	}

	transactionUseCase := factory.TransactionUseCaseFactory(kafkaProcessor.Database)

	createdTransaction, err := transactionUseCase.Register(
		transaction.AccountID,
		transaction.Amount,
		transaction.PixKeyTo,
		transaction.PixKeyKindTo,
		transaction.Description,
		transaction.ID,
	)

	if err != nil {
		return err
	}

	topic := "bank" + createdTransaction.PixKeyTo.Account.Bank.Code
	transaction.ID = createdTransaction.ID
	transaction.Status = model.TransactionPending
	transactionJson, err := transaction.ToJson()

	if err != nil {
		return err
	}

	err = Publish(string(transactionJson), topic, kafkaProcessor.Producer, kafkaProcessor.DeliveryChannel)

	if err != nil {
		return err
	}

	return nil
}

func (kafkaProcessor *KafkaProcessor) processTransactionConfirmation(message *ckafka.Message) error {
	transaction := appmodel.NewTransaction()
	err := transaction.ParseJson(message.Value)

	if err != nil {
		return err
	}

	transactionUseCase := factory.TransactionUseCaseFactory(kafkaProcessor.Database)

	if transaction.Status == model.TransactionConfirmed {
		err = kafkaProcessor.confirmTransaction(transaction, transactionUseCase)

		if err != nil {
			return err
		}
	} else if transaction.Status == model.TransactionCompleted {
		_, err := transactionUseCase.Complete(transaction.ID)

		if err != nil {
			return err
		}
	}

	return nil
}

func (kafkaProcessor *KafkaProcessor) confirmTransaction(transaction *appmodel.Transaction, transactionUseCase usecase.TransactionUseCase) error {
	confirmedTransaction, err := transactionUseCase.Confirm(transaction.ID)

	if err != nil {
		return err
	}

	topic := "bank" + confirmedTransaction.AccountFrom.Bank.Code
	transactionJson, err := transaction.ToJson()

	if err != nil {
		return err
	}

	err = Publish(string(transactionJson), topic, kafkaProcessor.Producer, kafkaProcessor.DeliveryChannel)

	if err != nil {
		return err
	}

	return nil
}
