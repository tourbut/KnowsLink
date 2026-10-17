#!/bin/sh
# Mixkit audio is not committed (license: no standalone redistribution). Re-download and verify.
set -e
cd "$(dirname "$0")/assets/audio"
get() { curl -fsSL -o "$1" "$2"; }
get mixkit-house-vibes-129.mp3 https://assets.mixkit.co/music/129/129.mp3
get mixkit-dry-gin-868.mp3     https://assets.mixkit.co/music/868/868.mp3
get mixkit-shame-553.mp3       https://assets.mixkit.co/music/553/553.mp3
get sfx-click-216.mp3          https://assets.mixkit.co/active_storage/sfx/216/216-preview.mp3
get sfx-key-2841.mp3           https://assets.mixkit.co/active_storage/sfx/2841/2841-preview.mp3
get sfx-switch-2579.mp3        https://assets.mixkit.co/active_storage/sfx/2579/2579-preview.mp3
get sfx-notify-2357.mp3        https://assets.mixkit.co/active_storage/sfx/2357/2357-preview.mp3
sha256sum -c ../SHA256SUMS
