#!/bin/sh

set -euo pipefail

repo_dir=$1
main_branch=$2
branch_name=$3

delete_branch() {
    git checkout -- .
    git clean -fd

    git checkout $main_branch

    set +e
        git rev-parse --verify "$branch_name" 2>/dev/null
        if [ $? -eq 0 ]; then
            # delete remote
            git push origin -d $branch_name

            # delete local
            git branch -D $branch_name
            git reflog expire --expire=now --all
            git gc --prune=now --quiet
        fi
    set -e
}

cd $repo_dir

delete_branch
