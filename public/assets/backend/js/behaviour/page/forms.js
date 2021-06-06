(function(){
var initFlag = 'data-initialized';
function initForms(formElem, urlPrefix) {
   if(formElem == null) formElem = 'form';
   if(urlPrefix == null) urlPrefix = BACKEND_URL;
   $(formElem).find('.form-selectpage:not(['+initFlag+'=true])').each(function(){
    $(this).attr(initFlag,'true');
    if(!$(this).data('url')) $(this).data('url',urlPrefix+$(this).attr('url'));
    App.editor.selectPage(this);
  });
  $(formElem).find('select.form-cascade:not(['+initFlag+'=true])').each(function(){
    $(this).attr(initFlag,'true');
    if(!$(this).data('url')) $(this).data('url',urlPrefix+$(this).attr('url'));
    App.editor.cascadeSelect(this);
  });
}
window.initForms = initForms;
})();