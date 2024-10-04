/*
 * HTMLParser - This implementation of parser assumes we are parsing HTML in browser
 * and user DOM methods available in browser for parsing HTML.
 *
 * @author Himanshu Gilani
 *
 */

/*
 Universal JavaScript Module, supports AMD (RequireJS), Node.js, and the browser.
 https://gist.github.com/kirel/1268753
*/
(function (name, definition) {
    if (typeof define === 'function') { // AMD
      define(definition);
    } else if (typeof module !== 'undefined' && module.exports) { // Node.js
      module.exports = definition();
    } else { // Browser
      var theModule = definition(), global = this, old = global[name];
      theModule.noConflict = function () {
        global[name] = old;
        return theModule;
      };
      global[name] = theModule;
    }
  })('markdownDOMParser', function() {
  
  var HTMLParser = function (html, handler, opts) {
      opts = opts || {};
  
      var e = document.createElement('div');
      e.innerHTML = html;
      var node = e;
      var nodesToIgnore = opts['nodesToIgnore'] || [];
      var parseHiddenNodes = opts['parseHiddenNodes'] || 'false';
  
      var c = node.childNodes;
      for (var i = 0; i < c.length; i++) {
          try {
              var ignore = false;
              for (var k=0; k< nodesToIgnore.length; k++) {
                  if (c[i].nodeName.toLowerCase() == nodesToIgnore[k]) {
                      ignore= true;
                      break;
                  }
              }
  
              //NOTE hidden node testing is expensive in FF.
              if (ignore || (!parseHiddenNodes && isHiddenNode(c[i]))  ){
                  continue;
              }
  
              if (c[i].nodeName.toLowerCase() != "#text" && c[i].nodeName.toLowerCase() != "#comment") {
                  var attrs = [];
  
                  if (c[i].hasAttributes()) {
                      var attributes = c[i].attributes;
                      for ( var a = 0; a < attributes.length; a++) {
                          var attribute = attributes.item(a);
  
                          attrs.push({
                              name : attribute.nodeName,
                              value : attribute.nodeValue,
                          });
                      }
                  }
  
                  if (handler.start) {
                      if (c[i].hasChildNodes()) {
                          handler.start(c[i].nodeName, attrs, false);
  
                          //if (c[i].nodeName.toLowerCase() == "pre" || c[i].nodeName.toLowerCase() == "code") {
                          //	handler.chars(c[i].innerHTML);
                          //} else
                          if (c[i].nodeName.toLowerCase() == "iframe" || c[i].nodeName.toLowerCase() == "frame") {
                              if (c[i].contentDocument && c[i].contentDocument.documentElement) {
                                  return HTMLParser(c[i].contentDocument.documentElement, handler, opts);
                              }
                          } else {
                              HTMLParser(c[i].innerHTML, handler, opts);
                          }
  
                          if (handler.end) {
                              handler.end(c[i].nodeName);
                          }
                      } else {
                          handler.start(c[i].nodeName, attrs, true);
                      }
                  }
              } else if (c[i].nodeName.toLowerCase() == "#text") {
                  if (handler.chars) {
                      handler.chars(c[i].nodeValue);
                  }
              } else if (c[i].nodeName.toLowerCase() == "#comment") {
                  if (handler.comment) {
                      handler.comment(c[i].nodeValue);
                  }
              }
  
          } catch (e) {
              //properly log error
              console.error(e);
              console.log("error while parsing node: " + c[i].nodeName.toLowerCase());
          }
      }
  };
  
  function isHiddenNode(node) {
      if (node.nodeName.toLowerCase() == "title"){
          return false;
      }
  
      if (window.getComputedStyle) {
          try {
              var style = window.getComputedStyle(node, null);
              if (style.getPropertyValue && style.getPropertyValue('display') == 'none') {
                  return true;
              }
          } catch (e) {
              // consume and ignore. node styles are not accessible
          }
          return false;
      }
  }
  
  return HTMLParser;
  });

/*!
 * editormd附件拖拽和图片粘贴上传插件
 *
 * @file   drop-or-paste-upload.js
 * @author admpub
 * @date   2024年10月04日
 * @link   https://github.com/admpub/nging
 * 引入方式：
 * settings.onload = function() {
 *  var _this=this;
    editormd.loadPlugin("../plugins/drop-or-paste-upload/drop-or-paste-upload", function(){
        _this.dropOrPasteUpload();
    });
 }
 */

(function () {

    function trim(value) {
        if(value===undefined||value==null) return '';
        return value.replace(/^\s+|\s+$/g,"");
    }
    
    function endsWith(value, suffix) {
        return value.match(suffix+"$") == suffix;
    }
    
    function startsWith(value, str) {
        return value.indexOf(str) == 0;
    }
    
    function html2markdown(html, opts) {
        opts = opts || {};
    
        var nodeList = [];
        var listTagStack = [];
        var linkAttrStack = [];
        var blockquoteStack = [];
        var preStack = [];
        var codeStack = [];
        var links = [];
        var inlineStyle = opts['inlineStyle'] || false;
        var parser = opts['parser'];
        var markdownTags = {
            "hr": "- - -\n\n",
            "br": "  \n",
            "title": "# ",
            "h1": "# ",
            "h2": "## ",
            "h3": "### ",
            "h4": "#### ",
            "h5": "##### ",
            "h6": "###### ",
            "b": "**",
            "strong": "**",
            "i": "_",
            "em": "_",
            "dfn": "_",
            "var": "_",
            "cite": "_",
            "span": " ",
            "ul": "* ",
            "ol": "1. ",
            "dl": "- ",
            "blockquote": "> "
        };
    
        if (!parser && typeof markdownDOMParser !== 'undefined') {
            parser = markdownDOMParser;
        }
    
        function getListMarkdownTag() {
            var listItem = "";
            if (listTagStack) {
                for (var i = 0; i < listTagStack.length - 1; i++) {
                    listItem += "  ";
                }
            }
            listItem += peek(listTagStack);
            return listItem;
        }
    
        function convertAttrs(attrs) {
            var attributes = {};
            for (var k in attrs) {
                var attr = attrs[k];
                attributes[attr.name] = attr;
            }
            return attributes;
        }
    
        function peek(list) {
            if (list && list.length > 0) {
                return list.slice(-1)[0];
            }
            return "";
        }
    
        function peekTillNotEmpty(list) {
            if (!list) {
                return "";
            }
    
            for (var i = list.length - 1; i >= 0; i--){
                if (list[i] != "") {
                    return list[i];
                }
            }
            return "";
        }
    
        function removeIfEmptyTag(start) {
            var cleaned = false;
            if (start == peekTillNotEmpty(nodeList)) {
                while (peek(nodeList) != start) {
                    nodeList.pop();
                }
                nodeList.pop();
                cleaned = true;
            }
            return cleaned;
        }
    
        function sliceText(start) {
            var text = [];
            while (nodeList.length > 0 && peek(nodeList) != start) {
                var t = nodeList.pop();
                text.unshift(t);
            }
            return text.join("");
        }
    
        function block(isEndBlock) {
            var lastItem = nodeList.pop();
            if (!lastItem) {
                return;
            }
    
            if (!isEndBlock) {
                var block;
                if (/\s*\n\n\s*$/.test(lastItem)) {
                    lastItem = lastItem.replace(/\s*\n\n\s*$/, "\n\n");
                    block = "";
                } else if (/\s*\n\s*$/.test(lastItem)) {
                    lastItem = lastItem.replace(/\s*\n\s*$/, "\n");
                    block = "\n";
                } else if (/\s+$/.test(lastItem)) {
                    block = "\n\n";
                } else {
                    block = "\n\n";
                }
    
                nodeList.push(lastItem);
                nodeList.push(block);
            } else {
                nodeList.push(lastItem);
                if (!endsWith(lastItem, "\n")) {
                    nodeList.push("\n\n");
                }
            }
        }
    
        function listBlock() {
            if (nodeList.length > 0) {
                var li = peek(nodeList);
    
                if (!endsWith(li, "\n")) {
                    nodeList.push("\n");
                }
            } else {
                nodeList.push("\n");
            }
        }
    
        parser(html, {
            start: function(tag, attrs, unary) {
                tag = tag.toLowerCase();
    
                if (unary && (tag != "br" && tag != "hr" && tag != "img")) {
                    return;
                }
    
                switch (tag) {
                case "br":
                    nodeList.push(markdownTags[tag]);
                    break;
                case "hr":
                    block();
                    nodeList.push(markdownTags[tag]);
                    break;
                case "title":
                case "h1":
                case "h2":
                case "h3":
                case "h4":
                case "h5":
                case "h6":
                    block();
                    nodeList.push(markdownTags[tag]);
                    break;
                case "b":
                case "strong":
                case "i":
                case "em":
                case "dfn":
                case "var":
                case "cite":
                    nodeList.push(markdownTags[tag]);
                    break;
                case "code":
                case "span":
                    if (preStack.length > 0) {
                        break;
                    } else if (!/\s+$/.test(peek(nodeList))) {
                        nodeList.push(markdownTags[tag]);
                    }
                    break;
                case "p":
                case "div":
                case "table":
                case "tbody":
                case "tr":
                case "td":
                    block();
                    break;
                case "ul":
                case "ol":
                case "dl":
                    listTagStack.push(markdownTags[tag]);
                    // lists are block elements
                    if (listTagStack.length > 1) {
                        listBlock();
                    } else {
                        block();
                    }
                    break;
                case "li":
                case "dt":
                    var li = getListMarkdownTag();
                    nodeList.push(li);
                    break;
                case "a":
                    var attribs = convertAttrs(attrs);
                    linkAttrStack.push(attribs);
                    nodeList.push("[");
                    break;
                case "img":
                    var attribs = convertAttrs(attrs);
                    var alt, title, url;
    
                    attribs["src"] ? url = attribs["src"].value : url = "";
                    if (!url) {
                        break;
                    }
    
                    attribs['alt'] ? alt = trim(attribs['alt'].value) : alt = "";
                    attribs['title'] ? title = trim(attribs['title'].value) : title = "";
    
                    // if parent of image tag is nested in anchor tag use inline style
                    if (!inlineStyle && !startsWith(peekTillNotEmpty(nodeList), "[")) {
                        var l = links.indexOf(url);
                        if (l == -1) {
                            links.push(url);
                            l=links.length-1;
                        }
    
                        block();
                        nodeList.push("![");
                        if (alt!= "") {
                            nodeList.push(alt);
                        } else if (title != null) {
                            nodeList.push(title);
                        }
    
                        nodeList.push("][" + l + "]");
                        block();
                    } else {
                        //if image is not a link image then treat images as block elements
                        if (!startsWith(peekTillNotEmpty(nodeList), "[")) {
                            block();
                        }
    
                        nodeList.push("![" + alt + "](" + url + (title ? " \"" + title + "\"" : "") + ")");
    
                        if (!startsWith(peekTillNotEmpty(nodeList), "[")) {
                            block(true);
                        }
                    }
                    break;
                case "blockquote":
                    //listBlock();
                    block();
                    blockquoteStack.push(markdownTags[tag]);
                    break;
                case "pre":
                    block();
                    preStack.push(true);
                    nodeList.push("    ");
                    break;
                case "table":
                    nodeList.push("<table>");
                    break;
                case "thead":
                    nodeList.push("<thead>");
                    break;
                case "tbody":
                    nodeList.push("<tbody>");
                    break;
                case "tr":
                    nodeList.push("<tr>");
                    break;
                case "td":
                    nodeList.push("<td>");
                    break;
                }
            },
            chars: function(text) {
                if (preStack.length > 0) {
                    text = text.replace(/\n/g,"\n    ");
                } else if (trim(text) != "") {
                    text = text.replace(/\s+/g, " ");
    
                    var prevText = peekTillNotEmpty(nodeList);
                    if (/\s+$/.test(prevText)) {
                        text = text.replace(/^\s+/g, "");
                    }
                } else {
                    nodeList.push("");
                    return;
                }
    
                //if(blockquoteStack.length > 0 && peekTillNotEmpty(nodeList).endsWith("\n")) {
                if (blockquoteStack.length > 0) {
                    nodeList.push(blockquoteStack.join(""));
                }
    
                nodeList.push(text);
            },
            end: function(tag) {
                tag = tag.toLowerCase();
    
            switch (tag) {
                case "title":
                case "h1":
                case "h2":
                case "h3":
                case "h4":
                case "h5":
                case "h6":
                    if(!removeIfEmptyTag(markdownTags[tag])) {
                        block(true);
                    }
                    break;
                case "p":
                case "div":
                case "table":
                case "tbody":
                case "tr":
                case "td":
                    while(nodeList.length > 0 && trim(peek(nodeList)) == "") {
                        nodeList.pop();
                    }
                    block(true);
                    break;
                case "b":
                case "strong":
                case "i":
                case "em":
                case "dfn":
                case "var":
                case "cite":
                    if (!removeIfEmptyTag(markdownTags[tag])) {
                        nodeList.push(trim(sliceText(markdownTags[tag])));
                        nodeList.push(markdownTags[tag]);
                    }
                    break;
                case "a":
                    var text = sliceText("[");
                    text = text.replace(/\s+/g, " ");
                    text = trim(text);
    
                    if (text == "") {
                        nodeList.pop();
                        break;
                    }
    
                    var attrs = linkAttrStack.pop();
                    var url;
                    attrs["href"] &&  attrs["href"].value != "" ? url = attrs["href"].value : url = "";
    
                    if (url == "") {
                        nodeList.pop();
                        nodeList.push(text);
                        break;
                    }
    
                    nodeList.push(text);
    
                    if (!inlineStyle && !startsWith(peek(nodeList), "!")){
                        var l = links.indexOf(url);
                        if (l == -1) {
                            links.push(url);
                            l=links.length-1;
                        }
                        nodeList.push("][" + l + "]");
                    } else {
                        if(startsWith(peek(nodeList), "!")){
                            var text = nodeList.pop();
                            text = nodeList.pop() + text;
                            block();
                            nodeList.push(text);
                        }
    
                        var title = attrs["title"];
                        nodeList.push("](" + url + (title ? " \"" + trim(title.value).replace(/\s+/g, " ") + "\"" : "") + ")");
    
                        if(startsWith(peek(nodeList), "!")){
                            block(true);
                        }
                    }
                    break;
                case "ul":
                case "ol":
                case "dl":
                    listBlock();
                    listTagStack.pop();
                    break;
                case "li":
                case "dt":
                    var li = getListMarkdownTag();
                    if (!removeIfEmptyTag(li)) {
                        var text = trim(sliceText(li));
    
                        if (startsWith(text, "[![")) {
                            nodeList.pop();
                            block();
                            nodeList.push(text);
                            block(true);
                        } else {
                            nodeList.push(text);
                            listBlock();
                        }
                    }
                    break;
                case "blockquote":
                    blockquoteStack.pop();
                    break;
                case "pre":
                    //uncomment following experimental code to discard line numbers when syntax highlighters are used
                    //notes this code thorough testing before production user
                    /*
                    var p=[];
                    var flag = true;
                    var count = 0, whiteSpace = 0, line = 0;
                    console.log(">> " + peek(nodeList));
                    while(peek(nodeList).startsWith("    ") || flag == true)
                    {
                        //console.log('inside');
                        var text = nodeList.pop();
                        p.push(text);
    
                        if(flag == true && !text.startsWith("    ")) {
                            continue;
                        } else {
                            flag = false;
                        }
    
                        //var result = parseInt(text.trim());
                        if(!isNaN(text.trim())) {
                            count++;
                        } else if(text.trim() == ""){
                            whiteSpace++;
                        } else {
                            line++;
                        }
                        flag = false;
                    }
    
                    console.log(line);
                    if(line != 0)
                    {
                        while(p.length != 0) {
                            nodeList.push(p.pop());
                        }
                    }
                    */
                    block(true);
                    preStack.pop();
                    break;
                case "code":
                case "span":
                    if (preStack.length > 0) {
                        break;
                    } else if (trim(peek(nodeList)) == "") {
                        nodeList.pop();
                        nodeList.push(markdownTags[tag]);
                    } else {
                        var text = nodeList.pop();
                        nodeList.push(trim(text));
                        nodeList.push(markdownTags[tag]);
                    }
                    break;
                case "table":
                    nodeList.push("</table>");
                    break;
                case "thead":
                    nodeList.push("</thead>");
                    break;
                case "tbody":
                    nodeList.push("</tbody>");
                    break;
                case "tr":
                    nodeList.push("</tr>");
                    break;
                case "td":
                    nodeList.push("</td>");
                    break;
                case "br":
                case "hr":
                case "img":
                    break;
                }
    
            }
        }, {"nodesToIgnore": ["script", "noscript", "object", "iframe", "frame", "head", "style", "label"]});
    
        if (!inlineStyle) {
            for (var i = 0; i < links.length; i++) {
                if (i == 0) {
                    var lastItem = nodeList.pop();
                    nodeList.push(lastItem.replace(/\s+$/g, ""));
                    nodeList.push("\n\n[" + i + "]: " + links[i]);
                } else {
                    nodeList.push("\n[" + i + "]: " + links[i]);
                }
            }
        }
    
        return nodeList.join("");
    
    }
    var factory = function (exports) {
        var $ = jQuery;           // if using module loader(Require.js/Sea.js).
        // ajax上传图片 可自行处理
        function _ajax(url, data, callback) {
            $.ajax({
                "type": 'post',
                "cache": false,
                "url": url,
                "data": data,
                "processData": false,
                "contentType": false,
                "dataType": 'json',
                "mimeType": "multipart/form-data",
                success: callback,
                error: function (err) {
                    console.error('请求失败:', err)
                }
            })
        };
        function _isImage(file) {
            return file.type.indexOf('image/') === 0;
        }
        function upload(file, callback) {
            // File { name: "mrvx5eugoqb.jpg", lastModified: 1533178146229, webkitRelativePath: "", size: 261784, type: "image/jpeg" }
            var _this = this, isImage = _isImage(file), 
                forms = new FormData(), fileName = new Date().getTime() + "." + file.name.split(".").pop();
            forms.append(_this.classPrefix + "image-file", file, fileName);
            _ajax(_this.settings.imageUploadURL, forms, function (ret) {
                if (ret.success == 1) {
                    var url = ret.url, content = '';
                    if (isImage) {
                        content = "![](" + url  + ")";
                    } else {
                        content = "[下载附件](" + url + ")";
                    }
                    if(callback){
                        callback(content)
                    }else{
                        _this.cm.replaceSelection(content);
                    }
                } else {
                    alert(ret.message);
                }
            })
        }
        exports.fn.dropOrPasteUploadDestory = function () {
            var _this = this;
            var id = _this.id;
            $('#' + id).off('paste dragover dragenter drop');
        }
        exports.fn.dropOrPasteUpload = function () {
            var _this = this;
            var settings = _this.settings;
            var id = _this.id;
            var makeHTMLCallback = function(content){
                content = html2markdown(content);
                _this.cm.replaceSelection(content);
            };

            if (!settings.imageUpload || !settings.imageUploadURL) {
                console.log('你还未开启图片上传或者没有配置上传地址');
                return false;
            }
            if(typeof(settings.parsePastedHTML)=='undefined'){
                settings.parsePasteHTML = _this.markdownTextarea.data('parse-pasted-html');
                if(settings.parsePasteHTML==undefined) settings.parsePasteHTML = true;
            }
            var parseHTML=settings.parsePasteHTML?true:false;
            //监听粘贴板事件
            var $textarea = $('#' + id);//<div id="id" class="editormd">
            $textarea.off('paste').on('paste', function (e) {
                var clipboardData = (e.clipboardData || e.originalEvent.clipboardData || window.clipboardData);
                var items = clipboardData.items;
                if (items && items.length) {
                    var did = false, stopEvent = function(){
                        if (did) return;
                        did = true
                        e.preventDefault()
                        e.stopPropagation()
                    };
                    for (var i = 0; i < items.length; i++) {
                        switch(items[i].kind){
                            case 'file':
                                if (_isImage(items[i])) {
                                    stopEvent()
                                    upload.call(_this, items[i].getAsFile());
                                }
                                break;
                            case 'string':
                                if(parseHTML && items[i].type=='text/html'){
                                    stopEvent()
                                    items[i].getAsString(makeHTMLCallback);
                                }
                                break;
                        }
                    }

                }
            });
            $textarea.off("dragover").off("dragenter").on("dragover dragenter", function (e) {
                e.preventDefault()
                e.stopPropagation()
            });
            $textarea.off("drop").on("drop", function (e) {
                e.preventDefault()
                e.stopPropagation()
                var files = this.files || e.originalEvent.dataTransfer.files;
                if(files && files.length){
                    for (var i = 0; i < files.length; i++) {
                        if (_isImage(files[i])) {
                            upload.call(_this, files[i]);
                        }
                    }
                }
            });
        };
    };

    // CommonJS/Node.js
    if (typeof require === "function" && typeof exports === "object" && typeof module === "object") {
        module.exports = factory;
    }
    else if (typeof define === "function")  // AMD/CMD/Sea.js
    {
        if (define.amd) { // for Require.js
            define(["editormd"], function (editormd) {
                factory(editormd);
            });

        } else { // for Sea.js
            define(function (require) {
                var editormd = require("../../editormd");
                factory(editormd);
            });
        }
    } else {
        factory(window.editormd);
    }

})();