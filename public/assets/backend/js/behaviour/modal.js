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
    function parseLangFieldName(name,prefix){
        name = name.substring((prefix + "[").length);
        return name.substring(0, name.length - 1);
    }
    function setFormFieldValue(input, val) {
        if (input.length < 1) return;
        switch (input[0].tagName.toLowerCase()) {
            case 'input':
                var type = input.attr('type');
                if (type === 'checkbox') {
                    input.prop('checked', false);
                    if(val === undefined || val === null) return;
                    input.filter('[value="' + val + '"]').prop('checked', true);
                } else if(type === 'radio'){
                    if(val === undefined || val === null) return;
                    input.filter('[value="' + val + '"]').prop('checked', true);
                }else {
                    if(val === undefined || val === null) return;
                    input.val(val);
                }
                break;
            case 'select':
                if(val === undefined || val === null) return;
                input.find('option[value="' + val + '"]').prop('selected', true);
                break;
            default:
                if(val === undefined || val === null) return;
                input.val(val);
                break;
        }
    }
    function getFormFieldValue(input) {
        if (input.length < 1) return;
        var type = input.attr('type');
        switch (input[0].tagName.toLowerCase()) {
            case 'input':
                if (type === 'checkbox') {
                    if (input.attr('name').endsWith('[]')) {
                        var values = [];
                        input.filter(':checked').each(function () {
                            values.push($(this).val());
                        });
                        return values;
                    }
                    return input.filter(':checked:last').val();
                } else if(type === 'radio'){
                    return input.filter(':checked').val();
                }
                return input.val();
            case 'select':
                return input.val();
            default:
                return input.val();
        }
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
                setFormFieldValue(form.find('[name="' + name + '"]'), formData[name]);
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
                form.find('[name]').each(function () {
                    var name = $(this).attr('name');
                    var field = name;
                    if(name.startsWith(multilingualFieldPrefix+'[')){
                        if(name.startsWith(prefix)){
                            field = parseLangFieldName(name,prefix);
                        }else{
                            return;
                        }
                    }
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
                    var fieldName = 'Language['+langDefault+']['+name+']';
                    if(form.find('[name="'+fieldName+'"]').length>0) return fieldName;
                    return name;
                }
            }else{
                getModalFieldName = function(name){
                    return name;
                }
            }
            var getFieldValue = afterOpenCallback(button, modal, {'formData':formData, 'multilingual':multilingual, 'langDefault':langDefault, 'getModalFieldName':getModalFieldName, 'getFields':getFields});
            if(getFieldValue){
                var fields = getFields();
                for(var i in fields){
                    var val = getFieldValue(fields[i]);
                    if(val && typeof(val) == 'object') val = getFormFieldValue(val);
                    setFormFieldValue(form.find('[name="'+getModalFieldName(fields[i])+'"]'), val);
                }
            }
        }
        var submitBtn = modal.find('.modal-footer .btn-primary');
        submitBtn.off('click').on('click', function () {
            var form = modal.find('form'), data = {};
            form.find('[name]').each(function () {
                var name = $(this).attr('name');
                if (!name) return;
                switch (this.tagName.toLowerCase()) {
                    case 'input':
                        switch(this.type){
                            case 'checkbox':
                                if(!this.checked) return;
                                if(name.endsWith('[]')){
                                    if(!data[name]) data[name] = [];
                                    data[name].push($(this).val());
                                }else{
                                    data[name] = $(this).val();
                                }
                                return;
                            case 'radio':
                                if(!this.checked) return;
                                data[name] = $(this).val();
                                return;
                            case 'button':
                                return;
                            default:
                                data[name] = $(this).val();
                        }
                        break;
                    default:
                        data[name] = $(this).val();
                        break;
                }
            });
            var values = { data: data, multilingual: multilingual };
            if (multilingual) {
                if(!multilingualFieldPrefix) multilingualFieldPrefix = 'Language';
                if (!fields) {
                    var prefix = multilingualFieldPrefix + "[" + langDefault + "]";
                    form.find('[name]').each(function () {
                        var name = $(this).attr('name');
                        var field = name;
                        if(name.startsWith(multilingualFieldPrefix+'[')){
                            if(name.startsWith(prefix)){
                                field = parseLangFieldName(name,prefix);
                            }else{
                                return;
                            }
                        }
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
        var genFieldName;
        if(translatePrefix) {
            genFieldName = function(name){
                return translatePrefix+'['+name+']';
            }
        }else{
            genFieldName = function(name){
                return name;
            }
        }
        var fieldName = genFieldName('translate'),
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
                var $e = parentForDefaultLang.find('[name="'+genFieldName(name)+'"]');
                if($e.length==0) continue;
                var type = $e.attr('type');
                switch(type){
                    case 'radio':
                        $e.filter('[value="'+values[name]+'"]').prop('checked',true);
                        break;
                    case 'checkbox':
                        $e.prop('checked',false);
                        if(typeof(values[name])=='array'){
                            for(var i = 0; i < values[name].length; i++){
                                $e.filter('[value="'+values[name][i]+'"]').prop('checked',true);
                            }
                        }else{
                            $e.filter('[value="'+values[name]+'"]').prop('checked',true);
                        }
                        break;
                    default:
                        if($e[0].tagName.toLowerCase()=='select') {
                            $e.find('option[value="'+values[name]+'"]').prop('selected',true);
                        } else {
                            $e.val(values[name]);
                        }
                        break;
                }
            }
        }
    }
    App.initModalBody = initModalBody;
    App.initModalBodyPagination = initModalBodyPagination;
    App.initModalForm = initModalForm;
    App.updateMultilingualFormByModal = updateMultilingualFormByModal;
    App.getFormFieldValue = getFormFieldValue;
    App.setFormFieldValue = setFormFieldValue;
})();