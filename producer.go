package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "strconv"
    "time"

    amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection
var queueName string

func init() {
    url := os.Getenv("RABBITMQ_URI")
    if url == "" {
        log.Fatal("missing environment variable RABBITMQ_URI")
    }

    queueName = os.Getenv("RABBITMQ_QUEUE_NAME")
    if queueName == "" {
        log.Fatal("missing environment variable RABBITMQ_QUEUE_NAME")
    }

    var err error

    conn, err = amqp.Dial(url)
    if err != nil {
        log.Fatal("dial failed ", err)
    }
}

func main() {

    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to create channel: %v", err)
    }
    defer ch.Close()

    q, err := ch.QueueDeclare(
        queueName,
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        log.Fatalf("Failed to declare queue: %v", err)
    }

    for i := 0; i <= 1000000; i++ {
        _i := strconv.Itoa(i)

        msg := "message-" + _i

        err = ch.PublishWithContext(
            context.Background(),
            "",
            q.Name,
            false,
            false,
            amqp.Publishing{
                ContentType: "text/plain",
                Body:        []byte(msg),
            },
        )
        if err != nil {
            log.Fatalf("failed to publish message: %v", err)
        }

        fmt.Println("message", msg, "sent to queue", q.Name)
        time.Sleep(1 * time.Second)
    }
}


