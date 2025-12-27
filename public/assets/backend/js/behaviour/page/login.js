$(function(){
    $('[data-toggle="tooltip"]').tooltip();
    submitEncryptedData('#login-form');
    $('#login-form').find('input[name=remember]').on('change', function(){
        var checked = $(this).prop('checked');
        if(checked){
            setCookie('RememberBackendLogin', $(this).val(), 30, `/`);
        }else{
            setCookie('RememberBackendLogin', '', -1, `/`);
        }
    }).trigger('change');
});