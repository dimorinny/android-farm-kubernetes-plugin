FROM golang:1.11-alpine as builder

RUN apk add --no-cache gcc musl-dev alpine-sdk libusb-dev

WORKDIR /plugin
COPY . .

RUN go build -mod vendor


FROM alpine

RUN apk add --no-cache libusb-dev

COPY --from=builder /plugin/android-devices-kubernetes-plugin /usr/bin/android-devices-kubernetes-plugin

CMD ["/usr/bin/android-devices-kubernetes-plugin"]