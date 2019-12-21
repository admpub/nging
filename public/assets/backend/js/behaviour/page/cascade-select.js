
(function(win){
function fetchOptions(elem,pid,selectedId,selectedIds,url){
  if(!url) url=BACKEND_URL+'/tool/area/index';
  $.get(url,{pid:pid},function(r){
      if(r.Code!=1){
        return App.message({text:r.Info,type:'error'});
      }
      if(r.Data.listData.length<1){
        while($(elem).next('select').length>0) $(elem).next('select').remove();
        $(elem).remove();
        return;
      }
      var exclude = $(elem).data('exclude');
      var h = '<option value=""> - '+(App.i18n.PLEASE_SELECT?App.i18n.PLEASE_SELECT:'请选择')+' - </option>';
      for(var i=0;i<r.Data.listData.length;i++){
        var v=r.Data.listData[i];
        var s=selectedId==v.id?' selected':'';
        if(exclude==v.id) s+=' disabled';
        h+='<option value="'+v.id+'"'+s+'>'+v.name+'</option>';
      }
      $(elem).html(h);
      bindEvent(elem,selectedIds,url).trigger('change');
  },'json');
}
function bindEvent(elem,selectedIds,url){
  if($(elem).data('bindEvent')){
    return $(elem);
  }
  $(elem).data('bindEvent',true);
  return $(elem).on('change',function(){
      var v=$(this).val();
      var n=$(this).attr('name');
      var c=$(this).attr('class');
      var p=Number($(this).attr('pos')||0);
      var target=$(this).data('target');
      var exclude=$(this).data('exclude');
      if(v==''||v=='0'){
        if(p==0 && target) $(target).val(v);
        while($(elem).next('select').length>0) $(elem).next('select').remove();
        return;
      }
      if(target) $(target).val(v);
      var index=p+1;
      if($(this).next('select').length<1){
        var props = ' pos="'+index+'"';
        if(n) props+=' name="'+n+'"';
        if(c) props+=' class="'+c+'"';
        if(target) props+=' data-target="'+target+'"';
        if(exclude) props+=' data-exclude="'+exclude+'"';
        $(this).after('<select'+props+'></select>');
      }
      var selectedId='';
      if($.isArray(selectedIds) && selectedIds.length>index) selectedId=selectedIds[index];
      fetchOptions($(this).next('select')[0],v,selectedId,selectedIds,url);
  });
}
win.CascadeSelect = {
  init: function(elem,selectedIds,url){
    fetchOptions(elem,0,selectedIds?selectedIds[0]:'',selectedIds,url);
  }
};
})(window);