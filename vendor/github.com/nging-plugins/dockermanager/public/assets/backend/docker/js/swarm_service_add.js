
$(function(){
  $('#endpointSpecPortMappingBtn').on('click',function(){
    $(this).closest('tr').before($('#endpointSpecPortMappingTmpl').html())
  });
  $('#taskTemplateContainerSpecMountsBtn').on('click',function(){
    $(this).closest('tr').before($('#taskTemplateContainerSpecMountsTmpl').html())
  });
  App.attachTurn('input[name="endpointSpec[mode]"]',{target:'#endpointSpecPorts'})
  App.attachTurn('input[name="storageVolumeMount"]',{target:'#storageVolumeMountRuleBox'})
  App.attachTurn('input[name="commandEnabled"]',{target:'#commandRuleBox'})
  App.editor.selectPage('#taskTemplateNetworkAttachmentSpecContainerID',{data:BACKEND_URL+'/docker/base/container/index?op=ajaxList&type=selectpage',eAjaxMethod:'GET'})
  App.editor.selectPage('#taskTemplateContainerSpecImage',{data:BACKEND_URL+'/docker/base/image/index?op=ajaxList&type=selectpage',eAjaxMethod:'GET'})
  App.editor.selectPage('#taskTemplateContainerSpecSecrets',{data:BACKEND_URL+'/docker/swarm/secret/index?op=ajaxList&type=selectpage',eAjaxMethod:'GET',multiple:true})
  App.editor.selectPage('#taskTemplateContainerSpecConfigs',{data:BACKEND_URL+'/docker/swarm/config/index?op=ajaxList&type=selectpage',eAjaxMethod:'GET',multiple:true})
  $('#service-mode').on('change',function(){
    var v=this.value;
    $('#mode-'+v).show().find('input').prop('disabled',false);
    $('#mode-'+v).siblings('.service-mode').hide();
    $('#mode-'+v).siblings('.service-mode').each(function(){
      $(this).find('input').prop('disabled',true)
    })
  }).trigger('change')
});