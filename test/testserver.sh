#!/bin/bash

while IFS= read -r line
do
	if [[ $line == server_port* ]]; then
		PORT=${line##* }
	fi
done < /config.ini

OUTPUT=$(echo ${1:-"HolaMundo"} | nc server $PORT)

if [[ $? -ne 0 ]]; then
	echo "Could not connect to server" >&2
else
	if [[ $OUTPUT == "Your Message has been received: ${1:-"HolaMundo"}" ]]; then
		echo "Server responded as expected: $OUTPUT" >&2
	else
		echo "Message different than expected: $OUTPUT" >&2
	fi
fi