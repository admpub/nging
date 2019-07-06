
$(function(){
	App.init();
	App.initLeftNavAjax(window.activeURL,null);
	App.bottomFloat('.pagination');
	App.bottomFloat('.form-submit-group',0,true);
    if(window.errorMSG) App.message({title: App.i18n.SYS_INFO, text: window.errorMSG, class_name: "danger"});
	if(window.successMSG) App.message({title: App.i18n.SYS_INFO, text: window.successMSG, class_name: "success"});
	App.notifyListen();
	$('#topnav a[data-project]').on('click',function(e){
		e.preventDefault();
		var ident=$(this).data('project');
		if(ident==$('#leftnav').data('project')) return;
		$('#leftnav').data('project',ident);
		var li=$(this).parent('li');
		$.get(window.BACKEND_URL+'/project/'+ident,{partial:1},function(r){
			if(r.Code!=1){
				App.message({title:App.i18n.SYS_INFO,text:r.Info,type:'error'});
				return;
			}
			$('#leftnav').html(r.Data.list);
			App.initLeftNav();
			App.initLeftNavAjax(window.activeURL,'#leftnav');
			li.siblings('li.active').removeClass('active');
			li.addClass('active');
		},'json');
	});
});