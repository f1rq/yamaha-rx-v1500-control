#!/bin/sh
echo "[HOOK] Event: $PLAYER_EVENT" >&2
case "$PLAYER_EVENT" in
  play|playing|track_changed)
    curl -s -X POST http://127.0.0.1:8080/hook/spotify/play > /dev/null
    ;;
  pause|paused|stop|stopped)
    curl -s -X POST http://127.0.0.1:8080/hook/spotify/pause > /dev/null
    ;;
esac
