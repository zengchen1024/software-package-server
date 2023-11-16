#!/bin/sh

set -euo pipefail

repo_dir=$1
git_token=$2
master_branch=$3
branch_name=$4
spec_url=$5
src_rpm_url=$6

new_branch() {
    git checkout -- .
    git clean -fd

    git checkout $master_branch

    set +e
    git rev-parse --verify "$branch_name" 2>/dev/null
    has=$?
    set -e

    if [ $has -eq 0 ]; then
        git checkout "$branch_name"
    else
        git checkout -b "$branch_name"
    fi
}

download() {
    url=$1

    if [[ $url == *"gitee.com"* ]]; then
        /opt/app/download "$url" "${git_token}"
    else
        curl -LO "$url"
    fi
}

modify() {
    ignore="-" 

    # download spec file
    if [ "$spec_url" != "$ignore" ]; then
        download $spec_url
    fi

    # download source rpm
    if [ "$src_rpm_url" != "ignore" ]; then
        download $src_rpm_url

	rpm2cpio *.rpm | cpio -div
    fi
}

commit() {
    git add .

    git commit -m "apply new package $branch_name"

    git push -f origin "$branch_name"

    git checkout $master_branch
}

cd $repo_dir

new_branch

modify

commit
