FROM golang:1.12-alpine as builder
RUN apk --no-cache add ca-certificates git
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -installsuffix 'static' -o /lottsvc ./cmd/lottsvc

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /lottsvc /lottsvc
ENV PROJECT_ID=lott-235905
EXPOSE 8080
ENTRYPOINT ["/lottsvc"]
