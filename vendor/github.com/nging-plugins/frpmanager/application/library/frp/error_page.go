package frp

import "github.com/admpub/frp/pkg/util/vhost"

const (
	NotFound = `<!doctype html>
<html lang='en'>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=0, minimum-scale=1.0, maximum-scale=1.0">
<meta name="apple-mobile-web-app-capable" content="yes">
<meta name="apple-mobile-web-app-status-bar-style" content="black">
<meta name="format-detection" content="telephone=no">
<title>404</title>
<style>html, body{height:95%;}body{background: #0f3854;background: -webkit-radial-gradient(center ellipse, #0a2e38 0%, #000000 70%);background: radial-gradient(ellipse at center, #0a2e38 0%, #000000 70%);background-size: 100%;}p{margin:0;padding:0;}#clock{font-family: 'Share Tech Mono', monospace;color: #ffffff;text-align: center;position: absolute;left: 50%;top: 50%;-webkit-transform: translate(-50%, -50%);transform: translate(-50%, -50%);color: #daf6ff;text-shadow: 0 0 20px #0aafe6, 0 0 20px rgba(10, 175, 230, 0);}#clock .time{letter-spacing: 0.05em;font-size: 60px;padding: 5px 0;}#clock .date{letter-spacing:0.1em;font-size:15px;}#clock .text{letter-spacing: 0.1em;font-size:12px;padding:20px 0 0;}</style>
</head>
<body>
<script>
var app=function(){"use strict";function t(){}function n(t){return t()}function e(){return Object.create(null)}function o(t){t.forEach(n)}function r(t){return"function"==typeof t}function c(t,n){return t!=t?n==n:t!==n||t&&"object"==typeof t||"function"==typeof t}function u(t,n){t.appendChild(n)}function s(t){t.parentNode.removeChild(t)}function i(t){return document.createElement(t)}function a(t){return document.createTextNode(t)}function f(){return a(" ")}function l(t,n,e){null==e?t.removeAttribute(n):t.getAttribute(n)!==e&&t.setAttribute(n,e)}function d(t,n){n=""+n,t.wholeText!==n&&(t.data=n)}let p;function g(t){p=t}const h=[],$=[],m=[],b=[],y=Promise.resolve();let _=!1;function x(t){m.push(t)}let k=!1;const v=new Set;function w(){if(!k){k=!0;do{for(let t=0;t<h.length;t+=1){const n=h[t];g(n),E(n.$$)}for(g(null),h.length=0;$.length;)$.pop()();for(let t=0;t<m.length;t+=1){const n=m[t];v.has(n)||(v.add(n),n())}m.length=0}while(h.length);for(;b.length;)b.pop()();_=!1,k=!1,v.clear()}}function E(t){if(null!==t.fragment){t.update(),o(t.before_update);const n=t.dirty;t.dirty=[-1],t.fragment&&t.fragment.p(t.ctx,n),t.after_update.forEach(x)}}const A=new Set;function N(t,n){-1===t.$$.dirty[0]&&(h.push(t),_||(_=!0,y.then(w)),t.$$.dirty.fill(0)),t.$$.dirty[n/31|0]|=1<<n%31}function O(c,u,i,a,f,l,d=[-1]){const h=p;g(c);const $=c.$$={fragment:null,ctx:null,props:l,update:t,not_equal:f,bound:e(),on_mount:[],on_destroy:[],on_disconnect:[],before_update:[],after_update:[],context:new Map(h?h.$$.context:u.context||[]),callbacks:e(),dirty:d,skip_bound:!1};let m=!1;if($.ctx=i?i(c,u.props||{},((t,n,...e)=>{const o=e.length?e[0]:n;return $.ctx&&f($.ctx[t],$.ctx[t]=o)&&(!$.skip_bound&&$.bound[t]&&$.bound[t](o),m&&N(c,t)),n})):[],$.update(),m=!0,o($.before_update),$.fragment=!!a&&a($.ctx),u.target){if(u.hydrate){const t=function(t){return Array.from(t.childNodes)}(u.target);$.fragment&&$.fragment.l(t),t.forEach(s)}else $.fragment&&$.fragment.c();u.intro&&((b=c.$$.fragment)&&b.i&&(A.delete(b),b.i(y))),function(t,e,c,u){const{fragment:s,on_mount:i,on_destroy:a,after_update:f}=t.$$;s&&s.m(e,c),u||x((()=>{const e=i.map(n).filter(r);a?a.push(...e):o(e),t.$$.on_mount=[]})),f.forEach(x)}(c,u.target,u.anchor,u.customElement),w()}var b,y;g(h)}function j(n){let e,o,r,c,p,g,h,$;return{c(){e=i("div"),o=i("p"),o.textContent="...Oops! Page Not Found...",r=f(),c=i("p"),p=a(n[0]),g=f(),h=i("p"),$=a(n[1]),l(o,"class","date"),l(c,"class","time"),l(h,"class","text"),l(e,"id","clock")},m(t,n){!function(t,n,e){t.insertBefore(n,e||null)}(t,e,n),u(e,o),u(e,r),u(e,c),u(c,p),u(e,g),u(e,h),u(h,$)},p(t,[n]){1&n&&d(p,t[0]),2&n&&d($,t[1])},i:t,o:t,d(t){t&&s(e)}}}function C(t,n){for(var e="",o=0;o<n;o++)e+="0";return(e+t).slice(-n)}function D(t,n,e){let o=["星期日","星期一","星期二","星期三","星期四","星期五","星期六"],r="",c="";function u(){let t=new Date;e(0,r=C(t.getHours(),2)+":"+C(t.getMinutes(),2)+":"+C(t.getSeconds(),2)),e(1,c=C(t.getFullYear(),4)+"-"+C(t.getMonth()+1,2)+"-"+C(t.getDate(),2)+" "+o[t.getDay()])}return setInterval(u,1e3),u(),[r,c]}return new class extends class{$destroy(){!function(t,n){const e=t.$$;null!==e.fragment&&(o(e.on_destroy),e.fragment&&e.fragment.d(n),e.on_destroy=e.fragment=null,e.ctx=[])}(this,1),this.$destroy=t}$on(t,n){const e=this.$$.callbacks[t]||(this.$$.callbacks[t]=[]);return e.push(n),()=>{const t=e.indexOf(n);-1!==t&&e.splice(t,1)}}$set(t){var n;this.$$set&&(n=t,0!==Object.keys(n).length)&&(this.$$.skip_bound=!0,this.$$set(t),this.$$.skip_bound=!1)}}{constructor(t){super(),O(this,t,D,j,c,{})}}({target:document.body,props:{}})}();
</script>
</body>
</html>`
)

func init() {
	vhost.NotFound = NotFound
}
