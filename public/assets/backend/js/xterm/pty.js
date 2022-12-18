
applyTerminalCommonAddon();
var term, socket, 
  terminalContainer = document.getElementById('terminal-container'), 
  urlPrefix = getQueryStringByName("urlPrefix")||'/server';

function connect() {
  var wsProtocol = window.location.protocol != 'https:' ? 'ws:' : 'wss:';
  var targetUrl = wsProtocol + "//" + document.location.host + urlPrefix
  createTerminal(targetUrl);
}

function createTerminal(targetUrl) {
  // Clean terminal
  while (terminalContainer.children.length) {
    terminalContainer.removeChild(terminalContainer.children[0]);
  }
  term = new Terminal({ cursorBlink: true, fontSize: 17 });
  term.on('resize', function (size) {
    //console.log(size)
    if (socket) socket.send('<RESIZE>' + size.cols + ',' + size.rows + "\n");
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

    // Set terminal size again to set the specific dimensions on the demo
    //setTermSize();

    socket = new WebSocket(targetUrl);
    socket.onopen = function () {
      term.attach(socket);
      term._initialized = true;
    };
    socket.onclose = function () {
      alert("连接已经关闭，如要重新连接，请刷新页面");
      //term.destroy();
    };
    socket.onerror = function () {
      alert("连接出错！请检查服务是否正常或是否拥有权限");
    };
  }, 0);
}
window.addEventListener('load', connect);
