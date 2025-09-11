if which fetchtext &> /dev/null; then
    echo "found fetchtext"
else
    echo "installing fetchtext"
    go install github.com/admpub/i18n/cmd/fetchtext@latest
fi

fetchtext --src=../ --dist=../config/i18n/messages --default=zh-CN --translate=true --clean=true --onlyExport=false --translator=tencent --translatorConfig="appid=&secret=" --envFile="$PWD/translator_tencent.env" --onlyTranslateIncr=true
