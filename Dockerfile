FROM golang:1.13.0

EXPOSE 8080
#Copy app
RUN mkdir /local_dir
COPY ./ /local_dir
RUN mv /local_dir /backend

WORKDIR /backend

# Create app directories
RUN mkdir /go/bin/config && mv ./config/.env /go/bin/config/.env && go mod download; cat /go/bin/config/.env

# Download dependencies
RUN go mod download

# Build main app
WORKDIR /backend/cmd/api
RUN go build -o /go/bin/backend main.go

RUN chmod +x /go/bin/backend

WORKDIR /go/bin
# Copy migrations
RUN mkdir -p /go/bin/module/db/migrations
COPY  /module/db/migrations /go/bin/module/db/migrations

#Copy email templates
RUN mkdir -p /go/bin/services/transport/messages/templates
COPY  /services/transport/messages/templates /go/bin/services/transport/messages/templates

USER 1001:1001

CMD ["backend"]
