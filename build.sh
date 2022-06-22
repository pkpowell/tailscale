#!/usr/bin/env sh

npm --prefix ipn/ipnlocal/embeds/html run build > /dev/null && {
    echo "npm build succeeded"
} || {
    echo "npm build failed"
    exit $?
}

sh ./build_dist.sh tailscale.com/cmd/tailscaled && {
    echo "tailscale.com/cmd/tailscaled built"

} || {
    echo "build tailscale.com/cmd/tailscaled failed"
    exit $?
}



sh ./build_dist.sh tailscale.com/cmd/tailscale && {
    echo "tailscale.com/cmd/tailscale built"

} || {
    echo "build tailscale.com/cmd/tailscale failed"
    exit $?
}