FROM golang:1.12 AS builder
WORKDIR /letterdrop
COPY . .
RUN CGO_ENABLED=0 go build -a -installsuffix cgo

FROM alpine
CMD ["/letterdrop"]
COPY --from=builder /letterdrop/public /public
COPY --from=builder /letterdrop/templates /templates
COPY --from=builder /letterdrop/letterdrop /letterdrop
