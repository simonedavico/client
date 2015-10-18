package deploy

import (

)

//TODO; decide if more to Rancher if I need load balancer services

//TODO:
//

//Manage the configuration starting from the docker images, building the go yaml representation of the docker-compose, then
//unsmarhall it to send to docker-compose
//https://github.com/go-yaml/yaml

//Also the variables and dynamics IPs etc has to be generated dinamycally...and the I may also be able to add two brokers on two server using swarm, eg:

// Solve problem with permission on writing to shared volume from kafka

//Note: add containers to the rancher network: io.rancher.container.network=true (soon there will be the docker custom network, so I can create the benchflow network)

//Write a method to parse the docker-compose.yml file and do the following customizations:
//modify the KAFKA_ADVERTISED_HOST_NAME in docker-compose.yml to match your docker host IP (Note: Do not use localhost or 127.0.0.1 as the host ip if you want to run multiple brokers.)
// = neha IP
// KAFKA_MESSAGE_MAX_BYTES: 1000000000 (1gb)
// Set the broker name (alias for the consumers)

//Docker-compose file:
//zookeeper:
//  image: jplock/zookeeper:3.4.6
//  hostname: zookeeper
//  ports:
//    - "2181:2181"
//  volumes:
//    #- /opt/zookeeper/conf:/opt/zookeeper/conf:rw
//    - /opt/zookeeper/tmp:/tmp/zookeeper:rw
//  environment:
//    - "constraint:node==SERVER_HOSTNAME"
//  labels:
//    - "io.rancher.container.network=true"
//kafka:
//  image: ches/kafka:0.8.1.1-1
//  hostname: kafka
//  #privileged: true
//  ports:
//    - "9092:9092"
//  links:
//    - zookeeper:zookeeper
//  #TODO, enable shared volumes. It has a permission error while trying to use the shared volumes
//  #volumes:
//    #- /opt/kafka/data:/data:rw
//    #- /opt/kafka/logs:/logs:rw
//  environment:
//    - "EXPOSED_HOST=SERVER_IP"
//    - "EXPOSED_PORT=9092"
//    - "BROKER_ID=0"
//    - "constraint:node==SERVER_HOSTNAME"
//  labels:
//    - "io.rancher.container.network=true"


//Here I need to configure also the topics, but maybe I can pass a ENV variable to let the image automatically create topics


//Usage Test:
//docker run --rm ches/kafka kafka-topics.sh --create --topic test --replication-factor 1 --partitions 1 --zookeeper $ZK_IP:2181
//
//
//docker run --rm --interactive ches/kafka kafka-console-producer.sh --topic test --broker-list $KAFKA_IP:9092
//<type some messages followed by newline>
//
//
//docker run --rm ches/kafka kafka-console-consumer.sh --topic test --from-beginning --zookeeper $ZK_IP:2181