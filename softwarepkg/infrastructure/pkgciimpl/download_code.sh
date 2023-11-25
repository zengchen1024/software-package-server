#!/bin/sh

set -euo pipefail

repo_dir=$1
git_token=$2
master_branch=$3
branch_name=$4
spec_url=$5
spec_file_name=$6
srpm_url=$7
srpm_file_name=$8

wanted_spec_name="${branch_name}.spec"
wanted_srpm_name="${branch_name}.src.rpm"

srpm_files_dir="./code"
files_in_srpm="./files_in_srpm.txt"
lfs_size="50M"

has_submitted=""

checkout_branch() {
    git checkout -- .
    git clean -fd

    git checkout $master_branch

    set +e
    git rev-parse --verify "$branch_name" 2>/dev/null
    local has=$?
    set -e

    if [ $has -eq 0 ]; then
        git checkout "$branch_name"
    else
        git checkout -b "$branch_name"
    fi

    has_submitted=$(git ls-files $wanted_spec_name)
}

download() {
    local url=$1

    if [[ $url == *"gitee.com"* ]]; then
        /opt/app/download "$url" "${git_token}"
    else
        curl -LO "$url"
    fi
}

handle_spec() {
    if [ "$spec_url" != "-" ]; then
        download $spec_url

        if [ "$spec_file_name" != "$wanted_spec_name" ]; then
            mv $spec_file_name $wanted_spec_name
        fi
    fi
}

handle_srpm() {
    if [ "$srpm_url" = "-" ]; then
        return
    fi

    # download
    download $srpm_url

    if [ "$srpm_file_name" != "$wanted_srpm_name" ]; then
        mv $srpm_file_name $wanted_srpm_name
    fi

    if [ -z "$(git ls-files -om $wanted_srpm_name)" ]; then
        # srpm does not change
        return
    fi

    # if srpm file is lfs, echo it
    find $wanted_srpm_name -size +$lfs_size -print

    # delete the files of last srpm
    if [ -n "$(ls $files_in_srpm)" ]; then
        while read f
        do
            if [[ "$i" ~= *".spec" ]]; then
                rm -f $f
            fi
        done < $files_in_srpm
    fi

    > $files_in_srpm

    test -d $srpm_files_dir || mkdir $srpm_files_dir

    rpm2cpio $wanted_srpm_name | cpio -div --quiet -D $srpm_files_dir > $files_in_srpm 2>&1

    if [ -n "$(ls $srpm_files_dir/*.spec)" ]; then
        rm $srpm_files_dir/*.spec
        mv $srpm_files_dir/* .
    fi
}

commit() {
    # nothing changed
    if [ -z "$(git ls-files -om)" ]; then
        return
    fi

    # track lfs, ignore .git dir
    lfs=$(find . -path '*/.git' -prune -o -type f -size +$lfs_size -print)
    if [ -n "$lfs" ]; then
        for item in "${lfs[@]}"
        do
            git lfs track --filename $item
        done
    fi

    git add .

    if [ -z "$has_submitted" ]; then
        git commit -m "apply new package $branch_name"
    else
        git commit --amend --no-edit --quiet
    fi

    git push -f origin "$branch_name"

    git checkout $master_branch
}

cd $repo_dir

checkout_branch

# must call handle_srpm before handle_spec
handle_srpm

handle_spec

commit
