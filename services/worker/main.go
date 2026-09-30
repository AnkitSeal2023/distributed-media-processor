package main

import (
	"context"
	"distributed-media-processing-platform/services/worker/repository"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Event struct {
	Records []struct {
		S3 struct {
			Object struct {
				Key string `json:"key"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

var (
	port                  = os.Getenv("WORKER_PORT")
	bucketName            = os.Getenv("BUCKET_NAME")
	minio_accessKeyID     = os.Getenv("MINIO_ACCESSKEYID")
	minio_secretAccessKey = os.Getenv("MINIO_SECRETACCESSKEY")
	minio_endpoint        = fmt.Sprintf("http://localhost:%v", os.Getenv("MINIO_SERVER_PORT"))
	rabbitmq_endpoint     = os.Getenv("RABBITMQ_URL")
	minioClient           *minio.Client
	ctx                   = context.Background()
)

func main() {
	if port == "" {
		fmt.Println("WORKER_PORT environment variable is not set")
		os.Exit(1)
	}
	if bucketName == "" {
		fmt.Println("BUCKET_NAME environment variable is not set")
		os.Exit(1)
	}
	if minio_accessKeyID == "" {
		fmt.Println("MINIO_ACCESSKEYID environment variable is not set")
		os.Exit(1)
	}
	if minio_secretAccessKey == "" {
		fmt.Println("MINIO_SECRETACCESSKEY environment variable is not set")
		os.Exit(1)
	}
	if minio_endpoint == "" {
		fmt.Println("MINIO_ENDPOINT environment variable is not set")
		os.Exit(1)
	}
	if rabbitmq_endpoint == "" {
		fmt.Println("RABBITMQ_URL environment variable is not set")
		os.Exit(1)
	}

	// init db
	db, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to create database connection pool: ", err)
	}
	defer db.Close()

	workerRepo := repository.NewWorkerRepository(db)

	//init minio client
	useSSL := false
	minioclient, err := minio.New(minio_endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minio_accessKeyID, minio_secretAccessKey, ""),
		Secure: useSSL,
	})
	minioClient = minioclient
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := amqp.Dial(rabbitmq_endpoint)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"media_events",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"videos-ids",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")
	fmt.Printf("Queue declared successfully: %v\n", q.Name)

	msgs, err := ch.Consume(
		"videos-ids",
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			var event Event
			err := json.Unmarshal(d.Body, &event)
			if err != nil {
				log.Printf("Error unmarshalling message: %s", err)
				continue
			}

			for _, record := range event.Records {
				key := record.S3.Object.Key
				//TODO: remove printf
				log.Printf("Key: %s", key)
				err, requeue := processVideo(ctx, key, *workerRepo)
				//TODO : implement 3 tries, then DLQ for failed videos
				if err != nil {
					log.Printf("Error processing video: %s", err)
				}
				if requeue {
					d.Nack(false, true)
				} else {
					d.Nack(false, true)
				}
			}
		}
	}()

	log.Println(" [*] Waiting for logs. To exit press CTRL+C")
	select {}
}
