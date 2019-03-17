FROM sorccu/adb

RUN apk add --no-cache bash

ADD app.apk .
ADD entrypoint.sh .

ENTRYPOINT ["./entrypoint.sh"]