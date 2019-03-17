#!/usr/bin/env bash

set -ex

SLEEP_DURATION=15

adb devices

adb install -t app.apk
adb shell am start -n "com.dimorinny.farm/com.dimorinny.farm.MainActivity" -a android.intent.action.MAIN -c android.intent.category.LAUNCHER -e title ${JOB_NAME:-Unknown} -e timer ${SLEEP_DURATION}
sleep ${SLEEP_DURATION}
adb shell pm uninstall -k com.dimorinny.farm