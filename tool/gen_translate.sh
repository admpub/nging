#go install github.com/admpub/i18n/cmd/fetchtext@latest
fetchtext --src=../ --dist=../config/i18n/messages --default=zh-cn --translate=true --onlyExport=false --translator=baidu --translatorConfig="appid=&secret="
