FROM sorccu/adb:latest

RUN apk --no-cache add socat bash

COPY entrypoint.sh /

EXPOSE 5037
EXPOSE 5554
EXPOSE 5555

ENTRYPOINT ["/entrypoint.sh"]