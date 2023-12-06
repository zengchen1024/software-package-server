#!/bin/sh

set -euo pipefail

work_dir=$1
main_branch=$2
pkg_name=$3

delete_branch() {
    git checkout -- .
    git clean -fd

    git checkout $main_branch

    set +e
        git rev-parse --verify "$pkg_name" 2>/dev/null
        if [ $? -eq 0 ]; then
            # delete remote
            git push origin -d $pkg_name

            # delete local
            git branch -D $pkg_name
            git reflog expire --expire=now --all
            git gc --prune=now --quiet
        fi
    set -e
}

cd $work_dir

delete_branch
