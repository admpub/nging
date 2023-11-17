function attachSearchNetworkModeField(){
    $('#networkMode').not('.tt-input').each(function(){
        searchNetworkModeField(this);
    });
}
function searchNetworkModeField(elem,size){
  if(size==null)size=500;
  $(elem).typeahead({
   hint: true, highlight: true, minLength: 0
  }, {source: function (query, sync, async) {
      $.ajax({
        url: BACKEND_URL+'/docker/base/network/index',
        type: 'get',
        data: {op:'ajaxList'},
        dataType: 'json',
        async: false,
        success: function (r) {
          if(r.Code!=1) return App.message({text:r.Info,type:'error'});
          if(!r.Data) return;
          sync(r.Data.listData);
        }
      });
  },limit: size});
}
function attachSearchVolumesFromField(){
  var $form=$('#container-add-form');
  $form.find('input.volumesFrom').not('.tt-input').each(function(){
    searchVolumesFromField(this);
  });
  $form.on('click','.btn-volumes-from-add',function(){
    $(this).closest('.input-group').next('.input-group').find('input.volumesFrom').not('.tt-input').each(function(){
      searchVolumesFromField(this);
    });
  })
}
function searchVolumesFromField(elem,size){
if(size==null)size=500;
$(elem).typeahead({
 hint: true, highlight: true, minLength: 0
}, {source: function (query, sync, async) {
    $.ajax({
      url: BACKEND_URL+'/docker/base/container/index',
      type: 'get',
      data: {op:'ajaxList'},
      dataType: 'json',
      async: false,
      success: function (r) {
        if(r.Code!=1) return App.message({text:r.Info,type:'error'});
        if(!r.Data) return;
        sync(r.Data.listData);
      }
    });
},limit: size});
}
$(function(){
  attachSearchNetworkModeField();attachSearchVolumesFromField();
  App.attachTurn('input[name="commandEnabled"]',{target:'#commandRuleBox'})
  App.attachTurn('input[name="portExport"]',{target:'#portExportRuleBox'})
  App.attachTurn('input[name="storageVolumeMount"]',{target:'#storageVolumeMountRuleBox'})
  $('#containerPortMappingBtn').on('click',function(){
    $(this).closest('tr').before($('#portMappingTmpl').html())
  })
  $('#containerPathMappingBtn').on('click',function(){
    $(this).closest('tr').before($('#pathMappingTmpl').html())
  })
  App.editor.selectPage('#dockerImage',{data:BACKEND_URL+'/docker/base/image/index?op=ajaxList&type=selectpage',eAjaxMethod:'GET'})
  $('#restartPolicy').on('change',function(){
    var v=this.value,$form=$('#container-add-form'),className='.restartPolicy-show-'+v;
    $form.find(className).show();
    $form.find('[class^="restartPolicy-show-"]:not('+className+')').hide();
    if(v!='on-failure')$('#restartMaxRetryCount').val('0');
  }).trigger('change')
});