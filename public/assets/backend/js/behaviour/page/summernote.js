function summernoteEditor(elem,minHeight){
  if(minHeight==null) minHeight=400;
  $(elem).summernote({lang:App.langTag(),
    minHeight:minHeight,
    callbacks:{
      onImageUpload:function(files, editor, $editable) {
        var $files = $(files);
        $files.each(function() {
        var file = this;
        var formdata = new FormData();  
        formdata.append("files[]", file);  
        $.ajax({  
            data : formdata,  
            type : "POST",  
            url : $(elem).attr('action'),
            cache : false,  
            contentType : false,  
            processData : false,  
            dataType : "json",  
            success: function(r) {
              if(r.Code!=1){
                return App.message({title:App.i18n.SYS_INFO,text:r.Info,time:5000,sticky:false,class_name:r.Code==1?'success':'error'});
              }
              $.each(r.Data.files, function (index, file) {
                $(elem).summernote('insertImage', file, function($image) {});
              });
            },
            error:function(){  
              alert(App.i18n.UPLOAD_ERR);  
            }  
        });
       });
      }
    }
  });
}