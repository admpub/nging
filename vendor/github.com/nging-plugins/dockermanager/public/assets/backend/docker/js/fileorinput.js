$(function(){
  $('#contentFrom-input,#contentFrom-file').on('click',function(){
    var input=$('#contentInput'),inputWrap=input.siblings('.CodeMirror');
    if(inputWrap.length<0) inputWrap=input;
    switch(this.value){
      case 'input': inputWrap.show(); input.trigger('initcodemirror'); $('#contentFile').prop('required',false).hide(); break;
      case 'file': inputWrap.hide(); $('#contentFile').prop('required',true).show(); break;
    }
  });
  $("#contentInput").on('initcodemirror',function(){
    if($(this).siblings('.CodeMirror').length>0) return;
    App.editor.codemirror("#contentInput", {theme:'ambiance'});
  });
  $("#contentFrom-input:checked,#contentFrom-file:checked").trigger('click');
})