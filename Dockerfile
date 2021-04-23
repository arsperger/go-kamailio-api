FROM golang:1.16.3 as builder
ARG PROJECT
ARG GITREPO
ARG COMMIT
ARG VERSION
COPY go.mod go.sum /go/src/code/
WORKDIR /go/src/code
RUN go mod download
COPY . /go/src/code
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X ${GITREPO}/internal/utils.commit=${COMMIT} \
	-X ${GITREPO}/internal/utils.version=${VERSION}" -o ${PROJECT}

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/code/${PROJECT} /usr/bin/${PROJECT}
COPY --from=builder /go/src/code/${PROJECT}/config.json /opt/
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/go-kamailio-api", "-config", "/opt/config.json"]
