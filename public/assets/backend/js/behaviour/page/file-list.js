function ownerTypeChange(a){}
function applySelected(){
	if(client=='xheditor'&&window.callback){
		window.callback('!'+getSelectedFiles().join(' '));
		return false;
	}
	if(callback){
		if(typeof(target[callback])=='function'){
			target[callback](getSelectedFiles());
		}
	}else{
		if(insertTo){
			$(insertTo).val(getSelectedFiles().join(','));
			return false;
		}
        App.message({title:App.i18n.SYS_INFO,text:App.i18n.NO_CALLBACK_NAME,type:'error'});
        return false;
	}
    return false;
}
function getSelectedFiles(){
	var files=[];
	$("input.check-table:checked").each(function(){
		files.push($(this).data('file-url'));
	});
	if(files.length<1){
		App.message({title:App.i18n.SYS_INFO,text:App.i18n.PLEASE_SELECT,type:'error'});
		return false;
	}
	return files;
}
$(function(){
	$('#timerange').on('focus',function(){
		if($(this).data('attached')) return false;
		$(this).data('attached',true);
		App.daterangepicker('#timerange',{
			showShortcuts: true,
			shortcuts: {
				'prev-days': [1,3,5,7],
				'next-days': [3,5,7],
				'prev' : ['week','month'],
				'next' : ['week','month']
			}
		});
	});
	function submitSearch(e){
		e.preventDefault();
		var data=$('#search-form').serializeArray();
		data.push({name:'partial',value:1});
		loadList($('#search-form').attr('action'),data);
	}
	if(dialogMode) {
		$('#search-form').on('submit',submitSearch);
		$('#timerange,#type,#table,#ownerId,#used').on('change', submitSearch);
	}else{
		$('#timerange,#type,#table,#ownerId,#used').on('change', function(){
			$('#search-form').submit();
		});
	}
	var myUploadInput;
	function initUploadButton(){
		if($("#input-file-upload").parent('.file-preview-shadow').length<1)
		myUploadInput = $("#input-file-upload").uploadPreviewer({
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

		$('#checkedAll,input[type=checkbox][name="id[]"]:checked').prop('checked',false);
		App.attachCheckedAll('#checkedAll','input[type=checkbox][name="id[]"]');
	}
	function loadList(url,args){
		$.get(url,args,function(r){
			$('#table-container').html(r);
			initTable();
			$('#table-container .pagination a').on('click',function(e){
				e.preventDefault();
				var url=$(this).attr('href');
				loadList(url,{});
			});
		},'html');
	}
	function initTable(){
		$('#table-container thead').data('sort-trigger',function(){
			var thead=$('#table-container thead');
			var url=thead.data('sort-url');
			loadList(url,{partial:1});
		});
		App.tableSorting('#table-container');
		App.float('#tbody-content img.previewable');
	}
	initUploadButton();
  	$(document).on("file-preview:changed", function(e) {
		$.LoadingOverlay("show", {
    		image : ASSETS_URL+"/images/nging-gear.png",//progress : true, image: "",//fontawesome : "fa fa-cog fa-spin",
			text  : App.i18n.UPLOADING
        });
		myUploadInput.submit(function(r){
			$.LoadingOverlay("hide");
			if(r.Code==1){
				App.message({text:App.i18n.UPLOAD_SUCCEED,type:'success'});
				loadList(listURL,{partial:1});
			}else{
				App.message({text:r.Info,type:'error'});
			}
		});
	});
	initTable();
});