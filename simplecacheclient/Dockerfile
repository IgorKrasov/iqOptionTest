FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/iqOptionTest/simplecacheclient

RUN go get -u github.com/kardianos/govendor
WORKDIR /go/src/github.com/iqOptionTest/simplecacheclient
RUN govendor sync
RUN govendor install

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/simplecacheclient

# Document that the service listens on port 8080.
EXPOSE 9004