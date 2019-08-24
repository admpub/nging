App.select2 = {
    i18n: {
        TAG_INPUT: '请输入或选择标签，如有多个用空格隔开',
        TAG_SELECT: '请输入关键词后从搜索列表中选择',
        SELECT: '请选择'
    },
    tags: function (element, tagsArray, ajax, sortable, onlySelect) {
        //tagsArray:[{id:1,text:'coscms',locked:true}] locked元素不是必须的，如果为true代表不可删除
        if (tagsArray == null) tagsArray = $(element).data('tags') || [];
        if (ajax == null) ajax = $(element).data('url');
        if (sortable == null) sortable = $(element).data('sortable') || false;
        if (onlySelect == null) onlySelect = $(element).data('onlyselect') || false;
        var single = $(element).data('single') || false;
        if (single) single = App.parseBool(single);
        var options = { multiple: !single, width: '100%' };
        if (onlySelect) {//仅仅可选择，不可新增选项
            options.minimumInputLength = 1
            options.placeholder = App.select2.i18n.TAG_SELECT;
            options.data = tagsArray;
            options.tags = true;
        } else {//支持新增选项(注意：ajax方式获取数据时，将不支持新增选项)
            options.placeholder = App.select2.i18n.TAG_INPUT;
            options.tags = tagsArray;
            options.tokenSeparators = [',', ' '];
        }
        if (ajax) {
            switch (typeof (ajax)) {
                case 'string':
                    var listKey = $(element).data('listkey') || 'list';
                    options.query = App.select2.buildQueryFunction(ajax, {}, listKey);
                    break;

                default:
                    if ($.isArray(ajax)) {
                        var listKey = ajax.length > 2 ? ajax[2] : ($(element).data('listkey') || 'list');
                        options.query = App.select2.buildQueryFunction(ajax[0], ajax.length>1?ajax[1]:{}, listKey);
                        break;
                    }
                    options.ajax = ajax;
                    break;
            }
        }
        var sel = $(element).select2(options);
        $(element).data('select2', sel);
        var initSelected = $(element).data('init');
        if (initSelected) $(element).val(initSelected.split(',')).trigger('change');
        if (!sortable) return;

        //拖动排序
        $(element).on('change', function () {
            var valElement;
            if (typeof (element) == 'string' && element.indexOf('#') === 0) {
                valElement = element + '_val';
            } else {
                valElement = $(element).attr('id') + '_val';
            }
            $(valElement).html($(element).val());
        });
        $(element).select2('container').find('ul.select2-choices').sortable({
            containment: 'parent',
            start: function () { $(element).select2('onSortStart'); },
            update: function () { $(element).select2('onSortEnd'); }
        });
    },
    select: function (element, options) {
        var defaults = {
            placeholder: App.select2.i18n.SELECT, width: '100%',
            minimumInputLength: 3,
            /*
            ajax: {},
            initSelection: function(element, callback) {
                var id = $(element).val();
                if (id !== "") {
                    $.ajax("https://api.github.com/repositories/" + id, {
                        dataType: "json"
                    }).done(function(data) { callback(data); });
                }
            },
            formatResult: function(row){return row.text}, // 格式化显示每一行数据
            formatSelection: function(row){return row.text}, // 格式化选中项在下拉框上显示的内容
            dropdownCssClass: "bigdrop", // 下拉列表样式
            escapeMarkup: function (m) { return m; } //如果要显示html内容，则不需要进行escape处理
            */
        };
        options = $.extend({}, defaults, options || {});
        var sel = $(element).select2(options);
        $(element).data('select2', sel);
    },
    update: function (element, preloadData) {
        $(element).select2('data', preloadData);
    },
    enable: function (element, enable) {
        $(element).prop('disabled', !enable);
    },
    readonly: function (element, readonly) {
        $(element).prop('disabled', readonly);
    },
    buildQueryFunction: function (url, params, listKey) {
        if (listKey == null) listKey = 'list';
        return function (query) {
            params.q = query.term;
            params.select2 = 1;
            $.get(url, params, function (r) {
                if (r.Code != 1) {
                    App.message({
                        title: App.i18n.SYS_INFO,
                        text: r.Info,
                        type: 'error'
                    }, false);
                    return;
                }
                //r.Data:[{id:1,text:'coscms',locked:true,disabled:true}] 
                //locked元素不是必须的，如果为true代表不可删除(用于tags)
                //disabled元素不是必须的，如果为true代表不可选择(用于select)
                var data = { results: r.Data[listKey] };
                query.callback(data);
            }, 'json');
        };
    },
    buildAjaxOptions: function (options, params, listKey) {
        if (listKey == null) listKey = 'list';
        var defaults = {
            url: "",
            dataType: 'json',
            quietMillis: 250,
            data: function (term, page) { // 基于页码构建查询数据
                return $.extend({}, {
                    q: term, //搜索词
                    page: page, //页码
                    select2: 1,
                }, params || {});
            },
            results: function (res, page) { // 重新组织ajax响应数据后供selec2使用
                if (res.Code != 1) {
                    App.message({
                        title: App.i18n.SYS_INFO,
                        text: res.Info,
                        type: 'error'
                    }, false);
                    return { results: [], more: false };
                }
                var data = res.Data;
                var more = page < data.pagination.pages; // more变量用于通知select2可以加载更多数据
                return {
                    results: data[listKey],
                    more: more
                };
            }
        }
        return $.extend({}, defaults, options || {});
    }
};