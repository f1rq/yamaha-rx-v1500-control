# Yamaha RX-V1500 RS-232 Controller

A lightweight Go microservice for controlling the **Yamaha RX-V1500** AV receiver over RS-232C and automating power management with **Spotify Connect ([librespot](https://github.com/librespot-org/librespot))**.

## Features
- **Spotify Connect Automation:** Automatically wakes the receiver from Standby, switches the input to `CD-R`, and enables `Speaker A` when playback starts.
- **Auto-Standby:** Powers down the receiver after 10 minutes of playback inactivity.
- **Reliable Serial Transport:** Thread-safe RS-232 communication (`9600 8N1`) with automatic controller wake-up frames (`DC1` / `\x11000\x03`).
- **REST API:** Simple HTTP endpoints for controlling power, Speaker A/B zones, and exact volume in dB.

## HTTP Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/power?on=true\|false` | Powers the receiver on or puts it into Standby |
| `POST` | `/api/speaker/a?on=true\|false` | Toggles Speaker A zone |
| `POST` | `/api/speaker/b?on=true\|false` | Toggles Speaker B zone |
| `POST` | `/api/volume?db=-40.0` | Sets master volume (`-80.0` to `0.0` dB, `0.5` dB steps) |
| `POST` | `/hook/spotify/play` | Webhook for `librespot` playback start |
| `POST` | `/hook/spotify/pause` | Webhook for `librespot` pause/stop (starts 10m sleep timer) |

## Docker Compose Setup

```yaml
services:
  yamaha-api:
    build: .
    container_name: yamaha-rx-v1500-api
    restart: unless-stopped
    network_mode: host
    devices:
      - "/dev/ttyUSB0:/dev/ttyUSB0"
    environment:
      - SERIAL_PORT=/dev/ttyUSB0

  librespot:
    image: giof71/librespot:latest
    container_name: librespot
    restart: unless-stopped
    network_mode: host
    depends_on:
      - yamaha-api
    devices:
      - "/dev/snd:/dev/snd"
    environment:
      - DEVICE_NAME=Yamaha RX-V1500
      - DEVICE_TYPE=avr
      - BITRATE=320
      - BACKEND=alsa
      - DEVICE=plughw:0,0
      - INITIAL_VOLUME=80
      - VOLUME_CTRL=softvol
      - ADDITIONAL_ARGUMENTS=--onevent /scripts/hook.sh
    volumes:
      - ./scripts/hook.sh:/scripts/hook.sh:ro
