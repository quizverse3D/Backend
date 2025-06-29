package common

import (
	"context"
	"log"

	"github.com/streadway/amqp"
)

type Consumer struct {
	channel *amqp.Channel
	queue   string
	handler func(amqp.Delivery)
}

func NewConsumer(channel *amqp.Channel, queue string, handler func(amqp.Delivery)) *Consumer {
	return &Consumer{
		channel: channel,
		queue:   queue,
		handler: handler,
	}
}

func (c *Consumer) DeclareQueue() error {
	_, err := c.channel.QueueDeclare(
		c.queue,
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}

func (c *Consumer) Listen(ctx context.Context) error {
	msgs, err := c.channel.Consume(
		c.queue,
		"",
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				// остановка обработки по контексту
				log.Printf("consumer for queue %s stopped", c.queue)
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Printf("channel closed for queue %s", c.queue)
				}
				c.safeHandle(msg)
			}
		}
	}()

	return nil
}

func (c *Consumer) safeHandle(msg amqp.Delivery) {
	defer func() {
		if r := recover(); r != nil {
			// отлов исключения и продолжение выполнения в случае panic()
			log.Printf("recovered in consumer for queue %s: %v", c.queue, r)
			msg.Nack(false, true) // retry
		}
	}()

	c.handler(msg)
}
