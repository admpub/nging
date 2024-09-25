go install github.com/daaku/cssdalek@latest
baseDir=../../
cssdalek \
  --css "${baseDir}public/assets/backend/js/bootstrap/dist/css/bootstrap.css"\
  --word "${baseDir}public/assets/backend/js/bootstrap/dist/js/bootstrap.js"\
  --word "${baseDir}template/backend/*.html"\
  --word "${baseDir}template/backend/*/*.html"\
  --word "${baseDir}template/backend/*/*/*.html"\
  --word "${baseDir}vendor/github.com/nging-plugins/*/template/backend/*.html"\
  --word "${baseDir}vendor/github.com/nging-plugins/*/template/backend/*/*.html"\
  --word "${baseDir}vendor/github.com/nging-plugins/*/template/backend/*/*/*.html"\
  --word "${baseDir}vendor/github.com/nging-plugins/*/template/backend/*/*/*/*.html"\
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
  --word "${baseDir}vendor/github.com/nging-plugins/*/template/backend/*.html"\
  --word "${baseDir}vendor/github.com/nging-plugins/*/template/backend/*/*.html"\
  --word "${baseDir}vendor/github.com/nging-plugins/*/template/backend/*/*/*.html"\
  --word "${baseDir}vendor/github.com/nging-plugins/*/template/backend/*/*/*/*.html"\
  --word "${baseDir}../webx/template/backend/official/*.html"\
  --word "${baseDir}../webx/template/backend/official/*/*.html"\
  --word "${baseDir}../webx/template/backend/official/*/*/*.html"\
  --word "${baseDir}public/assets/backend/js/behaviour/*.js"\
  --word "${baseDir}public/assets/backend/js/behaviour/*/*.js"\
  --word "${baseDir}public/assets/backend/js/bootstrap/dist/js/bootstrap.js"\
  --word "${baseDir}public/assets/backend/js/bootstrap.switch/bootstrap-switch.min.js"\
  --word "${baseDir}public/assets/backend/js/jquery.select2/select2.js"\
  --word "${baseDir}public/assets/backend/js/jquery.parsley/parsley.js"\
  --word "${baseDir}public/assets/backend/js/jquery.sparkline/jquery.sparkline.min.js"\
  --include-selector ".profile_menu .dropdown-toggle"\
  --include-selector ".code-cont .main-app"\
  --include-selector ".modal-body .dropzone .dz-preview"\
  --include-id "captchaImage"\
  --include-class "progress-bar-.*" > ${baseDir}public/assets/backend/css/style.lite.min.css
