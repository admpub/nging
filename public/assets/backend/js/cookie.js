function setCookie(name, value, days, path) {
    var exp = new Date();
    if(!days) days=30;
    exp.setTime(exp.getTime() + days * 24 * 60 * 60 * 1000);
    var cookie = name + '=' + escape(value) + ';path=' + (path?path:window.location.pathname) + ';expires=' + exp.toUTCString() + ';sameSite=Lax';
    if (window.location.protocol === 'https:') cookie += ';secure=true';
    document.cookie = cookie;
}
function getCookie(name) {
	var arr = document.cookie.match(new RegExp('(^| )' + name + '=([^;]*)(;|$)'));
	if (arr != null) return unescape(arr[2]);
}