go install github.com/daaku/cssdalek@latest
baseDir=../../
cssdalek \
  --css "${baseDir}public/assets/backend/js/bootstrap/dist/css/bootstrap.css"\
  --word "${baseDir}public/assets/backend/js/bootstrap/dist/js/bootstrap.js"\
  --word "${baseDir}template/backend/*.html"\
  --word "${baseDir}template/backend/*/*.html"\
  --word "${baseDir}template/backend/*/*/*.html"\
  --word "${baseDir}../webx/template/backend/official/*.html"\
  --word "${baseDir}../webx/template/backend/official/*/*.html"\
  --word "${baseDir}../webx/template/backend/official/*/*/*.html"\
  --word "${baseDir}public/assets/backend/js/behaviour/*.js"\
  --word "${baseDir}public/assets/backend/js/behaviour/*/*.js" > ${baseDir}public/assets/backend/js/bootstrap/dist/css/bootstrap.lite.min.css

cssdalek \
  --css "${baseDir}public/assets/backend/css/style.css"\
  --word "${baseDir}template/backend/*.html"\
  --word "${baseDir}template/backend/*/*.html"\
  --word "${baseDir}template/backend/*/*/*.html"\
  --word "${baseDir}../webx/template/backend/official/*.html"\
  --word "${baseDir}../webx/template/backend/official/*/*.html"\
  --word "${baseDir}../webx/template/backend/official/*/*/*.html"\
  --word "${baseDir}public/assets/backend/js/behaviour/*.js"\
  --word "${baseDir}public/assets/backend/js/behaviour/*/*.js"\
  --word "${baseDir}public/assets/backend/js/bootstrap/dist/js/bootstrap.js"\
  --word "${baseDir}public/assets/backend/js/bootstrap.switch/bootstrap-switch.min.js"\
  --include-selector ".profile_menu .dropdown-toggle"\
  --include-selector ".code-cont .main-app"\
  --include-selector ".sp_results>li"\
  --include-id "captchaImage" > ${baseDir}public/assets/backend/css/style.lite.min.css
