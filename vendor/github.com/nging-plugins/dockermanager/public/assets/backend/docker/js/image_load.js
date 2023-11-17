function recv_notice_dockerImageLoad(m){
	var percent=m.progress?m.progress.percent.toFixed(2)*1:0;
    if(m.status>0){
        console.info('['+percent+'/100] '+m.content);
    }else{
        console.error('['+percent+'/100] '+m.content);
    }
	if(!m.progress) return;
	showProgress(m.id,percent,m.content);
	if(percent>=100&&m.progress.complete){
		App.message({title:App.i18n.SYS_INFO,text:m.content,
		time:5000,sticky:false,class_name:'success'});
		window.setTimeout(function(){
			var el='docker-image-'+m.id,tr=$('#'+el);
			tr.next('tr.tr-progressbar').slideUp('slow',function(){$(this).remove()});
		},5000);
	}
}
function showProgress(id,percent,content, isHTML){
	var el='docker-image-'+id;
	var tr=$('#'+el);
	if(tr.length<1) return;
	var pb=tr.next('.tr-progressbar');
	if(pb.length<1){
		var html=$('#tr-progressbar').html();
		$('#'+el).after(html);
		pb=$('#'+el).next('.tr-progressbar');
	}
	var bar=pb.children('td').children('.progress').children('.progress-bar');
	bar.css('width',percent+'%').attr('aria-valuenow',percent);
	//bar.children('.sr-only').html('['+percent+'%] '+content);
	if(isHTML){
		pb.find('.progress-description').html(content);
	}else{
		pb.find('.progress-description').text(content);
	}
}
function imageLoad(data){
	var $form=$('#formImageLoad');
	var $a=$form.find('button:submit');
	var id=$a.data('id')||'0';
	var fa=$a.children('.fa');
	var cb=function(show){
		if(show){
			fa.addClass('fa-spin');
			$a.prop('disabled',true);
		}else{
			fa.removeClass('fa-spin');
			$a.prop('disabled',false);
		}
	}
	cb(true);
	showProgress(id,0,'<i class="fa fa-refresh fa-spin"></i> '+App.t('开始处理，请稍候...'),true);
	var postData={id:id,clientID:App.clientID['notify']};
	postData=$.extend(postData,data||{});
	App.postFormData($form[0],postData,function (r) {
		cb(false);
		App.message({title:App.i18n.SYS_INFO,text:r.Info,time:5000,sticky:false,class_name:r.Code==1?'success':'error'});
		if(r.Code==0) showProgress(id,0,'<i class="fa fa-times-circle text-danger"></i> '+App.htmlEncode(r.Info),true);
		else showProgress(id,100,'<i class="fa fa-check"></i> '+App.htmlEncode(r.Info),true);
    },function (xhr, status, info) {
		cb(false);
	})
}
$(function(){
$('#formImageLoad').on('submit',function(e){
    e.preventDefault();
    imageLoad({});
})
})