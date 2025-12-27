$(function(){
    $('[data-toggle="tooltip"]').tooltip();
    submitEncryptedData('#login-form');
    $('#login-form').find('input[name=remember]').on('click', function(){
        var checked = $(this).prop('checked');
        if(checked){
            setCookie('rememberBackendLogin', $(this).val(), 30, `/`);
        }else{
            setCookie('rememberBackendLogin', '', -1, `/`);
        }
    }).trigger('click');
});