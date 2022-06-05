(function(){
    function initModalBody($md, ajaxData, onloadCallback, showFooterTabIndexies){
        var tabs = $md.find('.modal-body .nav-tabs>li');
        $md.find('.modal-footer').removeClass('hidden');
        if(showFooterTabIndexies == null) showFooterTabIndexies = [0];
        tabs.on('click', function(){
            var index = tabs.index(this), show = false;
            for(var i = 0; i < showFooterTabIndexies.length;i++){
                if(index == showFooterTabIndexies[i]){
                    show = true;
                    break;
                }
            }
            if(show){
                $md.find('.modal-footer').removeClass('hidden');
            }else{
                $md.find('.modal-footer:not(.hidden)').addClass('hidden');
            }
        });
        $md.find('table[data-page-size]').each(function(){
            initModalBodyPagination($(this), ajaxData, onloadCallback);
        });
    }
    function initModalBodyPagination(that, ajaxData, onloadCallback){
        if(!that || that.length<1) return;
        var data=that.data(),onSwitchPage=null;
        if (onloadCallback){
            if(data.ajaxList && data.ajaxList in onloadCallback) onloadCallback[data.ajaxList](that);
            if('switchPage' in onloadCallback) onSwitchPage=onloadCallback['switchPage'];
            if('pageData' in onloadCallback) data=typeof(onloadCallback['pageData'])=='function'?onloadCallback['pageData'](that):onloadCallback['pageData'];
        }
        App.withPagination(data,null);
        that.after(data.pagination);
        if (!onSwitchPage) {
            onSwitchPage=function(page){
                var params = {ajaxList:data.ajaxList,page:page};
                params = $.extend(params, ajaxData||{});
                $.get(window.location.href,params,function(res){
                    var container = that.parent();
                    container.html(res);
                    initModalBodyPagination(container.children('table[data-page-size]'), ajaxData, onloadCallback);
                },'html');
            };
        }
        that.next('ul.pagination').find('li > a[page]').on('click',function(){
            if($(this).closest('li.disabled').length > 0) return;
            onSwitchPage($(this).attr('page'));
        });
        that.next('ul.pagination').on('refresh',function(){
            var page=1;
            if($(this).data('page')){
                page=$(this).data('page');
                $(this).data('page',false);
            }else{
                page=$(this).find('li.active > a[page]').data('page')||1;
            }
            onSwitchPage(page);
        });
    }
    App.initModalBody = initModalBody;
    App.initModalBodyPagination = initModalBodyPagination;
})();