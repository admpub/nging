function ownerTypeChange(a){}
function applySelected(){
	if(!callback){
        App.message({title:App.i18n.SYS_INFO,text:App.i18n.NO_CALLBACK_NAME,type:'error'});
        return false;
    }
	if(typeof(target[callback])=='function'){
		var files=[];
		$("input.check-table:checked").each(function(){
			files.push($(this).data('file-url'));
        });
        if(files.length<1){
            App.message({title:App.i18n.SYS_INFO,text:App.i18n.PLEASE_SELECT,type:'error'});
            return false;
        }
		target[callback](files);
    }
    return false;
}
$(function(){
	App.daterangepicker('#timerange',{
		showShortcuts: true,
		shortcuts: {
			'prev-days': [1,3,5,7],
			'next-days': [3,5,7],
			'prev' : ['week','month'],
			'next' : ['week','month']
		}
	});
	$('#timerange').on('change',function(event,obj){
		$('#search-form').submit();
	});
	App.float('#tbody-content img.previewable');
	var myUploadInput = $("#input-file-upload").uploadPreviewer({
		"buttonText":'<i class="fa fa-cloud-upload"></i> '+App.i18n.BUTTON_UPLOAD,
		"previewTableContainer":'#previewTableContainer',
		"url":uploadURL,
		"previewTableShow":false/*,
		"uploadProgress":function(progress){
			var count=progress*100;
			if(count>100){
				$.LoadingOverlay("hide");
				return;
			}
			$.LoadingOverlay("progress", count);
		}*/
	});
  	$(document).on("file-preview:changed", function(e) {
		$.LoadingOverlay("show", {
    		image : ASSETS_URL+"/images/nging-gear.png",//progress : true, image: "",//fontawesome : "fa fa-cog fa-spin",
			text  : App.i18n.UPLOADING
        });
		myUploadInput.submit(function(r){
			$.LoadingOverlay("hide");
			if(r.Code==1){
				App.message({text:App.i18n.UPLOAD_SUCCEED,type:'success'});
				window.setTimeout(function(){window.location=listURL;},2000);
			}else{
				App.message({text:r.Info,type:'error'});
			}
		});
	});
	App.tableSorting();
	$('#checkedAll,input[type=checkbox][name="id[]"]:checked').prop('checked',false);
	App.attachCheckedAll('#checkedAll','input[type=checkbox][name="id[]"]');
});