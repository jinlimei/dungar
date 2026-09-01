ARG BUILDPLATFORM

FROM --platform=$BUILDPLATFORM golang:1.21-bookworm AS build

WORKDIR /src

ARG TARGETOS
ARG TARGETARCH

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN target_os="${TARGETOS:-$(go env GOOS)}" && \
    target_arch="${TARGETARCH:-$(go env GOARCH)}" && \
    CGO_ENABLED=0 GOOS="$target_os" GOARCH="$target_arch" \
    go build -mod=mod -o /out/dungar -ldflags="-s -w" ./cmd/dungar

FROM scratch

WORKDIR /app

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /out/dungar /usr/local/bin/dungar

USER 65532:65532

ENTRYPOINT ["/usr/local/bin/dungar"]
CMD ["run"]
