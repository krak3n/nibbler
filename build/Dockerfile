FROM golang:1.12-alpine
RUN apk update && apk add build-base git
WORKDIR /nibbler
COPY ./ /nibbler
RUN make nibbler

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /nibbler
COPY --from=0 /nibbler/bin/nibbler .
CMD ["./nibbler"]
