(function(App){
    Loader={
        data:{},
        libs:{},
        staticURL:'',
        siteURL:'',
        assetsURL:ASSETS_URL,
    };
	Loader.getValue = function(key, data) {
		var keys = key.split(".");
		var v = data[keys.shift()];
		if (v === null) return "";
		for (var i = 0, l = keys.length; i < l; i++) {
			v = v[keys[i]];
			if (v === null) return "";
		}
		return typeof(v) !== "undefined" && v !== null ? v : "";
	};
    Loader.parseTmpl = function(template, data) {
		return template.replace(/\{=([\w\.]*)=\}/g, function(str, key) {
			return Loader.getValue(key, data);
		});
	};
    Loader.include = function(file, location,once) {
        if (location == null) location = "head";
        if (once == null) once = true;
        if (location == "head" && typeof(Loader.data["include"]) == "undefined") {
            var jsAfter = $("#js-lazyload-begin"),
                cssAfter = $("#css-lazyload-begin");
            Loader.data.include = {
                before: {},
                after: {}
            };
            if (jsAfter.length) {
                Loader.data.include.after.script = jsAfter;
            } else {
                var jsBefore = $("#js-lazyload-end");
                if (jsBefore.length) Loader.data.include.before.script = jsBefore;
            }
            if (cssAfter.length) {
                Loader.data.include.after.link = cssAfter;
            } else {
                var cssBefore = $("#css-lazyload-end");
                if (cssBefore.length) Loader.data.include.before.link = cssBefore;
            }
        }
        $.ajaxSetup({cache: true});
        var files = typeof(file) == "string" ? [file] : file;
        for (var i = 0; i < files.length; i++) {
            var name = files[i].replace(/^\s|\s$/g, ""),
                att = name.split('.');
            var ext = att[att.length - 1].toLowerCase(),
                isCSS = ext == "css";
            var tag = isCSS ? "link" : "script";
            var attr = isCSS ? ' type="text/css" rel="stylesheet"' : ' type="text/javascript"';
            attr += ' charset="utf-8" ';
            var link = (isCSS ? "href" : "src") + "='" + name + "'";
            if (once && $(tag + "[" + link + "]").length > 0) continue;
            var ej = $("<" + tag + attr + link + "></" + tag + ">");
            if (location == "head") {
                if (typeof(Loader.data.include.after[tag]) != 'undefined') {
                    Loader.data.include.after[tag].after(ej);
                    continue;
                } else if (typeof(Loader.data.include.before[tag]) != 'undefined') {
                    Loader.data.include.before[tag].before(ej);
                    continue;
                }
            }
            $(location).append(ej);
        }
        $.ajaxSetup({cache: false});
    };
    Loader.defined = function(vType, key, callback) {
        if (vType != 'undefined' || key == null) {
            if (key != null && callback != null) return callback();
            return;
        }
        if (typeof(key) == 'string' && typeof(Loader.libs[key]) != 'undefined') key = Loader.libs[key];
        Loader.includes(key);
        if (callback != null) return callback();
    };
    Loader.includes = function(js,once) {
        if (!js) return;
        switch (typeof(js)) {
        case 'string':
            var url=Loader.staticURL;
            if (js.substring(0,1)=='#') {
                url=Loader.assetsURL+'/js/';
                js=js.substring(1);
            }
            Loader.include(url + '/' + js,null,once);
            return;
        default:
            if (typeof(js.length) == 'undefined') return;
            var jss = [];
            for (var i = 0; i < js.length; i++) {
                var url=Loader.staticURL;
                var jsf=js[i];
                if (jsf.substring(0,1)=='#') {
                    url=Loader.assetsURL+'/js/';
                    jsf=jsf.substring(1);
                }
                jss.push(url + '/' + jsf);
            }
            Loader.include(jss,null,once);
        }
    };
    App.loader=Loader;
})(App);