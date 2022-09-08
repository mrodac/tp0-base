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
    networks:
      - testing_net
    depends_on:
      - server"""      
}

render_compose_yml ${1:-1}

  
