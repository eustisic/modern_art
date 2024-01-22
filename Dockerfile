# Base image
FROM golang:alpine

# current working directory
WORKDIR /app

# copy go mod
COPY go.mod go.sum ./

# download dependencies
RUN go mod download

# copy src code
COPY . .

# compile src code, output as s3-api
RUN go build -o s3-api .

# start application
CMD ["./modern_art"]
