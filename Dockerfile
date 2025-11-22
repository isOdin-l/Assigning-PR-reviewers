FROM golang:1.25.4

WORKDIR /

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /pr-assigning ./cmd/Assigning-PR/main.go

FROM ubuntu:24.04
COPY --from=0 /pr-assigning /bin/pr-assigning

CMD ["/bin/pr-assigning"]