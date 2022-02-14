
function getCropServerURL(){
  if(typeof(window.CropServerURL)!='undefined' && window.CropServerURL){
    return window.CropServerURL;
  }
  if(typeof(window.IS_BACKEND)!='undefined' && window.IS_BACKEND){
    return BACKEND_URL+'/manager/crop';
  }
  return FRONTEND_URL+'/user/file/crop';
}
function cropImage(uploadURL,thumbsnailInput,originalInput,subdir,width,height){
  var options = {
    uploadURL:uploadURL,
    croperURL:getCropServerURL(),
    fileElem:null,
    thumbsnailInput:thumbsnailInput,
    originalInput:originalInput,
    previewElem:null,
    subdir:subdir,
    width:width,
    height:height,
    prefix:''
  };
  if(typeof(uploadURL)=='object') {
    options = $.extend(options,uploadURL);
  }
  if(!options.fileElem) {
    if(options.prefix){
      options.fileElem = '#'+options.prefix+'-fileupload';
    }else{
      options.fileElem = '#fileupload';
    }
  }
  if(!options.prefix){
    options.prefix = 'noprefix';
  }
  if(!options.thumbsnailInput) {
    options.thumbsnailInput = $(options.fileElem).data('thumbsnail-input')||'#'+options.prefix+'-image';
  }
  if(!options.originalInput) {
    options.originalInput = $(options.fileElem).data('original-input')||'#'+options.prefix+'-image-original';
  }
  if(!options.previewElem) {
    options.previewElem = $(options.fileElem).data('preview')||'';
    if(!options.previewElem) options.previewElem=$(options.fileElem).siblings('img')[0];
  }
  if(!options.subdir) {
    options.subdir = $(options.fileElem).data('subdir')||'';
    if(!options.subdir){
      var matched = options.uploadURL.match(/subdir=([^&]+)/);
      if(matched && matched.length>0) options.subdir=matched[1];
    }
  }
  //alert(options.subdir)
  if(options.width==null) {
    options.width=$(options.fileElem).data('width')||200;
  }
  if(options.height==null) {
    options.height=$(options.fileElem).data('height')||options.width;
  }
  //console.dir(options);
  var modal=$('#'+options.prefix+'-crop-image'),
      progress=$('#'+options.prefix+'-progress'),
      saveBtn=$('#'+options.prefix+'-save-image'),
      width=options.width,height=options.height;

  var jcrop=null;
  var imageSet=function(img){
      App.getImgNaturalDimensions(img[0],function(natural){
        var ratio=1;
        if(natural.w>0){
          ratio=img.width()/natural.w;
        }
        var w=width*ratio,h=height*ratio;
        img.Jcrop({
          aspectRatio:width/height,
          //minSize:[w,h],
          setSelect:[0,0,w,h],
        },function(){
          jcrop=this;
        });
      });
  };
  var crop=function(fileList,type){
    var afterClose=function(){
      if(jcrop) jcrop.destroy();
      modal.find(".crop-image").empty();
      jcrop=null;
    };
    if(jcrop) afterClose();
    $.each(fileList, function (index, file) {
      $(options.originalInput).val(file);
      modal.find(".crop-image").html('<img src="' + file + '?_='+Math.random()+'" />');
      progress.fadeOut();
      saveBtn.data('image-file',file);
      var img=modal.find(".crop-image img");
      img.off().on('load',function(){
        imageSet(img);
      });
    });
    if(typeof($.fn.niftyModal)!='undefined'){
      modal.niftyModal('show',{afterClose:afterClose});
    }else{
      modal.modal('show');
      modal.off('hide.bs.modal').on('hide.bs.modal', function(){
        afterClose();
      });
    }
  };
  $(options.fileElem).fileupload({
      url: options.uploadURL,
      dataType: 'json',
      done: function (e, data) {
        var r=data.result;
        if(r.Code!=1){
          return App.message({title:App.i18n.SYS_INFO,text:r.Info,time:5000,sticky:false,class_name:r.Code==1?'success':'error'});
        }
        crop(r.Data.files);
      },
      progressall: function (e, data) {
          var progressValue = parseInt(data.loaded / data.total * 100, 10);
          progress.fadeIn();
          progress.css('width',progressValue + '%');
      }
  }).prop('disabled', !$.support.fileInput).parent().addClass($.support.fileInput ? undefined : 'disabled');

  var actions=$(options.fileElem).siblings('.avatar-actions');
  if(actions.length<1){
    actions=$('<div class="avatar-actions"></div');
    $(options.fileElem).parent('div').prepend(actions);
  }

  saveBtn.on('click',function(){
    var self=$(this),img=modal.find(".crop-image img");
    App.getImgNaturalDimensions(img[0],function(natural){
    var c = jcrop.tellSelect();
    if( c.w != 0 ){
      var ratio=natural.w/img.width();
      var timg=self.data('image-file');
      var token=actions.children('.avatar-resize').data('token');
      var data={
        src:timg,
        x:c.x*ratio,
        y:c.y*ratio,
        w:c.w*ratio,
        h:c.h*ratio,
        size:width+'x'+height
      };
      if(token!==undefined && token) data.token = token;
      $.get(options.croperURL,data,function(r){
        if(r.Code!=1){
          return App.message({title:App.i18n.SYS_INFO,text:r.Info,time:5000,sticky:false,class_name:r.Code==1?'success':'error'});
        }
        $(options.previewElem).attr("src", r.Data+'?_='+Math.random());
        $(options.thumbsnailInput).val(r.Data);
        if(typeof($.fn.niftyModal)!='undefined'){
          modal.niftyModal('hide');
        }else{
          modal.modal('hide');
        }
        if(jcrop) jcrop.destroy();
        modal.find(".crop-image").empty();
        jcrop=null;
      },'json');
    }else{
      alert("Please select a crop region then press save.");
    }
    });
  });
  var removeBtn=actions.children('.avatar-remove');
  if(removeBtn.length<1){
    removeBtn=$('<a class="label label-danger avatar-remove" href="javascript:;" title="'+ App.t('删除图片')+'"><i class="fa fa-times"></i></a>');
    actions.append(removeBtn);
  }
  removeBtn.on('click',function(){
    //if(!confirm(App.t('确定删除封面图吗？'))) return;
    $(options.previewElem).attr("src", ASSETS_URL+"/images/user_128.png");
    $(options.originalInput).val("");
    $(options.thumbsnailInput).val("");
    actions.children('.avatar-resize').hide();
  });
  if(typeof(App.editor)=="undefined"){
    return;
  }
  var browsing=$('<a class="label label-primary avatar-browsing" href="javascript:;" title="'+App.t('从服务器选择')+'"><i class="fa fa-folder-open-o"></i></a>');
  actions.append(browsing);
  var resizeBtn=actions.children('.avatar-resize');
  if(resizeBtn.length<1){
    resizeBtn=$('<a class="label label-warning avatar-resize" href="javascript:;" title="'+ App.t('裁剪图片')+'" style="display:none"><i class="fa fa-crop"></i></a>');
    actions.append(resizeBtn);
  }
  browsing.on('click',function(){
    App.editor.finderDialog(App.editor.browsingFileURL + '?from=parent&size=12&filetype=image&subdir='+options.subdir+'&multiple=0', function (fileList,infoList) {
        if (fileList.length <= 0) {
          return App.message({ type: 'error', text: App.t('没有选择任何选项！') });
        }
        $.post(options.uploadURL,{
          pipe:'_queryThumb',
          file:fileList[0],
          size:width+'x'+height
        },function(r){
          if(r.Code!=1) return App.message({ type: 'error', text: r.Info });
          if('thumb' in r.Data) {
            var thumb = r.Data.thumb;
            $(options.previewElem).attr("src", thumb+'?_='+Math.random());
            $(options.thumbsnailInput).val(thumb);
            $(options.originalInput).val(fileList[0]);
            actions.children('.avatar-resize').data('token',r.Data.token).show();
            return;
          }
          crop([fileList[0]]);
        },'json');
    }, 100000);
  });
  resizeBtn.on('click',function(){
    var file = $(options.originalInput).val();
    if (!file) return App.message({ type: 'error', text: App.t('请先选择图片！') });
    crop([file]);
  });
}