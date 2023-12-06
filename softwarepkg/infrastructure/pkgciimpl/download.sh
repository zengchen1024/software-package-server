#!/bin/sh

set -euo pipefail

work_dir=$1
git_token=$2
main_branch=$3
pkg_name=$4
spec_url=$5
spec_file_name=$6
srpm_url=$7
srpm_file_name=$8
code_changed_tag=$9
srpm_file_lfs_tag=${10}

wanted_spec_name="${pkg_name}.spec"
wanted_srpm_name="${pkg_name}.src.rpm"

srpm_files_dir="./code"
files_in_srpm="./files_in_srpm.txt"
lfs_size="50M"

has_submitted=""

checkout_branch() {
    git checkout -- .
    git clean -fd

    git checkout $main_branch

    set +e
    git rev-parse --verify "$pkg_name" 2>/dev/null
    local has=$?
    set -e

    if [ $has -eq 0 ]; then
        git checkout "$pkg_name"
    else
        git checkout -b "$pkg_name"
    fi

    has_submitted=$(git ls-files $wanted_spec_name)
}

download() {
    local url=$1

    if [[ "$url" =~ "gitee.com" ]]; then
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
    if [ -n "$(find $wanted_srpm_name -size +$lfs_size -print)" ]; then
        echo "$srpm_file_lfs_tag"
    fi

    # delete the files of last srpm
    if [ -f "$files_in_srpm" ]; then
        while read fn
        do
            if [[ ! "$fn" =~ ".spec" ]]; then
                rm -f $fn
            fi
        done < $files_in_srpm

        > $files_in_srpm
    fi

    # decompress
    test -d $srpm_files_dir || mkdir $srpm_files_dir

    rpm2cpio $wanted_srpm_name | cpio -div --quiet -D $srpm_files_dir > $files_in_srpm 2>&1

    # mv the files of srpm
    if [ -f "$files_in_srpm" ]; then
        while read fn
        do
            fn=${srpm_files_dir}/$fn

            if [[ "$fn" =~ ".spec" ]]; then
                rm -f $fn
            else
                mv $fn .
            fi
        done < $files_in_srpm
    fi
}

commit() {
    # nothing changed
    if [ -z "$(git ls-files -om)" ]; then
        return
    fi

    echo "$code_changed_tag"

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
        git commit -m "apply new package $pkg_name"
    else
        git commit --amend --no-edit --quiet
    fi

    git push -f origin "$pkg_name"

    git checkout $main_branch
}

cd $work_dir

checkout_branch

# must call handle_srpm before handle_spec
handle_srpm

handle_spec

commit
