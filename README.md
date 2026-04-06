# Canonflow Golang Domain-Driven Design (DDD)

A production-ready Go **REST API** boilerplate build on **Domain-Driven Design** with **Hexagonal Architecture** principles. This project provides a clean, scalable architecture with integrated messaging (Kafka and RabbitMQ), caching (Redis), database migration, JWT Authentication, and auto-generated Swagger documentation - all wired together with Docker Compose for a smooth local development experiences.

---

## Table of Contents

- [Canonflow Golang Domain-Driven Design (DDD)](#canonflow-golang-domain-driven-design-ddd)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Tech Stack](#tech-stack)
  - [Project Structure](#project-structure)

## Overview

`canonflow-go-ddd` is a structured Go backend starter that enforces a clear separation of conecerns via DDD. Each business domain lives in its own self-contained package under `internal/domain`, with distinct layers for `models`, `DTOs`, `repositories`, `usecases`, and `HTTP Messaging` delivery. This design makes it straightforward to scale teams and features without tangling responsibilities.

The project also ships with first-class support fo **event-driven messaging** - you can scaffold a new domain wiht optional `Kafka` producer/consumer stubs in a single `make` command.

## Tech Stack

| Category         | Technology                                                      |
| :--------------- | :-------------------------------------------------------------- |
| Language         | Go 1.25+                                                        |
| Web Framework    | Gin                                                             |
| ORM              | GORM                                                            |
| Database         | MySQL / PostgreSQL                                              |
| Cache            | Redis 8                                                         |
| Message Broker   | Apache Kafka (via Confluent) & RabbitMQ (Internal Queue System) |
| Migrations       | golang-migrate                                                  |
| Authentication   | JWT (golang-jwt)                                                |
| API Docs         | Swagger (swaggo/swag + gin-swagger)                             |
| Config           | Viper                                                           |
| Logging          | Logrus                                                          |
| Containerization | Docker                                                          |

## Project Structure

```
canonflow-go-ddd/
├── .env.example
├── .gitattributes
├── .gitignore
├── cmd/
│   ├── api/
│   │   ├── .gitignore
│   │   ├── docs/
│   │   │   ├── docs.go
│   │   │   ├── swagger.json
│   │   │   └── swagger.yaml
│   │   └── main.go
│   ├── migration/
│   │   ├── .gitignore
│   │   └── main.go
│   ├── queue/
│   │   ├── .gitignore
│   │   └── main.go
│   └── worker/
│       ├── .gitignore
│       └── main.go
├── dev.sh
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal/
│   ├── config/
│   │   ├── app.go
│   │   ├── gin.go
│   │   ├── gorm.go
│   │   ├── kafka.go
│   │   ├── logrus.go
│   │   ├── queue.go
│   │   ├── redis.go
│   │   ├── validator.go
│   │   └── viper.go
│   ├── contract/
│   │   ├── consumer.go
│   │   ├── database.go
│   │   ├── event.go
│   │   ├── producer.go
│   │   ├── queue.go
│   │   └── ratelimiter.go
│   ├── domain/
│   │   └── <domain>/
│   │       ├── delivery/
│   │       │   ├── http/
│   │       │   │   ├── handler.go
│   │       │   │   └── route.go
│   │       │   └── messaging/
│   │       │       └── <domain>_consumer.go
│   │       ├── dto/
│   │       │   └── <domain>.go
│   │       ├── gateway/
│   │       │   └── messaging/
│   │       │       └── <domain>_producer.go
│   │       ├── model/
│   │       │   ├── <domain>_event.go
│   │       │   └── <domain>.go
│   │       ├── queue/
│   │       │   └── <domain>_queue.go
│   │       ├── repository/
│   │       │   ├── <domain>_repository_impl.go
│   │       │   └── <domain>_repository.go
│   │       └── usecase/
│   │           ├── <domain>_usecase_impl.go
│   │           └── <domain>_usecase.go
│   ├── factory/
│   │   ├── database.go
│   │   └── ratelimiter.go
│   └── middleware/
│       ├── auth_middleware.go
│       └── ratelimiter_middleware.go
├── makefile
├── migrations/
│   └── mysql/
│       ├── 000001_create_users_table.down.sql
│       ├── 000001_create_users_table.up.sql
│       ├── 000002_create_queues_table.down.sql
│       └── 000002_create_queues_table.up.sql
├── pkg/
│   ├── broker/
│   │   ├── consumer.go
│   │   ├── producer.go
│   │   └── topic.go
│   ├── jwt/
│   │   └── jwt.go
│   ├── queue/
│   │   ├── queue.go
│   │   ├── repository.go
│   │   └── type.go
│   ├── response/
│   │   ├── ratelimiter.go
│   │   └── web_response.go
│   └── utils/
│       ├── password.go
│       └── utils.go
├── swag.sh
└── test/
    └── user/
        ├── handler_test.go
        ├── repository_test.go
        └── usecase_test.go
```
