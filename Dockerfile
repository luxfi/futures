FROM golang:1.27.1-alpine AS build
ENV GOTOOLCHAIN=auto
WORKDIR /src
RUN apk add --no-cache ca-certificates tzdata \
    && echo 'nonroot:x:65532:65532:nonroot:/home/nonroot:/sbin/nologin' >> /etc/passwd \
    && echo 'nonroot:x:65532:' >> /etc/group
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /futuresd ./cmd/futuresd

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /futuresd /futuresd
USER 65532:65532
EXPOSE 8090
ENTRYPOINT ["/futuresd"]
