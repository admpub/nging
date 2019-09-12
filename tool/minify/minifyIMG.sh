SHELL_FOLDER=$(cd "$(dirname "$0")";pwd)
node ${SHELL_FOLDER}/minifyIMG.js "$1" "$2"
