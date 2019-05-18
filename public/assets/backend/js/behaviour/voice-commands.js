$(function(){
function getParam(sParam)
{
    var sPageURL = window.location.search.substring(1);
    var sURLVariables = sPageURL.split('&');
    for (var i = 0; i < sURLVariables.length; i++) 
    {
        var sParameterName = sURLVariables[i].split('=');
        if (sParameterName[0] == sParam) 
        {
            return sParameterName[1];
        }
    }
}

if(('webkitSpeechRecognition' in window)){
  
  //Speech Recognition Options
  App.speech({
    lang: 'en'
  });//Initialize
  
  /*Goto Command*/
  App.speechCommand('go to',{
    action: function(){
      $.gritter.add({
      title: "Go to Page",
      text: 'Tell me where do you want to go?',
      image: 'images/mic-icon.png',
      class_name: 'clean',
      time: ''
      });          
    },
    listen: function(datos){
      switch(datos){
        case "dashboard": location.href = "index.html?listen=on"; break;
        case "sidebar": location.href = "layouts-sidebar.html?listen=on"; break;
        case "ui elements": location.href = "ui-elements.html?listen=on"; break;
        case "buttons": location.href = "ui-buttons.html?listen=on"; break;
        case "icons": location.href = "ui-icons.html?listen=on"; break;
        case "grid": location.href = "ui-grid.html?listen=on"; break;
        case "tabs": location.href = "ui-tabs-accordions.html?listen=on"; break;
        case "accordions": location.href = "ui-tabs-accordions.html?listen=on"; break;
        case "tabs and accordions": location.href = "ui-tabs-accordions.html?listen=on"; break;
        case "nestable lists": location.href = "ui-nestable-lists.html?listen=on"; break;
        case "form elements": location.href = "form-elements.html?listen=on"; break;
        case "validation": location.href = "form-validation.html?listen=on"; break;
        case "wizard": location.href = "form-wizard.html?listen=on"; break;
        case "input masks": location.href = "form-masks.html?listen=on"; break;
        case "text editor": location.href = "form-wysiwyg.html?listen=on"; break;
        case "tables": location.href = "tables-general.html?listen=on"; break;
        case "data tables": location.href = "tables-datatables.html?listen=on"; break;
        case "maps": location.href = "maps.html?listen=on"; break;
        case "typography": location.href = "typography.html?listen=on"; break;
        case "charts": location.href = "charts.html?listen=on"; break;
        case "blank page": location.href = "pages-blank.html?listen=on"; break;
        case "blank page header": location.href = "pages-blank-header.html?listen=on"; break;
        case "login": location.href = "pages-login.html?listen=on"; break;
        case "404 page": location.href = "pages-404.html?listen=on"; break;
        case "500 page": location.href = "pages-500.html?listen=on"; break;
        case "500 page": location.href = "pages-500.html?listen=on"; break;
        default:
          $.gritter.add({title: "Error",text: "Could not find: <strong>" + datos + "</strong> page, Please try again.",image: 'images/mic-icon.png',class_name: 'clean',time: ''});  
        break;
      }
    }
  });
  
  /*Change Title Command*/
  App.speechCommand('change title',{
    action: function(){
      $.gritter.add({
      title: "Change Title",
      text: 'Tell me the new title...',
      image: 'images/mic-icon.png',
      class_name: 'clean',
      time: ''
      });     
    },
    listen: function(r){
      $(".navbar-brand span").html(r);
    },
    interim: function(r){
      $(".navbar-brand span").html(r);
    }
  });
  
  /*Logout*/
  App.speechCommand('log out',{
    action: function(){
      location.href = "pages-login.html";
    }
  });
  
  /*Toggle Sidebar*/
  App.speechCommand('toggle sidebar',{
    action: function(){
      App.toggleSideBar();
    }
  });
  
  /*Scroll Down*/
  App.speechCommand('scroll down',{
    action: function(){
      var y = $(window).scrollTop();
     $("html, body").animate({
          scrollTop:  y + 500
     },'slow');
    }
  });  
  
  /*Scroll Up*/
  App.speechCommand('scroll up',{
    action: function(){
     var y = $(window).scrollTop();
     $("html, body").animate({
          scrollTop:  y - 500
     },'slow');
    }
  });
    
  /*Go Bottom*/
  App.speechCommand('go bottom',{
    action: function(){
     $("html, body").animate({
          scrollTop:  $(document).height()
     },'slow');
    }
  });  
  
  /*Go Top*/
  App.speechCommand('go top',{
    action: function(){
     $("html, body").animate({
          scrollTop:  0
     },'slow');
    }
  });  

  /*Hello*/
  var hello_actions = {
    action: function(){
      $.gritter.add({
      title: "Hello",
      text: 'Tell me what is your name?',
      image: 'images/user-icon.png',
      class_name: 'clean',
      time: ''
      });  
    },
    listen:function(r){
      function titleCase(str){return str.replace(/\w\S*/g, function(txt){return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();});}
      var name = titleCase(r);
      $.gritter.add({
      title: "Clean Zone",
      text: 'Welcome home <strong>' + name + '</strong>!' ,
      image: 'images/user-icon.png',
      class_name: 'clean',
      time: ''
      });    
      $(".side-user .info a").html(name);
      $(".profile_menu > a span").html(name);
    }
  };
  
  App.speechCommand('hello',hello_actions);
  App.speechCommand('hi',hello_actions);
  
  /*Thank you*/
  var thank_actions = {
    action: function(){
      $.gritter.add({
      title: "Clean Zone",
      text: 'Your welcome!',
      image: 'images/user-icon.png',
      class_name: 'clean',
      time: ''
      });  
    }
  };
  App.speechCommand('thanks',thank_actions);
  App.speechCommand('thank you',thank_actions);

  
  /*Compose Email*/
  App.speechCommand('email',{
    dictation: true,
    dictationEndCommand: 'end of email',
    dictationEnd: function(){
      var progress  = $('<div class="progress progress-striped active" style="display:none;"><div style="width: 0%" class="progress-bar progress-bar-info">0%</div></div>').css('margin','10px 0 0 0');
      $('#mod-info .interim').fadeOut(function(){ $(this).html(""); });     
      progress.appendTo("#mod-info .modal-body").fadeIn();
      progress.find(".progress-bar").animate({width:900},{
        duration: 5000,
        step:function(now, fx){
          var percent = ((100*now)/fx.end).toFixed(0);
          $(this).html( percent + "%");
        },
        complete:function(){
          $("#mod-info").modal('hide');
        }
      });
      $("#mod-info .modal-body h4").html("Thank you!");
      $("#mod-info .modal-body p").addClass("text-center").html("We are sending a new e-mail...");
    },
    action: function(){
      var modal = $('<div role="dialog" tabindex="-1" id="mod-info" class="modal fade"><div class="modal-dialog"><div class="modal-content"><div class="modal-header"><button aria-hidden="true" data-dismiss="modal" class="close" type="button">×</button></div><div class="modal-body"><div class="text-center"><div class="i-circle primary"><i class="fa fa-envelope"></i></div><h4>Tell me your message...</h4><div contenteditable="true"><p class="text-left"><span class="msg"></span><span class="interim color-primary"></span></p></div></div></div><div class="modal-footer"><button data-dismiss="modal" class="btn btn-default" type="button">Cancel</button><button data-dismiss="modal" class="btn btn-primary" type="button">Send</button></div></div></div></div>');
      modal.appendTo("body");
      $('#mod-info').modal();
      $('#mod-info').on('hidden.bs.modal', function () {
          $(this).removeData('bs.modal');
          $(this).remove();
      });
    },
    listen: function(r){
      $('#mod-info .msg').append(" " + r);
      $('#mod-info .interim').fadeOut(function(){ $(this).html(""); });
    },
    interim: function(r){
      $('#mod-info .interim').show().html(r);
    }
  });  
  
  /*Stop Recognition*/
  App.speechCommand('stop',{
    action: function(){
      App.speech("stop");
    }
  });
  
  if(getParam("listen")=="on"){
    App.speech("start");
  }
  
  $(".speech-button").click(function(e){
    var modal = $('<div role="dialog" tabindex="-1" id="mod-sound" data-backdrop="false" class="modal fade"><div class="modal-dialog"><div class="modal-content"><div class="modal-header"><button aria-hidden="true" data-dismiss="modal" class="close" type="button">×</button><h2 class="hthin"><img src="images/mic-icon.png" /> Speech API</h2></div><div class="modal-body" style="padding:0 25px;"><div><h4>Voice Recognition</h4><div><p class="text-left">Thanks to Google Speech API we can do <a href="https://dvcs.w3.org/hg/speech-api/raw-file/tip/speechapi.html">Speech Recognition</a> in our web sites, initially Chrome 25+ and up versions support this, but don&#39;t worry! browsers are working on a <a href="https://wiki.mozilla.org/HTML5_Speech_API">standard</a> and soon we will see this working on our favorites browsers. </p><h4 class="spacer2">Let the party begin</h4><p>First you must to allow microphone access clicking on <strong>"Allow"</strong> option above this modal. After this you&#39;ll see the Microphone icon with a blur effect, this means that voice recognition starts.</p><h4 class="spacer2">Voice Commands</h4><p>After that, try to say <strong>"Hello"</strong> at your mic. Do you see what happens? things in template start to change, now try these commands:</p><ul><li><strong>"Go to"</strong>: wait for a message and then say a page title like "tables"</li><li><strong>"Change title"</strong> - Sets template title</li><li><strong>"Scroll down"</strong> and <strong>"Scroll up"</strong></li><li><strong>"Go top"</strong> and <strong>"Go bottom"</strong></li><li><strong>"Toggle sidebar"</strong></li><li><strong>"log out"</strong></li><li><strong>"Thank you"</strong></li><li><strong>"Stop"</strong> - Stops recognition</li><li><strong>"Email"</strong> - Compose an e-mail with your voice, to end dictation just say "end of email"</li></ul><p>Do you want more commands? you can add a voice command easily with our own API.</p></div></div></div><div class="modal-footer"><button data-dismiss="modal" class="btn btn-primary" type="button">DONE!</button></div></div></div></div>');
    modal.appendTo("body");
    $('#mod-sound').modal();
    $('#mod-sound').on('hidden.bs.modal', function () {
        $(this).removeData('bs.modal');
        $(this).remove();
    });
    App.speech("start");
    e.preventDefault();
    e.stopPropagation();
  });

}

});