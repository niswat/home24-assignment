FROM golang:1.16.5 as builder

# Setting the necessary Environment variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /src

# Copying and downloading dependency using go mod
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# Copying code into the container
COPY . .

# Build the application
##RUN go build -o main . 
RUN go build -o webtool main.go

# Move to /dest directory as the place for resulting binary folder
WORKDIR /dest

# Copy binary from build to main folder
##RUN cp /src/main .
RUN cp /src/webtool .

FROM alpine:latest

WORKDIR /
##COPY --from=builder /dest/main .
COPY --from=builder /dest/webtool .

# Adding ssl certificates
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY ./certificate.crt /usr/local/share/ca-certificates/certificate.crt
RUN update-ca-certificates

# copying the html templates
COPY index.html index.html
COPY parse.html parse.html

EXPOSE 8082

##ENTRYPOINT ["/main"]
ENTRYPOINT ["/webtool"]






