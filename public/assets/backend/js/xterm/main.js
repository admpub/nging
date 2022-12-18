applyTerminalCommonAddon();
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

var id = getQueryStringByName("id"), urlPrefix = getQueryStringByName("urlPrefix"),
  protocol = getQueryStringByName("protocol"), hostname = getQueryStringByName("hostname"),
  file = getQueryStringByName("file"), port = getQueryStringByName("port"),
  cmd = getQueryStringByName("cmd"), is_debug = getQueryStringByName("debug"),
  user = getQueryStringByName("user"), password = getQueryStringByName("password"),
  casename = getQueryStringByName("name");

if (hostname) document.title = hostname + ' - ' + document.title;

function toggleLogin() {
  var loginEl = document.getElementById("login"), optionsEl = document.getElementById("options");
  changeClassList(optionsEl, "hide", "active")
  var klsList = loginEl.classList;
  if (klsList.contains("hide")) {
    changeClassList(loginEl, "active", "hide")
  } else {
    changeClassList(loginEl, "hide", "active")
  }
}

function toggleOptions() {
  var loginEl = document.getElementById("login"), optionsEl = document.getElementById("options");
  changeClassList(loginEl, "hide", "active")
  var klsList = optionsEl.classList;
  if (klsList.contains("hide")) {
    changeClassList(optionsEl, "active", "hide")
  } else {
    changeClassList(optionsEl, "hide", "active")
  }
}

if (actionElements.findNext) actionElements.findNext.addEventListener('click', function () {
  term.findNext(actionElements.findText.value);
});
if (actionElements.findPrevious) actionElements.findPrevious.addEventListener('click', function () {
  term.findPrevious(actionElements.findText.value);
});
if (actionElements.toggleOptions) actionElements.toggleOptions.addEventListener('click', function () {
  toggleOptions();
});
if (loginElements.login) loginElements.login.addEventListener('click', function () {
  user = loginElements.user.value;
  password = loginElements.password.value;

  toggleLogin();
  connect();
});

function setTerminalSize() {
  var cols = colsElement ? parseInt(colsElement.value, 10) : 80;
  var rows = rowsElement ? parseInt(rowsElement.value, 10) : 32;
  setTermSize(cols, rows);
}

if (colsElement) colsElement.addEventListener('change', setTerminalSize);
if (rowsElement) rowsElement.addEventListener('change', setTerminalSize);

if (optionElements.cursorBlink) optionElements.cursorBlink.addEventListener('change', function () {
  term.setOption('cursorBlink', optionElements.cursorBlink.checked);
});
if (optionElements.cursorStyle) optionElements.cursorStyle.addEventListener('change', function () {
  term.setOption('cursorStyle', optionElements.cursorStyle.value);
});
if (optionElements.bellStyle) optionElements.bellStyle.addEventListener('change', function () {
  term.setOption('bellStyle', optionElements.bellStyle.value);
});
if (optionElements.scrollback) optionElements.scrollback.addEventListener('change', function () {
  term.setOption('scrollback', parseInt(optionElements.scrollback.value, 10));
});
if (optionElements.tabstopwidth) optionElements.tabstopwidth.addEventListener('change', function () {
  term.setOption('tabStopWidth', parseInt(optionElements.tabstopwidth.value, 10));
});
function makeFullUrl() {
  var wsProtocol = window.location.protocol != 'https:' ? 'ws:' : 'wss:';
  var url = wsProtocol + "//" + document.location.host + urlPrefix + "/" + protocol + "?";
  switch (protocol) {
    case "replay":
      url += "file=" + file + "&user=" + user
      break;
    case "ssh_exec":
      url += "dump_file=" + file + "&hostname=" + hostname + "&port=" + port + "&user=" + user + "&cmd=" + cmd
      break;
    default:
      url += "hostname=" + hostname + "&port=" + port + "&user=" + user
  }
  if (id) url += 'id=' + id;
  else if (password) url += '&password=' + password;
  if (is_debug) url += '&debug=' + is_debug;
  return url;
}
function connect() {
  if (protocol == "ssh") {
    if (!id && !password) {
      toggleLogin()
      return
    }
  }
  document.getElementById('title').innerText = casename + ' (' + user + '@' + hostname + ':' + port + ')';
  createTerminal(makeFullUrl());
}

function createTerminal(targetUrl) {
  // Clean terminal
  while (terminalContainer.children.length) {
    terminalContainer.removeChild(terminalContainer.children[0]);
  }
  term = new Terminal({
    cursorBlink: optionElements.cursorBlink ? optionElements.cursorBlink.checked : true,
    scrollback: optionElements.scrollback ? parseInt(optionElements.scrollback.value, 10) : 1000,
    tabStopWidth: optionElements.tabstopwidth ? parseInt(optionElements.tabstopwidth.value, 10) : 8,
    fontSize: 17
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
  if(zmodem) term.zmodem();

  // fit is called within a setTimeout, cols and rows need this.
  setTimeout(function () {
    if (colsElement) colsElement.value = term.cols;
    if (rowsElement) rowsElement.value = term.rows;

    // Set terminal size again to set the specific dimensions on the demo
    setTerminalSize();

    socket = new WebSocket(targetUrl + '&columns=' + term.cols + '&rows=' + term.rows);
    socket.onopen = function () {
      term.attach(socket);
      term._initialized = true;
    };
    socket.onclose = function () {
      //term.destroy();
    };
    socket.onerror = function () {
      alert("连接出错！");
    };
    if(zmodem) zmodemEventBind(term, socket);
  }, 0);
}
function autoDetectPort() {
  if (!protocol) protocol = "ssh";
  switch (protocol) {
    case "telnet":
      if (!port) port = '23';
      break;
    case "ssh":
      if (!port) port = '22';
      break;
    case "replay":
      if (!file) {
        alert("file is empty.")
        return false;
      }
      break;
    default:
      if (!hostname) {
        alert("hostname is empty.")
        return false;
      }
  }
  return true;
}
window.addEventListener('load', function () {
  if (!autoDetectPort()) return;
  if (urlPrefix) {
    if (urlPrefix.charAt(urlPrefix.length - 1) == '/') {
      urlPrefix = urlPrefix.substring(0, urlPrefix.length - 1)
    }
    if (urlPrefix.charAt(0) != '/') {
      urlPrefix = '/' + urlPrefix
    }
  }
  connect();
}, false);