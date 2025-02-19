#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source $SCRIPT_DIR/../.env

function healthcheck {
    container_name=$1
    bootstrap_server=$2
    if docker exec -it $container_name kafka-cluster \
        cluster-id \
        --bootstrap-server $bootstrap_server > /dev/null;
    then
        echo "$container_name is healthy."
    else
        echo "$container_name is unhealthy."
    fi
}

function create_topic {
    topic_name=$1
    partitions_count=$2
    docker exec -it $BROKER_1 kafka-topics \
        --create --topic $topic_name --partitions $partitions_count --replication-factor 3 \
        --bootstrap-server $BOOTSTRAP_SERVER_ALL
}

function delete_topic {
    topic_name=$1
    docker exec -it $BROKER_1 kafka-topics \
        --delete --topic $topic_name \
        --bootstrap-server $BOOTSTRAP_SERVER_ALL
}

function list_topics {
	docker exec -it $BROKER_1 kafka-topics \
		--list \
		--bootstrap-server $BOOTSTRAP_SERVER_ALL
}

function describe_topic {
    topic_name=$1
    docker exec -it $BROKER_1 kafka-topics \
        --describe --topic $topic_name \
        --bootstrap-server $BOOTSTRAP_SERVER_ALL
}

function delete_records {
    json_file=$1
    cat $json_file | \
    docker exec -i $BROKER_1 kafka-delete-records \
        --offset-json-file /dev/stdin \
        --bootstrap-server $BOOTSTRAP_SERVER_ALL
}


export -f healthcheck create_topic list_topics describe_topic delete_records