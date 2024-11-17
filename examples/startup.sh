#!/bin/sh
set -e

REPO=vordimous/gohlay
RELEASE_URL="https://github.com/$REPO/releases/download"
GOHLAY_VERSION=""
EXAMPLE_FOLDER=""
AUTO_TEARDOWN=false
WORKDIR=$(pwd)

# help text
HELP_TEXT=$(cat <<EOF
Usage: ${CMD:=${0##*/}} [-m][-d WORKDIR][-v GOHLAY_VERSION] example.name

Operand:
    example.name          The name of the example to use                                 [default: quickstart][string]

Options:
    -d | --workdir        Sets the directory used to download and run the example                             [string]
    -m | --use-main       Download the head of the main branch                                               [boolean]
    -v | --gohlay-version Sets the gohlay version to use                                     [default: latest][string]
         --auto-teardown  Executes the teardown script immediately after setup                               [boolean]
         --help           Print help                                                                         [boolean]

Report a bug: github.com/aklivity/gohlay/issues/new/choose
EOF
)
export USAGE="$HELP_TEXT"
exit2 () { printf >&2 "%s:  %s: '%s'\n%s\n" "$CMD" "$1" "$2" "$USAGE"; exit 2; }
check () { { [ "$1" != "$EOL" ] && [ "$1" != '--' ]; } || exit2 "missing argument" "$2"; } # avoid infinite loop

# parse command-line options
set -- "$@" "${EOL:=$(printf '\1\3\3\7')}"  # end-of-list marker
while [ "$1" != "$EOL" ]; do
  opt="$1"; shift
  case "$opt" in

    #defined options
    -d | --workdir        ) check "$1" "$opt"; WORKDIR="$1"; shift;;
    -v | --gohlay-version ) check "$1" "$opt"; GOHLAY_VERSION="$1"; shift;;
         --auto-teardown  ) AUTO_TEARDOWN=true;;
         --help           ) printf "%s\n" "$USAGE"; exit 0;;

    # process special cases
    --) while [ "$1" != "$EOL" ]; do set -- "$@" "$1"; shift; done;;   # parse remaining as positional
    --[!=]*=*) set -- "${opt%%=*}" "${opt#*=}" "$@";;                  # "--opt=arg"  ->  "--opt" "arg"
    -[A-Za-z0-9] | -*[!A-Za-z0-9]*) exit2 "invalid option" "$opt";;    # anything invalid like '-*'
    -?*) other="${opt#-?}"; set -- "${opt%"$other"}" "-${other}" "$@";;  # "-abc"  ->  "-a" "-bc"
    *) set -- "$@" "$opt";;                                            # positional, rotate to the end
  esac
done; shift

# check ability to run Gohlay with docker
! [ -x "$(command -v docker)" ] && echo "WARN: Docker is required to run this setup."
! [ -x "$(command -v docker compose)" ] && echo "WARN: Docker Compose is required to run this setup."

# pull the example folder from the end of the params and set defaults
operand=$*
EXAMPLE_FOLDER=$(echo "$operand" | sed 's/\///g')
[ -z "$EXAMPLE_FOLDER" ] && EXAMPLE_FOLDER="quickstart"
[ -z "$GOHLAY_VERSION" ] && GOHLAY_VERSION=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep -i "tag_name" | awk -F '"' '{print $4}')

echo "==== Starting Gohlay $GOHLAY_VERSION Example $EXAMPLE_FOLDER at $WORKDIR ===="

! [ -d "$WORKDIR" ] && echo "Error: WORKDIR must be a valid directory." && exit2
if [ -d "$WORKDIR" ] && ! [ -d "$WORKDIR/$EXAMPLE_FOLDER" ]; then
  echo "==== Downloading $RELEASE_URL/$GOHLAY_VERSION/$EXAMPLE_FOLDER.tar.gz to $WORKDIR ===="
  wget -qO- "$RELEASE_URL"/"$GOHLAY_VERSION"/"$EXAMPLE_FOLDER.tar.gz" | tar -xf - -C "$WORKDIR"
fi

export GOHLAY_VERSION="$GOHLAY_VERSION"
export EXAMPLE_DIR="$WORKDIR/$EXAMPLE_FOLDER"

TEARDOWN_SCRIPT=""
cd "$WORKDIR"/"$EXAMPLE_FOLDER"

chmod u+x teardown.sh
TEARDOWN_SCRIPT="$(pwd)/teardown.sh"
printf "\n\n"
echo "==== Use this script to teardown ===="
printf '%s\n' "$TEARDOWN_SCRIPT"
sh setup.sh

printf "\n\n"
echo "==== Check out the README to see how to use this example ==== "
echo "cd $WORKDIR/$EXAMPLE_FOLDER"
echo "cat README.md"
head -n 4 "$WORKDIR"/"$EXAMPLE_FOLDER"/README.md | tail -n 3

printf "\n\n"
echo "==== Finished, use the teardown script(s) to clean up ===="
printf '%s\n' "$TEARDOWN_SCRIPT"

if [ $AUTO_TEARDOWN = true ]; then
    printf "\n\n"
    echo "==== Auto teardown ===="
    printf '%s\n' "$TEARDOWN_SCRIPT"
    [ -n "$TEARDOWN_SCRIPT" ] && bash -c "$TEARDOWN_SCRIPT"
fi
printf "\n"
