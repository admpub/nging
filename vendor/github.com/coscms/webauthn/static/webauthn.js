(function () {
  function check() {
    // check whether current browser supports WebAuthn
    if (!window.PublicKeyCredential) {
      alert("Error: this browser does not support WebAuthn");
      return false;
    }
    return true
  }

  function isSupported() {
    return typeof(window.PublicKeyCredential)!='undefined';
  }

  // Base64 to ArrayBuffer
  function bufferDecode(value) {
    try {
      value = String(value).replace(/_/g,'/').replace(/-/g,'+');
      return Uint8Array.from(atob(value), c => c.charCodeAt(0));
    } catch (error) {
      console.error(error+": "+value);
    }
  }

  // ArrayBuffer to URLBase64
  function bufferEncode(value) {
    return btoa(String.fromCharCode.apply(null, new Uint8Array(value)))
      .replace(/\+/g, "-")
      .replace(/\//g, "_")
      .replace(/=/g, "");
  }

  function webAuthn(options) {
    var $this = this;
    this.options = {
      urlPrefix: '/webauthn',
      debug: false,
      getRegisterData: function () { return {} },
      getLoginData: function () { return {} },
      getUnbindData: function () { return {} },
      onRegisterSuccess: function (response) {$this.options.debug && console.log(response)},
      onRegisterError: function (error) {$this.options.debug && console.error(error)},
      onLoginSuccess: function (response) {$this.options.debug && console.log(response)},
      onLoginError: function (error) {$this.options.debug && console.error(error)},
      onUnbindSuccess: function (response) {$this.options.debug && console.log(response)},
      onUnbindError: function (error) {$this.options.debug && console.error(error)},
      checkResponseBeginLogin: function(data) {return data},
      checkResponseFinishLogin: function(data) {return data},
      checkResponseBeginRegister: function(data) {return data},
      checkResponseFinishRegister: function(data) {return data},
      checkResponseBeginUnbind: function(data) {return data},
      checkResponseFinishUnbind: function(data) {return data},
    }
    $.extend(this.options, options || {});
  }

  webAuthn.prototype.check = check;
  webAuthn.prototype.isSupported = isSupported();
  webAuthn.prototype.register = function (username) {
    if (username === "") {
      alert("Please enter a username");
      return;
    }
    var $this = this;

    $.post(
      $this.options.urlPrefix + '/register/begin/' + username,
      $this.options.getRegisterData(),
      function (data) {
        data = $this.options.checkResponseBeginRegister(data);
        return data;
      },'json')
      .then((credentialCreationOptions) => {
        $this.options.debug && console.log(credentialCreationOptions);

        if(typeof credentialCreationOptions.Code != 'undefined') throw new Error(credentialCreationOptions.Info);
        
        credentialCreationOptions.publicKey.challenge = bufferDecode(credentialCreationOptions.publicKey.challenge);
        credentialCreationOptions.publicKey.user.id = bufferDecode(credentialCreationOptions.publicKey.user.id);
        if (credentialCreationOptions.publicKey.excludeCredentials) {
          for (var i = 0; i < credentialCreationOptions.publicKey.excludeCredentials.length; i++) {
            credentialCreationOptions.publicKey.excludeCredentials[i].id = bufferDecode(credentialCreationOptions.publicKey.excludeCredentials[i].id);
          }
        }

        return navigator.credentials.create({
          publicKey: credentialCreationOptions.publicKey
        })
      })
      .then((credential) => {
        $this.options.debug && console.log(credential);
        let attestationObject = credential.response.attestationObject;
        let clientDataJSON = credential.response.clientDataJSON;
        let rawId = credential.rawId;

        $.ajax({
          type:'post',
          url:$this.options.urlPrefix + '/register/finish/' + username,
          dataType:'json',
          async:false,
          data:JSON.stringify({
            id: credential.id,
            rawId: bufferEncode(rawId),
            type: credential.type,
            response: {
              attestationObject: bufferEncode(attestationObject),
              clientDataJSON: bufferEncode(clientDataJSON),
            },
          }),
          success:function (data) {
            data = $this.options.checkResponseFinishRegister(data);
            return data;
          }});
      })
      .then((response) => {
        $this.options.debug && alert("successfully registered " + username + "!");
        $this.options.onRegisterSuccess.call(this, response);
      })
      .catch((error) => {
        console.log("failed to register " + username + ": " + error);
        $this.options.onRegisterError.call(this, error);
      })
  }

  webAuthn.prototype.auth = function (username, type) {
    if (username === "") {
      alert("Please enter a username");
      return;
    }
    var $this = this;

    $.post(
      $this.options.urlPrefix + '/'+type+'/begin/' + username,
      type=='login'?$this.options.getLoginData():$this.options.getUnbindData(),
      function (data) {
        if(type=='login'){
          data = $this.options.checkResponseBeginLogin(data);
        }else{
          data = $this.options.checkResponseBeginUnbind(data);
        }
        return data;
      },'json')
      .then((credentialRequestOptions) => {
        $this.options.debug && console.log(credentialRequestOptions);

        if(typeof credentialRequestOptions.Code != 'undefined') throw new Error(credentialRequestOptions.Info);

        credentialRequestOptions.publicKey.challenge = bufferDecode(credentialRequestOptions.publicKey.challenge);
        credentialRequestOptions.publicKey.allowCredentials.forEach(function (listItem) {
          listItem.id = bufferDecode(listItem.id)
        });

        return navigator.credentials.get({
          publicKey: credentialRequestOptions.publicKey
        })
      })
      .then((assertion) => {
        $this.options.debug && console.log(assertion);
        let authData = assertion.response.authenticatorData;
        let clientDataJSON = assertion.response.clientDataJSON;
        let rawId = assertion.rawId;
        let sig = assertion.response.signature;
        let userHandle = assertion.response.userHandle;

        $.ajax({
          type:'post',
          url:$this.options.urlPrefix + '/'+type+'/finish/' + username,
          dataType: 'json',
          async: false,
          data:JSON.stringify({
            id: assertion.id,
            rawId: bufferEncode(rawId),
            type: assertion.type,
            response: {
              authenticatorData: bufferEncode(authData),
              clientDataJSON: bufferEncode(clientDataJSON),
              signature: bufferEncode(sig),
              userHandle: bufferEncode(userHandle),
            },
          }),
          success:function (data) {
            if(type=='login'){
              data = $this.options.checkResponseFinishLogin(data);
            }else{
              data = $this.options.checkResponseFinishUnbind(data);
            }
            return data;
          }});
      })
      .then((response) => {
        $this.options.debug && alert("successfully "+type+" " + username + "!");
        if(type=='login'){
          $this.options.onLoginSuccess.call(this, response);
        }else{
          $this.options.onUnbindSuccess.call(this, response);
        }
      })
      .catch((error) => {
        console.log("failed to "+type+" " + username + ": " +error);
        if(type=='login'){
          $this.options.onLoginError.call(this, error);
        }else{
          $this.options.onUnbindError.call(this, error);
        }
      })
  }

  webAuthn.prototype.login = function (username) {
    this.auth(username,'login');
  }

  webAuthn.prototype.unbind = function (username) {
    this.auth(username,'unbind');
  }
  window.WebAuthn = webAuthn;
})();