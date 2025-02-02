#!/bin/bash
set -euo pipefail

# Go to the project root.
cd ..

# Base output directory.
BUILD_DIR="./bin/flex6k-discovery-util-go-build"

build_target() {
  local os="$1"
  local arch="$2"
  local extra="$3"
  local out="$4"
  local build_cmd="env GOOS=${os} GOARCH=${arch}"
  
  # Append extra environment variable if provided.
  if [[ -n "$extra" ]]; then
    build_cmd+=" ${extra}"
  fi

  # Complete the build command.
  build_cmd+=" go build -o \"${BUILD_DIR}/${out}\" ."
  echo "Building for ${os}/${arch} ${extra:+with ${extra}} -> ${BUILD_DIR}/${out}"
  eval "${build_cmd}"
}

targets=(
  "linux;amd64;;linux64/flexi"
  "linux;386;;linux32/flexi"
  "linux;arm;GOARM=5;raspberryPi/flexi"
  "windows;amd64;;Win64/flexi.exe"
  "windows;386;;Win32/flexi.exe"
  "freebsd;amd64;;pfSense64/flexi"
  "freebsd;386;;pfSense32/flexi"
)

# Create necessary output directories.
for target in "${targets[@]}"; do
  IFS=";" read -r os arch extra out <<< "$target"
  mkdir -p "$(dirname "${BUILD_DIR}/${out}")"
done

# Build for each target.
for target in "${targets[@]}"; do
  IFS=";" read -r os arch extra out <<< "$target"
  build_target "$os" "$arch" "$extra" "$out"
done
