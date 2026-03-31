FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /futuresd ./cmd/futuresd

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /futuresd /futuresd
EXPOSE 8090
ENTRYPOINT ["/futuresd"]
