#!/bin/bash

render_compose_yml () {
  echo """version: '3'
services:"""
  render_server
  for i in `seq $1`
  do
    render_client $i
  done
  
  echo """networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24""" 
}

render_server () {
  echo """  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    volumes:
      - ./config/server.ini:/config.ini
    environment:
      - PYTHONUNBUFFERED=1
      - SERVER_PORT=12345
      - SERVER_LISTEN_BACKLOG=7
      - LOGGING_LEVEL=DEBUG
    networks:
      - testing_net"""      
}

render_client () {
  echo """
  client${1}:
    container_name: client${1}
    image: client:latest
    entrypoint: /client
    volumes:
      - ./config/client.yaml:/config.yaml
    environment:
      - CLI_ID=${1}
      - CLI_SERVER_ADDRESS=server:12345
      - CLI_LOOP_LAPSE=1m2s
      - CLI_LOG_LEVEL=DEBUG
    networks:
      - testing_net
    depends_on:
      - server"""      
}

render_compose_yml ${1:-1}

  
