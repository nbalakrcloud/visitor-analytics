# Latest golang image on apline linux
FROM golang:1.19-alpine

# Work directory
WORKDIR /app

# Copying all the files
COPY . ./

# Installing dependencies
COPY go.mod go.sum ./
RUN go mod download

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

RUN go build -o visitor-analytics ./cmd
# Exposing server port
EXPOSE 8080
# Starting our application
CMD ["./visitor-analytics"]

