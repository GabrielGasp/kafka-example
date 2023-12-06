### Basic Kafka example using Golang

#### Requirements to run this example:
- Docker
- Docker Compose
- Golang 1.21 or higher

#### Steps to run this example:

1. Bootstrap Kafka cluster and create a topic
```bash
make setup
```

2. Run consumer
```bash
make run-consumer
```

3. Run producer (in a different terminal tab/window)
```bash
make run-producer
```

4. Tear down Kafka cluster after you are done.
```bash
make teardown
```
