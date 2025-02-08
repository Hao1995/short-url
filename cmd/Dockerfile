FROM golang:1.23.6-alpine

# Specify the working path
WORKDIR /user/src/app

# Download dependency
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY pkg ./pkg
COPY database ./database
COPY cmd ./cmd
COPY internal ./internal

RUN go build -o app ./cmd

# Run APP
CMD ["./app"]