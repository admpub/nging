App.loader.libs.editormdPreview = ['#editor/markdown/lib/marked.min.js', '#editor/markdown/lib/prettify.min.js', '#editor/markdown/lib/raphael.min.js', '#editor/markdown/lib/underscore.min.js', '#editor/markdown/css/editormd.preview.min.css', '#editor/markdown/editormd.min.js'];
App.loader.libs.editormd = ['#editor/markdown/css/editormd.min.css', '#editor/markdown/editormd.min.js'];
App.loader.libs.flowChart = ['#editor/markdown/lib/flowchart.min.js', '#editor/markdown/lib/jquery.flowchart.min.js'];
App.loader.libs.sequenceDiagram = ['#editor/markdown/lib/sequence-diagram.min.js'];
App.loader.libs.xheditor = ['#editor/xheditor/xheditor.min.js', '#editor/xheditor/xheditor_lang/' + App.lang + '.js'];
App.loader.libs.ueditor = ['#editor/ueditor/ueditor.config.js', '#editor/ueditor/ueditor.all.min.js'];
App.loader.libs.summernote = ['#editor/summernote/summernote.css', '#editor/summernote/summernote.min.js', '#editor/summernote/lang/summernote-' + App.langTag() + '.js'];
App.loader.libs.summernote_bs4 = ['#editor/summernote/summernote-bs4.css', '#editor/summernote/summernote-bs4.min.js', '#editor/summernote/lang/summernote-' + App.langTag() + '.js'];
App.loader.libs.markdownit = ['#markdown/it/markdown-it.min.js', '#markdown/it/plugins/emoji/markdown-it-emoji.min.js'];
App.loader.libs.codehighlight = ['#markdown/it/plugins/highlight/loader/prettify.js', '#markdown/it/plugins/highlight/loader/run_prettify.js?skin=sons-of-obsidian'];
window.UEDITOR_HOME_URL = ASSETS_URL + '/js/editor/ueditor/';

App.editor = {
	browsingFileURL: App.loader.siteURL + '/manager/select_file'
};

/* 解析markdown为html */
App.editor.markdownToHTML = function (viewZoneId, markdownData, options) {
	var defaults = {
		markdown: markdownData,
		//htmlDecode    : true,  // 开启HTML标签解析，为了安全性，默认不开启
		htmlDecode: "style,script,iframe",  // you can filter tags decode
		//toc           : false,
		tocm: true,  // Using [TOCM]
		//gfm           : false,
		//tocDropdown   : true,
		emoji: true,
		taskList: true,
		tex: true,  // 默认不解析
		flowChart: true,  // 默认不解析
		sequenceDiagram: true,  // 默认不解析
	};
	var params = $.extend({}, defaults, options || {});
	if (params.flowChart) App.loader.defined(typeof ($.fn.flowChart), 'flowChart');
	if (params.sequenceDiagram) App.loader.defined(typeof ($.fn.sequenceDiagram), 'sequenceDiagram');
	App.loader.defined(typeof (editormd), 'editormdPreview');
	var EditormdView = editormd.markdownToHTML(viewZoneId, params);
	return EditormdView;
};

// =================================================================
// ueditor
// =================================================================

App.editor.ueditors = function (editorElement, uploadUrl, options) {
	$(editorElement).each(function () {
		App.editor.ueditor(this, uploadUrl, options);
	});
};
/* 初始化UEditor编辑器 */
App.editor.ueditor = function (editorElement, uploadUrl, options) {
	if ($(editorElement).hasClass('form-control')) $(editorElement).removeClass('form-control');
	if (!uploadUrl) uploadUrl = $(editorElement).attr('action');
	if (uploadUrl.substr(0, 1) == '!') uploadUrl = uploadUrl.substr(1);
	if (options == null) options = {};
	App.loader.defined(typeof (window.UE), 'ueditor');
	if (uploadUrl != null) {
		if (uploadUrl.indexOf('?') >= 0) {
			uploadUrl += '&';
		} else {
			uploadUrl += '?';
		}
		uploadUrl += 'format=json&';
		uploadUrl += 'client=webuploader';
		options.serverUrl = uploadUrl;
	}
	var idv = $(editorElement).attr('id');
	if (!idv) {
		idv = 'ueditor-instance-' + Math.random();
		$(editorElement).attr('id', idv);
	}
	var editor = UE.getEditor(idv, options);
	$(editorElement).data('editor-name', 'ueditor');
	$(editorElement).data('editor-object', editor);
};

// =================================================================
// editormd
// =================================================================

App.editor.markdowns = function (editorElement, uploadUrl, options) {
	$(editorElement).each(function () {
		App.editor.markdown(this, uploadUrl, options);
	});
};
/* 初始化Markdown编辑器 */
App.editor.markdown = function (editorElement, uploadUrl, options) {
	var isManager = false;
	if (!uploadUrl) uploadUrl = $(editorElement).attr('action');
	if (uploadUrl.substr(0, 1) == '!') {
		uploadUrl = uploadUrl.substr(1);
		isManager = true;
	}
	App.loader.defined(typeof (editormd), 'editormd');
	if (uploadUrl != null) {
		if (uploadUrl.indexOf('?') >= 0) {
			uploadUrl += '&';
		} else {
			uploadUrl += '?';
		}
		if (!isManager) uploadUrl += 'format=json&';
		uploadUrl += 'filetype=image&client=markdown';
	}
	var container = $(editorElement).parent(),
		containerId = container.attr('id');
	if (containerId === undefined) {
		containerId = 'webx-md-' + window.location.href.replace(/[^\w]+/g, '-');
		container.attr('id', containerId);
	};
	var path = BACKEND_URL + '/public/assets/backend/js/editor/markdown/';
	var defaults = {
		width: "100%",
		height: container.height(),
		path: path + 'lib/',
		markdown: $(editorElement).text(),
		codeFold: true,
		saveHTMLToTextarea: true,			// 保存HTML到Textarea
		searchReplace: true,
		watch: true,						// 关闭实时预览
		htmlDecode: "style,script,iframe",	// 开启HTML标签解析，为了安全性，默认不开启
		emoji: true,
		taskList: true,
		tocm: true,					 // Using [TOCM]
		tex: true,                   // 开启科学公式TeX语言支持，默认关闭
		flowChart: true,             // 开启流程图支持，默认关闭
		sequenceDiagram: true,       // 开启时序/序列图支持，默认关闭,
		imageUpload: true,
		imageFormats: ["jpg", "jpeg", "gif", "png", "bmp"],
		imageUploadURL: uploadUrl,
		crossDomainUpload: true,
		uploadCallbackURL: path + 'plugins/image-dialog/upload_callback.htm',
		onload: function () { }
	};
	var params = $.extend({}, defaults, options || {});
	if (isManager) {
		params.toolbarIcons = function () {
			// Or return editormd.toolbarModes[name]; // full, simple, mini
			return [
				"undo", "redo", "|",
				"bold", "del", "italic", "quote", "ucwords", "uppercase", "lowercase", "|",
				"h1", "h2", "h3", "h4", "h5", "h6", "|",
				"list-ul", "list-ol", "hr", "|",
				"link", "reference-link", "browsing-image", "code", "preformatted-text", "code-block", "table", "datetime", "emoji", "html-entities", "pagebreak", "|",
				"goto-line", "watch", "preview", "fullscreen", "clear", "search", "|",
				"help", "info"
			];
		};
		params.toolbarIconsClass = {
			'browsing-image': "fa-image"
		};
		params.toolbarIconTexts = {
			'browsing-image': App.loader.t('选择图片')
		};
		params.toolbarHandlers = {
			'browsing-image': function (cm, icon, cursor, selection) {
				Coscms.Dialog.Modal(App.editor.browsingFileURL + '?pagerows=12&filetype=image&multiple=1', {
					title: App.loader.t('选择图片'),
					width: '600px',
					submit: function (dialog) {
						var ck = dialog.find('input[type=checkbox][name="id[]"]:checked');
						if (ck.length <= 0) {
							App.loader.noty({ type: 'error', text: T('没有选择任何选项！') });
						} else {
							var urls = [];
							ck.each(function () {
								var v = $(this).data('raw');
								urls.push('![' + v.Name + '](' + v.ViewUrl + ')');
							});
							//var linenum=urls.length>0?urls.length-1:0;
							urls = urls.join('\n') + '\n';
							cm.replaceSelection(urls);
							//if(selection === "") cm.setCursor(cm.line+linenum, cm.ch+1);
							dialog.modal('hide');
						}
					},
					cancel: function (dialog) {
					}
				}, null).css('z-index', 20030902);

			}
		};
		params.lang = {
			toolbar: {
				'browsing-image': App.loader.t("从服务器选择图片")
			}
		};
	}
	if (!uploadUrl) params.imageUpload = false;
	var editor = editormd(containerId, params);
	$(editorElement).data('editor-name', 'markdown');
	$(editorElement).data('editor-object', editor);
	return editor;
};
App.editor.md = App.editor.markdown;

// =================================================================
// XHEditor
// =================================================================

App.editor.xheditors = function (editorElement, uploadUrl, options) {
	$(editorElement).each(function () {
		App.editor.xheditor(this, uploadUrl, options);
	});
};
/* 初始化xheditor */
App.editor.xheditor = function (editorElement, uploadUrl, settings) {
	App.loader.defined(typeof ($.fn.xheditor), 'xheditor');
	if (!uploadUrl) uploadUrl = $(editorElement).attr('action');
	var editor, editorRoot = BACKEND_URL + '/public/assets/backend/js/editor/xheditor/';
	if (!uploadUrl) { editor = $(editorElement).xheditor({ 'editorRoot': editorRoot }); } else {
		if (uploadUrl.indexOf('?') >= 0) {
			uploadUrl += '&';
		} else {
			uploadUrl += '?';
		}
		if (uploadUrl.substr(0, 1) == '!') {
			settings = $.extend({
				'modalWidth': 620,
				'modalHeight': 635,
				'upBtnText': App.loader.t('浏览')
			}, settings || {});
		} else {
			uploadUrl += 'format=json&';
		}

		uploadUrl += 'client=xheditor';
		var plugins = {
			Code: {
				c: 'xhe_btnCode', t: '插入代码', h: 1, e: function () {
					var that = this;
					var lang = ["erlang", "go", "html", "javascript", "php", "scala", "sql", "xquery", "xml", "yaml", "yml"];
					var htmlCode = '<div><select id="xheCodeType">';
					for (var i = 0; i < lang.length; i++) {
						var s = lang[i] == 'go' ? ' selected="selected"' : '';
						htmlCode += '<option value="' + lang[i] + '"' + s + '>' + lang[i] + '</option>';
					}
					htmlCode += '<option value="">其它</option></select></div><div><textarea id="xheCodeValue" wrap="soft" spellcheck="false" style="width:300px;height:100px;" /></div><div style="text-align:right;"><input type="button" id="xheSave" value="确定" /></div>';
					var jCode = $(htmlCode), jType = $('#xheCodeType', jCode),
						jValue = $('#xheCodeValue', jCode), jSave = $('#xheSave', jCode);
					jSave.click(function () {
						that.loadBookmark();
						that.pasteHTML('<pre class="prettyprint linenums lang-' + jType.val() + '">' + that.domEncode(jValue.val()) + '</pre>');
						that.hidePanel();
						return false;
					});
					that.saveBookmark();
					that.showDialog(jCode);
				}
			},
			EndInput: {
				c: 'xhe_btnEndInput', t: '末尾新行 (Shift+End)', s: 'shift+end', e: function () {
					this.appendHTML('<p><br /></p>');/*解决光标无法移出容器的问题*/
				}
			}
		};
		var option = {
			'skin': 'default',//'shortcuts':{'ctrl+enter':submitForm},'loadCSS':'<style></style>',
			'plugins': plugins,
			'upLinkUrl': uploadUrl + '&filetype=file',
			'upLinkExt': "zip,rar,7z,tar,gz,txt,xls,doc,docx,ppt,pptx,et,wps,rtf,dps",
			'upImgUrl': uploadUrl + '&filetype=image',
			'upImgExt': "jpg,jpeg,gif,png",
			'upFlashUrl': uploadUrl + '&filetype=flash',
			'upFlashExt': "swf",
			'upMediaUrl': uploadUrl + '&filetype=media',
			'upMediaExt': "avi,wmv,wma,mp3,mp4,mpeg,mkv,rm,rmv,mid",
			'editorRoot': editorRoot
		};
		option = $.extend(option, settings || {});
		/* IE10以下不支持HTML5中input:file域的mutiple属性，采用iframe加载swfupload实现批量选择上传 */
		if ($.browser.msie && parseFloat($.browser.version) < 10.0) {
			uploadUrl = '!{editorRoot}xheditor_plugins/multiupload/multiupload.html?uploadurl=' + encodeURIComponent(uploadUrl);
			if (option.upLinkUrl) {
				option.upLinkUrl = uploadUrl + '&ext=Attachment(' + '*.' + option.upLinkExt.replace(/,/g, ';*.') + ')';
				option.upLinkExt = '';
			}
			if (option.upImgUrl) {
				option.upImgUrl = uploadUrl + '&ext=Image(' + '*.' + option.upImgExt.replace(/,/g, ';*.') + ')';
				option.upImgExt = '';
			}
			if (option.upFlashUrl) {
				option.upFlashUrl = uploadUrl + '&ext=Flash(' + '*.' + option.upFlashExt.replace(/,/g, ';*.') + ')';
				option.upFlashExt = '';
			}
			if (option.upMediaUrl) {
				option.upMediaUrl = uploadUrl + '&ext=Media(' + '*.' + option.upMediaUrl.replace(/,/g, ';*.') + ')';
				option.upMediaExt = '';
			}
		}
		editor = $(editorElement).xheditor(option);
	}
	$(editorElement).data('editor-name', 'xheditor');
	$(editorElement).data('editor-object', editor);
	return editor;
};

// =================================================================
// summernote
// =================================================================

App.editor.summernotes = function (elem, minHeight) {
	$(editorElement).each(function () {
		App.editor.summernote(this, minHeight);
	});
};
App.editor.summernote = function (elem, minHeight) {
	if (minHeight == null) minHeight = 400;
	App.loader.defined(typeof ($.fn.summernote), 'summernote');
	$(elem).summernote({
		lang: App.langTag(),
		minHeight: minHeight,
		callbacks: {
			onImageUpload: function (files, editor, $editable) {
				var $files = $(files);
				$files.each(function () {
					var file = this;
					var formdata = new FormData();
					formdata.append("files[]", file);
					$.ajax({
						data: formdata,
						type: "POST",
						url: $(elem).attr('action'),
						cache: false,
						contentType: false,
						processData: false,
						dataType: "json",
						success: function (r) {
							if (r.Code != 1) {
								return App.message({ title: App.i18n.SYS_INFO, text: r.Info, time: 5000, sticky: false, class_name: r.Code == 1 ? 'success' : 'error' });
							}
							$.each(r.Data.files, function (index, file) {
								$(elem).summernote('insertImage', file, function ($image) { });
							});
						},
						error: function () {
							alert(App.i18n.UPLOAD_ERR);
						}
					});
				});
			}
		}
	});
};

// =================================================================
// markdownit
// =================================================================

App.editor.markdownItToHTML = function markdownParse(box, isContainer) {
	App.loader.defined(typeof (window.markdownit), 'markdownit');
	App.loader.defined(typeof (window.prettyPrint), 'codehighlight');
	if (isContainer != false) box = box.find('.markdown-code');
	var md = App.editor.markdownItInstance();
	box.each(function () {
		$(this).html(md.render($.trim($(this).html())));
		$(this).find("pre > code").each(function () {
			$(this).parent("pre").addClass("prettyprint linenums");
		});

		if (typeof (prettyPrint) !== "undefined") prettyPrint();
	});
};
App.editor.markdownItInstance = function () {
	App.loader.defined(typeof (window.markdownit), 'markdownit');
	var md = window.markdownit();
	if (typeof (window.markdownitEmoji) != 'undefined') md.use(window.markdownitEmoji);
	return md;
};


//例如：App.editor.switch($('textarea'))
App.editor.switch = function (texta, cancelFn, tips) {
	var upurl = texta.attr('action') || texta.data("upload-url") || '!' + App.editor.browsingFileURL + '?pagerows=12&multiple=1',
		etype = texta.data("editor"),
		ename = texta.data("editor-name"),
		eobject = texta.data("editor-object"),
		ctype = texta.data("current-editor");
	if (ctype == etype) return;
	var className = texta.data("class");
	if (className === undefined) {
		className = texta.attr("class");
		if (!className) className = '';
		texta.data("class", className);
	}
	var content = texta.data("content-elem"), cElem = content;
	if (content) cElem = App.loader.parseTmpl(content, { type: etype });
	var obj = texta.get(0);
	switch (etype) {
		case 'markdown':
			if (tips) {
				var cc = App.loader.parseTmpl(content, { type: ctype });
				if (cc && $(cc).length > 0) {
					if (texta.val() != $(cc).val() && !confirm('确定要切换吗？切换编辑器将会丢失您当前所做的修改。')) {
						if (cancelFn != null) cancelFn();
						return false;
					}
				}
			}
			switch (ename) {
				case 'xheditor':
					if (typeof (texta.xheditor) != 'undefined') {
						texta.xheditor(false);
					}
					break;
				case 'ueditor':
					eobject.destroy();
					break;
			}
			if (cElem && $(cElem).length > 0) {
				texta.text($(cElem).val());
				texta.val($(cElem).val());
			}

			App.editor.markdown(obj, upurl);
			texta.data("current-editor", etype);
			break;
		default:
			if (cElem && $(cElem).length > 0) {
				var cc = App.loader.parseTmpl(content, { type: ctype });
				if (cc && $(cc).length > 0) {
					$(cc).text(texta.val());
					var ht = $('textarea[name="' + texta.parent().attr('id') + '-html-code"]');
					if (ht.length > 0 && ht.val() != "") {
						$(cElem).text(ht.val());
						$(cElem).val(ht.val());
					}
				};
				texta.val($(cElem).val());
				texta.text($(cElem).val());
			}
			texta.parent().removeAttr('class');
			texta.attr('class', className).siblings().remove();
			App.editor.xheditor(obj, upurl);
			//App.editor.initUE(obj, upurl);
			texta.data("current-editor", "html");
	};
	return true;
};
