var CaptchaSlideRegion = function (options) {
    var getCaptDataApi = options.dataApi || "/api/go-captcha-data/slide-basic"
    var checkCaptDataApi = options.verifyApi || "/api/go-captcha-check-data/slide-basic"
    if(!options.success) options.success=function(){alert("验证成功")}
    if(!options.error) options.error=function(){alert("验证失败")}
    var idSuffix = options.idSuffix || '';

    var captchaKey = ""
    var pointX = 0
    var pointY = 0

    var hiddenClassName = "wg-cap-wrap__hidden"
    var dialogActiveClassName = "wg-cap-dialog__active"
    var activeDefaultClassName = "wg-cap-active__default"
    var activeOverClassName = "wg-cap-active__over"
    var activeErrorClassName = "wg-cap-active__error"
    var activeSuccessClassName = "wg-cap-active__success"

    var captchaDragWrapDom    = document.querySelector("#wg-cap-wrap-drag"+idSuffix)
    var captchaTileWrapDom    = document.querySelector("#wg-cap-tile"+idSuffix)
    var captchaTileImageDom   = document.querySelector("#wg-cap-tile-picture"+idSuffix)
    var captchaImageDom       = document.querySelector("#wg-cap-image"+idSuffix)
    var captchaBtnControlDom  = document.querySelector("#wg-cap-btn-control"+idSuffix)
    var captchaCloseBtnDom    = document.querySelector("#wg-cap-close-btn"+idSuffix)
    var captchaDialogBtnDom   = document.querySelector("#wg-cap-dialog"+idSuffix)
    var captchaRefreshBtnDom  = document.querySelector("#wg-cap-refresh-btn"+idSuffix)
    var captchaDefaultBtnDom  = document.querySelector("#wg-cap-btn-default"+idSuffix)
    var captchaErrorBtnDom    = document.querySelector("#wg-cap-btn-error"+idSuffix)
    var captchaOverBtnDom     = document.querySelector("#wg-cap-btn-over"+idSuffix)
    var dialogDom             = document.querySelector("#wg-cap-container"+idSuffix)
    var retryBtnDom           = document.querySelector("#wg-cap-btn-retry"+idSuffix)

    function __initialize() {
        // requestCaptchaData()
        handleEvent()

        document.addEventListener('touchstart', (event) => {
            if (event.touches.length > 1) {
                event.preventDefault()
            }
        })
        document.addEventListener('gesturestart', (event) => {
            event.preventDefault()
        })
        document.body.addEventListener('touchend', () => { })
    }

    function handleEvent() {
        Helper.addEventListener(captchaTileWrapDom, "mousedown", handleDragEvent, false)
        Helper.addEventListener(captchaTileWrapDom, "dragstart", function (e){e.preventDefault()}, false)
        Helper.addEventListener(captchaTileWrapDom, "touchstart", handleDragEvent, false)
        Helper.addEventListener(captchaCloseBtnDom, "click", handleClickClose, false)
        Helper.addEventListener(captchaDialogBtnDom, "click", handleClickClose, false)
        Helper.addEventListener(captchaRefreshBtnDom, "click", handleClickRefresh, false)
        Helper.addEventListener(captchaDefaultBtnDom, "click", handleClickDefault, false)
        Helper.addEventListener(captchaErrorBtnDom, "click", handleClickDefault, false)
        Helper.addEventListener(captchaOverBtnDom, "click", handleClickDefault, false)
        Helper.addEventListener(retryBtnDom, "click", handleSucceedRetry, false)
    }

    function handleSucceedRetry() {
        captchaBtnControlDom.classList.remove(activeSuccessClassName)
        captchaBtnControlDom.classList.add(activeDefaultClassName)
        if(options.input) document.querySelector(options.input).value='';
        handleClickDefault.apply(this,arguments)
    }

    function resetCaptcha() {
        captchaKey = ""
        captchaTileImageDom.setAttribute("src", "")
        captchaTileWrapDom.style.left = 0
        captchaTileWrapDom.style.top = 0
        captchaTileWrapDom.style.display = "none"
    }

    function clearImage() {
        captchaImageDom.setAttribute("src", "")
    }

    function handleDragEvent(ev){
        var ee = ev || window.event;
        var touch = ev.touches && ev.touches[0];
        var offsetLeft = captchaTileWrapDom.offsetLeft
        var offsetTop = captchaTileWrapDom.offsetTop
        var width = captchaDragWrapDom.offsetWidth
        var height = captchaDragWrapDom.offsetHeight
        var tileWidth = captchaTileWrapDom.offsetWidth
        var tileHeight = captchaTileWrapDom.offsetHeight
        var maxWidth = width - tileWidth
        var maxHeight = height - tileHeight

        var isMoving = false
        var startX = 0;
        var startY = 0;
        if (touch) {
            startX = touch.pageX - offsetLeft
            startY = touch.pageY - offsetTop
        } else {
            startX = ee.clientX - offsetLeft
            startY = ee.clientY - offsetTop
        }

        var handleMove = function(e) {
            isMoving = true

            var mTouche = e.touches && e.touches[0];
            var me = e || window.event;

            var left = 0;
            var top = 0;
            if (mTouche) {
                left = mTouche.pageX - startX
                top = mTouche.pageY - startY
            } else {
                left = me.clientX - startX
                top = me.clientY - startY
            }

            captchaTileWrapDom.style.left = left + "px";
            captchaTileWrapDom.style.top = top + "px";

            if (left <= 0) {
                captchaTileWrapDom.style.left = 0 + "px";
            }

            if (top <= 0) {
                captchaTileWrapDom.style.top = 0 + "px";
            }

            if (left >= maxWidth) {
                captchaTileWrapDom.style.left = maxWidth + "px";
            }

            if (top >= maxHeight) {
                captchaTileWrapDom.style.top = maxHeight + "px";
            }

            me.cancelBubble = true
            me.preventDefault()
        }

        var handleUp = function(e) {
            var ue = e || window.event;
            if (!isMoving) {
                return
            }

            var uTouche = e.changedTouches && e.changedTouches[0];
            var mouseX = 0;
            var mouseY = 0;
            if (uTouche) {
                mouseX = uTouche.pageX
                mouseY = uTouche.pageY
            } else {
                mouseX = ue.pageX || ue.clientX
                mouseY = ue.pageY || ue.clientY
            }

            isMoving = false
            Helper.removeEventListener(captchaDragWrapDom, "mousemove", handleMove, false)
            Helper.removeEventListener(captchaDragWrapDom, "mouseup", handleUp, false)
            Helper.removeEventListener(captchaDragWrapDom, "mouseout", handleUp, false)

            Helper.removeEventListener(captchaDragWrapDom, "touchmove", handleMove, false)
            Helper.removeEventListener(captchaDragWrapDom, "touchend", handleUp, false)

            var xy = Helper.getDomXY(captchaDragWrapDom)

            var domX = xy.domX
            var domY = xy.domY

            var txy = Helper.getDomXY(captchaTileWrapDom)
            var ox = mouseX - txy.domX
            var oy = mouseY - txy.domY

            var xPos = mouseX - domX - ox;
            var yPos = mouseY - domY - oy;

            handleClickCheck(xPos, yPos)

            ue.cancelBubble = true
            ue.preventDefault()
        }

        Helper.addEventListener(captchaDragWrapDom, "mousemove", handleMove, false);
        Helper.addEventListener(captchaDragWrapDom, "mouseout", handleUp, false);
        Helper.addEventListener(captchaDragWrapDom, "mouseup", handleUp, false);

        Helper.addEventListener(captchaDragWrapDom, "touchmove", handleMove, false);
        Helper.addEventListener(captchaDragWrapDom, "touchend", handleUp, false);

        ee.cancelBubble = true
        ee.preventDefault()
    }

    function handleClickRefresh() {
        requestCaptchaData()
    }

    function handleClickClose() {
        dialogDom.classList.remove(dialogActiveClassName)
    }

    function handleClickCheck(x, y) {
        if (x === pointX || y === pointY) {
            alert("请点拖动图案进行验证")
            return
        }

        requestCheckCaptchaData({'response': [x, y].join(','), 'key': captchaKey})
    }

    function handleClickDefault() {
        requestCaptchaData()
        dialogDom.classList.add(dialogActiveClassName)
    }

    function requestCaptchaData() {
        resetCaptcha()
        clearImage()
        captchaImageDom.classList.add(hiddenClassName)

        Ajax.get(getCaptDataApi, {}, function(data){
            if (data['code'] === 0) {
                captchaImageDom.classList.remove(hiddenClassName)
                captchaImageDom.setAttribute("src", data['image'])
                captchaTileImageDom.setAttribute("src", data['tile']['image'])

                captchaTileWrapDom.style.left = data['tile']['x'] + "px"
                captchaTileWrapDom.style.top = data['tile']['y'] + "px"
                captchaTileWrapDom.style.width = data['tile']['width'] + "px"
                captchaTileWrapDom.style.height = data['tile']['height'] + "px"
                captchaTileWrapDom.style.display = "block"

                captchaKey = data['key']
                pointX = data['tile']['x']
                pointY = data['tile']['y']
                if(options.input) document.querySelector(options.input).value=captchaKey;
            } else {
                alert("请求验证码数据失败：" + data['message'])
            }
        }, function(e){
            console.log("请求验证码数据失败：" + e['message']);
        })
    }

    function requestCheckCaptchaData(point) {
        Ajax.post(checkCaptDataApi, point, function(data){
            captchaBtnControlDom.classList.remove(activeDefaultClassName)
            captchaBtnControlDom.classList.remove(activeOverClassName)
            if (data['code'] === 0) {
                captchaBtnControlDom.classList.remove(activeErrorClassName)
                captchaBtnControlDom.classList.add(activeSuccessClassName)
                setTimeout(function () {
                    handleClickClose()
                }, 200)
                options.success && options.success.apply(this,arguments)
            } else {
                captchaBtnControlDom.classList.remove(activeSuccessClassName)
                captchaBtnControlDom.classList.add(activeErrorClassName)
                requestCaptchaData()
                options.error && options.error.apply(this,arguments)
            }
        }, function(e){
            captchaBtnControlDom.classList.remove(activeDefaultClassName)
            captchaBtnControlDom.classList.remove(activeOverClassName)
            captchaBtnControlDom.classList.remove(activeSuccessClassName)
            captchaBtnControlDom.classList.add(activeErrorClassName)
            requestCaptchaData()
            options.networkError && options.networkError.apply(this,arguments)
        }, function () {
            captchaKey = ""
            options.complete && options.complete.apply(this,arguments)
        })
    }

    __initialize()
    return {}
};