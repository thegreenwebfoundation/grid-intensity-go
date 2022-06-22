#!/bin/sh
# Based on Deno installer: Copyright 2019 the Deno authors. All rights reserved. MIT license.

set -e

case $(uname -sm) in
"Darwin x86_64") target="Darwin_x86_64" ;;
"Darwin arm64") target="Darwin_arm64" ;;
"Linux aarch64") target="Linux_arm64" ;;
*) target="Linux_i386" ;;
esac

version=${1}

if [ $# -eq 0 ]; then
    version="$(curl -s https://api.github.com/repos/thegreenwebfoundation/grid-intensity-go/releases/latest | grep tag_name | cut -d '"' -f 4)"
fi

version="$(echo "$version" | cut -d 'v' -f 2)"
grid_intensity_uri="https://github.com/thegreenwebfoundation/grid-intensity-go/releases/download/v${version}/grid-intensity_${version}_${target}.tar.gz"

bin_dir="/usr/local/bin"
binary="grid-intensity"
exe="$bin_dir/grid-intensity"

curl --fail --location --progress-bar --output "$exe.tar.gz" "$grid_intensity_uri"
tar xzf "$exe.tar.gz" -C $bin_dir $binary
chmod +x "$exe"
rm "$exe.tar.gz"
