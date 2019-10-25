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
function DESEncryptECB(data, key, iv) {//加密
    var key = CryptoJS.enc.Utf8.parse(key);
    var encrypted = CryptoJS.DES.encrypt(data, key, {
        mode: CryptoJS.mode.ECB,
        padding: CryptoJS.pad.Pkcs7
    });
    return encrypted.toString(); //返回的是base64格式的密文
}
function DESDecryptECB(encrypted, key, iv) {//解密
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