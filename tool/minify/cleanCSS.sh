go install github.com/daaku/cssdalek@latest
cssdalek \
  --css '../../public/assets/backend/js/bootstrap/dist/css/bootstrap.css'\
  --word '../../template/backend/*.html'\
  --word '../../template/backend/*/*.html'\
  --word '../../template/backend/*/*/*.html'\
  --word '../../public/assets/backend/js/behaviour/*.js'\
  --word '../../public/assets/backend/js/behaviour/*/*.js' > ../../public/assets/backend/js/bootstrap/dist/css/bootstrap.lite.min.css

cssdalek \
  --css '../../public/assets/backend/css/style.css'\
  --word '../../template/backend/*.html'\
  --word '../../template/backend/*/*.html'\
  --word '../../template/backend/*/*/*.html'\
  --word '../../public/assets/backend/js/behaviour/*.js'\
  --word '../../public/assets/backend/js/behaviour/*/*.js'\
  --include-id 'captchaImage' > ../../public/assets/backend/css/style.lite.min.css
