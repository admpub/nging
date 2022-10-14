//根据QueryString参数名称获取值
function getQueryStringByName(name) {
    var result = location.search.match(new RegExp("[\?\&]" + name + "=([^\&]+)", "i"));
    if (result == null || result.length < 1) {
        return "";
    }
    var value = decodeURIComponent(result[1]);
    //console.log(name+':'+value);
    return value;
}

function startsWith(s, prefix) {
    return s.indexOf(prefix) == 0;
}

function changeClassList(ele, add, del) {
    var klsList = ele.classList;
    klsList.add(add);
    klsList.remove(del);
}

function toggleFullscreen() {
  term.toggleFullScreen();
}

function setTermSize(cols, rows) {
  if(!term) return;
  var viewportElement = document.querySelector('.xterm-viewport');
  var scrollBarWidth = viewportElement.offsetWidth - viewportElement.clientWidth;
  var width = (cols * term.charMeasure.width + 20 /*room for scrollbar*/).toString() + 'px';
  var height = (rows * term.charMeasure.height).toString() + 'px';

  terminalContainer.style.width = width;
  terminalContainer.style.height = height;
  term.resize(cols, rows);
}
