
var Ajax = (function () {
    var xhr

    function __initialize() {
        if(window.XMLHttpRequest){
            xhr = new XMLHttpRequest()
        } else if(window.ActiveObject){
            xhr = new ActiveXobject('Microsoft.XMLHTTP')
        }
    }

    function requestJQuery(options){
        $.ajax(options);
    }

    function requestHandler(options){
        options = options || {}
        options.type = (options.type || "GET").toUpperCase()
        options.dataType = options.dataType || "json"
        var params = ajaxFormatParams(options.data)

        if(options.type === "GET"){
            xhr.open("GET", options.url+"?" + params,true);
            xhr.send(null);
        }else if(options.type === "POST"){
            xhr.open("post",options.url,true);
            xhr.setRequestHeader("Content-type","application/x-www-form-urlencoded");
            xhr.send(params);
        }

        xhr.timeout = options.timeout;

        xhr.onreadystatechange = function(){
            if(xhr.readyState === 4){
                var status = xhr.status;
                if(status >= 200 && status < 300 || status === 304){
                    var data = xhr.responseText, contentType = String(xhr.getResponseHeader("Content-Type")).split(";")[0];
                    if(contentType=="application/json"){
                        try {
                            data = JSON.parse(data);
                        } catch(e) {
                            console.error(e)
                        }
                    }
                    options.success&&options.success(data, xhr.responseXML);
                }else{
                    options.error&&options.error(status);
                }
            }
            options.complete&&options.complete(status);
        }
    }

    function ajaxFormatParams(data){
        var arr = [];
        for(var name in data){
            arr.push(encodeURIComponent(name) + "=" + encodeURIComponent(data[name]))
        }
        arr.push(("v=" + Math.random()).replace(".",""))
        return arr.join("&")
    }

    function handlePost (url, data, success, error, complete) {
        requestHandler({
            url: url,
            type: 'POST',
            data: data,
            dataType:'json',
            timeout: 10000,
            contentType: "application/json",
            success: success,
            error: error,
            complete: complete
        })
    }

    function handleGet (url, data, success, error, complete) {
        requestHandler({
            url: url,
            type: 'GET',
            data: data,
            dataType:'json',
            timeout: 10000,
            success: success,
            error: error,
            complete: complete
        })
    }

    __initialize()
    return {
        post: handlePost,
        get: handleGet
    }
})();

var Helper = (function () {
    function addEventListener(el,type,fn, c) {
        if(el.addEventListener){
            el.addEventListener(type,fn, c);
        }else{
            el["on" + type]=fn;
        }
    }

    function removeEventListener(el,type,fn, c) {
        if(el.removeEventListener){
            el.removeEventListener(type,fn, c);
        }else{
            el["on" + type]=null;
        }
    }

    function calcLocationLeft(el){
        var tmp = el.offsetLeft
        var val = el.offsetParent
        while(val != null){
            tmp += val.offsetLeft
            val = val.offsetParent
        }
        return tmp
    }

    function calcLocationTop(el){
        var tmp = el.offsetTop
        var val = el.offsetParent
        while(val != null){
            tmp += val.offsetTop
            val = val.offsetParent
        }
        return tmp
    }

    function getDomXY(dom){
        var x = 0
        var y = 0
        if (dom.getBoundingClientRect) {
            var box = dom.getBoundingClientRect();
            var D = document.documentElement;
            x = box.left + Math.max(D.scrollLeft, document.body.scrollLeft) - D.clientLeft;
            y = box.top + Math.max(D.scrollTop, document.body.scrollTop) - D.clientTop
        }
        else{
            while (dom !== document.body) {
                x += dom.offsetLeft
                y += dom.offsetTop
                dom = dom.offsetParent
            }
        }
        return {
            domX: x,
            domY: y
        }
    }

    var checkTargetFather = function (that, e) {
        var parent = e.relatedTarget
        try{
            while(parent && parent !== that) {
                parent = parent.parentNode
            }
        }catch (e){
            console.warn(e)
        }
        return parent !== that
    }

    return {
        addEventListener: addEventListener,
        removeEventListener: removeEventListener,
        getDomXY: getDomXY,
        calcLocationTop: calcLocationTop,
        calcLocationLeft: calcLocationLeft,
        checkTargetFather: checkTargetFather,
    }
})();