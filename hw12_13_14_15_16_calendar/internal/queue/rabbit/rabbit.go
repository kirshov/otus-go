package rabbit

import (
	"errors"
	"io"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RbtCon struct {
	Conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbit(dsn string) *RbtCon {
	conn, err := amqp.Dial(dsn)
	if err != nil {
		log.Fatal(err)
	}

	return &RbtCon{
		Conn: conn,
	}
}

func (r *RbtCon) Stop() error {
	if r.Conn != nil {
		err := r.Conn.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RbtCon) InitQueue(name string) (io.Closer, error) {
	ch, err := r.Conn.Channel()
	if err != nil {
		return nil, err
	}

	r.channel = ch

	_, err = ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)

	return ch, err
}

func (r *RbtCon) Publish(queue string, data []byte) error {
	if r.channel == nil {
		return errors.New("channel is not init")
	}

	return r.channel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json", // Указываем тип контента
			Body:        data,
		},
	)
}

func (r *RbtCon) Consume(queue string) (<-chan []byte, error) {
	_ = queue
	return nil, nil
}
