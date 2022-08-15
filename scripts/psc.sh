#!/usr/bin/env bash

OS="$(uname -s)"

usage() {
  cat <<EOT
Usage: $0 [option] command

Options:
   -d          download all the binaries
   -h          print this help
Note: Arguments must be passed at the end for ./psc to work correctly.
For instance: ./psc.sh balances <hex-address> --node=https://api.posichain.org/

EOT
}

set_download() {
  local rel='mainnet'
  case "$OS" in
  Darwin)
    BIN=(psc-arm64 libbls384_256.dylib libcrypto.1.0.0.dylib libgmp.10.dylib libgmpxx.4.dylib libmcl.dylib)
    ;;
  Linux)
    BIN=(psc)
    ;;
  *)
    echo "${OS} not supported."
    exit 2
    ;;
  esac
}

do_download() {
  # download all the binaries
  for bin in "${BIN[@]}"; do
    rm -f ${bin}
    curl https://download.posichain.org/latest/${bin} -o ${bin}
  done
  chmod +x psc*
}

while getopts "dh" opt; do
  case ${opt} in
  d)
    set_download
    do_download
    exit 0
    ;;
  h | *)
    usage
    exit 1
    ;;
  esac
done

shift $((OPTIND - 1))

if [ "$OS" = "Linux" ]; then
  ./psc "$@"
else
  DYLD_FALLBACK_LIBRARY_PATH="$(pwd)" ./psc "$@"
fi
