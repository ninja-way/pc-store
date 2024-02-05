package rabbitmq

import (
	"context"
	"encoding/json"
	audit "github.com/ninja-way/mq-audit-log/pkg/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
}

func NewClient(mqURI string) (*Client, error) {
	conn, err := amqp.Dial(mqURI)
	if err != nil {
		return nil, err
	}
	//defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	//defer ch.Close()

	q, err := ch.QueueDeclare(
		"logs", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
		ch:   ch,
		q:    q,
	}, nil
}

func (c *Client) CloseConnection() error {
	return c.conn.Close()
}

func (c *Client) SendLogRequest(ctx context.Context, req audit.LogItem) error {
	item, err := json.Marshal(req)
	if err != nil {
		return err
	}

	return c.ch.PublishWithContext(
		ctx,
		"",       // exchange
		c.q.Name, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "json",
			Body:        item,
		})
}
