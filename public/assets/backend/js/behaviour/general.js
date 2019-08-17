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
  
  /*Form Wizard*/
  var wizard = function(){
    //Fuel UX
    $('.wizard-ux').wizard();

    $('.wizard-ux').on('changed',function(){
      //delete $.fn.slider;
      $('.bslider').slider();
    });
    
    $(".wizard-next").click(function(e){
      $('.wizard-ux').wizard('next');
      e.preventDefault();
    });
    
    $(".wizard-previous").click(function(e){
      $('.wizard-ux').wizard('previous');
      e.preventDefault();
    });
  };//End of wizard

      function toggleSideBar(_this){
        var b = $("#sidebar-collapse")[0];
        var w = $("#cl-wrapper");
        if(w.hasClass("sb-collapsed")){
          $(".fa",b).addClass("fa-angle-left").removeClass("fa-angle-right");
          w.removeClass("sb-collapsed");
        }else{
          $(".fa",b).removeClass("fa-angle-left").addClass("fa-angle-right");
          w.addClass("sb-collapsed");
        }
      }
      
      function pageAside(){
        var pageKey=window.location.pathname.replace(/^\//,'').replace(/[^\w]/g,'-')+'-'+window.location.search.replace(/[^\w]/g,'_');
        $('.page-aside').each(function(index){
          var aside=$(this);
          $(this).find('.header > .collapse-button').on('click',function(){
              aside.addClass('collapsed');
              store.set(pageKey+'.page-aside-'+index,'collapsed');
              aside.trigger('collapsed');
              tableReponsive(0);
          }); 
          $(this).children('.collapsed-button').on('click',function(){
              aside.removeClass('collapsed');
              store.set(pageKey+'.page-aside-'+index,'');
              aside.trigger('expanded');
              tableReponsive(280);
          }); 
          if(store.get(pageKey+'.page-aside-'+index)=='collapsed'){
            $(this).trigger('click');
          }
        });
      }

      function tableReponsiveInit(){
        var aside=$('#main-container > .page-aside');
        var asideWidth=aside.width();
        //if($('#pcont').width()>bodyWidth){
        if(asideWidth<270){
          aside.find('.header > .collapse-button').trigger('click');
          return;
        }
        tableReponsive(asideWidth);
      }
      
      function tableReponsive(asideWidth){
        var sidebarWidth=$('#cl-wrapper > .cl-sidebar').width();
        var windowWidth=$(window).width();
        var bodyWidth=windowWidth-sidebarWidth-asideWidth;
        bodyWidth-=80;
        $('#cl-wrapper > .cl-body .table-responsive').each(function(){
          if($(this).children('table').width()>bodyWidth){
            $(this).addClass('overflow').css('max-width',bodyWidth);
          }else{
            $(this).removeClass('overflow');
          }
        });
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

      var cachedLang=null;
  return {
    clientID: {},
    i18n: {SYS_INFO:'System Information',UPLOAD_ERR:'Upload Error'},
    lang: 'en',
    langInfo: function(){
      if(cachedLang!=null) return cachedLang;
      var _lang=App.lang.split('-',2);
      cachedLang={encoding:_lang[0],country:''};
      if(_lang.length>1)cachedLang.country=_lang[1].toUpperCase();
      return cachedLang;
    },
    langTag: function(){
      var l=App.langInfo();
      if(l.country) return l.encoding+'-'+l.country;
      return l.encoding;
    },
    initTool: function(){
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
    },
    initLeftNav: function(){
      /*VERTICAL MENU*/
      $(".cl-vnavigation li ul").each(function(){
        $(this).parent().addClass("parent");
      });
      
      $(".cl-vnavigation li ul li.active").each(function(){
        $(this).parent().show().parent().addClass("open");
        //setTimeout(function(){updateHeight();},200);
      });
      if(!$(".cl-vnavigation").data('initclick')){
        $(".cl-vnavigation").data('initclick',true);
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
      
      $(".cl-vnavigation li").click(function(e){
        if((($("#cl-wrapper").hasClass("sb-collapsed") || ($(window).width() > 755 && $(window).width() < 963)) && $("ul",this).length > 0) && !($(window).width() < 755)){
          showMenu(this, e);
          e.stopPropagation();
        }
      });
    },
    initLeftNavAjax: function(activeURL,elem){
      App.markNavByURL(activeURL);
      App.attachPjax(elem,{
      	 onclick: function(obj){
      		//console.log($(obj).data('marknav'))
      		if($(obj).data('marknav')){
      			App.unmarkNav($(obj),$(obj).data('marknav'));
      			App.markNav($(obj),$(obj).data('marknav'));
      		}
      	 },
      	 onend: function(evt,xhr,opt){
      		 opt.container.find('[data-popover="popover"]').popover();
      		 opt.container.find('.ttip, [data-toggle="tooltip"]').tooltip();
      	 }
      });
      App.attachAjaxURL(elem);
    },
    init: function (options) {
      //Extends basic config with options
      $.extend( config, options );
      App.initLeftNav();
      App.initTool();
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
      }

      /*Return to top*/
      var offset = 220;
      var duration = 500;
      var button = $('<a href="#" class="back-to-top"><i class="fa fa-angle-up"></i></a>');
      button.appendTo("body");
      
      $(window).scroll(function() {
        if ($(this).scrollTop() > offset) {
            $('.back-to-top').fadeIn(duration);
        } else {
            $('.back-to-top').fadeOut(duration);
        }
      });
    
      $('.back-to-top').click(function(event) {
          event.preventDefault();
          $('html, body').animate({scrollTop: 0}, duration);
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
        if($('.switch:not(.has-switch)').length>0)$('.switch:not(.has-switch)').bootstrapSwitch();
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
      App.autoFixedThead();
    },
    autoFixedThead:function (prefix){
      if(prefix==null) prefix=''; 
      App.topFloatThead(prefix+'thead.auto-fixed',$('#head-nav').height());
    },
    pageAside: function(options){
      pageAside(options);
    },
    tableReponsiveInit: function(options){
      tableReponsiveInit(options);
    },
    tableReponsive: function(options){
      tableReponsive(options);
    },
    toggleSideBar: function(){
      toggleSideBar();
    },
    wizard: function(){
      wizard();
    },
    markNavByURL:function(url){
      if(url==null)url=window.location.pathname;
      App.markNav($('#leftnav a[href="'+url+'"]'),'left');
      App.markNav($('#topnav a[href="'+url+'"]'),'top');
    },
    markNav:function(curNavA,position){
	    if(curNavA.length>0){
        var li=curNavA.parent('li').addClass('active');
        switch(position){
          case 'left':
          li.parent('.from-left').show().parent('li').addClass("open");
          break;

          case 'top':
          li.parent('.from-top').parent('li').addClass("active").siblings('li.active').removeClass('active');
          break;
        }
	    }
    },
    unmarkNav:function(curNavA,position){
      var li=curNavA.parent('li');
      var siblings=li.siblings('li.active');
      if(siblings.length>0){
        siblings.removeClass('active');
      }else{
        switch(position){
          case 'left':
          // 点击的左侧边栏菜单
          if(li.parent('ul.sub-menu').length>0){
            var op2=$('.col-menu-2').children('li.open');
            if(op2.length>0) op2.removeClass('open').find('li.active').removeClass('active');
          }
          $('#leftnav > .open').removeClass('open').children('ul.sub-menu').hide().children('li.active').removeClass('active');
          $('#leftnav .active').removeClass('active');
          break;

          case 'top':
          var topnavDropdown=li.parent('ul.dropdown-menu.from-top');
          if(topnavDropdown.length>0){
            siblings=topnavDropdown.parent('li').addClass('active').siblings('li.active');
          }
          if(siblings.length>0){
            siblings.removeClass('active').children('ul.dropdown-menu.from-top').children('li.active').removeClass('active');
          }
          break;
        }
      }
    },
    message:function(options,sticky){
      var defaults={title: App.i18n.SYS_INFO,
        text: '',
        image: '',
        class_name: 'clean',//primary|info|danger|warning|success|dark
        sticky: false
        //,time: 1000,speed: 500,position: 'bottom-right'
      };
      if(typeof(options)!="object")options={text:options};
      if(typeof(options.type)!="undefined"&&options.type)options.class_name=options.type;
      options=$.extend({},defaults,options||{});
      switch(options.class_name){
        case 'dark':
        case 'primary':
        case 'clean':
        case 'info':
        if(options.title) options.title='<i class="fa fa-info-circle"></i> '+options.title;break;

        case 'error':
          options.class_name='danger';
        case 'danger':
        if(options.title) options.title='<i class="fa fa-comment-o"></i> '+options.title;break;

        case 'warning':
        if(options.title) options.title='<i class="fa fa-warning"></i> '+options.title;break;
          
        case 'success':
        if(options.title) options.title='<i class="fa fa-check"></i> '+options.title;break;
      }
      if(sticky!=null)options.sticky=sticky;
	    $.gritter.add(options);
    },
    attachAjaxURL:function(elem){
      if(elem==null) elem=document;
      $(elem).on('click','[data-ajax-url]',function(){
        var url=$(this).data('ajax-url'),title=$(this).attr('title');
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
        var container = $(this).data('pjax'),keepjs=$(this).data('keepjs');
        var onclick=$(this).data('onclick');
        $.pjax.click(event, $(container),{timeout:timeout,keepjs:keepjs});
        if(options.onclick)options.onclick(this);
        if(onclick && typeof(window[onclick])=='function')window[onclick](this);
      }).on('pjax:send',function(evt,xhr,option){
        App.loading('show');
        if(options.onsend)options.onsend(evt,xhr,option);
      }).on('pjax:complete',function(evt, xhr, textStatus, option){
        App.loading('hide');
        if(options.oncomplete)options.oncomplete(evt, xhr, textStatus, option);
      }).on('pjax:timeout',function(evt,xhr,option){
        console.log('timeout');
        App.loading('hide');
        if(options.ontimeout)options.ontimeout(evt,xhr,option);
      }).on('pjax:start',function(evt,xhr,option){
        if(options.onstart)options.onstart(evt,xhr,option);
      }).on('pjax:end',function(evt,xhr,option){
        App.loading('hide');
        if(options.onend)options.onend(evt,xhr,option);
        //console.debug(option);
        var id=option.container.attr('id');
        if(id){
          App.bottomFloat('#'+id+' .pagination');
          App.bottomFloat('#'+id+' .form-submit-group',0,true);
          $('#'+id+' .switch:not(.has-switch)').bootstrapSwitch();
          App.autoFixedThead('#'+id+' ');
        }
        if(option.type=='GET') $('#global-search-form').attr('action',option.url);
      });
    },
    wsURL:function(url) {
    	var protocol='ws:';
      if(window.location.protocol=='https:')protocol='wss:';
      var p=String(url).indexOf('//');
      if(p==-1){
        url=protocol+"//"+window.location.host+url;
      }else{
        url=protocol+String(url).substring(p);
      }
      return url;
    },
    websocket:function(showmsg,url,onopen){
      url = App.wsURL(url);
    	var ws = new WebSocket(url);
    	ws.onopen = function(evt) {
    	    console.log('Websocket Server is connected');
    		  if(onopen!=null&&$.isFunction(onopen))onopen.apply(this,arguments);
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
    notifyListen:function(){
      var messageCount={
        notify:0,
        element:0,
        modal:0,
      },
      messageMax={
        notify:20,
        element:50,
        modal:50,
      };
      App.websocket(function(message){
        //console.dir(message);
        var m=$.parseJSON(message);
        if(typeof(App.clientID['notify'])=='undefined'){
          App.clientID['notify']=m.client_id;
        }
        switch(m.mode){
          case '-':
          break;
          case 'element':
          var c=$('#notify-element-'+m.type);
          if(c.length<1){
            var callback='recv_notice_'+m.type;
            if(typeof(window[callback])!='undefined'){
              return window[callback](m);
            }
            if(m.status>0){
              console.info(m.content);
            }else{
              console.error(m.content);
            }
            return;
          }
          if(messageCount[m.mode]>=messageMax[m.mode]){
            c.find('li:first').remove();
          }
          if(m.title){
            var badge='badge-danger';
            if(m.status>0) badge='badge-success';
            message='<span class="badge '+badge+'">'+App.text2html(m.title)+'</span> '+App.text2html(m.content);
          }else{
            message=App.text2html(m.content);
          }
          c.append('<li>'+message+'</li>');
          messageCount[m.mode]++;
          break;
          case 'modal':
          var c=$('#notify-modal-'+m.type);
          if(c.length<1){
            var callback='recv_notice_'+m.type;
            if(typeof(window[callback])!='undefined'){
              return window[callback](m);
            }
            if(m.status>0){
              console.info(m.content);
            }else{
              console.error(m.content);
            }
            return;
          }
          if(m.title){
            var badge='badge-danger';
            if(m.status>0) badge='badge-success';
            message='<span class="badge '+badge+'">'+App.text2html(m.title)+'</span> '+App.text2html(m.content);
          }else{
            message=App.text2html(m.content);
          }
          if(!c.data('shown')){
            messageCount[m.mode]=0;
            c.data('shown',true);
            var mbody=c.find('.modal-body'),mbodyUL=mbody.children('ul.modal-body-ul');
            if(mbodyUL.length<1){
              mbody.html('<ul class="modal-body-ul" id="notify-modal-'+m.type+'-container"><li>'+message+'</li></ul>');
            }else{
              mbodyUL.html('<li>'+message+'</li>');
            }
            c.niftyModal('show',{
              afterOpen: function(modal) {},
              afterClose: function(modal) {
                c.data('shown',false);
              }
            });
          }else{
            var cc=$('#notify-modal-'+m.type+'-container');
            if(messageCount[m.mode]>=messageMax[m.mode]){
              cc.find('li:first').remove();
            }
            cc.append('<li>'+message+'</li>');
          }
          messageCount[m.mode]++;
          break;
          case 'notify':
          default:
            if('notify'!=m.mode) m.mode='notify';
            var c=$('#notice-message-container');
            if(c.length<1){
              App.message({title: App.i18n.SYS_INFO, text: '<ul id="notice-message-container" class="no-list-style" style="max-height:500px;overflow-y:auto;overflow-x:hidden"></ul>'});
              c=$('#notice-message-container');
            }
            if(messageCount[m.mode]>=messageMax[m.mode]){
              c.find('li:first').remove();
            }
            if(m.title){
              var badge='badge-danger';
              if(m.status>0) badge='badge-success';
              message='<span class="badge '+badge+'">'+App.text2html(m.title)+'</span>'+App.text2html(m.content);
            }else{
              message=App.text2html(m.content);
            }
            c.append('<li>'+message+'</li>');
            messageCount[m.mode]++;
          break;
        }
      },BACKEND_URL+'/user/notice');
    },
    text2html:function(text){
      return String(text).replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/\n/g,'<br />').replace(/  /g,'&nbsp; ').replace(/\t/g,'&nbsp; &nbsp; ')
    },
    checkedAll:function(ctrl,target){
      return $(target).prop('checked',$(ctrl).prop('checked'));
    },
    attachCheckedAll:function(ctrl,target,showNumElem){
      $(ctrl).on('ifChecked ifUnchecked click',function(){
        App.checkedAll(this,target);
        if(showNumElem) $(showNumElem).text($(target+':checked').length);
      });
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
    insertAtCursor: function(myField, myValue,posStart,posEnd) { 
      if (typeof TextAreaEditor != 'undefined') {
        TextAreaEditor.setSelectText(myField, myValue,posStart,posEnd);
        return;
      } 
      /* IE support */
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
    },
    searchFS:function(elem,size,type,url,before){
      if(size==null)size=10;
      if(url==null)url=BACKEND_URL+'/user/autocomplete_path';
      $(elem).typeahead({
       hint: true, highlight: true, minLength: 1
      }, {source: function (query, sync, async) {
          var data={query:query,size:size,type:type};
          $.ajax({
            url: url,
            type: 'get',
            data: before?before(data):data,
            dataType: 'json',
            async: false,
            success: function (data) {
              var arr = [];
              if(!data.Data) return;
              $.each(data.Data, function (index, val) {
                arr.push(val);
              });
              sync(arr);
            }
          });
      },limit: size});
    },
    randomString:function (len) {
      　　len = len || 32;
      　　var $chars = 'ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678'; //默认去掉了容易混淆的字符oOLl,9gq,Vv,Uu,I1
      　　var maxPos = $chars.length;
      　　var pwd = '';
      　　for (i = 0; i < len; i++) {
      　　　　pwd += $chars.charAt(Math.floor(Math.random() * maxPos));
      　　}
      　　return pwd;
      },
      bottomFloat:function (elems, top, autoWith) {
      	if ($(elems).length<1) return;
        if (top == null) top = 0;
        $(elems).not('[disabled-fixed]').each(function(){
        $(this).attr('disabled-fixed','fixed');
        var elem=this;
      	var _offset = $(elem).height() + top;
      	var offsetY = $(elem).offset().top + _offset;
      	var w = $(elem).outerWidth(), h = $(elem).outerHeight();
      	if(!autoWith) autoWith=$(elem).data('auto-width');
      	if(autoWith) $(elem).css('width', w);
      	$(window).on('scroll',function () {
      		var scrollH = $(this).scrollTop() + $(window).height();
      		if (scrollH >= offsetY) {
      			if ($(elem).hasClass('always-bottom')) {
      				$(elem).removeClass('always-bottom');
      				$(elem).next('.fixed-placeholder').hide();
      			}
      			return;
      		}
      		if (!$(elem).hasClass('always-bottom')) {
      			$(elem).addClass('always-bottom');
      			if ($(elem).next('.fixed-placeholder').length > 0) {
      				$(elem).next('.fixed-placeholder').show();
      			} else {
      				$(elem).after('<div style="width:' + w + 'px;height:' + h + 'px" class="fixed-placeholder"></div>');
      			}
      		}
      	});//on-scroll
        });//end-each
        $(window).trigger('scroll');
      },
      topFloat:function (elems, top, autoWith) {
      	if ($(elems).length<1) return;
        if (top == null) top = 0;
        $(elems).not('[disabled-fixed]').each(function(){
        $(this).attr('disabled-fixed','fixed');
        var elem=this;
      	var offsetY = $(elem).offset().top;
        var w = $(elem).outerWidth(), h = $(elem).outerHeight();
        if(!autoWith) autoWith=$(elem).data('auto-width');
      	if(autoWith) $(elem).css('width', w);
      	$(window).on('scroll',function () {
      		var scrollH = $(this).scrollTop();
      		if (scrollH <= offsetY) {
      			if ($(elem).hasClass('always-top')) {
      				$(elem).removeClass('always-top');
      				$(elem).next('.fixed-placeholder').hide();
      			}
      			return;
      		}
      		if (!$(elem).hasClass('always-top')) {
      			$(elem).addClass('always-top').css('top',top);
      			if ($(elem).next('.fixed-placeholder').length > 0) {
      				$(elem).next('.fixed-placeholder').show();
      			} else {
      				$(elem).after('<div style="width:' + w + 'px;height:' + h + 'px" class="fixed-placeholder"></div>');
      			}
      		}
      	});//on-scroll
        });//end-each
        $(window).trigger('scroll');
      },
      topFloatRawThead:function (elems, top) {
      	if ($(elems).length<1) return;
        if (top==null) top = 0;
        $(elems).not('[disabled-fixed]').each(function(){
        $(this).attr('disabled-fixed','fixed');
        var elem=this,table=$(elem).parent('table');
        var offsetY=$(elem).offset().top,maxOffsetY=table.height()+offsetY-$(elem).outerHeight()*2;
        $(elem).css({'background-color':'white'});
        var setSize=function(init){
          if(init==null) init=false;
          var width=$(elem).outerWidth(),ratio=1;
          if(!init){
            if(Math.abs(table.data('width')-width)>1){//避免抖动
              ratio=width/table.data('width');
            }
            if(table.data('offset-left')!=$(table).offset().left){
            }else if(table.data('scroll-left')!=$(window).scrollLeft()){
              $(elem).css({'left':table.offset().left - $(window).scrollLeft()});
              if(ratio==1) return;
            }
          }
          table.data('width',width);//记录宽度
          table.data('offset-left',$(table).offset().left);//记录左侧偏移
          table.data('scroll-left',$(window).scrollLeft());//记录左侧滚动条
          var cols=table.children('col'),tds=$(elem).find('td,th');
          if(cols.length<1){
            var html='';
            tds.each(function(){
              var w=$(this).outerWidth()*ratio;
              html+='<col style="min-width:'+w+'px;max-width:auto" />';
              $(this).css({'min-width':w,'max-width':'auto'});
            });
            table.prepend(html);
            return;
          }
          tds.each(function(index){
            var col=cols.eq(index);
            var w=$(this).outerWidth()*ratio;
            col.css({'width':w});
            $(this).css({'width':w});
          });
        }
        setSize(true);
      	$(window).on('scroll resize',function () {
          setSize();
          var scrollH=$(this).scrollTop();
          if (scrollH<=offsetY || scrollH>=maxOffsetY) {
      			if ($(elem).hasClass('always-top')) {
              $(elem).removeClass('always-top');
      			}
      			return;
          }
      		if (table.height()>$(window).height() && !$(elem).hasClass('always-top')) {
      			$(elem).addClass('always-top').css('top',top);
      		}
        });
      });
      $(window).trigger('scroll');
      },
      topFloatThead:function (elems, top, clone) {
        if (!clone) return App.topFloatRawThead(elems,top);
      	if ($(elems).length<1) return;
        if (top == null) top = 0;
        $(elems).not('[disabled-fixed]').each(function(){
        $(this).attr('disabled-fixed','fixed');
        var elem=this,table=$(elem).parent('table');
      	var offsetY = $(elem).offset().top,maxOffsetY=table.height()+offsetY-$(elem).outerHeight()*2, cid = $(elem).data('copy');
        if(cid) {
          $('#tableCopy'+cid).remove();
        }else{
          cid = Math.random();
          $(elem).data('copy',cid);
        }
        var eCopy=$('<table class="'+table.attr('class')+' always-top" style="background-color:white" id="tableCopy'+cid+'"></table>');
        var hCopy=$(elem).clone();
        eCopy.append(hCopy);
        var setSize=function(init){
          if(init==null) init=false;
          if(!init){
            if(eCopy.data('offset-left')!=$(elem).offset().left){
            }else if(eCopy.data('scroll-left')!=$(window).scrollLeft()){
              eCopy.css({'left':$(elem).offset().left - $(window).scrollLeft()});
              return;
            }else{
              return;
            }
          }
          eCopy.data('offset-left',$(elem).offset().left);//记录左侧偏移
          eCopy.data('scroll-left',$(window).scrollLeft());//记录左侧滚动条
          var cols=hCopy.find('td,th'),rawCols=$(elem).find('td,th');
          rawCols.each(function(index){
            var col=cols.eq(index);
            col.css('width',$(this).outerWidth());
            if(!init)return;
            var chk=col.find('input:checkbox');
            if(chk.length<1)return;
            var rawChk=rawCols.find('input:checkbox');
            chk.each(function(idx){
              rawChk.eq(idx).on('click change',function(){
                chk.prop('checked',$(this).prop('checked'));
              });
            });
          });
          var offsetX = $(elem).offset().left - $(window).scrollLeft();
          var w = $(elem).outerWidth(), h = $(elem).outerHeight()
          eCopy.css({'top':top,'left':offsetX,'width':w,'height':h});
        }
        setSize(true);
        eCopy.hide();
        table.after(eCopy);
      	$(window).on('scroll',function () {
          setSize();
          var scrollH = $(this).scrollTop();
          if (scrollH <= offsetY||scrollH>=maxOffsetY) return eCopy.hide();
      		eCopy.show();
        });
      });
      $(window).trigger('scroll');
    },
    getImgNaturalDimensions:function (oImg, callback) {
      if (!oImg.naturalWidth) { // 现代浏览器
      　　callback({w: oImg.naturalWidth, h:oImg.naturalHeight});
        return;
      } 
      // IE6/7/8
      var nImg = new Image();
      nImg.onload = function() {
        callback({w: nImg.width, h:nImg.height});
      }
      nImg.src = oImg.src;
    },
    reportBug:function(url){
	    $.post(url,{"panic":$('#panic-content').html(),"url":window.location.href},function(r){},'json');
    },
    replaceURLParam:function(name,value,url){
      if(url==null) url=window.location.href;
      value=encodeURIComponent(value);
      var pos=String(url).indexOf('?');
      if(pos<0) return url+'?'+name+'='+value;
      var q=url.substring(pos),r=new RegExp('([\\?&]'+name+'=)[^&]*(&|$)');
      if(!r.test(q)) return url+'&'+name+'='+value;
      url=url.substring(0,pos);
	    q=q.replace(r,'$1'+value+'$2');
	    return url+q;
    },
    switchLang:function(lang){
      window.location=App.replaceURLParam('lang',lang);
    },
    extends:function(child,parent){
      //parent.call(this);
      var obj = function(){};
      obj.prototype = parent.prototype;
      child.prototype = new obj();
      child.prototype.constructor = child;
    },
    formatBytes:function(cellValue,precision) {
      if(precision==null) precision=1
      var intVal = parseInt(cellValue),ras = " B";
      if (intVal > 1024) {
          intVal /= 1024
          ras = " KB"
      }
      if (intVal > 1024) {
          intVal /= 1024
          ras = " MB"
      }
      if (intVal > 1024) {
          intVal /= 1024
          ras = " GB"
      }
      if (intVal > 1024) {
          intVal /= 1024
          ras = " TB"
      }
      if (intVal > 1024) {
          intVal /= 1024
          ras = " PB"
      }
      return (intVal).toFixed(precision) + ras;
    },
    format:function(raw,xy,formatters){
      if(formatters && typeof(formatters[xy.y])=='function') return formatters[xy.y](raw,xy.x);
      return raw;
    },
    genTable:function(rows,formatters){
      var h='<table class="table table-bordered no-margin">';
      var th='<thead>',bd='<tbody>';
      for(var i=0;i<rows.length;i++){
        var v=rows[i];
        if(i==0){
          for(var k in v) th+='<th class="'+k+'"><strong>'+k+'</strong></th>';
        }
        bd+='<tr>';
        for(var k in v) bd+='<td class="'+k+'">'+App.format(v[k],{'x':i,'y':k},formatters)+'</td>';
        bd+='</tr>';
      }
      th+='</thead>';
      bd+='</tbody>';
      h+=th+bd+'</table>';
      return h;
    },
    httpStatusColor:function(code){
      if(code>=500)return 'danger';
      if(code>=400)return 'warning';
      if(code>=300)return 'info';
      return 'success';
    },
    logShow:function(elem,trigger,pipe){
      if(!$('#log-show-modal').data('init')){
      $('#log-show-modal').data('init',true);
      $(window).off().on('resize',function(){
        $('#log-show-modal').css({height:$(window).height(),width:'100%', 'max-width':'100%',left:0,top:0,transform:'none'});
		    $('#log-show-modal').find('.md-content').css('height',$(window).height());
		    $('#log-show-content').css('height',$(window).height()-200);
      });
      $('#log-show-last-lines').on('change',function(r){
        var target=$(this).data('target');
        if(!target)return;
        var lastLines=$(this).val();
        target.data('last-lines',lastLines);
        target.trigger('click');
      });
      $('#log-show-modal .modal-footer .btn-refresh').on('click',function(r){
        var target=$('#log-show-last-lines').data('target');
        if(!target)return;
        target.trigger('click');
      });
      $(window).trigger('resize');
      }
      if(pipe==null) pipe='';
      var done=function(a){
	    	var url=$(a).data('url');
	    	var lastLines=$(a).data('last-lines');
        if(lastLines==null) lastLines=100;
        $('#log-show-last-lines').data('target',$(a));
        var contentID='log-show-content',contentE='#'+contentID;
        $('#log-show-modal').niftyModal('show',{
          afterOpen: function(modal) {
	    			$.get(url,{lastLines:lastLines,pipe:pipe},function(r){
	    				if(r.Code==1){
                var subTitle=$('#log-show-modal .modal-header .modal-subtitle');
                if(typeof(r.Data.title)!='undefined'){
                  if(r.Data.title)r.Data.title=' ('+r.Data.title+')';
                  subTitle.html(r.Data.title);
                }else{
                  subTitle.empty();
                }
                if(typeof(r.Data.list)!='undefined'){
                  var h='<div class="table-responsive" id="'+contentID+'">'+App.genTable(r.Data.list,{'StatusCode':function(raw,index){
                    return '<span class="label label-'+App.httpStatusColor(raw)+'">'+raw+'</span>';
                  }})+'</div>';
                  $(contentE).parent('.modal-body').css('padding',0);
                  $(contentE).replaceWith(h);
                }else{
                  if($(contentE)[0].tagName.toUpperCase()!='TEXTAREA'){
                    $(contentE).replaceWith("<textarea name='content' class='form-control' id='"+contentID+"'></textarea>");
                  }
                  $(contentE).text(r.Data.content);
                }
                $(window).trigger('resize');
	    				  var textarea=$(contentE)[0];
	    				  textarea.scrollTop = textarea.scrollHeight;
	    				}else{
	    					$(contentE).text(r.Info);
	    				}
	    			},'json');
          },
          afterClose: function(modal) {}
        });
      };
      if(trigger) return done(elem);
	    $(elem).on('click',function(){
        done(this);
      });
    },
    tableSorting:function(table){
      table=table==null?'':table+' ';
       $(table+'[sort-current!=""]').each(function(){//<thead sort-current="created">
        var current=$(this).attr('sort-current');
        var sortObj=$(this).find('[sort="'+current+'"]');//<th sort="-created">
        var newCls='fa-arrow-down',oldCls='fa-arrow-up';
        if(sortObj.length<1){
          newCls='fa-arrow-up';
          oldCls='fa-arrow-down';
          sortObj=$(this).find('[sort="-'+current+'"]');
        }
        if(sortObj.length>0){
          var icon=sortObj.children('.fa');
          if(icon.length<1){
            sortObj.append('<i class="fa '+newCls+'"></i>');
          }else{
            icon.removeClass(oldCls).addClass(newCls);
          }
        }
       $(table+'[sort-current] [sort]').css('cursor','pointer').on('click',function(e){
         var thead=$(this).parents('[sort-current]');
         var current=thead.attr('sort-current');
         var url=thead.attr('sort-url')||window.location.href;
         var trigger=thead.attr('sort-trigger');
         var sort=$(this).attr('sort');
         if(current && (current==sort || current=='-'+sort)){
          var reg=/^\-/;
          current=reg.test(current)?current.replace(reg,''):'-'+current;
         }else{
          current=sort;
         }
         url=App.replaceURLParam('sort',current,url);
         if(trigger){
           thead.trigger('sort');
         }else{
           var setto=thead.attr('sort-setto');
           if(setto){
             $(setto).load(url);
           }else{
             window.location=url;
           }
         }
       });
      });
    },
    resizeModalHeight:function(el){
      var h=$(window).height()-200;
      if(h<200)h=200;
      var bh=h-150;
      $(el).css({"max-height":h+'px'});
      $(el).find('.modal-body').css({"max-height":bh+'px'});
    },
    switchStatus:function(a,type,editURL){
      if(type==null)type=$(a).data('type');
      if(editURL==null)editURL=$(a).data('url');
      var that=$(a),status=that.data(type)=='Y'?'N':'Y',data={id:that.data('id')};
      var v=that.val();
      data[type]=status;
      if(editURL.charAt(0)!='/') editURL=BACKEND_URL+'/'+editURL;
	    $.get(editURL,data,function(r){
	    	if(r.Code==1){
	    		that.data(type,status);
	    		that.prop('checked',status==v);
	    	}
	    	App.message({title:App.i18n.SYS_INFO,text:r.Info,time:5000,sticky:false,   class_name:r.Code==1?'success':'error'});
	    },'json');
    },
    bindSwitch:function(elem,eventName,editURL,type){
      if(eventName==null)eventName='click';
      var re = new RegExp('switch-([\\w\\d]+)');
	    $(elem).on(eventName,function(){
        if(type==null){
          var matches=String($(this).attr('class')).match(re);
          type=matches[1];
        }
        App.switchStatus(this,type,editURL);
      });
    }
  };
 
}();

$(function(){
  //$("body").animate({opacity:1,'margin-left':0},500);
  $("body").css({opacity:1,'margin-left':0});
});
