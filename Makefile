include .env
export

#! Go local
_mkdir_bin:
	mkdir -p bin
_build_producer:
	@go build -o ./bin/producer ./cmd/producer/main.go
_build_consumer:
	@go build -o ./bin/consumer ./cmd/consumer/main.go
build: _mkdir_bin _build_producer _build_consumer

run_producer: _build_producer
	./bin/producer
run_consumer: _build_consumer
	./bin/consumer

#! Docker-compose
up:
	docker-compose up -d
	make status
down:
	docker-compose down
status:
	docker-compose ps

#! Kafka
healthcheck:
	@./bash_scripts/kafka_healthcheck.sh

init_kafka:
	@./bash_scripts/kafka_init.sh

reinit_kafka:
	@./bash_scripts/kafka_delete_topic.sh
	@./bash_scripts/kafka_init.sh

list_topics:
	@./bash_scripts/kafka_list_topics.sh

describe_topic:
	@./bash_scripts/kafka_describe_topic.sh

delete_records:
	@./bash_scripts/kafka_delete_records.sh

#! Mount persistant volumes
MNT_1=./mnt/kafka-data-1
MNT_2=./mnt/kafka-data-2
MNT_3=./mnt/kafka-data-3
MNT_4=./mnt/kafka-data-3
MNT_5=./mnt/kafka-data-3
DOCKER_VOLUMES_DIR=/var/lib/docker/volumes
DOCKER_VOLUME_PREFIX=go-kafka_events_kafka-data
_mkdir_mnt:
	mkdir -p ${MNT_1} ${MNT_2} ${MNT_3} ${MNT_4} ${MNT_5}
mount_volumes: _mkdir_mnt
	sudo mount -t none -o bind ${DOCKER_VOLUMES_DIR}/${DOCKER_VOLUME_PREFIX}-1/_data ${MNT_1}
	sudo mount -t none -o bind ${DOCKER_VOLUMES_DIR}/${DOCKER_VOLUME_PREFIX}-2/_data ${MNT_2}
	sudo mount -t none -o bind ${DOCKER_VOLUMES_DIR}/${DOCKER_VOLUME_PREFIX}-3/_data ${MNT_3}
	sudo mount -t none -o bind ${DOCKER_VOLUMES_DIR}/${DOCKER_VOLUME_PREFIX}-4/_data ${MNT_4}
	sudo mount -t none -o bind ${DOCKER_VOLUMES_DIR}/${DOCKER_VOLUME_PREFIX}-5/_data ${MNT_5}
unmount_volumes:
	sudo umount .${MNT_1}
	sudo umount .${MNT_2}
	sudo umount .${MNT_3}
	sudo umount .${MNT_4}
	sudo umount .${MNT_5}