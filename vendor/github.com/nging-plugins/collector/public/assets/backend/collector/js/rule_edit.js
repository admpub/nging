function addRule(a,idx){
  var content=$('#rule-tmpl').html();
  if(idx) content=content.replace(/ name\="rule\[/g,' name="extra[rule]['+idx+'][');
  content=content.replace(/\{\=idx\=\}/g,idx||'');
  $(a).parents('div.form-group').after(content);
}
function removeRule(a){
  var g=$(a).parents('div.form-group');
  var rule=g.find('[name$="[rule][]"]').text();
  var vname=g.find('[name$="[var][]"]').val();
  var filter=g.find('[name$="[filter][]"]').val();
  if(rule||vname||filter){
    if(!confirm('确定要删除规则吗？'))return;
  }
  g.remove();
}
function removePage(idx){
  if(!confirm('确定要删除页面规则吗？'))return;
  $('#extra-page-'+idx).remove();
}
function addPage(a){
  var idx=$('hr.extra-page').length, html = $('#page-tmpl').html().replace(/\{\=idx\=\}/g,idx);
  $(a).parents('div.form-group').before('<div class="extra-page-container" id="extra-page-'+idx+'"><a href="javascript:;" class="label label-danger extra-page-remove" onclick="removePage('+idx+');" data-toggle="tooltip" title="'+App.t("删除页面规则")+'"><i class="fa fa-times"></i></a>'+html+'</div>');
}
function verifyRule(a){
  var prefix='',v=$(a).val();
  var reset=function(){
    if(!$(a).hasClass('parsley-error')) return;
    $(a).removeClass('parsley-error');
    if($(a).next('ul.parsley-errors-list').length>0) $(a).next('ul.parsley-errors-list').remove();
  };
  if(!v||v.length<9)return reset();
  if(v.substring(0,7)=='regexp:'){
    prefix='regexp';
    v=v.substring(7);
  }else if(v.substring(0,8)=='regexp2:'){
    prefix='regexp2';
    v=v.substring(8);
  }else{
    return reset();
  } 
  var data={type:prefix,regexp:v};
  $.post(BACKEND_URL+'/collector/regexp_test',data,function(r){
    if(!r.Error) return reset();
    if(!$(a).hasClass('parsley-error')) $(a).addClass('parsley-error');
    if($(a).next('ul.parsley-errors-list').length<1) $(a).after('<ul class="parsley-errors-list"></ul>');
    $(a).next('ul.parsley-errors-list').html('<li class="required" style="display:list-item">'+App.t("正则表达式错误")+': '+r.Error+'</li>');
    $(a).focus();
  },'json');
}
function showRegexpTester(a){
  var modalBody=$('#regexp-test-modal').find('.modal-body'),prefix='';
  var v=$(a).val();
  if(!v||v.length<9)return;
  if(v.substring(0,7)=='regexp:'){
    prefix='regexp';
    v=v.substring(7);
  }else if(v.substring(0,8)=='regexp2:'){
    prefix='regexp2';
    v=v.substring(8);
  }else{
    return;
  } 
  var funcShow=function(){
    $('#regexp-test-modal').niftyModal('show',{
      afterOpen: function(modal) {
        $(window).trigger('resize');
        modalBody.find('[name="regexp"]').val(v);
        modalBody.find('[name="type"][value="'+prefix+'"]').prop('checked',true);
      },
      afterClose: function(modal) {
        $(a).val(modalBody.find('[name="type"]:checked').val()+':'+modalBody.find('[name="regexp"]').val());
      }
    });
  };
  if(!modalBody.data('ready')){
    $.get(BACKEND_URL+'/collector/regexp_test',{type:prefix},function(r){
      modalBody.html(r);
      modalBody.data('ready',true);
      funcShow();
    },'html');
  }else{
    funcShow();
  }
}
var filterInput=null;
function showFilterModal(a){
  filterInput=$(a).parent('.input-group-btn').siblings('input:first');
  $('#filter-modal').niftyModal('show');
}
$(function(){
  $('#filter-modal').find('.modal-body').html($('#filter-list-tmpl').html());
  $('#filter-modal').on('click','a[data-placeholder]',function(){
    var value=filterInput.val();
    if(value){
      filterInput.val(value+'|'+$(this).data('placeholder'));
    }else{
      filterInput.val($(this).data('placeholder'));
    }
    $('#filter-modal').niftyModal('hide');
    filterInput.focus();
  });
  $(window).on('resize',function(){
    $('#regexp-test-modal').css({height:$(window).height(),width:'100%','max-width':'100%',left:0,top:0,transform:'none','z-index':'9999'});
    $('#regexp-test-form').css({height:$(window).height()-180,overflow:'auto'});
  });
  $('#regexp-test-modal .modal-footer .btn-default.md-close:last').html('<i class="fa fa-times"></i> '+App.t("关闭"));
  var modalSubmitBtn=$('#regexp-test-modal .modal-footer .btn-primary:last');
  modalSubmitBtn.html('<i class="fa fa-stethoscope"></i> '+App.t("测试"));
  modalSubmitBtn.removeClass('md-close').click(function(){
        var $form=$('#regexp-test-form');
        App.loading('show');
        $.post(BACKEND_URL+'/collector/regexp_test',$form.serialize(),function(r){
            App.loading('hide');
            if(r.Error)return App.message({text: r.Error},false);
            if(!r.Result){
              $form.find('.result tbody').empty();
              return App.message({text:App.t("没有匹配到任何结果"),type:'warning'});
            }
            var tbody='';
            for(var i=0; i<r.Result.length; i++){
              var captures=r.Result[i];
              tbody+='<tr>';
              tbody+='<td>'+i+'</td>';
              tbody+='<td><table class="table table-bordered no-margin"><thead><tr><th style="width:50px">'+App.t("编号")+'</th><th>'+App.t("结果")+'</th></tr></thead><tbody>';
              for(var j=0; j<captures.length; j++){
                tbody+='<tr>';
                tbody+='<td>'+j+'</td>';
                tbody+='<td>'+App.htmlEncode(captures[j])+'</td>';
                tbody+='</tr>';
              }
              tbody+='</tbody></table></td>';
              tbody+='</tr>';
            }
            $form.find('.result tbody').html(tbody);
        },'json');
    });
    $('#collector-rule-form').on('submit',function(event){
      event.preventDefault();
      var $form=$(this);
      $.post($form.attr('action'),$form.serialize(),function(r){
        if(r.Code==1){
          App.message({text: r.Info, class_name: "success", after_open: function(e){window.location=BACKEND_URL+"/collector/rule"}});
          return;
        }
        if(r.Zone){
          $form.find('[name="'+r.Zone+'"]').focus();
        }
	      App.message({text: r.Info, class_name: "danger"});
      },'json');
    });
});