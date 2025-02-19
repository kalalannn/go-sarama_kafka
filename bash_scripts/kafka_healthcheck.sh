#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source $SCRIPT_DIR/kafka_base.sh

healthcheck $BROKER_1 $BROKER_1_INTERNAL_ADDR
healthcheck $BROKER_2 $BROKER_2_INTERNAL_ADDR
healthcheck $BROKER_3 $BROKER_3_INTERNAL_ADDR
healthcheck $BROKER_4 $BROKER_4_INTERNAL_ADDR
healthcheck $BROKER_5 $BROKER_5_INTERNAL_ADDR