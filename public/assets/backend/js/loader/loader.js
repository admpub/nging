(function(App){
    var Loader={
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
    Loader.include = function(file,location,once,successCallback,failureCallback) {
        if (location == null) location = "head";
        if (once == null) once = true;
        if (location == "head" && typeof(Loader.data["include"]) == "undefined") {
            var jsAfter = $("#js-lazyload-begin"),
                cssAfter = $("#css-lazyload-begin");
            Loader.data.include = {
                before: {},
                after: {}
            };
            if (jsAfter.length>0) {
                Loader.data.include.after.script = jsAfter;
            } else {
                var jsBefore = $("#js-lazyload-end");
                if (jsBefore.length>0) Loader.data.include.before.script = jsBefore;
            }
            if (cssAfter.length>0) {
                Loader.data.include.after.link = cssAfter;
            } else {
                var cssBefore = $("#css-lazyload-end");
                if (cssBefore.length>0) Loader.data.include.before.link = cssBefore;
            }
        }
        $.ajaxSetup({cache: true});
        var files = typeof(file) == "string" ? [file] : file;
        var loaded = {success:0,failure:0,total:files.length};
        if(successCallback||failureCallback){
            var timer = setInterval(function(){
                //console.log(loaded)
                if (loaded.success+loaded.failure < files.length) return;
                clearInterval(timer);
                if (loaded.success >= files.length) {
                    successCallback && successCallback();
                } else {
                    failureCallback && failureCallback();
                }
            }, 200);
        }
        for (var i = 0; i < files.length; i++) {
            var name = files[i].replace(/^\s|\s$/g, ""),
                att = name.split('.');
            var ext = att[att.length - 1].toLowerCase(),
                isCSS = ext == "css";
            var tag = isCSS ? "link" : "script";
            var attr = isCSS ? ' type="text/css" rel="stylesheet"' : ' type="text/javascript"';
            attr += ' charset="utf-8" ';
            var link = (isCSS ? "href" : "src") + "='" + name + "'";
            if (once && $(tag + "[" + link + "]").length > 0) {
                loaded.success++;
                continue;
            }
            var ej = $("<" + tag + attr + link + "></" + tag + ">");
            var script = ej[0];
            if (script.readyState) {  // IE
                script.onreadystatechange = function() {
                    if (script.readyState === 'loaded' || script.readyState === 'complete') {
                        script.onreadystatechange = null;
                        loaded.success++;
                    }
                };
            } else {  // Other Browsers
                script.onload = function() {
                    loaded.success++;
                };
            }
            if (location == "head") {
                if (typeof(Loader.data.include.after[tag]) != 'undefined') {
                    Loader.data.include.after[tag].after(ej);
                    loaded.success++;
                    continue;
                } 
                if (typeof(Loader.data.include.before[tag]) != 'undefined') {
                    Loader.data.include.before[tag].before(ej);
                    loaded.success++;
                    continue;
                }
            }
            try{
                $(location).append(ej);
                loaded.success++;
            }catch(err){
                loaded.failure++;
                console.error(err.message);
                console.log(name);
            }
        }
        $.ajaxSetup({cache: false});
    };
    Loader.defined = function(vType, key, callback, onloadCallback, failureCallback) {
        if (vType != 'undefined' || key == null) {
            if (key != null && callback != null) return callback();
            return;
        }
        if (typeof(key) == 'string' && typeof(Loader.libs[key]) != 'undefined') key = Loader.libs[key];
        var successCallback = onloadCallback;
        if (callback != null) {
            if ($.isPlainObject(callback)) {
                if (typeof callback['success'] != "undefined") successCallback = callback['success'];
                if (typeof callback['error'] != "undefined") failureCallback = callback['error'];
                else if (typeof callback['failure'] != "undefined") failureCallback = callback['failure'];
            } else if (onloadCallback != null) {
                successCallback = function() { 
                    onloadCallback();
                    callback();
                };
            } else {
                successCallback = callback;
            }
        }
        Loader.includes(key, true, successCallback, failureCallback);
    };
    Loader.fullURL = function(file) {
        var url=Loader.staticURL;
        if (file.substring(0,1)=='#') {
            url=Loader.assetsURL+'/js/';
            file=file.substring(1);
        }
        return url+file;
    };
    Loader.includes = function(js,once,onloadCallback,failureCallback) {
        if (!js) return;
        switch (typeof(js)) {
        case 'string':
            Loader.include(Loader.fullURL(js),null,once,onloadCallback,failureCallback);
            break;
        default:
            if (typeof(js.length) == 'undefined') return;
            var jss = [];
            for (var i = 0; i < js.length; i++) {
                jss.push(Loader.fullURL(js[i]));
            }
            Loader.include(jss,null,once,onloadCallback,failureCallback);
        }
    };
    App.loader=Loader;
})(App);