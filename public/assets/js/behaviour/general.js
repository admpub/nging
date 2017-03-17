var App = function () {

  var config = {//Basic Config
    tooltip: true,
    popover: true,
    nanoScroller: true,
    nestableLists: true,
    hiddenElements: true,
    bootstrapSwitch:true,
    dateTime:true,
    select2:true,
    tags:true,
    slider:true
  }; 
  
  var voice_methods = [];

  /*Widgets*/
  var widgets = function(){
    var skycons = new Skycons({"color": "#FFFFFF"});
    skycons.add($("#sun-icon")[0], Skycons.PARTLY_CLOUDY_DAY);
    skycons.play();
    
  };//End of widgets
  
  
  /*Speech Recognition*/
  var speech_commands = [];
  if(('webkitSpeechRecognition' in window)){
    var rec = new webkitSpeechRecognition();  
  }
  
  var speech = function(options){
   
    if(('webkitSpeechRecognition' in window)){
    
      if(options=="start"){
        rec.start();
      }else if(options=="stop"){
        rec.stop();
      }else{
        var config = {
          continuous: true,
          interim: true,
          lang: false,
          onEnd: false,
          onResult: false,
          onNoMatch: false,
          onSpeechStart: false,
          onSpeechEnd: false
        };
        $.extend( config, options );
        
        if(config.continuous){rec.continuous = true;}
        if(config.interim){rec.interimResults = true;}
        if(config.lang){rec.lang = config.lang;}        
        var values = false;
        var val_command = "";
        
        rec.onresult = function(event){
          for (var i = event.resultIndex; i < event.results.length; ++i) {
            if (event.results[i].isFinal) {
              var command = event.results[i][0].transcript;//Return the voice command
              command = command.toLowerCase();//Lowercase
              command = command.replace(/^\s+|\s+$/g,'');//Trim spaces
              console.log(command);
              if(config.onResult){
                config.onResult(command);
              }   
              
              $.each(speech_commands,function(i, v){
                if(values){//Second command
                  if(val_command == v.command){
                    if(v.dictation){
                      if(command == v.dictationEndCommand){
                        values = false;
                        v.dictationEnd(command);
                      }else{
                        v.listen(command);
                      }
                    }else{
                      values = false;
                      v.listen(command);
                    }
                  }
                }else{//Primary command
                  if(v.command == command){
                    v.action(command);
                    if(v.listen){
                      values = true;
                      val_command = v.command;
                    }
                  }
                }
              });
            }else{
              var res = event.results[i][0].transcript;//Return the interim results
              $.each(speech_commands,function(i, v){
                if(v.interim !== false){
                  if(values){                
                    if(val_command == v.command){
                      v.interim(res);
                    }
                  }
                }
              });
            }
          }
        };      
        
        
        if(config.onNoMatch){rec.onnomatch = function(){config.onNoMatch();};}
        if(config.onSpeechStart){rec.onspeechstart = function(){config.onSpeechStart();};}
        if(config.onSpeechEnd){rec.onspeechend = function(){config.onSpeechEnd();};}
        rec.onaudiostart = function(){ $(".speech-button i").addClass("blur"); }
        rec.onend = function(){
          $(".speech-button i").removeClass("blur");
          if(config.onEnd){config.onEnd();}
        };
        
        
      }    
      
    }else{ 
      alert("Only Chrome25+ browser support voice recognition.");
    }
   
    
  };
  
  var speechCommand = function(command, options){
    var config = {
      action: false,
      dictation: false,
      interim: false,
      dictationEnd: false,
      dictationEndCommand: 'final.',
      listen: false
    };
    
    $.extend( config, options );
    if(command){
      if(config.action){
        speech_commands.push({
          command: command, 
          dictation: config.dictation,
          dictationEnd: config.dictationEnd,
          dictationEndCommand: config.dictationEndCommand,
          interim: config.interim,
          action: config.action, 
          listen: config.listen 
        });
      }else{
        alert("Must have an action function");
      }
    }else{
      alert("Must have a command text");
    }
  };
  
      function toggleSideBar(_this){
        var b = $("#sidebar-collapse")[0];
        var w = $("#cl-wrapper");
        var s = $(".cl-sidebar");
        
        if(w.hasClass("sb-collapsed")){
          $(".fa",b).addClass("fa-angle-left").removeClass("fa-angle-right");
          w.removeClass("sb-collapsed");
        }else{
          $(".fa",b).removeClass("fa-angle-left").addClass("fa-angle-right");
          w.addClass("sb-collapsed");
        }
        //updateHeight();
      }
      
      function updateHeight(){
        if(!$("#cl-wrapper").hasClass("fixed-menu")){
          var button = $("#cl-wrapper .collapse-button").outerHeight();
          var navH = $("#head-nav").height();
          //var document = $(document).height();
          var cont = $("#pcont").height();
          var sidebar = ($(window).width() > 755 && $(window).width() < 963)?0:$("#cl-wrapper .menu-space .content").height();
          var windowH = $(window).height();
          
          if(sidebar < windowH && cont < windowH){
            if(($(window).width() > 755 && $(window).width() < 963)){
              var height = windowH;
            }else{
              var height = windowH - button - navH;
            }
          }else if((sidebar < cont && sidebar > windowH) || (sidebar < windowH && sidebar < cont)){
            var height = cont + button + navH;
          }else if(sidebar > windowH && sidebar > cont){
            var height = sidebar + button;
          }  
          
          // var height = ($("#pcont").height() < $(window).height())?$(window).height():$(document).height();
          $("#cl-wrapper .menu-space").css("min-height",height);
        }else{
          $("#cl-wrapper .nscroller").nanoScroller({ preventPageScrolling: true });
        }
      }
        
  return {
   
    init: function (options) {
      //Extends basic config with options
      $.extend( config, options );
      
      /*VERTICAL MENU*/
      $(".cl-vnavigation li ul").each(function(){
        $(this).parent().addClass("parent");
      });
      
      $(".cl-vnavigation li ul li.active").each(function(){
        $(this).parent().show().parent().addClass("open");
        //setTimeout(function(){updateHeight();},200);
      });
      
      $(".cl-vnavigation").delegate(".parent > a","click",function(e){
        $(".cl-vnavigation .parent.open > ul").not($(this).parent().find("ul")).slideUp(300, 'swing',function(){
           $(this).parent().removeClass("open");
        });
        
        var ul = $(this).parent().find("ul");
        ul.slideToggle(300, 'swing', function () {
          var p = $(this).parent();
          if(p.hasClass("open")){
            p.removeClass("open");
          }else{
            p.addClass("open");
          }
          //var menuH = $("#cl-wrapper .menu-space .content").height();
          // var height = ($(document).height() < $(window).height())?$(window).height():menuH;
          //updateHeight();
         $("#cl-wrapper .nscroller").nanoScroller({ preventPageScrolling: true });
        });
        e.preventDefault();
      });
      
      /*Small devices toggle*/
      $(".cl-toggle").click(function(e){
        var ul = $(".cl-vnavigation");
        ul.slideToggle(300, 'swing', function () {
        });
        e.preventDefault();
      });
      
      /*Collapse sidebar*/
      $("#sidebar-collapse").click(function(){
          toggleSideBar();
      });
      
      
      if($("#cl-wrapper").hasClass("fixed-menu")){
        var scroll =  $("#cl-wrapper .menu-space");
        scroll.addClass("nano nscroller");
 
        function update_height(){
          var button = $("#cl-wrapper .collapse-button");
          var collapseH = button.outerHeight();
          var navH = $("#head-nav").height();
          var height = $(window).height() - ((button.is(":visible"))?collapseH:0) - navH;
          scroll.css("height",height);
          $("#cl-wrapper .nscroller").nanoScroller({ preventPageScrolling: true });
        }
        
        $(window).resize(function() {
          update_height();
        });    
            
        update_height();
        $("#cl-wrapper .nscroller").nanoScroller({ preventPageScrolling: true });
        
      }else{
        $(window).resize(function(){
          //updateHeight();
        }); 
        //updateHeight();
      }

      
      /*SubMenu hover */
        var tool = $("<div id='sub-menu-nav' style='position:fixed;z-index:9999;'></div>");
        
        function showMenu(_this, e){
          if(($("#cl-wrapper").hasClass("sb-collapsed") || ($(window).width() > 755 && $(window).width() < 963)) && $("ul",_this).length > 0){   
            $(_this).removeClass("ocult");
            var menu = $("ul",_this);
            if(!$(".dropdown-header",_this).length){
              var head = '<li class="dropdown-header">' +  $(_this).children().html()  + "</li>" ;
              menu.prepend(head);
            }
            
            tool.appendTo("body");
            var top = ($(_this).offset().top + 8) - $(window).scrollTop();
            var left = $(_this).width();
            
            tool.css({
              'top': top,
              'left': left + 8
            });
            tool.html('<ul class="sub-menu">' + menu.html() + '</ul>');
            tool.show();
            
            menu.css('top', top);
          }else{
            tool.hide();
          }
        }

        $(".cl-vnavigation li").hover(function(e){
          showMenu(this, e);
        },function(e){
          tool.removeClass("over");
          setTimeout(function(){
            if(!tool.hasClass("over") && !$(".cl-vnavigation li:hover").length > 0){
              tool.hide();
            }
          },500);
        });
        
        tool.hover(function(e){
          $(this).addClass("over");
        },function(){
          $(this).removeClass("over");
          tool.fadeOut("fast");
        });
        
        
        $(document).click(function(){
          tool.hide();
        });
        $(document).on('touchstart click', function(e){
          tool.fadeOut("fast");
        });
        
        tool.click(function(e){
          e.stopPropagation();
        });
     
        $(".cl-vnavigation li").click(function(e){
          if((($("#cl-wrapper").hasClass("sb-collapsed") || ($(window).width() > 755 && $(window).width() < 963)) && $("ul",this).length > 0) && !($(window).width() < 755)){
            showMenu(this, e);
            e.stopPropagation();
          }
        });
        
        $(".cl-vnavigation li").on('touchstart click', function(){
          //alert($(window).width());
        });
        
      $(window).resize(function(){
        //updateHeight();
      });

      var domh = $("#pcont").height();
      $(document).bind('DOMSubtreeModified', function(){
        var h = $("#pcont").height();
        if(domh != h) {
          //updateHeight();
        }
      });
      
      /*Return to top*/
      var offset = 220;
      var duration = 500;
      var button = $('<a href="#" class="back-to-top"><i class="fa fa-angle-up"></i></a>');
      button.appendTo("body");
      
      jQuery(window).scroll(function() {
        if (jQuery(this).scrollTop() > offset) {
            jQuery('.back-to-top').fadeIn(duration);
        } else {
            jQuery('.back-to-top').fadeOut(duration);
        }
      });
    
      jQuery('.back-to-top').click(function(event) {
          event.preventDefault();
          jQuery('html, body').animate({scrollTop: 0}, duration);
          return false;
      });
      
      /*Datepicker UI*/
      if($(".ui-datepicker").length>0)$(".ui-datepicker").datepicker();
      
      /*Tooltips*/
      if(config.tooltip){
        $('.ttip, [data-toggle="tooltip"]').tooltip();
      }
      
      /*Popover*/
      if(config.popover){
        $('[data-popover="popover"]').popover();
      }

      /*NanoScroller*/      
      if(config.nanoScroller){
        $(".nscroller").nanoScroller();     
      }
      
      /*Nestable Lists*/
      if(config.nestableLists&&$('.dd').length>0){
        $('.dd').nestable();
      }
      
      /*Switch*/
      if(config.bootstrapSwitch){
        if($('.switch').length>0)$('.switch').bootstrapSwitch();
      }
      
      /*DateTime Picker*/
      if(config.dateTime){
        if($(".datetime").length>0)$(".datetime").datetimepicker({autoclose: true});
      }
      
      /*Select2*/
      if(config.select2){
         if($(".select2").length>0)$(".select2").select2({
          width: '100%'
         });
      }
      
       /*Tags*/
      if(config.tags){
        if($(".tags").length>0)$(".tags").select2({tags: 0,width: '100%'});
      }
      
       /*Slider*/
      if(config.slider){
        if($('.bslider').length>0)$('.bslider').slider();     
      }
      
      /*Input & Radio Buttons*/
      if(jQuery().iCheck){
        $('.icheck').iCheck({
          checkboxClass: 'icheckbox_square-blue checkbox',
          radioClass: 'iradio_square-blue'
        });
      }
      
      /*Bind plugins on hidden elements*/
      if(config.hiddenElements){
      	/*Dropdown shown event*/
        $('.dropdown').on('shown.bs.dropdown', function () {
          $(".nscroller").nanoScroller();
        });
          
        /*Tabs refresh hidden elements*/
        $('.nav-tabs').on('shown.bs.tab', function (e) {
          $(".nscroller").nanoScroller();
        });
      }
      
    },
      
    speech: function(options){
      speech(options);
    },
    
    speechCommand: function(com, options){
      speechCommand(com, options);
    },
    
    toggleSideBar: function(){
      toggleSideBar();
    },
    
    widgets: function(){
      widgets();
    },

    markNavByURL:function(url){
      if(url==null)url=window.location.pathname;
      App.markNav($('.cl-vnavigation a[href="'+url+'"]'));
    },
    markNav:function(curNavA){
	    if(curNavA.length>0){
        curNavA.parent('li').addClass('active').parent().show().parent().addClass("open");
	    }
    },
    unmarkNav:function(){
      $('.cl-vnavigation > .open').removeClass('open').children('.sub-menu').hide().children('li.active').removeClass('active');
      $('.cl-vnavigation > .active').removeClass('active');
    },
    message:function(options,sticky){
      var defaults={title: '', text: "", image: '', class_name: 'clean',//primary|info|danger|warning|success|dark
        sticky: true};
      if(typeof(options)!="object"){
        options={text:options};
      }
      options=$.extend({},defaults,options||{});
      switch(options.class_name){
        case 'dark':
        case 'primary':
        case 'clean':
        case 'info':
          options.title='<i class="fa fa-info-circle"></i> '+options.title;break;
        case 'danger':
          options.title='<i class="fa fa-comment-o"></i> '+options.title;break;
        case 'warning':
          options.title='<i class="fa fa-warning"></i> '+options.title;break;
        case 'success':
          options.title='<i class="fa fa-check"></i> '+options.title;break;
      }
      if(sticky!=null)options.sticky=sticky;
	    $.gritter.add(options);
    },
    attachAjaxURL:function(){
      $(document).on('click','[data-ajax-url]',function(){
        var url=$(this).data('ajax-url');
        var title=$(this).attr('title');
        if(!title)title=$(this).text();
        $.get(url,{},function(r){
          App.message({title:title,text:r,time:5000,sticky:false});
        },'html');
      });
    },
    attachPjax:function(elem,callbacks,timeout){
      if(!$.support.pjax)return;
      if(elem==null)elem='a';
      if(timeout==null)timeout=5000;
      var defaults={onclick:null,onsend:null,oncomplete:null,ontimeout:null,onstart:null,onend:null};
      var options=$.extend({},defaults,callbacks||{});
      $(document).on('click', elem+'[data-pjax]', function(event) {
        var container = $(this).data('pjax');
        $.pjax.click(event, $(container),{timeout:timeout});
        if(options.onclick)options.onclick(this);
      }).on('pjax:send',function(evt,xhr,option){
        App.loading('show');
        if(options.onsend)options.onsend(evt,xhr,option);
      }).on('pjax:complete',function(evt, xhr, textStatus, option){
        App.loading('hide');
        if(options.oncomplete)options.oncomplete(evt, xhr, textStatus, option);
      }).on('pjax:timeout',function(evt,xhr,option){
        console.log('timeout');
        if(options.ontimeout)options.ontimeout(evt,xhr,option);
      }).on('pjax:start',function(evt,xhr,option){
        if(options.onstart)options.onstart(evt,xhr,option);
      }).on('pjax:end',function(evt,xhr,option){
        if(options.onend)options.onend(evt,xhr,option);
      });
    },
    
    websocket:function(showmsg,url,onopen){
    	var protocol='ws:';
    	if(window.location.protocol=='https:')protocol='wss:';
    	var ws = new WebSocket(protocol+"//"+window.location.host+url);
    	ws.onopen = function(evt) {
    	    console.log('Websocket Server is connected');
    		  if(onopen!=null&&typeof(onopen)=='function')onopen();
    	};
    	ws.onclose = function(evt) {
    	    console.log('Websocket Server is disconnected');
    	};
    	ws.onmessage = function(evt) {
    	    showmsg(evt.data);
    	};
    	ws.onerror = function(evt) {
    	    console.dir(evt);
    	};
      if(onopen!=null&&typeof(onopen)=='object'){
        ws=$.extend({},ws,onopen);
      }
      return ws;
    },

    text2html:function(text){
      return String(text).replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/\n/g,'<br />').replace(/  /g,'&nbsp; ').replace(/\t/g,'&nbsp; &nbsp; ')
    },
    
    iCheck:function(elem,on,callback){
        var icOn='';
        switch(on){
            case 'click':
            icOn='ifClicked';
            break;
            case 'change':
            icOn='ifChanged';
            break;
            default:
            alert('unsupported '+on);
            return;
        }
        if($(elem).first().next('.iCheck-helper').length<1){
          $(elem).iCheck({
            checkboxClass: 'icheckbox_square-blue checkbox',
            radioClass: 'iradio_square-blue'
          });
        }
        $(elem).on(icOn,function(){
            $(this).trigger(on);
        }).on(on,callback);
    },
    alertBlock:function(content,title,type){
      switch(type){
        case 'info':
        if(title==null)title='Info!';
        return '<div class="alert alert-info">\
								<button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>\
								<i class="fa fa-info-circle sign"></i><strong>'+title+'</strong> '+content+'</div>';
        case 'warn':
        if(title==null)title='Alert!';
        return '<div class="alert alert-warning">\
								<button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>\
								<i class="fa fa-warning sign"></i><strong>'+title+'</strong> '+content+'</div>';
        case 'error':
        if(title==null)title='Error!';
        return '<div class="alert alert-danger">\
								<button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>\
								<i class="fa fa-times-circle sign"></i><strong>'+title+'</strong> '+content+'</div>';
        default:
        if(title==null)title='Success!';
        return '<div class="alert alert-success">\
								<button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>\
								<i class="fa fa-check sign"></i><strong>'+title+'</strong> '+content+'</div>';
      }
    },
    alertBlockx:function(content,title,type){
      switch(type){
        case 'info':
        if(title==null)title='Info!';
        return '<div class="alert alert-info alert-white rounded">\
								<button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>\
								<div class="icon"><i class="fa fa-info-circle"></i></div>\
								<strong>'+title+'</strong> '+content+'</div>';
        case 'warn':
        if(title==null)title='Alert!';
        return '<div class="alert alert-warning alert-white rounded">\
								<button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>\
								<div class="icon"><i class="fa fa-warning"></i></div>\
								<strong>'+title+'</strong> '+content+'</div>';
        case 'error':
        if(title==null)title='Error!';
        return '<div class="alert alert-danger alert-white rounded">\
								<button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>\
								<div class="icon"><i class="fa fa-times-circle"></i></div>\
								<strong>'+title+'</strong> '+content+'</div>';
        default:
        if(title==null)title='Success!';
        return '<div class="alert alert-success alert-white rounded">\
								<button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>\
								<div class="icon"><i class="fa fa-check"></i></div>\
								<strong>'+title+'</strong> '+content+'</div>';
      }
    },
    showDbLog:function(results,container){
      if(container==null)container='.block-flat:first';
      if(typeof(results.length)=='undefined')results=[results];
      for(var j=0;j<results.length;j++){
        var result=results[j];
        var s=result.Started;
        if(result.SQLs && result.SQLs.length>0){
          for(var i=0;i<result.SQLs.length;i++) s+='<code class="wrap">'+result.SQLs[i]+'</code>';
        }else{
          s+='<code class="wrap">'+result.SQL+'</code>';
        }
        var t='success';
        if(result.Error){
          s+='('+result.Error+')';
          t='error'
        }else{
          s+='('+result.Elapsed+')';
        }
        $(container).before(App.alertBlockx(s,null,t));
      }
    },
    loading:function(op){
      var obj=$('#loading-status');
      switch(op){
        case 'show':
          if(obj.length>0){
            obj.show();
          }else{
            $('body').append('<div id="loading-status"><i class="fa fa-spinner fa-spin fa-3x"></i></div>');
          }
          break;
        case 'hide':
          if(obj.length>0){
            obj.hide();
          }
      }
    },
    insertAtCursor: function(myField, myValue) { /* IE support */
			if (document.selection) {
				myField.focus();
				sel = document.selection.createRange();
				sel.text = myValue;
				sel.select();
			} /* MOZILLA/NETSCAPE support */
			else if (myField.selectionStart || myField.selectionStart == '0') {
				var startPos = myField.selectionStart;
				var endPos = myField.selectionEnd; /* save scrollTop before insert */
				var restoreTop = myField.scrollTop;
				myField.value = myField.value.substring(0, startPos) + myValue + myField.value.substring(endPos, myField.value.length);
				if (restoreTop > 0) myField.scrollTop = restoreTop;
				myField.focus();
				myField.selectionStart = startPos + myValue.length;
				myField.selectionEnd = startPos + myValue.length;
			} else {
				myField.value += myValue;
				myField.focus();
			}
		}
  };
 
}();

$(function(){
  //$("body").animate({opacity:1,'margin-left':0},500);
  $("body").css({opacity:1,'margin-left':0});
});