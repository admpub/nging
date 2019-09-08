/*
 * 设计器私有的配置说明 
 * 一
 * UE.FormDesignUrl  插件路径
 * 
 * 二
 *UE.getEditor('myFormDesign',{
 *          toolFormDesign:true,//是否显示，设计器的清单 tool
 */
UE.FormDesignUrl = UE.FormDesignUrl || 'third-party/formdesign';

function FormDesign(idName, options) {
    var defaults = {
        //allowDivTransToP: false,//阻止转换div 为p
        toolFormDesign: true, //是否显示，设计器的 toolbars
        textarea: 'design_content',
        //这里可以选择自己需要的工具按钮名称,此处仅选择如下五个
        toolbars: [
            [
                'fullscreen', 'source', '|', 'undo', 'redo', '|', 'bold', 'italic', 'underline', 'fontborder', 'strikethrough', 'removeformat', '|', 'forecolor', 'backcolor', 'insertorderedlist', 'insertunorderedlist', '|', 'fontfamily', 'fontsize', '|', 'indent', '|', 'justifyleft', 'justifycenter', 'justifyright', 'justifyjustify', '|', 'link', 'unlink', '|', 'horizontal', 'spechars', 'wordimage', '|', 'inserttable', 'deletetable', 'mergecells', 'splittocells'
            ]
        ],
        //focus时自动清空初始化时的内容
        //autoClearinitialContent:true,
        //关闭字数统计
        wordCount: false,
        //关闭elementPath
        elementPathEnabled: false,
        //默认的编辑区域高度
        initialFrameHeight: 300
            //,iframeCssUrl:"css/bootstrap/css/bootstrap.css" //引入自身 css使编辑器兼容你网站css
            //更多其他参数，请参考ueditor.config.js中的配置项
    };
    if (options) {
        for (var k in options) {
            defaults[k] = options[k];
        }
    }
    var editor = UE.getEditor(idName, defaults);

    window.webxFormDesign = {
        editor: editor,
        /*执行控件*/
        exec: function(method) {
            this.editor.execCommand(method);
        },
        showmsg: function(content,timeout){
            return this.editor.trigger('showmessage',{
                content : content,
                timeout : timeout||2000
            });
        },
        hidemsg:function(id){
            this.editor.trigger('hidemessage',id);
        },
        /*
            Javascript 解析表单
            template 表单设计器里的Html内容
            fields 字段总数
        */
        parse: function(template, fields) {
            var preg = /(<span(((?!<span).)*form-widget=\"(radios|checkboxs|select)\".*?)>(.*?)<\/span>|<(img|input|textarea|select).*?(<\/select>|<\/textarea>|\/>))/gi,
                preg_attr = /(\w+)=\"(.?|.+?)\"/gi,
                preg_group = /<input.*?\/>/gi;
            if (!fields) fields = 0;

            var template_parse = template,
                template_data = new Array(),
                add_fields = new Object(),
                checkboxs = 0;

            var pno = 0;
            template.replace(preg, function(plugin, p1, p2, p3, p4, p5, p6) {
                var attr_arr_all = new Object(),
                    name = '',select_dot = '',
                    is_new = false, p0 = plugin, tag = p6 ? p6 : p4;
                if (tag == 'radios' || tag == 'checkboxs') {
                    plugin = p2;
                } 
                plugin.replace(preg_attr, function(str0, attr, val) {
                        if (attr == 'name') {
                            if (val == 'FormNewField') {
                                is_new = true;
                                fields++;
                                val = 'data_' + fields;
                            }
                            name = val;
                        }
                        if (tag == 'select' && attr == 'value') {
                            if (!attr_arr_all[attr]) attr_arr_all[attr] = '';
                            attr_arr_all[attr] += select_dot + val;
                            select_dot = ',';
                        } else {
                            attr_arr_all[attr] = val;
                        }
                });
                attr_arr_all["widget"] = tag;
                /*复选组  多个字段 */
                if (tag == 'checkboxs') {
                    plugin = p0;
                    var name = 'checkboxs_' + checkboxs;
                    attr_arr_all['parseName'] = name;
                    attr_arr_all['name'] = '';
                    attr_arr_all['value'] = '';

                    attr_arr_all['content'] = '<span form-widget="checkboxs" title="' + attr_arr_all['title'] + '">';
                    var dot_name = '',
                        dot_value = '';
                    p5.replace(preg_group, function(parse_group) {
                        var is_new = false,
                            option = new Object();
                        parse_group.replace(preg_attr, function(str0, k, val) {
                            if (k == 'name') {
                                if (val == 'FormNewField') {
                                    is_new = true;
                                    fields++;
                                    val = 'data_' + fields;
                                }
                                attr_arr_all['name'] += dot_name + val;
                                dot_name = ',';
                            } else if (k == 'value') {
                                attr_arr_all['value'] += dot_value + val;
                                dot_value = ',';
                            }
                            option[k] = val;
                        });

                        if (!attr_arr_all['options']) attr_arr_all['options'] = new Array();
                        attr_arr_all['options'].push(option);
                        //if(!option['checked']) option['checked'] = '';
                        var checked = option['checked'] != undefined ? 'checked="checked"' : '';
                        attr_arr_all['content'] += '<input type="checkbox" name="' + option['name'] + '" value="' + option['value'] + '"  ' + checked + '/>' + option['value'] + '&nbsp;';

                        if (is_new) {
                            var arr = new Object();
                            arr['name'] = option['name'];
                            arr['widget'] = attr_arr_all['widget'];
                            add_fields[option['name']] = arr;
                        }
                    });
                    attr_arr_all['content'] += '</span>';

                    //parse
                    template = template.replace(plugin, attr_arr_all['content']);
                    template_parse = template_parse.replace(plugin, '{' + name + '}');
                    template_data[pno] = attr_arr_all;
                    checkboxs++;

                } else if (name) {
                    /* 单选组 一个字段 */ 
                    if (tag == 'radios'){
                        plugin = p0;
                        attr_arr_all['value'] = '';
                        attr_arr_all['content'] = '<span form-widget="radios" name="' + attr_arr_all['name'] + '" title="' + attr_arr_all['title'] + '">';
                        var dot = '';
                        p5.replace(preg_group, function(parse_group) {
                            var option = new Object();
                            parse_group.replace(preg_attr, function(str0, k, val) {
                                if (k == 'value') {
                                    attr_arr_all['value'] += dot + val;
                                    dot = ',';
                                }
                                option[k] = val;
                            });
                            option['name'] = attr_arr_all['name'];
                            if (!attr_arr_all['options']) attr_arr_all['options'] = new Array();
                            attr_arr_all['options'].push(option);
                            //if(!option['checked']) option['checked'] = '';
                            var checked = option['checked'] != undefined ? 'checked="checked"' : '';
                            attr_arr_all['content'] += '<input type="radio" name="' + attr_arr_all['name'] + '" value="' + option['value'] + '"  ' + checked + '/>' + option['value'] + '&nbsp;';

                        });
                        attr_arr_all['content'] += '</span>';
                    } else {
                        attr_arr_all['content'] = is_new ? plugin.replace(/FormNewField/, name) : plugin;
                    }
                    //attr_arr_all['itemid'] = fields;
                    //attr_arr_all['tag'] = tag;
                    template = template.replace(plugin, attr_arr_all['content']);
                    template_parse = template_parse.replace(plugin, '{' + name + '}');
                    if (is_new) {
                        var arr = new Object();
                        arr['name'] = name;
                        arr['widget'] = attr_arr_all['widget'];
                        add_fields[arr['name']] = arr;
                    }
                    template_data[pno] = attr_arr_all;
                }
                pno++;
            });
            var form = new Object({
                'fields': fields, //总字段数
                'template': template, //完整html
                'parse': template_parse, //控件替换为{data_1}的html
                'data': template_data, //控件属性
                'add_fields': add_fields //新增控件
            });
            return JSON.stringify(form);
        },
        /*type  =  save 保存设计 versions 保存版本  close关闭 */
        fnCheckForm: function(type) {
            if (this.editor.queryCommandState('source'))
                this.editor.execCommand('source'); //切换到编辑模式才提交，否则有bug

            if (this.editor.hasContents()) {
                this.editor.sync(); /*同步内容*/
                if(defaults.fnSave) defaults.fnSave.call(this);
            } else {
                this.showmsg('表单内容不能为空！');
                return false;
            }
        },
        /*预览表单*/
        fnReview: function() {
            if (this.editor.queryCommandState('source'))
                this.editor.execCommand('source'); /*切换到编辑模式才提交，否则部分浏览器有bug*/

            if (this.editor.hasContents()) {
                this.editor.sync(); /*同步内容*/
                if(defaults.fnView) defaults.fnView.call(this);
            } else {
                this.showmsg('表单内容不能为空！');
                return false;
            }
        }
    };
}
/**
 * 文本框
 * @command textfield
 * @method execCommand
 * @param { String } cmd 命令字符串
 * @example
 * ```javascript
 * editor.execCommand( 'textfield');
 * ```
 */
UE.plugins['text'] = function() {
    var me = this,
        thePlugins = 'text';
    me.commands[thePlugins] = {
        execCommand: function() {
            var dialog = new UE.ui.Dialog({
                iframeUrl: this.options.UEDITOR_HOME_URL + UE.FormDesignUrl + '/text.html',
                name: thePlugins,
                editor: this,
                title: '文本框',
                cssRules: "width:600px;height:310px;",
                buttons: [{
                        className: 'edui-okbutton',
                        label: '确定',
                        onclick: function() {
                            dialog.close(true);
                        }
                    },
                    {
                        className: 'edui-cancelbutton',
                        label: '取消',
                        onclick: function() {
                            dialog.close(false);
                        }
                    }
                ]
            });
            dialog.render();
            dialog.open();
        }
    };
    var popup = new baidu.editor.ui.Popup({
        editor: this,
        content: '',
        className: 'edui-bubble',
        _edittext: function() {
            baidu.editor.plugins[thePlugins].editdom = popup.anchorEl;
            me.execCommand(thePlugins);
            this.hide();
        },
        _delete: function() {
            if (window.confirm('确认删除该控件吗？')) {
                baidu.editor.dom.domUtils.remove(this.anchorEl, false);
            }
            this.hide();
        }
    });
    popup.render();
    me.addListener('mouseover', function(t, evt) {
        evt = evt || window.event;
        var el = evt.target || evt.srcElement;
        var widget = el.getAttribute('form-widget');
        if (/input/ig.test(el.tagName) && widget == thePlugins) {
            var html = popup.formatHtml(
                '<nobr>文本框: <span onclick=$$._edittext() class="edui-clickable">编辑</span>&nbsp;&nbsp;<span onclick=$$._delete() class="edui-clickable">删除</span></nobr>');
            if (html) {
                popup.getDom('content').innerHTML = html;
                popup.anchorEl = el;
                popup.showAnchor(popup.anchorEl);
            } else {
                popup.hide();
            }
        }
    });
};
/**
 * 宏控件
 * @command macros
 * @method execCommand
 * @param { String } cmd 命令字符串
 * @example
 * ```javascript
 * editor.execCommand( 'macros');
 * ```
 */
UE.plugins['macros'] = function() {
    var me = this,
        thePlugins = 'macros';
    me.commands[thePlugins] = {
        execCommand: function() {
            var dialog = new UE.ui.Dialog({
                iframeUrl: this.options.UEDITOR_HOME_URL + UE.FormDesignUrl + '/macros.html',
                name: thePlugins,
                editor: this,
                title: '宏控件',
                cssRules: "width:600px;height:270px;",
                buttons: [{
                        className: 'edui-okbutton',
                        label: '确定',
                        onclick: function() {
                            dialog.close(true);
                        }
                    },
                    {
                        className: 'edui-cancelbutton',
                        label: '取消',
                        onclick: function() {
                            dialog.close(false);
                        }
                    }
                ]
            });
            dialog.render();
            dialog.open();
        }
    };
    var popup = new baidu.editor.ui.Popup({
        editor: this,
        content: '',
        className: 'edui-bubble',
        _edittext: function() {
            baidu.editor.plugins[thePlugins].editdom = popup.anchorEl;
            me.execCommand(thePlugins);
            this.hide();
        },
        _delete: function() {
            if (window.confirm('确认删除该控件吗？')) {
                baidu.editor.dom.domUtils.remove(this.anchorEl, false);
            }
            this.hide();
        }
    });
    popup.render();
    me.addListener('mouseover', function(t, evt) {
        evt = evt || window.event;
        var el = evt.target || evt.srcElement;
        var widget = el.getAttribute('form-widget');
        if (/input/ig.test(el.tagName) && widget == thePlugins) {
            var html = popup.formatHtml(
                '<nobr>宏控件: <span onclick=$$._edittext() class="edui-clickable">编辑</span>&nbsp;&nbsp;<span onclick=$$._delete() class="edui-clickable">删除</span></nobr>');
            if (html) {
                popup.getDom('content').innerHTML = html;
                popup.anchorEl = el;
                popup.showAnchor(popup.anchorEl);
            } else {
                popup.hide();
            }
        }
    });
};
/**
 * 单选框组
 * @command radios
 * @method execCommand
 * @param { String } cmd 命令字符串
 * @example
 * ```javascript
 * editor.execCommand( 'radio');
 * ```
 */
UE.plugins['radios'] = function() {
    var me = this,
        thePlugins = 'radios';
    me.commands[thePlugins] = {
        execCommand: function() {
            var dialog = new UE.ui.Dialog({
                iframeUrl: this.options.UEDITOR_HOME_URL + UE.FormDesignUrl + '/radios.html',
                name: thePlugins,
                editor: this,
                title: '单选框组',
                cssRules: "width:590px;height:370px;",
                buttons: [{
                        className: 'edui-okbutton',
                        label: '确定',
                        onclick: function() {
                            dialog.close(true);
                        }
                    },
                    {
                        className: 'edui-cancelbutton',
                        label: '取消',
                        onclick: function() {
                            dialog.close(false);
                        }
                    }
                ]
            });
            dialog.render();
            dialog.open();
        }
    };
    var popup = new baidu.editor.ui.Popup({
        editor: this,
        content: '',
        className: 'edui-bubble',
        _edittext: function() {
            baidu.editor.plugins[thePlugins].editdom = popup.anchorEl;
            me.execCommand(thePlugins);
            this.hide();
        },
        _delete: function() {
            if (window.confirm('确认删除该控件吗？')) {
                baidu.editor.dom.domUtils.remove(this.anchorEl, false);
            }
            this.hide();
        }
    });
    popup.render();
    me.addListener('mouseover', function(t, evt) {
        evt = evt || window.event;
        var el = evt.target || evt.srcElement;
        var widget = el.getAttribute('form-widget');
        if (/span/ig.test(el.tagName) && widget == thePlugins) {
            var html = popup.formatHtml(
                '<nobr>单选框组: <span onclick=$$._edittext() class="edui-clickable">编辑</span>&nbsp;&nbsp;<span onclick=$$._delete() class="edui-clickable">删除</span></nobr>');
            if (html) {
                var elInput = el.getElementsByTagName("input");
                var rEl = elInput.length > 0 ? elInput[0] : el;
                popup.getDom('content').innerHTML = html;
                popup.anchorEl = el;
                popup.showAnchor(rEl);
            } else {
                popup.hide();
            }
        }
    });
};
/**
 * 复选框组
 * @command checkboxs
 * @method execCommand
 * @param { String } cmd 命令字符串
 * @example
 * ```javascript
 * editor.execCommand( 'checkboxs');
 * ```
 */
UE.plugins['checkboxs'] = function() {
    var me = this,
        thePlugins = 'checkboxs';
    me.commands[thePlugins] = {
        execCommand: function() {
            var dialog = new UE.ui.Dialog({
                iframeUrl: this.options.UEDITOR_HOME_URL + UE.FormDesignUrl + '/checkboxs.html',
                name: thePlugins,
                editor: this,
                title: '复选框组',
                cssRules: "width:600px;height:400px;",
                buttons: [{
                        className: 'edui-okbutton',
                        label: '确定',
                        onclick: function() {
                            dialog.close(true);
                        }
                    },
                    {
                        className: 'edui-cancelbutton',
                        label: '取消',
                        onclick: function() {
                            dialog.close(false);
                        }
                    }
                ]
            });
            dialog.render();
            dialog.open();
        }
    };
    var popup = new baidu.editor.ui.Popup({
        editor: this,
        content: '',
        className: 'edui-bubble',
        _edittext: function() {
            baidu.editor.plugins[thePlugins].editdom = popup.anchorEl;
            me.execCommand(thePlugins);
            this.hide();
        },
        _delete: function() {
            if (window.confirm('确认删除该控件吗？')) {
                baidu.editor.dom.domUtils.remove(this.anchorEl, false);
            }
            this.hide();
        }
    });
    popup.render();
    me.addListener('mouseover', function(t, evt) {
        evt = evt || window.event;
        var el = evt.target || evt.srcElement;
        var widget = el.getAttribute('form-widget');
        if (/span/ig.test(el.tagName) && widget == thePlugins) {
            var html = popup.formatHtml(
                '<nobr>复选框组: <span onclick=$$._edittext() class="edui-clickable">编辑</span>&nbsp;&nbsp;<span onclick=$$._delete() class="edui-clickable">删除</span></nobr>');
            if (html) {
                var elInput = el.getElementsByTagName("input");
                var rEl = elInput.length > 0 ? elInput[0] : el;
                popup.getDom('content').innerHTML = html;
                popup.anchorEl = el;
                popup.showAnchor(rEl);
            } else {
                popup.hide();
            }
        }
    });
};
/**
 * 多行文本框
 * @command textarea
 * @method execCommand
 * @param { String } cmd 命令字符串
 * @example
 * ```javascript
 * editor.execCommand( 'textarea');
 * ```
 */
UE.plugins['textarea'] = function() {
    var me = this,
        thePlugins = 'textarea';
    me.commands[thePlugins] = {
        execCommand: function() {
            var dialog = new UE.ui.Dialog({
                iframeUrl: this.options.UEDITOR_HOME_URL + UE.FormDesignUrl + '/textarea.html',
                name: thePlugins,
                editor: this,
                title: '多行文本框',
                cssRules: "width:600px;height:330px;",
                buttons: [{
                        className: 'edui-okbutton',
                        label: '确定',
                        onclick: function() {
                            dialog.close(true);
                        }
                    },
                    {
                        className: 'edui-cancelbutton',
                        label: '取消',
                        onclick: function() {
                            dialog.close(false);
                        }
                    }
                ]
            });
            dialog.render();
            dialog.open();
        }
    };
    var popup = new baidu.editor.ui.Popup({
        editor: this,
        content: '',
        className: 'edui-bubble',
        _edittext: function() {
            baidu.editor.plugins[thePlugins].editdom = popup.anchorEl;
            me.execCommand(thePlugins);
            this.hide();
        },
        _delete: function() {
            if (window.confirm('确认删除该控件吗？')) {
                baidu.editor.dom.domUtils.remove(this.anchorEl, false);
            }
            this.hide();
        }
    });
    popup.render();
    me.addListener('mouseover', function(t, evt) {
        evt = evt || window.event;
        var el = evt.target || evt.srcElement;
        if (/textarea/ig.test(el.tagName)) {
            var html = popup.formatHtml(
                '<nobr>多行文本框: <span onclick=$$._edittext() class="edui-clickable">编辑</span>&nbsp;&nbsp;<span onclick=$$._delete() class="edui-clickable">删除</span></nobr>');
            if (html) {
                popup.getDom('content').innerHTML = html;
                popup.anchorEl = el;
                popup.showAnchor(popup.anchorEl);
            } else {
                popup.hide();
            }
        }
    });
};
/**
 * 下拉菜单
 * @command select
 * @method execCommand
 * @param { String } cmd 命令字符串
 * @example
 * ```javascript
 * editor.execCommand( 'select');
 * ```
 */
UE.plugins['select'] = function() {
    var me = this,
        thePlugins = 'select';
    me.commands[thePlugins] = {
        execCommand: function() {
            var dialog = new UE.ui.Dialog({
                iframeUrl: this.options.UEDITOR_HOME_URL + UE.FormDesignUrl + '/select.html',
                name: thePlugins,
                editor: this,
                title: '下拉菜单',
                cssRules: "width:590px;height:370px;",
                buttons: [{
                        className: 'edui-okbutton',
                        label: '确定',
                        onclick: function() {
                            dialog.close(true);
                        }
                    },
                    {
                        className: 'edui-cancelbutton',
                        label: '取消',
                        onclick: function() {
                            dialog.close(false);
                        }
                    }
                ]
            });
            dialog.render();
            dialog.open();
        }
    };
    var popup = new baidu.editor.ui.Popup({
        editor: this,
        content: '',
        className: 'edui-bubble',
        _edittext: function() {
            baidu.editor.plugins[thePlugins].editdom = popup.anchorEl;
            me.execCommand(thePlugins);
            this.hide();
        },
        _delete: function() {
            if (window.confirm('确认删除该控件吗？')) {
                baidu.editor.dom.domUtils.remove(this.anchorEl, false);
            }
            this.hide();
        }
    });
    popup.render();
    me.addListener('mouseover', function(t, evt) {
        evt = evt || window.event;
        var el = evt.target || evt.srcElement;
        var widget = el.getAttribute('form-widget');
        if (/select|span/ig.test(el.tagName) && widget == thePlugins) {
            var html = popup.formatHtml(
                '<nobr>下拉菜单: <span onclick=$$._edittext() class="edui-clickable">编辑</span>&nbsp;&nbsp;<span onclick=$$._delete() class="edui-clickable">删除</span></nobr>');
            if (html) {
                if (el.tagName == 'SPAN') {
                    var elInput = el.getElementsByTagName("select");
                    el = elInput.length > 0 ? elInput[0] : el;
                }
                popup.getDom('content').innerHTML = html;
                popup.anchorEl = el;
                popup.showAnchor(popup.anchorEl);
            } else {
                popup.hide();
            }
        }
    });

};
/**
 * 进度条
 * @command progressbar
 * @method execCommand
 * @param { String } cmd 命令字符串
 * @example
 * ```javascript
 * editor.execCommand( 'progressbar');
 * ```
 */
UE.plugins['progressbar'] = function() {
    var me = this,
        thePlugins = 'progressbar';
    me.commands[thePlugins] = {
        execCommand: function() {
            var dialog = new UE.ui.Dialog({
                iframeUrl: this.options.UEDITOR_HOME_URL + UE.FormDesignUrl + '/progressbar.html',
                name: thePlugins,
                editor: this,
                title: '进度条',
                cssRules: "width:600px;height:450px;",
                buttons: [{
                        className: 'edui-okbutton',
                        label: '确定',
                        onclick: function() {
                            dialog.close(true);
                        }
                    },
                    {
                        className: 'edui-cancelbutton',
                        label: '取消',
                        onclick: function() {
                            dialog.close(false);
                        }
                    }
                ]
            });
            dialog.render();
            dialog.open();
        }
    };
    var popup = new baidu.editor.ui.Popup({
        editor: this,
        content: '',
        className: 'edui-bubble',
        _edittext: function() {
            baidu.editor.plugins[thePlugins].editdom = popup.anchorEl;
            me.execCommand(thePlugins);
            this.hide();
        },
        _delete: function() {
            if (window.confirm('确认删除该控件吗？')) {
                baidu.editor.dom.domUtils.remove(this.anchorEl, false);
            }
            this.hide();
        }
    });
    popup.render();
    me.addListener('mouseover', function(t, evt) {
        evt = evt || window.event;
        var el = evt.target || evt.srcElement;
        var widget = el.getAttribute('form-widget');
        if (/img/ig.test(el.tagName) && widget == thePlugins) {
            var html = popup.formatHtml(
                '<nobr>进度条: <span onclick=$$._edittext() class="edui-clickable">编辑</span>&nbsp;&nbsp;<span onclick=$$._delete() class="edui-clickable">删除</span></nobr>');
            if (html) {
                popup.getDom('content').innerHTML = html;
                popup.anchorEl = el;
                popup.showAnchor(popup.anchorEl);
            } else {
                popup.hide();
            }
        }
    });
};
/**
 * 二维码
 * @command qrcode
 * @method execCommand
 * @param { String } cmd 命令字符串
 * @example
 * ```javascript
 * editor.execCommand( 'qrcode');
 * ```
 */
UE.plugins['qrcode'] = function() {
    var me = this,
        thePlugins = 'qrcode';
    me.commands[thePlugins] = {
        execCommand: function() {
            var dialog = new UE.ui.Dialog({
                iframeUrl: this.options.UEDITOR_HOME_URL + UE.FormDesignUrl + '/qrcode.html',
                name: thePlugins,
                editor: this,
                title: '二维码',
                cssRules: "width:600px;height:370px;",
                buttons: [{
                        className: 'edui-okbutton',
                        label: '确定',
                        onclick: function() {
                            dialog.close(true);
                        }
                    },
                    {
                        className: 'edui-cancelbutton',
                        label: '取消',
                        onclick: function() {
                            dialog.close(false);
                        }
                    }
                ]
            });
            dialog.render();
            dialog.open();
        }
    };
    var popup = new baidu.editor.ui.Popup({
        editor: this,
        content: '',
        className: 'edui-bubble',
        _edittext: function() {
            baidu.editor.plugins[thePlugins].editdom = popup.anchorEl;
            me.execCommand(thePlugins);
            this.hide();
        },
        _delete: function() {
            if (window.confirm('确认删除该控件吗？')) {
                baidu.editor.dom.domUtils.remove(this.anchorEl, false);
            }
            this.hide();
        }
    });
    popup.render();
    me.addListener('mouseover', function(t, evt) {
        evt = evt || window.event;
        var el = evt.target || evt.srcElement;
        var widget = el.getAttribute('form-widget');
        if (/img/ig.test(el.tagName) && widget == thePlugins) {
            var html = popup.formatHtml(
                '<nobr>二维码: <span onclick=$$._edittext() class="edui-clickable">编辑</span>&nbsp;&nbsp;<span onclick=$$._delete() class="edui-clickable">删除</span></nobr>');
            if (html) {
                popup.getDom('content').innerHTML = html;
                popup.anchorEl = el;
                popup.showAnchor(popup.anchorEl);
            } else {
                popup.hide();
            }
        }
    });
};
/**
 * 列表控件
 * @command listctrl
 * @method execCommand
 * @param { String } cmd 命令字符串
 * @example
 * ```javascript
 * editor.execCommand( 'qrcode');
 * ```
 */
UE.plugins['listctrl'] = function() {
    var me = this,
        thePlugins = 'listctrl';
    me.commands[thePlugins] = {
        execCommand: function() {
            var dialog = new UE.ui.Dialog({
                iframeUrl: this.options.UEDITOR_HOME_URL + UE.FormDesignUrl + '/listctrl.html',
                name: thePlugins,
                editor: this,
                title: '列表控件',
                cssRules: "width:800px;height:400px;",
                buttons: [{
                        className: 'edui-okbutton',
                        label: '确定',
                        onclick: function() {
                            dialog.close(true);
                        }
                    },
                    {
                        className: 'edui-cancelbutton',
                        label: '取消',
                        onclick: function() {
                            dialog.close(false);
                        }
                    }
                ]
            });
            dialog.render();
            dialog.open();
        }
    };
    var popup = new baidu.editor.ui.Popup({
        editor: this,
        content: '',
        className: 'edui-bubble',
        _edittext: function() {
            baidu.editor.plugins[thePlugins].editdom = popup.anchorEl;
            me.execCommand(thePlugins);
            this.hide();
        },
        _delete: function() {
            if (window.confirm('确认删除该控件吗？')) {
                baidu.editor.dom.domUtils.remove(this.anchorEl, false);
            }
            this.hide();
        }
    });
    popup.render();
    me.addListener('mouseover', function(t, evt) {
        evt = evt || window.event;
        var el = evt.target || evt.srcElement;
        var widget = el.getAttribute('form-widget');
        if (/input/ig.test(el.tagName) && widget == thePlugins) {
            var html = popup.formatHtml(
                '<nobr>列表控件: <span onclick=$$._edittext() class="edui-clickable">编辑</span>&nbsp;&nbsp;<span onclick=$$._delete() class="edui-clickable">删除</span></nobr>');
            if (html) {
                popup.getDom('content').innerHTML = html;
                popup.anchorEl = el;
                popup.showAnchor(popup.anchorEl);
            } else {
                popup.hide();
            }
        }
    });
};
UE.plugins['error'] = function() {
    var me = this,
        thePlugins = 'error';
    me.commands[thePlugins] = {
        execCommand: function() {
            var dialog = new UE.ui.Dialog({
                iframeUrl: this.options.UEDITOR_HOME_URL + UE.FormDesignUrl + '/error.html',
                name: thePlugins,
                editor: this,
                title: '异常提示',
                cssRules: "width:400px;height:130px;",
                buttons: [{
                    className: 'edui-okbutton',
                    label: '确定',
                    onclick: function() {
                        dialog.close(true);
                    }
                }]
            });
            dialog.render();
            dialog.open();
        }
    };
};
UE.plugins['widgets'] = function() {
    var me = this,
        thePlugins = 'widgets';
    me.commands[thePlugins] = {
        execCommand: function() {
            var dialog = new UE.ui.Dialog({
                iframeUrl: this.options.UEDITOR_HOME_URL + UE.FormDesignUrl + '/widgets.html',
                name: thePlugins,
                editor: this,
                title: '表单设计器 - 清单',
                cssRules: "width:620px;height:220px;",
                buttons: [{
                    className: 'edui-okbutton',
                    label: '确定',
                    onclick: function() {
                        dialog.close(true);
                    }
                }]
            });
            dialog.render();
            dialog.open();
        }
    };
};
UE.plugins['form_template'] = function() {
    var me = this,
        thePlugins = 'form_template';
    me.commands[thePlugins] = {
        execCommand: function() {
            var dialog = new UE.ui.Dialog({
                iframeUrl: this.options.UEDITOR_HOME_URL + UE.FormDesignUrl + '/template.html',
                name: thePlugins,
                editor: this,
                title: '表单模板',
                cssRules: "width:640px;height:380px;",
                buttons: [{
                    className: 'edui-okbutton',
                    label: '确定',
                    onclick: function() {
                        dialog.close(true);
                    }
                }]
            });
            dialog.render();
            dialog.open();
        }
    };
};

UE.registerUI('button_form_design', function(editor, uiName) {
    if (!this.options.toolFormDesign) {
        return false;
    }
    //注册按钮执行时的command命令，使用命令默认就会带有回退操作
    editor.registerCommand(uiName, {
        execCommand: function() {
            editor.execCommand('widgets');
        }
    });
    //创建一个button
    var btn = new UE.ui.Button({
        //按钮的名字
        name: uiName,
        //提示
        title: "表单设计器",
        //需要添加的额外样式，指定icon图标，这里默认使用一个重复的icon
        cssRules: 'background-position: -401px -40px;',
        //点击时执行的命令
        onclick: function() {
            //这里可以不用执行命令,做你自己的操作也可
            editor.execCommand(uiName);
        }
    });
    return btn;
});
UE.registerUI('button_template', function(editor, uiName) {
    if (!this.options.toolFormDesign) {
        return false;
    }
    //注册按钮执行时的command命令，使用命令默认就会带有回退操作
    editor.registerCommand(uiName, {
        execCommand: function() {
            try {
                webxFormDesign.exec('form_template');
                //webxFormDesign.fnCheckForm('save');
            } catch (e) {
                webxFormDesign.showmsg('打开模板异常');
            }

        }
    });
    //创建一个button
    var btn = new UE.ui.Button({
        //按钮的名字
        name: uiName,
        //提示
        title: "表单模板",
        //需要添加的额外样式，指定icon图标，这里默认使用一个重复的icon
        cssRules: 'background-position: -339px -40px;',
        //点击时执行的命令
        onclick: function() {
            //这里可以不用执行命令,做你自己的操作也可
            editor.execCommand(uiName);
        }
    });
    return btn;
});
UE.registerUI('button_preview', function(editor, uiName) {
    if (!this.options.toolFormDesign) {
        return false;
    }
    //注册按钮执行时的command命令，使用命令默认就会带有回退操作
    editor.registerCommand(uiName, {
        execCommand: function() {
            try {
                webxFormDesign.fnReview();
            } catch (e) {
                webxFormDesign.showmsg('webxFormDesign.fnReview 预览异常');
            }
        }
    });
    //创建一个button
    var btn = new UE.ui.Button({
        //按钮的名字
        name: uiName,
        //提示
        title: "预览",
        //需要添加的额外样式，指定icon图标，这里默认使用一个重复的icon
        cssRules: 'background-position: -420px -19px;',
        //点击时执行的命令
        onclick: function() {
            //这里可以不用执行命令,做你自己的操作也可
            editor.execCommand(uiName);
        }
    });
    return btn;
});

UE.registerUI('button_save', function(editor, uiName) {
    if (!this.options.toolFormDesign) {
        return false;
    }
    //注册按钮执行时的command命令，使用命令默认就会带有回退操作
    editor.registerCommand(uiName, {
        execCommand: function() {
            try {
                webxFormDesign.fnCheckForm('save');
            } catch (e) {
                webxFormDesign.showmsg('webxFormDesign.fnCheckForm("save") 保存异常');
            }
        }
    });
    //创建一个button
    var btn = new UE.ui.Button({
        //按钮的名字
        name: uiName,
        //提示
        title: "保存表单",
        //需要添加的额外样式，指定icon图标，这里默认使用一个重复的icon
        cssRules: 'background-position: -481px -20px;',
        //点击时执行的命令
        onclick: function() {
            //这里可以不用执行命令,做你自己的操作也可
            editor.execCommand(uiName);
        }
    });
    return btn;
});