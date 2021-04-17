FROM golang:1.16-alpine

WORKDIR /sms
ADD . /sms

RUN cd /sms && go build

EXPOSE 7080 7090
CMD ./sms
