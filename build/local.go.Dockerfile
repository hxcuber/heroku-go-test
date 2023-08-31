FROM golang:1.18

RUN apt-get update

RUN GO111MODULE=on go install golang.org/x/tools/cmd/goimports@v0.1.10

RUN GO111MODULE=on go install github.com/volatiletech/sqlboiler/v4@v4.8.6 && \
    GO111MODULE=on go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@v4.8.6 \
