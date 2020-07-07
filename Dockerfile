FROM golang:alpine as builder

RUN mkdir /build
RUN apk add --update gcc musl-dev
ADD . /build/
WORKDIR /build
RUN go build -o terraform-server-state main.go

FROM alpine
RUN adduser -S -D -H tss
RUN mkdir -p /var/terraform-server-state/registry
RUN chown tss -R /var/terraform-server-state
USER tss
COPY --from=builder /build/terraform-server-state /usr/local/bin
COPY --from=builder /build/config.yaml /etc/

VOLUME /var/terraform-server-state
EXPOSE 8080

CMD ["terraform-server-state"]
