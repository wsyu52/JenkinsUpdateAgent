FROM golang:1.18 as builder
WORKDIR /root/
COPY go.mod go.sum ./
COPY src ./src
RUN go env -w GOPROXY=https://goproxy.cn && go mod tidy && CGO_ENABLED=0 go build src/JenkinsUpdateAgent.go

FROM scratch
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /root/JenkinsUpdateAgent ./
EXPOSE 8888
CMD ["/JenkinsUpdateAgent"]
