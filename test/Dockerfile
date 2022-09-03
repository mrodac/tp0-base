FROM debian:stable-slim
RUN apt update -y && apt install -y netcat
COPY test/testserver.sh .
ENTRYPOINT ["/testserver.sh"]