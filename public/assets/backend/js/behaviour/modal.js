(function () {
    function initModalBody($md, ajaxData, onloadCallback, showFooterTabIndexies) {
        var tabs = $md.find('.modal-body .nav-tabs>li');
        $md.find('.modal-footer').removeClass('hidden');
        if (showFooterTabIndexies == null) showFooterTabIndexies = [0];
        tabs.on('click', function () {
            var index = tabs.index(this), show = false;
            for (var i = 0; i < showFooterTabIndexies.length; i++) {
                if (index == showFooterTabIndexies[i]) {
                    show = true;
                    break;
                }
            }
            if (show) {
                $md.find('.modal-footer').removeClass('hidden');
            } else {
                $md.find('.modal-footer:not(.hidden)').addClass('hidden');
            }
        });
        $md.find('table[data-page-size]').each(function () {
            initModalBodyPagination($(this), ajaxData, onloadCallback);
        });
    }
    function initModalBodyPagination(that, ajaxData, onloadCallback) {
        if (!that || that.length < 1) return;
        var data = that.data(), onSwitchPage = null;
        if (onloadCallback) {
            if (data.ajaxList && data.ajaxList in onloadCallback) onloadCallback[data.ajaxList](that);
            if ('switchPage' in onloadCallback) onSwitchPage = onloadCallback['switchPage'];
            if ('pageData' in onloadCallback) data = typeof (onloadCallback['pageData']) == 'function' ? onloadCallback['pageData'](that) : onloadCallback['pageData'];
        }
        App.withPagination(data, null);
        that.after(data.pagination);
        if (!onSwitchPage) {
            onSwitchPage = function (page) {
                var params = { ajaxList: data.ajaxList, page: page };
                params = $.extend(params, ajaxData || {});
                var url = typeof (data['pageUrl']) != 'undefined' ? data['pageUrl'] : window.location.href;
                $.get(url, params, function (res) {
                    var container = that.parent();
                    container.html(res);
                    initModalBodyPagination(container.children('table[data-page-size]'), ajaxData, onloadCallback);
                }, 'html');
            };
        }
        var $paging = that.next('ul.pagination');
        $paging.find('li > a[page]').on('click', function () {
            if ($(this).closest('li.disabled').length > 0) return;
            onSwitchPage($(this).attr('page'));
        });
        $paging.on('refresh', function () {
            var page = 1;
            if ($(this).data('page')) {
                page = $(this).data('page');
                $(this).data('page', false);
            } else {
                page = $(this).find('li.active > a[page]').data('page') || 1;
            }
            onSwitchPage(page);
        });
    }
    function initModalForm(button, modal, fields, afterOpenCallback, onSubmitCallback, multilingualFieldPrefix) {
        if(typeof(fields) == 'object' && fields !== null) {
            var options = fields;
            var fields = null;
            if('fields' in options) fields = options.fields;
            if('afterOpenCallback' in options) afterOpenCallback = options.afterOpenCallback;
            if('onSubmitCallback' in options) onSubmitCallback = options.onSubmitCallback;
            if('multilingualFieldPrefix' in options) multilingualFieldPrefix = options.multilingualFieldPrefix;
        }
        var title = button.data('modal-title'), formData = button.data('form-data'), form = modal.find('form');
        var titleE = modal.find('.modal-header h3'), originalTitle = titleE.data('original-title');
        if (title) {
            if(!originalTitle) titleE.data('original-title',titleE.html());
            titleE.html(title);
        }else if(originalTitle){
            titleE.html(originalTitle);
            titleE.data('original-title',null);
        }
        form[0].reset();
        var multilingual = modal.find('.langset').length > 0, langDefault = '';
        if (multilingual) {
            var firstTabLi = modal.find('.langset > .nav-tabs > li:first');
            if (!firstTabLi.hasClass('active')) {
                firstTabLi.children('a').trigger('click');
            }
            langDefault = firstTabLi.children('a').data('lang');
        }
        if (formData) {
            for (var name in formData) {
                var val = formData[name], input = form.find('[name="' + name + '"]:first');
                if (input.length < 1) continue;
                var type = input.attr('type');
                switch (input[0].tagName.toLowerCase()) {
                    case 'input':
                        if (type === 'checkbox' || type === 'radio') {
                            form.find('[name="' + name + '"][value="' + val + '"]').prop('checked', true);
                        } else {
                            input.val(val);
                        }
                        break;
                    default:
                        input.val(val);
                        break;
                }
            }
        }
        var getFields;
        if(fields && fields.length > 0) {
            getFields = function(){return fields;};
        }else if(multilingual){
            if(!multilingualFieldPrefix) multilingualFieldPrefix = 'Language';
            getFields = function(){
                var prefix = multilingualFieldPrefix + "[" + langDefault + "]";
                var fields = [], fieldsNames = {};
                form.find('[name^="' + prefix + '"]').each(function () {
                    var name = $(this).attr('name'), field = name.replace(prefix + "[", '').replace(']', '');
                    if (fieldsNames[field]) return;
                    fieldsNames[field] = true;
                    fields.push(field);
                });
                return fields;
            }
        }else{
            getFields = function(){
                var fields = [], fieldsNames = {};
                form.find('[name]').each(function () {
                    var name = $(this).attr('name');
                    if (fieldsNames[name]) return;
                    fieldsNames[name] = true;
                    fields.push(name);
                });
                return fields;
            }
        }
        if (afterOpenCallback) {
            var getModalFieldName;
            if(multilingual&&langDefault){
                getModalFieldName = function(name){
                    return 'Language['+langDefault+']['+name+']';
                }
            }else{
                getModalFieldName = function(name){
                    return name;
                }
            }
            afterOpenCallback(button, modal, {'formData':formData, 'multilingual':multilingual, 'langDefault':langDefault, 'getModalFieldName':getModalFieldName, 'getFields':getFields});
        }
        var submitBtn = modal.find('.modal-footer .btn-primary');
        submitBtn.off('click').on('click', function () {
            var form = modal.find('form'), data = {};
            form.find('input,select,textarea').each(function () {
                var name = $(this).attr('name'), val = '';
                switch (this.tagName.toLowerCase()) {
                    case 'input':
                        if (this.type === 'checkbox' || this.type === 'radio') {
                            if (!this.checked) return;
                        } else if (this.type === 'button') {
                            return;
                        }
                        val = $(this).val();
                        break;
                    default:
                        val = $(this).val();
                        break;
                }
                if (name) data[name] = val;
            });
            var values = { data: data, multilingual: multilingual };
            if (multilingual) {
                if(!multilingualFieldPrefix) multilingualFieldPrefix = 'Language';
                if (!fields) {
                    var prefix = multilingualFieldPrefix + "[" + langDefault + "]";
                    form.find('[name^="' + prefix + '"]').each(function () {
                        var name = $(this).attr('name'), field = name.replace(prefix + "[", '').replace(']', '');
                        values[field] = data[name];
                    });
                }else{
                    for (var i = 0; i < fields.length; i++) {
                        var field = fields[i];
                        values[field] = data[multilingualFieldPrefix + "[" + langDefault + "][" + field + "]"];
                    }
                }
                values.langDefault = langDefault;
            } else {
                if (!fields) {
                    for (var name in data) {
                        values[name] = data[name];
                    }
                }else{
                    for (var i = 0; i < fields.length; i++) {
                        var field = fields[i];
                        values[field] = data[field];
                    }
                }
            }
            if (onSubmitCallback) onSubmitCallback(button, modal, values);
        });
    }
    function updateMultilingualFormByModal(parent, values, prefixNames, nameFixer, parentForDefaultLang){
        var prefix = 'Language[', langPrefix = '', translatePrefix = '';
        if(prefixNames){
            for(var i = 0; i < prefixNames.length; i++){
                langPrefix += '['+prefixNames[i]+']';
                if(i==0) {
                    translatePrefix += prefixNames[i];
                }else{
                    translatePrefix += '['+prefixNames[i]+']';
                }
            }
        }
        for(var name in values.data){
            if(!name.startsWith(prefix)) continue;
            var cleanedName = name.substring(prefix.length);
            cleanedName = cleanedName.substring(0, cleanedName.length - 1); // Language[zh-CN][value]
            var params = cleanedName.split(']['); // [zh-CN, value]
            if(params.length!=2) continue;
            if(params[0]==values.langDefault) continue;
            if(nameFixer) params[1] = nameFixer(params[1]);
            var fieldName = 'Language'+langPrefix+'['+params[0]+']['+params[1]+']',field = parent.find('input[type=hidden][name="'+fieldName+'"]');
            if(field.length>0){
                field.val(values.data[name]);
                continue;
            }
            parent.prepend('<input type="hidden" name="'+fieldName+'" value="'+values.data[name]+'" />');
        }
        var fieldName = translatePrefix+'[translate]',
            field = parent.find('input[type=hidden][name="'+fieldName+'"]'),
            value = ('forceTranslate' in values.data)?values.data.forceTranslate:'';
        if(field.length>0){
            field.val(value);
        }else{
            parent.prepend('<input type="hidden" name="'+fieldName+'" value="'+value+'" />');
        }
        if(parentForDefaultLang){
            for(var name in values){
                if(name=='data'||name=='langDefault'||name=='multilingual') continue;
                if(nameFixer) name = nameFixer(name);
                $h.find('input[name="'+translatePrefix+'['+name+']"]').val(values[name]);
            }
        }
    }
    App.initModalBody = initModalBody;
    App.initModalBodyPagination = initModalBodyPagination;
    App.initModalForm = initModalForm;
    App.updateMultilingualFormByModal = updateMultilingualFormByModal;
})();