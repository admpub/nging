$(function(){
function c(elem,show){
  if(show){
    $(elem).show();
    $(elem).children('input').prop('required',true);
  }else{
    $(elem).hide();
	$(elem).val('');
    $(elem).children('input').prop('required',false);
  }
}
$('#formContainerFileImport').find('input[name="contentFrom"]').on('click',function(){
    switch(this.value){
    case 'input': 
		c('#contentInput',true);
		c('#contentFile',false);
		break;
    case 'file': 
		c('#contentInput',false);
		c('#contentFile',true);
		break;
    }
});
$('#formImageImport').find('input[name="contentFrom"]:checked').trigger('click')
})