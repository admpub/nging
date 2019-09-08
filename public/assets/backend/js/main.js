require.config({
    baseUrl: ASSETS_URL + '/js',
    shim: {
        'bootstrap': ['jquery'],
        'bootstrap-switch': ['bootstrap', 'jquery'],
        'pjax': ['jquery'],
        'modalEffects': ['jquery'],
        'gritter': ['jquery'],
        'typeahead': ['bootstrap', 'jquery'],
        'parsley': ['jquery'],
        'flot': ['flot', 'flot.time', 'flot.canvas', 'flot.categories', 'flot.labels', 'jquery'],
        'general': ['sprintf', 'bootstrap', 'jquery']
    },
    map: {
        '*': {
            'css': 'require/css.min',
            'text': 'require/text.min',
            'i18n': 'require/i18n.min'
        }
    },
    paths: {
        //plugins
        'domReady': 'require/domReady.min',
        //app
        'general': 'behaviour/general.min',
        'jquery': 'jquery',
        'modernizr': 'modernizr',
        'nanoscroller': 'jquery.nanoscroller/jquery.nanoscroller',
        'bootstrap-switch': 'bootstrap.switch/bootstrap-switch.min',
        'pjax': 'jquery.pjax.min',
        'modalEffects': 'jquery.niftymodals/js/jquery.modalEffects.min',
        'sprintf': 'behaviour/sprintf.min',
        'gritter': 'jquery.gritter/js/jquery.gritter.min',
        'typeahead': 'bootstrap.typeahead/typeahead.bundle.min',
        'bootstrap': 'bootstrap/dist/js/bootstrap.min',
        'storeWithJson2': 'storeWithJson2.min',
        'parsley': 'jquery.parsley/dist/parsley.min',//jquery.parsley/i18n/messages.
        'footer': 'behaviour/page/footer.min',
        'flot':'jquery.flot/jquery.flot.min',
        'excanvas':'jquery.flot/excanvas.min',
        'flot.time':'jquery.flot/jquery.flot.time.min',
        'flot.canvas':'jquery.flot/jquery.flot.canvas.min',
        'flot.categories':'jquery.flot/jquery.flot.categories.min',
        'flot.labels':'jquery.flot/jquery.flot.labels.min',
        'real-time-chart':'behaviour/page/real-time-chart.min'
    }
});
require(['general'],function(_,jQuery){
});