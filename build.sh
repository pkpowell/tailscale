#!/usr/bin/env sh

npm --prefix ipn/ipnlocal/embeds/html i > /dev/null && {
    echo "npm install ok"
} || {
    echo "npm install failed"
    exit $?
}

npm --prefix ipn/ipnlocal/embeds/html run build > /dev/null && {
    echo "npm build ok"
} || {
    echo "npm build failed"
    exit $?
}

sh ./build_dist.sh tailscale.com/cmd/tailscaled && {
    echo "tailscale.com/cmd/tailscaled ok"

} || {
    echo "build tailscale.com/cmd/tailscaled failed"
    exit $?
}



sh ./build_dist.sh tailscale.com/cmd/tailscale && {
    echo "tailscale.com/cmd/tailscale ok"

} || {
    echo "build tailscale.com/cmd/tailscale failed"
    exit $?
}