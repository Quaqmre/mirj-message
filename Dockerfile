FROM golang:1.12

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY ./go.mod .
COPY ./go.sum .

# Download all the dependencies
RUN go mod download
COPY . .
RUN  CGO_ENABLED=0 go build -o ./mirjmessage 

FROM alpine
COPY --from=0 /app/mirjmessage /bin/mirjmessage
EXPOSE 9001
ENTRYPOINT [ "/bin/mirjmessage" ]