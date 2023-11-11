var keys=['name','filter'];
function selectDestType(a){
    if(!a.value) return $('.dest-value').addClass('hidden');
    var elem='#'+a.value+'Input';
    $(elem).removeClass('hidden').siblings('.dest-value').addClass('hidden');
}
function selectRuleList(pageId){
    $.post(window.location.href,{pageId:pageId,op:'ruleList'},function(r){
        if(r.Code!=1){
            return App.message({title: App.i18n.SYS_INFO, text: r.Info, class_name: 'error'});
        }
        var html='',tmpl=$('#tmpl-field-settings').html();
        for(var i=0;i<r.Data.length;i++){
            var v=r.Data[i],t=tmpl;
            for(var j=0;j<keys.length;j++){
                var k=keys[j],expr=new RegExp('{='+k+'=}','g');
                if(k=='filter'){
                  t=t.replace(expr,'');
                }else{
                  t=t.replace(expr,v[k]);
                }
            }
            html+=t;
        }
        $('#export-field-settings').html(html);
        attachSearchCollectField();
    },'json');
    return true;
}
function selectPageRule(a){
    var idv=$(a).attr('id');
    var pageId=$(a).val();
    if(idv=='pageRoot'){
        if(pageId<1){
            $('#pageId').html('<option value="0">'+App.t("请选择")+'</option>');
            return true;
        }
        $.post(window.location.href,{pageId:pageId,op:'childrenPageList'},function(r){
            if(r.Code!=1){
                return App.message({title: App.i18n.SYS_INFO, text: r.Info, class_name: 'error'});
            }
            var options='',addnum=2;
            for(var i=0;i<r.Data.length;i++){
                var v=r.Data[i];
                if(!v.name) v.name=App.t("%d级页面",i+addnum);
                options+='<option value="'+v.id+'">'+v.name+'</option>';
            }
            options='<option value="'+pageId+'">'+App.t("顶级页面")+'</option>'+options;
            $('#pageId').html(options);
            var subId=$('#pageId').val()*1;
            if(subId>0) pageId=subId;
            selectRuleList(pageId);
        },'json');
        return true;
    }
    if(!pageId) return $('#export-field-settings').empty();
    return selectRuleList(pageId);
}
function addRow(a) {
    var t=$('#tmpl-field-settings').html();
    for(var i=0;i<keys.length;i++){
        var k=keys[i],r=new RegExp('{='+k+'=}','g');
        t=t.replace(r,'');
    }
    if(a==null){
        $('#export-field-settings').prepend(t);
    }else{
        $(a).parent().parent().after(t);
    }
    attachSearchCollectField();
}
function delRow(a) {
    $(a).parent().parent().remove();
}
function moveUpRow(obj) {
    var currentTR=$(obj).parent().parent();
    if(currentTR.prev('tr').length>0){
        currentTR.prev('tr').before(currentTR.clone());
        currentTR.remove();
    }
}
function moveDownRow(obj) {
    var currentTR=$(obj).parent().parent();
    if(currentTR.next('tr').length>0){
        currentTR.next('tr').after(currentTR.clone());
        currentTR.remove();
    }
}
function attachSearchCollectField(){
    $('#export-field-settings').find('input[name="mapping[name][]"]').not('.tt-input').each(function(){
        searchCollectField(this);
    });
}
function searchCollectField(elem,size){
  if(size==null)size=500;
  $(elem).typeahead({
   hint: true, highlight: true, minLength: 1
  }, {source: function (query, sync, async) {
      var pageId=$('#pageId').val();
      $.ajax({
        url: window.location.href,
        type: 'post',
        data: {pageId:pageId,op:'ruleList'},
        dataType: 'json',
        async: false,
        success: function (data) {
          var arr = [];
          if(!data.Data) return;
          $.each(data.Data, function (index, val) {
            arr.push(val.name);
          });
          sync(arr);
        }
      });
  },limit: size});
}
var filterInput=null;
function showFilterModal(a){
  var $parent=$(a).parent('.input-group-btn');
  filterInput=$parent.siblings('input:first');
  $('#filter-modal').niftyModal('show');
}
$(function(){
  if($('#destType').val()) selectDestType($('#destType')[0]);
  attachSearchCollectField();
  $('#filter-modal').find('.modal-body').html($('#filter-list-tmpl').html());
  $('#filter-modal').on('click','a[data-placeholder]',function(){
    var value=filterInput.val(), newVal=$(this).data('placeholder');
    if(value) newVal=value+'|'+newVal;
    filterInput.val(newVal);
    $('#filter-modal').niftyModal('hide');
    filterInput.focus();
  });
});