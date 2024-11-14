
var CaptchaClickBasic = function (options) {
    var getCaptDataApi = options.dataApi || "/api/go-captcha-data/click-basic"
    var checkCaptDataApi = options.verifyApi || "/api/go-captcha-check-data/click-basic"
    if(!options.success) options.success=function(){alert("验证成功")}
    if(!options.error) options.error=function(){alert("验证失败")}
    var idSuffix = options.idSuffix || '';

    var captchaKey = ""
    var maxDot = 8
    var dots = []

    var hiddenClassName = "wg-cap-wrap__hidden"
    var dialogActiveClassName = "wg-cap-dialog__active"
    var activeDefaultClassName = "wg-cap-active__default"
    var activeOverClassName = "wg-cap-active__over"
    var activeErrorClassName = "wg-cap-active__error"
    var activeSuccessClassName = "wg-cap-active__success"

    var captchaWrapDom        = document.querySelector("#wg-cap-dots"+idSuffix)
    var captchaImageDom       = document.querySelector("#wg-cap-image"+idSuffix)
    var captchaThumbDom       = document.querySelector("#wg-cap-thumb"+idSuffix)
    var captchaBtnControlDom  = document.querySelector("#wg-cap-btn-control"+idSuffix)
    var captchaCheckBtnDom    = document.querySelector("#wg-cap-check-btn"+idSuffix)
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
        Helper.addEventListener(captchaImageDom, "click", handleClickPos, false)
        Helper.addEventListener(captchaCheckBtnDom, "click", handleClickCheck, false)
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

    function appendDotIcon(event, x, y) {
        var dot = document.createElement('div')
        dot.setAttribute('class', 'wg-cap-wrap__dot')
        dot.setAttribute('style', 'top:' + (y - 11) + 'px; left:' + (x - 11) + 'px;')
        dot.innerHTML = '<span>'+ (dots.length + 1) +'</span>'
        captchaWrapDom.appendChild(dot)
        dots.push([x, y])
    }

    function resetCaptcha() {
        captchaKey = ""
        dots = []
        captchaWrapDom.innerHTML = ""
    }

    function clearImage() {
        captchaImageDom.setAttribute("src", "")
        captchaThumbDom.setAttribute("src", "")
    }

    function handleClickPos(ev){
        if (dots.length >= maxDot || captchaKey === "") {
            return
        }

        var e = ev || window.event;
        var dom = e.currentTarget

        var xy = Helper.getDomXY(dom)

        var mouseX = e.pageX || e.clientX
        var mouseY = e.pageY || e.clientY

        var domX = xy.domX
        var domY = xy.domY

        var xPos = mouseX - domX;
        var yPos = mouseY - domY;

        appendDotIcon(e, parseInt(xPos.toString()), parseInt(yPos.toString()))
        e.cancelBubble = true
        e.preventDefault()
        return false
    }

    function handleClickRefresh() {
        requestCaptchaData()
    }

    function handleClickClose() {
        dialogDom.classList.remove(dialogActiveClassName)
    }

    function handleClickCheck() {
        var dotsA = []
        dots.forEach(function (value, key) {
            dotsA.push(value.join(","))
        })

        if (dotsA.length <= 0) {
            alert("请点选图案进行验证")
            return
        }

        requestCheckCaptchaData({'response': dotsA.join(','), 'key': captchaKey})
    }

    function handleClickDefault() {
        requestCaptchaData()
        dialogDom.classList.add(dialogActiveClassName)
    }

    function requestCaptchaData() {
        resetCaptcha()
        clearImage()
        captchaImageDom.classList.add(hiddenClassName)
        captchaThumbDom.classList.add(hiddenClassName)

        Ajax.get(getCaptDataApi, {}, function(data){
            if (data['code'] === 0) {
                captchaImageDom.classList.remove(hiddenClassName)
                captchaThumbDom.classList.remove(hiddenClassName)
                captchaImageDom.setAttribute("src", data['image'])
                captchaThumbDom.setAttribute("src", data['thumb'])
                captchaKey = data['key'];
                if(options.input) document.querySelector(options.input).value=captchaKey;
            } else {
                alert("请求验证码数据失败：" + data['message'])
            }
        }, function(e){
            console.log("请求验证码数据失败：" + e['message']);
        })
    }

    function requestCheckCaptchaData(dots) {
        Ajax.post(checkCaptDataApi, dots, function(data){
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
var CaptchaClickShape=CaptchaClickBasic;