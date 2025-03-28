# build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/gfstat
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/gfstat -v .

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/gfstat /gfstat
COPY views /views
WORKDIR /
ENV GH_BASIC_CLIENT_ID=${GH_BASIC_CLIENT_ID}
ENV GH_BASIC_SECRET_ID=${GH_BASIC_SECRET_ID}
ENTRYPOINT ["/gfstat"]
LABEL Name=gfstat Version=0.0.1
EXPOSE 3639