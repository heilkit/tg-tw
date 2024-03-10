FROM golang:1.22-alpine3.19 as BUILD

WORKDIR /build

COPY go.mod go.sum ./
RUN go get -u -v -f all
COPY telegram ./telegram
COPY twitter ./twitter
COPY cmd ./cmd
RUN go build -o tgtw ./cmd


FROM alpine:3.19

WORKDIR /app

COPY --from=BUILD /build/tgtw ./tgtw

ENTRYPOINT ["/app/tgtw"]