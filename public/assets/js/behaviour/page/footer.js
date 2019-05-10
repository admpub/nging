
$(function(){
	App.init();
	App.markNavByURL(window.activeURL);
	App.attachPjax(null,{
		 onclick: function(obj){
			//console.log($(obj).data('marknav'))
			if($(obj).data('marknav')){
				App.unmarkNav($(obj),$(obj).data('marknav'));
				App.markNav($(obj),$(obj).data('marknav'));
			}
		 },
		 onend: function(evt,xhr,opt){
			 opt.container.find('[data-popover="popover"]').popover();
			 opt.container.find('.ttip, [data-toggle="tooltip"]').tooltip();
		 }
	});
	App.attachAjaxURL();
	App.bottomFloat('.pagination');
	App.bottomFloat('.form-submit-group',0,true);
    if(window.errorMSG) App.message({title: App.i18n.SYS_INFO, text: window.errorMSG, class_name: "danger"});
	if(window.successMSG) App.message({title: App.i18n.SYS_INFO, text: window.successMSG, class_name: "success"});
	App.notifyListen();
});