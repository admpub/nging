
Terminal.applyAddon(attach);
Terminal.applyAddon(fit);
Terminal.applyAddon(fullscreen);
Terminal.applyAddon(search);
Terminal.applyAddon(webLinks);
Terminal.applyAddon(winptyCompat);

var term, socket;

var terminalContainer = document.getElementById('terminal-container'),
    actionElements = {
      findText: document.getElementById('find-text'),
      findNext: document.getElementById('find-next'),
      findPrevious: document.getElementById('find-previous'),
      toggleOptions: document.getElementById('toggle-options'),
    },
    loginElements = {
      user: document.getElementById('userName'),
      password: document.getElementById('password'),
      login: document.getElementById('ssh-login'),
    },
    optionElements = {
      cursorBlink: document.getElementById('option-cursor-blink'),
      cursorStyle: document.getElementById('option-cursor-style'),
      scrollback: document.getElementById('option-scrollback'),
      tabstopwidth: document.getElementById('option-tabstopwidth'),
      bellStyle: document.getElementById('option-bell-style')
    },
    colsElement = document.getElementById('cols'),
    rowsElement = document.getElementById('rows');


var id = getQueryStringByName("id")
var urlPrefix = getQueryStringByName("url_prefix")
var protocol = getQueryStringByName("protocol")
var hostname = getQueryStringByName("hostname")
var file = getQueryStringByName("file")
var port = getQueryStringByName("port")
var cmd = getQueryStringByName("cmd")
var is_debug = getQueryStringByName("debug")
var user = getQueryStringByName("user")
var password = getQueryStringByName("password")
var name = getQueryStringByName("name")
if(hostname) document.title=hostname+' - '+document.title;

function toggleLogin() {
    var loginEl = document.getElementById("login");
    var optionsEl = document.getElementById("options");

    changeClassList(optionsEl, "hide", "active")
    
    var klsList = loginEl.classList;
    if (klsList.contains("hide")) {
      changeClassList(loginEl, "active", "hide")
    } else {
      changeClassList(loginEl, "hide", "active")
    }
}

function toggleOptions() {
    var loginEl = document.getElementById("login");
    var optionsEl = document.getElementById("options");

    changeClassList(loginEl, "hide", "active")

    var klsList = optionsEl.classList;
    if (klsList.contains("hide")) {
      changeClassList(optionsEl, "active", "hide")
    } else {
      changeClassList(optionsEl, "hide", "active")
    }
}

if(actionElements.findNext) actionElements.findNext.addEventListener('click', function() {
  term.findNext(actionElements.findText.value);
});
if(actionElements.findPrevious) actionElements.findPrevious.addEventListener('click', function() {
  term.findPrevious(actionElements.findText.value);
});
if(actionElements.toggleOptions) actionElements.toggleOptions.addEventListener('click',  function() {
  toggleOptions();
});
if(loginElements.login) loginElements.login.addEventListener('click', function() {
    user = loginElements.user.value;
    password = loginElements.password.value;

    toggleLogin();
    connect();
});

function setTerminalSize() {
  var cols = colsElement?parseInt(colsElement.value, 10):80;
  var rows = rowsElement?parseInt(rowsElement.value, 10):32;
  setTermSize(cols,rows);
}

if(colsElement) colsElement.addEventListener('change', setTerminalSize);
if(rowsElement) rowsElement.addEventListener('change', setTerminalSize);

if(optionElements.cursorBlink) optionElements.cursorBlink.addEventListener('change', function () {
  term.setOption('cursorBlink', optionElements.cursorBlink.checked);
});
if(optionElements.cursorStyle) optionElements.cursorStyle.addEventListener('change', function () {
  term.setOption('cursorStyle', optionElements.cursorStyle.value);
});
if(optionElements.bellStyle) optionElements.bellStyle.addEventListener('change', function () {
  term.setOption('bellStyle', optionElements.bellStyle.value);
});
if(optionElements.scrollback) optionElements.scrollback.addEventListener('change', function () {
  term.setOption('scrollback', parseInt(optionElements.scrollback.value, 10));
});
if(optionElements.tabstopwidth) optionElements.tabstopwidth.addEventListener('change', function () {
  term.setOption('tabStopWidth', parseInt(optionElements.tabstopwidth.value, 10));
});
function connect() {
    if(protocol == "ssh") {
      if (!id && (undefined == password || null == password || "" == password)) {
        toggleLogin()
        return
      }
    }
    var wsProtocol=window.location.protocol!='https:'?'ws:':'wss:';
    var target_url = wsProtocol+"//" + document.location.host + urlPrefix + "/" + protocol + "?id=" + id + "&hostname=" + hostname + "&port=" + port + "&user=" + user + "&password=" + password + "&debug=" + is_debug
    if ("replay" == protocol) {
        target_url = wsProtocol+"//" + document.location.host + urlPrefix + "/" + protocol + "?id=" + id + "&file=" + file + "&user=" + user + "&password=" + password
    } else if ("ssh_exec" == protocol) {
        target_url = wsProtocol+"//" + document.location.host + urlPrefix + "/" + protocol + "?id=" + id + "&dump_file=" + file + "&hostname=" + hostname + "&port=" + port + "&user=" + user + "&password=" + password + "&cmd=" + cmd + "&debug=" + is_debug
    }
    document.getElementById('title').innerText=name+' ('+user+'@'+hostname+':'+port+')';
    createTerminal(target_url);
}

function createTerminal(targetUrl) {
  // Clean terminal
  while (terminalContainer.children.length) {
    terminalContainer.removeChild(terminalContainer.children[0]);
  }
  term = new Terminal({
    cursorBlink: optionElements.cursorBlink?optionElements.cursorBlink.checked:true,
    scrollback: optionElements.scrollback?parseInt(optionElements.scrollback.value, 10):1000,
    tabStopWidth: optionElements.tabstopwidth?parseInt(optionElements.tabstopwidth.value, 10):8,
    fontSize:17
  });
  term.on('resize', function (size) {
    //if (!pid) return;
    //var cols = size.cols,rows = size.rows,url = '/terminals/' + pid + '/size?cols=' + cols + '&rows=' + rows;
    //fetch(url, {method: 'POST'});
  });
  //键入字符
  //term.on('data',function(data){console.log('data:'+data)});

  term.open(terminalContainer);
  term.winptyCompatInit();
  term.webLinksInit();
  term.fit();
  term.focus();

  // fit is called within a setTimeout, cols and rows need this.
  setTimeout(function () {
    if(colsElement) colsElement.value = term.cols;
    if(rowsElement) rowsElement.value = term.rows;

    // Set terminal size again to set the specific dimensions on the demo
    setTerminalSize();

    socket = new WebSocket(targetUrl + '&columns=' + term.cols + '&rows=' + term.rows);
    socket.onopen = function() {
      term.attach(socket);
      term._initialized = true;
    };
    socket.onclose = function() {
      //term.destroy();
    };
    socket.onerror = function() {
      alert("连接出错！");
    };
  }, 0);
}

window.addEventListener('load', function () {
    if (undefined == protocol || null == protocol || "" == protocol) {
        protocol = "ssh"
        if (undefined == port || null == port || "" == port) {
            port = "22"
        }
    } else if ("telnet" == protocol) {
        if (undefined == port || null == port || "" == port) {
            port = "23"
        }
    } else if ("ssh" == protocol) {
        if (undefined == port || null == port || "" == port) {
            port = "22"
        }
    }

    if ("replay" == protocol) {
        if (undefined == file || null == file || "" == file) {
            alert("file is empty.")
            return
        }
    } else {
        if (undefined == hostname || null == hostname || "" == hostname) {
            alert("hostname is empty.")
            return
        }
    }

    if(undefined != urlPrefix && null != urlPrefix && "" != urlPrefix) {
      if (urlPrefix[urlPrefix.length-1] == "/") {
        urlPrefix = urlPrefix.substring(0, urlPrefix.length-1)
      }
    }

    if(undefined != urlPrefix && null != urlPrefix && "" != urlPrefix) {
      if (urlPrefix.indexOf("/") != 0) {
        urlPrefix = "/" + urlPrefix
      }
    }
    connect();
}, false);