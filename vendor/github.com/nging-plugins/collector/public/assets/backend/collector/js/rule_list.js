function recv_notice_collector(m){
	var percent=m.progress?m.progress.percent.toFixed(2)*1:0;
    if(m.status>0){
      console.info('['+percent+'/100] '+m.content);
    }else{
      console.error('['+percent+'/100] '+m.content);
	  App.message({title:App.i18n.SYS_INFO,text:m.content,
		time:5000,sticky:false,class_name:'error'});
    }
	if(!m.progress) return;
	var id=m.id;
	var el='collect-p-'+id;
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
	//bar.children('.sr-only').html('['+percent+'%] '+m.content);
	pb.find('.progress-description').text(' ['+percent+'%] '+m.content);
	if(percent>=100&&m.progress.complete){
		$('#collect-stop-'+id).addClass('hidden');
		$('#collect-start-'+id).removeClass('hidden');
		App.message({title:App.i18n.SYS_INFO,text:m.content,
		time:5000,sticky:false,class_name:'success'});
	}
}
function collecting(a,op){
	var id=$(a).data('id');
	if(op=='start'){
		$(a).addClass('hidden');
		$('#collect-stop-'+id).removeClass('hidden');
	}
	$.post(BACKEND_URL+'/collector/rule_collect',{id:id,op:op,clientID:App.clientID['notify']},function(r){
		if(r.Code<=0){
			op='stop';
		}
		if(op=='stop'){
			$('#collect-stop-'+id).addClass('hidden');
			$('#collect-start-'+id).removeClass('hidden');
			//$('#collect-p-'+id).next('.tr-progressbar').remove();
		}
		App.message({title:App.i18n.SYS_INFO,text:r.Info,time:5000,sticky:false,class_name:r.Code==1?'success':'error'});
	},'json');
}