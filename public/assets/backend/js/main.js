var ASSETS_BASE_URL = BACKEND_URL + '/public/assets/backend/', I18N_MODULE_URL = 'i18n!' + ASSETS_BASE_URL + 'js/nls/messages.js'; 
require.config({
    baseUrl: ASSETS_BASE_URL+'js',
    shim: {
        'bootstrap': ['jquery'],
        'bootstrap-switch': ['bootstrap', 'jquery'],
        'pjax': ['jquery'],
        'modalEffects': [
            'css!'+ASSETS_BASE_URL+'js/jquery.niftymodals/css/component.css',
            'jquery'
        ],
        'typeahead': [
            'css!'+ASSETS_BASE_URL+'js/bootstrap.typeahead/typeahead.css',
            'bootstrap', 'jquery'
        ],
        'parsley': ['jquery','jquery.parsley/dist/parsley.min','jquery.parsley/i18n/'+LANGUAGE.replace('-','_')],
        'flot': ['flot', 'flot.time', 'flot.canvas', 'flot.categories', 'flot.labels', 'jquery'],
        'general': ['sprintf', 'bootstrap', 'jquery'],
        'bootstrap-dialog': ['css!'+ASSETS_BASE_URL+'js/dialog/bootstrap-dialog.min.css'],
        'gritter': [
            'css!'+ASSETS_BASE_URL+'js/jquery.gritter/css/jquery.gritter.min.css',
            'css!'+ASSETS_BASE_URL+'js/jquery.gritter/css/custom.min.css',
            'jquery'
        ],
        'powerFloat': ['css!'+ASSETS_BASE_URL+'js/float/powerFloat.min.css'],
        'bootstrap-switch': ['css!'+ASSETS_BASE_URL+'js/bootstrap.switch/bootstrap-switch.min.css']
    },
    map: {
        '*': {
            'domReady': 'require/domReady.min',
            'css': 'require/css.min',
            'text': 'require/text.min',
            'i18n': 'require/i18n.min'
        }
    },
    paths: {
        'general': 'behaviour/general.min',
        'jquery': 'jquery',
        'modernizr': 'modernizr',
        'nanoscroller': 'jquery.nanoscroller/jquery.nanoscroller.min',
        'bootstrap-switch': 'bootstrap.switch/bootstrap-switch.min',
        'pjax': 'jquery.pjax.min',
        'modalEffects': 'jquery.niftymodals/js/jquery.modalEffects.min',
        'sprintf': 'behaviour/sprintf.min',
        'gritter': 'jquery.gritter/js/jquery.gritter.min',
        'typeahead': 'bootstrap.typeahead/typeahead.bundle.min',
        'bootstrap': 'bootstrap/dist/js/bootstrap.min',
        'nprogress': 'nprogress/nprogress.min',
        'gritter': 'jquery.gritter/js/jquery.gritter.min',
        'bootstrap-dialog': 'dialog/bootstrap-dialog.min',
        'storeWithJson2': 'storeWithJson2.min',
        'parsley': '',
        'powerFloat': 'float/powerFloat.min',
        'template': 'template',
        'cascadeSelect': 'behaviour/page/cascade-select.min',
        'footer': 'behaviour/page/footer.min',
        'flot':'jquery.flot/jquery.flot.min',
        'excanvas':'jquery.flot/excanvas.min',
        'flot.time':'jquery.flot/jquery.flot.time.min',
        'flot.canvas':'jquery.flot/jquery.flot.canvas.min',
        'flot.categories':'jquery.flot/jquery.flot.categories.min',
        'flot.labels':'jquery.flot/jquery.flot.labels.min',
        'real-time-chart':'behaviour/page/real-time-chart.min'
    },
    i18n: {
        locale: LANGUAGE //zh-cn:中文; en:英文
    },
    waitSeconds: 60
});
require(['general'],function(_,jQuery){
});