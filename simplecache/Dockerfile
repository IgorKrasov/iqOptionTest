FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/iqOptionTest/simplecache

RUN go get -u github.com/kardianos/govendor
WORKDIR /go/src/github.com/iqOptionTest/simplecache
RUN govendor sync
RUN govendor install

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/simplecache

# Document that the service listens on port 8080.
EXPOSE 9003