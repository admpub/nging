function AESEncrypt(data, key, iv) {//加密
    var key = CryptoJS.enc.Utf8.parse(key);
    var iv = CryptoJS.enc.Utf8.parse(iv);
    var encrypted = CryptoJS.AES.encrypt(data, key, {
        iv: iv,
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7
    });
    return encrypted.toString(); //返回的是base64格式的密文
}
function AESDecrypt(encrypted, key, iv) {//解密
    var key = CryptoJS.enc.Utf8.parse(key);
    var iv = CryptoJS.enc.Utf8.parse(iv);
    var decrypted = CryptoJS.AES.decrypt(encrypted, key, {
        iv: iv,
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7
    });
    return decrypted.toString(CryptoJS.enc.Utf8);
}
function DESEncrypt(data, key, iv) {//加密
    var key = CryptoJS.enc.Utf8.parse(key);
    var iv = CryptoJS.enc.Utf8.parse(iv);
    var encrypted = CryptoJS.DES.encrypt(data, key, {
        iv: iv,
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7
    });
    return encrypted.toString(); //返回的是base64格式的密文
}
function DESDecrypt(encrypted, key, iv) {//解密
    var key = CryptoJS.enc.Utf8.parse(key);
    var iv = CryptoJS.enc.Utf8.parse(iv);
    var decrypted = CryptoJS.DES.decrypt(encrypted, key, {
        iv: iv,
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7
    });
    return decrypted.toString(CryptoJS.enc.Utf8);
}
function DESEncryptECB(data, key) {//加密
    var key = CryptoJS.enc.Utf8.parse(key);
    var encrypted = CryptoJS.DES.encrypt(data, key, {
        mode: CryptoJS.mode.ECB,
        padding: CryptoJS.pad.Pkcs7
    });
    return encrypted.toString(); //返回的是base64格式的密文
}
function DESDecryptECB(encrypted, key) {//解密
    var key = CryptoJS.enc.Utf8.parse(key);
    var decrypted = CryptoJS.DES.decrypt(encrypted, key, {
        mode: CryptoJS.mode.ECB,
        padding: CryptoJS.pad.Pkcs7
    });
    return decrypted.toString(CryptoJS.enc.Utf8);
}
function SHA256(data){
    var encrypted = CryptoJS.SHA256(data).toString(CryptoJS.enc.Hex);
    return encrypted;
}
function MD5(data){
    var encrypted = CryptoJS.MD5(data).toString(CryptoJS.enc.Hex);
    return encrypted;
}
function SHA1(data){
    var encrypted = CryptoJS.SHA1(data).toString(CryptoJS.enc.Hex);
    return encrypted;
}
function EncryptFormPassword(formElem){
	var data=$(formElem).serializeArray();
	var pwdn=$(formElem).find('input[type="password"]').length;
    var pwdi=0,secret=$(formElem).data('secret');
    if(!secret) return data;
	for(var i=0;i<data.length;i++){
		var v=data[i];
		if($(formElem).find('input[name="'+v.name+'"][type="password"]').length>0){
			data[i].value=DESEncryptECB(v.value,secret);
			pwdi++;
			if(pwdi>=pwdn) break;
		}
    }
    return data;
}