package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/younesbeheshti/chatapp-backend/cmd/utils"
	"log"
	"os"
	"time"
)

type Service struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitService() *Service {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	dsn := os.Getenv("RABBITMQ_URL")
	fmt.Println(dsn)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	err = ch.ExchangeDeclare(
		"private_messages",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.ExchangeDeclare(
		"channel_messages",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	return &Service{conn: conn, ch: ch}
}

func (s *Service) ConsumeChannelMessages(messageHandler func(event *utils.Event)) error {

	q, err := s.ch.QueueDeclare(
		"",
		false,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		return failOnError(err, "Failed to declare a queue")
	}

	err = s.ch.QueueBind(
		q.Name,
		"",
		"channel_messages",
		false,
		nil,
	)
	if err != nil {
		return failOnError(err, "Failed to bind a queue")
	}

	msgs, err := s.ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return failOnError(err, "Failed to register a consumer")
	}

	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)
			var event *utils.Event
			err := json.Unmarshal(msg.Body, &event)
			if err != nil {
				continue
			}
			messageHandler(event)
		}
	}()

	return nil
}

func (s *Service) PublishChannelMessages(event *utils.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := json.Marshal(event)
	if err != nil {
		return failOnError(err, "Failed to marshal event")
	}

	return s.ch.PublishWithContext(
		ctx,
		"",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        req,
		},
	)

}

func (s *Service) ConsumePrivateMessages(userID uint, messageHandler func(event *utils.Event)) error {

	queueName := fmt.Sprintf("user_%d_queue", userID)

	q, err := s.ch.QueueDeclare(
		queueName,
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return failOnError(err, "Failed to declare a queue")
	}

	routingKey := fmt.Sprintf("user_%d", userID)

	err = s.ch.QueueBind(
		q.Name,
		routingKey,
		"private_messages",
		false,
		nil,
	)
	if err != nil {
		return failOnError(err, "Failed to bind a queue")
	}

	msgs, err := s.ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return failOnError(err, "Failed to register a consumer")
	}

	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)
			var event *utils.Event
			err := json.Unmarshal(msg.Body, &event)
			if err != nil {
				continue
			}
			messageHandler(event)
		}
	}()

	return nil

}

func (s *Service) PublishPrivateMessages(event *utils.Event) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := json.Marshal(event)
	if err != nil {
		return failOnError(err, "Failed to marshal event")
	}

	routingKey := fmt.Sprintf("user_%d", event.MessageRequest.ReceiverID)

	return s.ch.PublishWithContext(
		ctx,
		"private_messages",
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        req,
		})

}

func (s *Service) Close() {
	if s.ch != nil {
		s.ch.Close()
	}
	if s.conn != nil {
		s.conn.Close()
	}
}

func failOnError(err error, msg string) error {
	return fmt.Errorf("%s: %s", msg, err)
}
