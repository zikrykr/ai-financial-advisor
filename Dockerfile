FROM golang:1.22-alpine AS builder

WORKDIR /app
RUN export GO111MODULE=on

COPY go.mod go.sum ./

# install modules
RUN go mod download

COPY . .
RUN go build -o ai-financial-advisor .

FROM alpine:3.18
WORKDIR /app

COPY --from=builder /app/ai-financial-advisor .

ENTRYPOINT ["./ai-financial-advisor"]