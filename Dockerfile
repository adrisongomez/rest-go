ARG GO_VERSION=1.16.6

FROM golang:${GO_VERSION}-alpine AS builder

RUN go env -w GOPROXY=direct
RUN apk add --no-cache git
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY ./ /app/

WORKDIR /app/

RUN go mod download

RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /server-golang .

####################### [ builder end ] ##################################3

FROM scratch AS runner

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /ect/ssl/certs/

COPY .env ./

COPY --from=builder /server-golang /

EXPOSE 5050

ENTRYPOINT ["/server-golang"]



