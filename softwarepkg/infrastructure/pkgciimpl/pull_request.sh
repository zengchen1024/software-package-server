#!/bin/sh

set -euo pipefail

repo_dir=$1
git_token=$2
branch_name=$3
pkg_info_file=$4
spec_url=$5
src_rpm_url=$6

new_branch() {
    git checkout -- .
    git clean -fd

    git checkout master

    git fetch origin master
    git rebase origin/master

    set +e
    git rev-parse --verify $branch 2>/dev/null
    has=$?
    set -e

    if [ $has -eq 0 ]; then
        git branch -D "$branch_name"
    fi

    git checkout -b "$branch_name"
}

modify() {
    # add pkginfo.yaml
    mv $pkg_info_file .

    # download spec file
    curl -LO "$spec_url"

    # download source rpm
    if [[ $src_rpm_url == *"gitee.com"* ]]; then
        /opt/app/download "$src_rpm_url" "${git_token}"
    else
        curl -LO "$src_rpm_url"
    fi
}

commit() {
    git add .

    git commit -m 'apply new package ci pull request'

    git push origin "$branch_name"

    git checkout master

    git branch -D "$branch_name"
}

cd $repo_dir

new_branch

modify

commit
