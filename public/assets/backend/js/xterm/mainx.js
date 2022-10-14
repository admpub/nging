
Terminal.applyAddon(attach);
Terminal.applyAddon(fit);
Terminal.applyAddon(fullscreen);
Terminal.applyAddon(search);
Terminal.applyAddon(webLinks);
Terminal.applyAddon(winptyCompat);
Terminal.applyAddon(zmodem);

var term,
    socket

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


var urlPrefix = getQueryStringByName("url_prefix")
var protocol = getQueryStringByName("protocol")
var hostname = getQueryStringByName("hostname")
var file = getQueryStringByName("file")
var port = getQueryStringByName("port")
var cmd = getQueryStringByName("cmd")
var is_debug = getQueryStringByName("debug")
var user = getQueryStringByName("user")
var password = getQueryStringByName("password")

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
  var viewportElement = document.querySelector('.xterm-viewport');
  var scrollBarWidth = viewportElement.offsetWidth - viewportElement.clientWidth;
  var width = (cols * term.charMeasure.width + 20 /*room for scrollbar*/).toString() + 'px';
  var height = (rows * term.charMeasure.height).toString() + 'px';

  terminalContainer.style.width = width;
  terminalContainer.style.height = height;
  term.resize(cols, rows);
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
      if (undefined == password || null == password || "" == password) {
        toggleLogin()
        return
      }
    }

    var target_url = "ws://" + document.location.host + urlPrefix + "/" + protocol + "?hostname=" + hostname + "&port=" + port + "&user=" + user + "&password=" + password + "&debug=" + is_debug
    if ("replay" == protocol) {
        target_url = "ws://" + document.location.host + urlPrefix + "/" + protocol + "?file=" + file + "&user=" + user + "&password=" + password
    } else if ("ssh_exec" == protocol) {
        target_url = "ws://" + document.location.host + urlPrefix + "/" + protocol + "?dump_file=" + file + "&hostname=" + hostname + "&port=" + port + "&user=" + user + "&password=" + password + "&cmd=" + cmd + "&debug=" + is_debug
    }

    createTerminal(target_url);
}

function createTerminal(targetUrl) {
  // Clean terminal
  while (terminalContainer.children.length) {
    terminalContainer.removeChild(terminalContainer.children[0]);
  }
  term = new Terminal({
    cursorBlink: optionElements.cursorBlink?optionElements.cursorBlink.checked:false,
    scrollback: optionElements.scrollback?parseInt(optionElements.scrollback.value, 10):0,
    tabStopWidth: optionElements.tabstopwidth?parseInt(optionElements.tabstopwidth.value, 10):0
  });
  term.on('resize', function (size) {
    //if (!pid) {
    //  return;
    //}
    //var cols = size.cols,
    //    rows = size.rows,
    //    url = '/terminals/' + pid + '/size?cols=' + cols + '&rows=' + rows;

    //fetch(url, {method: 'POST'});
  });

  term.open(terminalContainer);
  term.winptyCompatInit();
  term.webLinksInit();
  term.fit();
  term.focus();
  term.zmodem();


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


  term.zmodemAttach(socket, {
      noTerminalWriteOutsideSession: true,
  } );

  term.on("zmodemRetract", () => {
      start_form.style.display = "none";
      start_form.onsubmit = null;
  });

  term.on("zmodemDetect", (detection) => {
      function do_zmodem() {
          term.detach();
          let zsession = detection.confirm();

          var promise;

          if (zsession.type === "receive") {
              promise = _handle_receive_session(zsession);
          }
          else {
              promise = _handle_send_session(zsession);
          }

          promise.catch( console.error.bind(console) ).then( () => {
              term.attach(socket);
          } );
      }

      if (_auto_zmodem()) {
          do_zmodem();
      }
      else {
          start_form.style.display = "";
          start_form.onsubmit = function(e) {
              start_form.style.display = "none";

              if (document.getElementById("zmstart_yes").checked) {
                  do_zmodem();
              }
              else {
                  detection.deny();
              }
          };
      }
  });


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


//----------------------------------------------------------------------
// UI STUFF

function _show_file_info(xfer) {
  var file_info = xfer.get_details();

  document.getElementById("name").textContent = file_info.name;
  document.getElementById("size").textContent = file_info.size;
  document.getElementById("mtime").textContent = file_info.mtime;
  document.getElementById("files_remaining").textContent = file_info.files_remaining;
  document.getElementById("bytes_remaining").textContent = file_info.bytes_remaining;

  document.getElementById("mode").textContent = "0" + file_info.mode.toString(8);

  var xfer_opts = xfer.get_options();
  ["conversion", "management", "transport", "sparse"].forEach( (lbl) => {
      document.getElementById(`zfile_${lbl}`).textContent = xfer_opts[lbl];
  } );

  document.getElementById("zm_file").style.display = "";
}
function _hide_file_info() {
  document.getElementById("zm_file").style.display = "none";
}

function _save_to_disk(xfer, buffer) {
  return Zmodem.Browser.save_to_disk(buffer, xfer.get_details().name);
}

var skipper_button = document.getElementById("zm_progress_skipper");
var skipper_button_orig_text = skipper_button.textContent;

function _show_progress() {
  skipper_button.disabled = false;
  skipper_button.textContent = skipper_button_orig_text;

  document.getElementById("bytes_received").textContent = 0;
  document.getElementById("percent_received").textContent = 0;

  document.getElementById("zm_progress").style.display = "";
}

function _update_progress(xfer) {
  var total_in = xfer.get_offset();

  document.getElementById("bytes_received").textContent = total_in;

  var percent_received = 100 * total_in / xfer.get_details().size;
  document.getElementById("percent_received").textContent = percent_received.toFixed(2);
}

function _hide_progress() {
  document.getElementById("zm_progress").style.display = "none";
}

var start_form = document.getElementById("zm_start");

function _auto_zmodem() {
  return document.getElementById("zmodem-auto").checked;
}

// END UI STUFF
//----------------------------------------------------------------------

function _handle_receive_session(zsession) {
  zsession.on("offer", function(xfer) {
      current_receive_xfer = xfer;

      _show_file_info(xfer);

      var offer_form = document.getElementById("zm_offer");

      function on_form_submit() {
          offer_form.style.display = "none";

          //START
          //if (offer_form.zmaccept.value) {
          if (_auto_zmodem() || document.getElementById("zmaccept_yes").checked) {
              _show_progress();

              var FILE_BUFFER = [];
              xfer.on("input", (payload) => {
                  _update_progress(xfer);
                  FILE_BUFFER.push( new Uint8Array(payload) );
              });
              xfer.accept().then(
                  () => {
                      _save_to_disk(xfer, FILE_BUFFER);
                  },
                  console.error.bind(console)
              );
          }
          else {
              xfer.skip();
          }
          //END
      }

      if (_auto_zmodem()) {
          on_form_submit();
      }
      else {
          offer_form.onsubmit = on_form_submit;
          offer_form.style.display = "";
      }
  } );

  var promise = new Promise( (res) => {
      zsession.on("session_end", () => {
          _hide_file_info();
          _hide_progress();
          res();
      } );
  } );

  zsession.start();

  return promise;
}

function _handle_send_session(zsession) {
  var choose_form = document.getElementById("zm_choose");
  choose_form.style.display = "";

  var file_el = document.getElementById("zm_files");

  var promise = new Promise( (res) => {
      file_el.onchange = function(e) {
          choose_form.style.display = "none";

          var files_obj = file_el.files;

          Zmodem.Browser.send_files(
              zsession,
              files_obj,
              {
                  on_offer_response(obj, xfer) {
                      if (xfer) _show_progress();
                      //console.log("offer", xfer ? "accepted" : "skipped");
                  },
                  on_progress(obj, xfer) {
                      _update_progress(xfer);
                  },
                  on_file_complete(obj) {
                      //console.log("COMPLETE", obj);
                      _hide_progress();
                  },
              }
          ).then(_hide_progress).then(
              zsession.close.bind(zsession),
              console.error.bind(console)
          ).then( () => {
              _hide_file_info();
              _hide_progress();
              res();
          } );
      };
  } );

  return promise;
}

//This is here to allow canceling of an in-progress ZMODEM transfer.
var current_receive_xfer;

//Called from HTML directly.
function skip_current_file() {
  current_receive_xfer.skip();

  skipper_button.disabled = true;
  skipper_button.textContent = "Waiting for server to acknowledge skip …";
}
