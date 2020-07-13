# Stage 1 - Start from the latest golang base image
FROM golang:latest as builder1

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2 - build frontend assets
FROM node:lts-alpine as builder2

WORKDIR /app

COPY --from=builder1 /app/web .

COPY /web/package*.json ./

RUN npm install
COPY . .
RUN npm run build

# Stage 3 - start a new stage from scratch #######
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder1 /app/main .
COPY --from=builder2 /app/dist ./web/dist

# Expose port 8080 to the outside world
EXPOSE 8080

# pass env vars???
# Command to run the executable
CMD [ "./main" ]
