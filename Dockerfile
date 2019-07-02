FROM golang:alpine as build-env
ADD . /src
RUN cd /src/ && go build -o poe-market

EXPOSE 9092
FROM golang:alpine
COPY --from=build-env /src/poe-market /opt/poe-market
ENTRYPOINT [ "/opt/poe-market" ]
