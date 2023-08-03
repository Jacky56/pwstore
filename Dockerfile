FROM golang:1.20-alpine as builder 

WORKDIR /build

COPY . .
RUN go mod download

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o webserver . 

# RUN mkdir ./statics
# RUN mkdir ./templates
# RUN mkdir ./configs
# COPY ./statics ./statics
# COPY ./templates ./templates
# COPY ./configs ./configs
# COPY ./webserver ./webserver
FROM alpine:latest as release

COPY --from=builder "/build/webserver" "/"
COPY --from=builder "/build/configs" "/configs/"
COPY --from=builder "/build/templates" "/templates/"
COPY --from=builder "/build/statics" "/statics/"


EXPOSE 8080
ENTRYPOINT ["/webserver", "-env=prod", "-port=8080"]


