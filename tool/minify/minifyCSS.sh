SHELL_FOLDER=$(cd "$(dirname "$0")";pwd)
node ${SHELL_FOLDER}/minifyCSS.js "$1" "$2"
