FROM golang:latest as build-env
ADD . /tmp/poe
RUN cd /tmp/poe && CGO_ENABLED=0 go build -o poe-market cmd/*.go

EXPOSE 9092
FROM golang:alpine
COPY --from=build-env /tmp/poe/poe-market /opt/poe-market
RUN chmod 0755 /opt/poe-market
ENTRYPOINT [ "/opt/poe-market" ]
