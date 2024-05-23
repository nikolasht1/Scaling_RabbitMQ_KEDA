package main

import (
    "fmt"
    "log"
    "os"
    "time"

    "github.com/google/uuid"
    amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection
var queueName string
var instanceName string

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
        log.Fatal(err)
    }

    instanceName = os.Getenv("INSTANCE_NAME")
    if instanceName == "" {
        instanceName = "rabbitmq-consumer-" + uuid.NewString()
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

    err = ch.Qos(
        1,
        0,
        false,
    )

    msgs, err := ch.Consume(
        q.Name,
        "",
        false,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        log.Fatalf("failed to consume messages from queue: %v", err)
    }

    fmt.Println("consumer instance", instanceName, "waiting for messages.....")

    for msg := range msgs {
        fmt.Println("Instance", instanceName, "received message", string(msg.Body), "from queue", q.Name)
        msg.Ack(false)
        time.Sleep(3 * time.Second)
    }
}


