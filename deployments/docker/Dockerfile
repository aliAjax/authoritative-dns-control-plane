FROM golang:1.23 AS build
WORKDIR /src
COPY . .
RUN go build -trimpath -ldflags='-s -w' -o /out/control-plane ./cmd/control-plane
FROM gcr.io/distroless/base-debian12
COPY --from=build /out/control-plane /control-plane
EXPOSE 8080 8053/udp
ENTRYPOINT ["/control-plane"]

