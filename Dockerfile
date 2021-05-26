FROM golang:1.16-alpine AS build

WORKDIR /src/
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
COPY main.go .env go* /src/
RUN CGO_ENABLED=0 go build -o /bin/sync

FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/
COPY --from=build /bin/sync /bin/sync
ENTRYPOINT ["/bin/sync"]