include .env
export

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

list_topics:
	@./bash_scripts/kafka_list_topics.sh

#! Mount persistant volumes
MNT_1=./mnt/kafka-data-1
MNT_2=./mnt/kafka-data-2
MNT_3=./mnt/kafka-data-3
DOCKER_VOLUMES_DIR=/var/lib/docker/volumes
DOCKER_VOLUME_PREFIX=go-kafka_events_kafka-data
_mkdir_mnt:
	mkdir -p ${MNT_1} ${MNT_2} ${MNT_3}
mount_volumes: _mkdir_mnt
	sudo mount -t none -o bind ${DOCKER_VOLUMES_DIR}/${DOCKER_VOLUME_PREFIX}-1/_data ${MNT_1}
	sudo mount -t none -o bind ${DOCKER_VOLUMES_DIR}/${DOCKER_VOLUME_PREFIX}-2/_data ${MNT_2}
	sudo mount -t none -o bind ${DOCKER_VOLUMES_DIR}/${DOCKER_VOLUME_PREFIX}-3/_data ${MNT_3}
unmount_volumes:
	sudo umount .${MNT_1}
	sudo umount .${MNT_2}
	sudo umount .${MNT_3}