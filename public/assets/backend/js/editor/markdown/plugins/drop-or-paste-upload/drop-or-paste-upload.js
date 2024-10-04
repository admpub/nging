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

    var factory = function (exports) {
        var $ = jQuery;           // if using module loader(Require.js/Sea.js).
        var pluginName = "drop-or-paste-upload";  // 定义插件名称
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
        function upload(file) {
            // File { name: "mrvx5eugoqb.jpg", lastModified: 1533178146229, webkitRelativePath: "", size: 261784, type: "image/jpeg" }
            var _this = this, isImage = _isImage(file), 
                forms = new FormData(), fileName = new Date().getTime() + "." + file.name.split(".").pop();
            forms.append(_this.classPrefix + "image-file", file, fileName);
            _ajax(_this.settings.imageUploadURL, forms, function (ret) {
                if (ret.success == 1) {
                    var url = ret.url;
                    if (isImage) {
                        _this.cm.replaceSelection("![](" + url  + ")");
                    } else {
                        _this.cm.replaceSelection("[下载附件](" + url + ")");
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

            if (!settings.imageUpload || !settings.imageUploadURL) {
                console.log('你还未开启图片上传或者没有配置上传地址');
                return false;
            }

            //监听粘贴板事件
            var $textarea = $('#' + id);
            $textarea.off('paste').on('paste', function (e) {
                e.preventDefault()
                e.stopPropagation()
                var items = (e.clipboardData || e.originalEvent.clipboardData || window.clipboardData).items;
                //判断图片类型
                if (items && items.length) {
                    var file = null;
                    for (var i = 0; i < items.length; i++) {
                        if (_isImage(items[i])) {
                            file = items[i].getAsFile();
                            break;
                        }
                    }

                    if (!file) {
                        console.log("粘贴内容非图片");
                        return;
                    }
                    upload.call(_this, file);
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
                upload.call(_this, files[0]);
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