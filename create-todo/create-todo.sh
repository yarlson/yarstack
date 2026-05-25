#!/usr/bin/env bash
set -euo pipefail

die() {
  printf 'create-todo: %s\n' "$*" >&2
  exit 1
}

usage() {
  die 'usage: create-todo.sh [--cwd <source_dir>] "<title>"'
}

has_git_dir() {
  [ -d "$1/.git" ] || [ -f "$1/.git" ]
}

init_git_repo() {
  local repo_dir=$1

  if git init -b main "$repo_dir" >/dev/null 2>&1; then
    return
  fi

  git init "$repo_dir" >/dev/null
  git -C "$repo_dir" symbolic-ref HEAD refs/heads/main
}

has_commits() {
  git -C "$1" rev-parse --verify HEAD >/dev/null 2>&1
}

ensure_todo_store() {
  local todo_root=$1

  mkdir -p "$todo_root"

  if ! has_git_dir "$todo_root"; then
    init_git_repo "$todo_root"
  fi

  if ! has_commits "$todo_root"; then
    git -C "$todo_root" symbolic-ref HEAD refs/heads/main

    printf '# create-todo\n\nThis repository stores generated todo Markdown files.\n' >"$todo_root/README.md"
    printf '.DS_Store\n*.swp\n*~\n' >"$todo_root/.gitignore"

    git -C "$todo_root" add README.md .gitignore
    git -C "$todo_root" commit -m "initialize todo store" -- README.md .gitignore >/dev/null
  fi
}

clean_remote_url() {
  local url=$1

  url=${url%%#*}
  url=${url%%\?*}
  url=${url%/}
  printf '%s\n' "$url"
}

repo_name_from_remote() {
  local url
  url=$(clean_remote_url "$1")
  url=${url%.git}
  printf '%s\n' "${url##*/}"
}

github_web_url() {
  local url path owner repo rest
  url=$(clean_remote_url "$1")

  case "$url" in
  https://github.com/*)
    path=${url#https://github.com/}
    ;;
  http://github.com/*)
    path=${url#http://github.com/}
    ;;
  git@github.com:*)
    path=${url#git@github.com:}
    ;;
  ssh://git@github.com/*)
    path=${url#ssh://git@github.com/}
    ;;
  *)
    return 1
    ;;
  esac

  path=${path%/}
  path=${path%.git}
  owner=${path%%/*}
  rest=${path#*/}
  repo=${rest%%/*}

  if [ -z "$owner" ] || [ -z "$repo" ] || [ "$owner" = "$path" ]; then
    return 1
  fi

  printf 'https://github.com/%s/%s\n' "$owner" "$repo"
}

slugify() {
  printf '%s' "$1" |
    LC_ALL=C tr '[:upper:]' '[:lower:]' |
    LC_ALL=C sed -E 's/[^a-z0-9]+/_/g; s/^_+//; s/_+$//'
}

random_hex() {
  if command -v openssl >/dev/null 2>&1; then
    openssl rand -hex 2
    return
  fi

  if [ -r /dev/urandom ] && command -v od >/dev/null 2>&1; then
    od -An -N2 -tx1 /dev/urandom | tr -d ' \n'
    return
  fi

  die 'openssl or /dev/urandom with od is required to generate a random suffix'
}

detect_project() {
  local source_dir=$1
  local git_root origin_url project_name source_url

  source_url=

  if git -C "$source_dir" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    git_root=$(git -C "$source_dir" rev-parse --show-toplevel)
    origin_url=$(git -C "$source_dir" remote get-url origin 2>/dev/null || true)

    if [ -n "$origin_url" ]; then
      project_name=$(repo_name_from_remote "$origin_url")
      source_url=$(github_web_url "$origin_url" || true)
    else
      project_name=$(basename "$git_root")
    fi
  else
    project_name=$(basename "$source_dir")
  fi

  printf '%s\n%s\n' "$project_name" "$source_url"
}

main() {
  local source_dir=$PWD
  local title

  if [ "$#" -eq 3 ] && [ "$1" = "--cwd" ]; then
    source_dir=$2
    title=$3
  elif [ "$#" -eq 1 ]; then
    title=$1
  else
    usage
  fi

  if [ ! -d "$source_dir" ]; then
    die "source directory does not exist: $source_dir"
  fi

  if [ -z "${title//[[:space:]]/}" ]; then
    usage
  fi

  local todo_root=$HOME/.create-todo
  local project_name source_url project_slug title_slug todo_dir todo_file
  local date_part time_part hex

  ensure_todo_store "$todo_root"

  {
    IFS= read -r project_name
    IFS= read -r source_url
  } < <(detect_project "$source_dir")

  project_slug=$(slugify "$project_name")
  title_slug=$(slugify "$title")

  if [ -z "$project_slug" ]; then
    die 'detected project name does not contain any alphanumeric characters'
  fi

  if [ -z "$title_slug" ]; then
    die 'title does not contain any alphanumeric characters'
  fi

  date_part=$(date '+%Y-%m-%d')
  time_part=$(date '+%H%M%S')
  hex=$(random_hex)

  todo_dir=$todo_root/$project_slug
  mkdir -p "$todo_dir"

  todo_file=$todo_dir/$date_part-$title_slug-$time_part-$hex.md
  : >"$todo_file"

  printf '%s\n%s\n' "$todo_file" "$source_url"
}

main "$@"
