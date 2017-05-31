FROM ubuntu:16.04
RUN apt-get update && apt-get install -y --force-yes golang
COPY . /src/github.com/ps/balancer
ENV GOPATH /
EXPOSE 9005
RUN go install github.com/ps/balancer
WORKDIR /
ENTRYPOINT ["./bin/balancer"]