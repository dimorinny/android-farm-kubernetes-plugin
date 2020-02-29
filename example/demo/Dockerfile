FROM sorccu/adb:latest

RUN apk add --no-cache bash

ADD app.apk .
ADD entrypoint.sh .

ENTRYPOINT ["./entrypoint.sh"]