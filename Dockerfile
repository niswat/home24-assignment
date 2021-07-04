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
RUN go build -o webtool main.go

# Move to /dest directory as the place for resulting binary folder
WORKDIR /dest

# Copy binary to dest folder
RUN cp /src/webtool .

FROM scratch

WORKDIR /
COPY --from=builder /dest/webtool .

# copying the html templates
COPY index.html index.html

EXPOSE 8082

ENTRYPOINT ["/webtool"]






