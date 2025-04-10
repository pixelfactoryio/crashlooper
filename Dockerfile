FROM alpine:3.13

COPY crashlooper_*.apk /tmp/

RUN apk add --allow-untrusted /tmp/crashlooper_*.apk \
  && rm -fr /tmp/crashlooper_*.apk

ENTRYPOINT ["/usr/bin/crashlooper"]
