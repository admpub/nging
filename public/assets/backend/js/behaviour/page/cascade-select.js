
(function (factory) {
  if (typeof define === 'function' && define.amd) {
      // AMD. Register as an anonymous module.
      define(['jquery'], factory);
  } else if (typeof exports === 'object') {
      // Node/CommonJS style for Browserify
      module.exports = factory;
  } else {
      // Browser globals
      factory(jQuery);
  }
}(function ($) {
function fetchOptions(elem,pid,selectedId,selectedIds,url,idKey,callback){
  if(!url) url=BACKEND_URL+'/tool/area/index';
  $.get(typeof url === 'function'?url():url,{pid:pid},function(r){
      if(r.Code!=1){
        return App.message({text:r.Info,type:'error'});
      }
      if(!r.Data || r.Data.listData.length<1){
        while($(elem).next('select').length>0) $(elem).next('select').remove();
        if(pid!=''&&pid!='0') $(elem).remove();
        else $(elem).html('<option value=""> - '+(App.i18n.PLEASE_SELECT?App.i18n.PLEASE_SELECT:'请选择')+' - </option>');
        return;
      }
      $(elem).show(h);
      var exclude = $(elem).data('exclude');
      var h = '<option value=""> - '+(App.i18n.PLEASE_SELECT?App.i18n.PLEASE_SELECT:'请选择')+' - </option>';
      for(var i=0;i<r.Data.listData.length;i++){
        var v=r.Data.listData[i];
        var s='';
        if(selectedId==v[idKey]) s+=' selected';
        if(exclude==v[idKey]) s+=' disabled';
        h+='<option value="'+v[idKey]+'"'+s+'>'+v.name+'</option>';
      }
      $(elem).html(h);
      callback && callback(pid,r);
      bindEvent(elem,selectedIds,url,idKey,callback).trigger('change');
  },'json');
}
function bindEvent(elem,selectedIds,url,idKey,callback){
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
        $(this).after('<select'+props+' style="display:none"></select>');
      }
      var selectedId='';
      if($.isArray(selectedIds) && selectedIds.length>index) selectedId=selectedIds[index];
      fetchOptions($(this).next('select')[0],v,selectedId,selectedIds,url,idKey,callback);
  });
}
function init(elem,selectedIds,url,idKey,callback){
  if(selectedIds == null){
    selectedIds = $(elem).data('selected') || [];
    if(selectedIds && !$.isArray(selectedIds)) selectedIds = String(selectedIds).split(",");
  }
  if(!url) url = $(elem).data('url');
  if(!idKey) idKey = $(elem).data('idKey') || 'id';
  var selected = selectedIds && selectedIds.length > 0 ? selectedIds[0] : '';
  var country = $(elem).data('country');
  var countryURL = $(elem).data('country-url');
  if(country || countryURL){
    if(!countryURL) countryURL = typeof url == 'function' ? url() : url;
    $.get(countryURL,{country:true,id:selected},function(r){
        if(r.Code!=1){
          return App.message({text:r.Info,type:'error'});
        }
        if(!r.Data || r.Data.listData.length<1){
          return;
        }
        var h = '<option value=""> - '+(App.i18n.PLEASE_SELECT?App.i18n.PLEASE_SELECT:'请选择')+' - </option>';
        var selectedAbbr = r.Data.countryAbbr||'';
        for(var i=0;i<r.Data.listData.length;i++){
          var v=r.Data.listData[i];
          var s='';
          if(selectedAbbr==v.abbr) s+=' selected';
          h+='<option value="'+v.abbr+'"'+s+'>'+v.name+'</option>';
        }
        var wrap=$(elem).prev('span.select-country-wrap');
        if(wrap.length<1){
          var $h=$('<span class="select-country-wrap"><select id="select-country">'+h+'</select></span>');
          $(elem).before($h);
          var $select=$h.children('select');
          $select.attr('class',$(elem).attr('class'));
          $select.on('change',function(){
            $(elem).trigger('reload')
          }).trigger('change');
        }
    },'json');
    var oldURL = typeof url == 'function' ? url() : url;
    var sep = oldURL.indexOf('?') > 0 ? '&' : '?';
    url = function(){
      return oldURL + sep + 'countryAbbr='+($('#select-country').val()||'');
    };
  }
  fetchOptions(elem,0,selected,selectedIds,url,idKey,callback);
  $(elem).on('reload',function(){
    fetchOptions(elem,0,selected,selectedIds,url,idKey,callback);
  });
}
window.CascadeSelect = {
  init: init
};
$.fn.cascadeSelect = function() {
  $(this).each(function(){
    init(this);
  });
};
return window.CascadeSelect;
}));