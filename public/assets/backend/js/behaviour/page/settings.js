
$(function(){
  App.editor.xheditors('#pcont .html-editor');
  App.editor.tinymces('#pcont .html-editor-tinymce');
  App.editor.markdowns('#pcont .html-editor-markdown');
  App.editor.fileInput();
  $("#body-left-navigate.nscroller:not(.has-scrollbar)").nanoScroller();
  function getTmplTag(name){
    if(!/\[value\]$/.test(name)){
      name=name.replace(/\]\[value\]\[/,'.');
    }else{
      name=name.replace(/\]\[value\]$/,'');
    }
    name=name.replace(/\]\[/g,'.');
    name=name.replace(/\]|\[/g,'.');
    name=name.replace(/\.$/g,'');
    return '{'+'{Config.'+name+'}'+'}';
  }
  if(focusInputName&&focusInputName.indexOf('"')==-1 && $('[name="'+focusInputName+'"]').length>0) {
    $('[name="'+focusInputName+'"]').focus();
    var steps = [{element:'[name="'+focusInputName+'"]',intro:App.t("设置它")}];
    introJs().setOptions({showButtons:false,showBullets:false,showStepNumbers:false,steps:steps}).start();
  }
  $('#pcont .cl-mcont form [name*="[value]"]').each(function(){
    if($(this).parent('.input-group').length>0){
      var e = $(this).parent('.input-group');
      if(e.prev('.input-group-addon').length>0){
        e.prev('.input-group-addon').attr('data-toggle','tooltip').attr('title',getTmplTag($(this).attr('name'))).tooltip();
        return;
      }
    }
    if($(this).closest('div[class*="col-"]').length>0){
      var e = $(this).closest('div[class*="col-"]');
      if(e.prev('.control-label').length>0){
        e.prev('.control-label').attr('data-toggle','tooltip').attr('title',getTmplTag($(this).attr('name'))).tooltip();
        return;
      }
    }
    if($(this).closest('.form-group').length>0){
      var e = $(this).closest('.form-group');
      if(e.children('.control-label').length>0){
        e.children('.control-label').attr('data-toggle','tooltip').attr('title',getTmplTag($(this).attr('name'))).tooltip();
        return;
      }
    }
  });
});