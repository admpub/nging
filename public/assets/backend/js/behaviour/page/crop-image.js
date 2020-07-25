
function getCropServerURL(){
  if(typeof(window.CropServerURL)!='undefined' && window.CropServerURL){
    return window.CropServerURL;
  }
  if(typeof(window.IS_BACKEND)!='undefined' && window.IS_BACKEND){
    return BACKEND_URL+'/manager/crop';
  }
  return FRONTEND_URL+'/user/file/crop';
}
function cropImage(uploadURL,thumbsnailElem,originalElem,type,width,height){
  var jcrop=null;
  if(!type) {
    type=uploadURL.split('?',2)[0];
    type=type.substring(type.lastIndexOf('/')+1);
  }
  if(width==null) width=200;
  if(height==null) height=width;
  $('#fileupload').fileupload({
      url: uploadURL,
      dataType: 'json',
      done: function (e, data) {
        var r=data.result;
        if(r.Code!=1){
          return App.message({title:App.i18n.SYS_INFO,text:r.Info,time:5000,sticky:false,class_name:r.Code==1?'success':'error'});
        }
        var afterClose=function(){
          jcrop.destroy();
          $(".crop-image").empty();
          jcrop=null;
        };
        if(jcrop) afterClose();
        $.each(r.Data.files, function (index, file) {
          $(originalElem).val(file);
          $(".crop-image").html('<img src="' + file + '?_='+Math.random()+'" />');
          $('#progress').fadeOut();
          $("#save-image").data('image-file',file);
          //Crop Image
          var img=$(".crop-image img");
          img.off().on('load',function(){
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
          });
        });
        if(typeof($.fn.niftyModal)!='undefined'){
          $("#crop-image").niftyModal('show',{afterClose:afterClose});
        }else{
          $("#crop-image").modal('show');
          $('#crop-image').off('hide.bs.modal').on('hide.bs.modal', function(){
            afterClose();
          });
        }
      },
      progressall: function (e, data) {
          var progress = parseInt(data.loaded / data.total * 100, 10);
          $('#progress').fadeIn();
          $('#progress').css('width',progress + '%');
      }
  }).prop('disabled', !$.support.fileInput).parent().addClass($.support.fileInput ? undefined : 'disabled');

  $("#save-image").on('click',function(){
    var self=$(this),img=$(".crop-image img");
    App.getImgNaturalDimensions(img[0],function(natural){
    var c = jcrop.tellSelect();
    if( c.w != 0 ){
      var ratio=natural.w/img.width();
      var timg=self.data('image-file');
      $.get(getCropServerURL(),{
        src:timg,
        x:c.x*ratio,
        y:c.y*ratio,
        w:c.w*ratio,
        h:c.h*ratio,
        type:type,
        size:width+'x'+height
      },function(r){
        if(r.Code!=1){
          return App.message({title:App.i18n.SYS_INFO,text:r.Info,time:5000,sticky:false,class_name:r.Code==1?'success':'error'});
        }
        $(".profile-avatar").attr("src", r.Data+'?_='+Math.random());
        $(thumbsnailElem).val(r.Data);
        if(typeof($.fn.niftyModal)!='undefined'){
          $("#crop-image").niftyModal('hide');
        }else{
          $("#crop-image").modal('hide');
        }
        jcrop.destroy();
        $(".crop-image").empty();
        jcrop=null;
      },'json');
    }else{
      alert("Please select a crop region then press save.");
    }
    });
  });
  $('.avatar-upload .avatar-remove').on('click',function(){
    $(".profile-avatar").attr("src", "");
    $(originalElem).val("");
    $(thumbsnailElem).val("");
  });
}