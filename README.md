# Kafka & Sarama Demo parallel processing (exactly-once consumer group)

### Get started
* Build:
```bash
make up           # Start kafka cluster
make healthcheck  # Wait for healthy
make init_kafka   # Create topics as it is defined in kafka_init.sh
make build        # Build producer and consumer binaries
```

* Run: 
```bash
./bin/producer    # Producer binary
./bin/consumer    # Consumer binary
```

Refer to `bash_scripts` and `Makefile` for more commands