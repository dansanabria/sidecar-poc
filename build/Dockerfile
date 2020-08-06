FROM golang

ADD ./cmd/client_sidecar /go/src/github.com/dansanabria/sidecar-poc/cmd/client_sidecar

RUN go install -v github.com/dansanabria/sidecar-poc/cmd/client_sidecar

CMD ["client_sidecar"]
