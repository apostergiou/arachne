FROM golang:1.11.2 as builder

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

# Copy the code from the host
RUN mkdir /app
RUN mkdir /go/src/app
ADD . /go/src/app
WORKDIR /go/src/app
COPY Gopkg.toml Gopkg.lock /go/src/app/
RUN dep ensure --vendor-only -v

EXPOSE 53/tcp

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app/arachne .
CMD ["/app/arachne"]
