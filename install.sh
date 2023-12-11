#!/bin/bash

set -e

owner="fitz7"
repo="tfnew"
bin_name="tfnew"


get_arch() {
    # darwin/amd64: Darwin axetroydeMacBook-Air.local 20.5.0 Darwin Kernel Version 20.5.0: Sat May  8 05:10:33 PDT 2021; root:xnu-7195.121.3~9/RELEASE_X86_64 x86_64
    # linux/amd64: Linux test-ubuntu1804 5.4.0-42-generic #46~18.04.1-Ubuntu SMP Fri Jul 10 07:21:24 UTC 2020 x86_64 x86_64 x86_64 GNU/Linux
    a=$(uname -m)
    case ${a} in
        "x86_64" | "amd64" )
            echo "amd64"
        ;;
        "aarch64" | "arm64" | "arm")
            echo "arm64"
        ;;
        *)
            echo ${NIL}
        ;;
    esac
}

get_os(){
    # darwin: Darwin
    echo $(uname -s | awk '{print tolower($0)}')
}

downloadFolder="${TMPDIR:-/tmp}"
mkdir -p ${downloadFolder}
os=$(get_os)
arch=$(get_arch)
file_name="${bin_name}_${os}_${arch}"
downloaded_file="${downloadFolder}/${file_name}" #
bin_folder="/usr/local/bin"

githubUrl="https://github.com"
githubApiUrl="https://api.github.com"

asset_path=$(
    command curl -sL \
        -H "Accept: application/vnd.github+json" \
        -H "X-GitHub-Api-Version: 2022-11-28" \
        ${githubApiUrl}/repos/${owner}/${repo}/releases |
    command grep -o "/${owner}/${repo}/releases/download/.*/${file_name}" |
    command head -n 1
)
if [[ ! "$asset_path" ]]; then
    echo "ERROR: unable to find a release asset called ${file_name}"
    exit 1
fi
asset_uri="${githubUrl}${asset_path}"

rm -f "${downloaded_file}"
sudo curl -s --fail --location --output "${bin_folder}/${bin_name}" "${asset_uri}"

bin="${bin_folder}/${bin_name}"
sudo chmod +x ${bin}

echo "${bin_name} was installed successfully to ${bin}"
