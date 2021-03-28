var dropzone,dropzoneZIP,editor;
function resetCheckedbox() {
    $('#checkedAll:checked').prop('checked', false);
    $('#tbody-content input[type=checkbox][name="path[]"]:checked').prop('checked', false);
}
function refreshList() {
    if($('#tbody-content').length<1){
        window.location.reload();
        return;
    }
    App.loading('show');
    $.get(window.location.href,{'_pjax':'tbody-content'},function(r){
        var e=$(r);
        $('#tbody-content').html(e.find('#tbody-content').html());
        App.float('#tbody-content img.previewable');
        App.loading('hide');
        $('#tbody-content').trigger('refresh');
        resetCheckedbox();
    },'html');
}
function initCodeMirrorEditor() {
    editor = CodeMirror.fromTextArea($("#file-edit-content")[0], {
      lineNumbers: true,
      theme: "night",
      extraKeys: {
        "F11": function(cm) {
          cm.setOption("fullScreen", !cm.getOption("fullScreen"));
        },
        "Esc": function(cm) {
          if (cm.getOption("fullScreen")) cm.setOption("fullScreen", false);
        }
      }
    });
    editor.setOption('lineWrapping', true);
    editor.setSize('auto', 'auto');
    $('#file-edit-modal .modal-footer .btn-success').on('click',function(){
        var url=$(this).data('url');
        var enc=$('#use-encoding-open').val();
        if(!enc)enc='';
        $.post(url,{content:editor.getValue(),encoding:enc},function(r){
            if(r.Code!=1)return App.message({title: App.i18n.SYS_INFO, text: r.Info},false);
            return App.message({title: App.i18n.SYS_INFO, text: App.i18n.SAVE_SUCCEED},false);
        },'json');
    });
    $('#file-edit-modal .modal-body').css('padding',0);
    $('#use-encoding-open').on('change',function(){
        var enc=$(this).val();
        fileReopen(enc);
    });
    
}

function fileReopen(encoding,url) {
    App.loading('show');
    if(url==null)url=$('#file-edit-modal .modal-footer .btn-success').data('url');
    $.get(url,{encoding:encoding},function(r){
        App.loading('hide');
        if(r.Code!=1)return App.message({title: App.i18n.SYS_INFO, text: r.Info},false);
        editor.setValue(r.Data);
    },'json');
}

function fileEdit(obj,file) {
    var url=$(obj).data('url');
    $('#file-edit-modal .modal-footer .btn-success').data('url',url);
    App.loading('show');
    $.get(url,{},function(r){
        App.loading('hide');
        if(r.Code!=1)return App.message({title: App.i18n.SYS_INFO, text: r.Info},false);
        $('#file-edit-modal .modal-header h3').html(App.i18n.EDIT_FILE+': '+file);
        $('#file-edit-modal').niftyModal('show',{
            afterOpen: function(modal) {
                editor.setValue(r.Data);
                codeMirrorChangeMode(editor,file);
            },
            afterClose: function(modal) {
                $('#use-encoding-open').find('option:selected').prop('selected',false);
            }
        });
    },'json');
}

function fileRename(obj,file,isDir) {
    var url=$(obj).data('url');
    $('#file-rename-modal .modal-footer .btn-primary:last').data('url',url);
    $('#file-rename-modal .modal-header h3').html((isDir ? App.i18n.MODIFY_DIRNAME : App.i18n.MODIFY_FILENAME)+': '+file);
    $('#file-rename-modal').niftyModal('show',{afterOpen:function(modal){
        $('#file-rename-input').val(file);
    }});
}

function fileMkdir(obj) {
    var url=$(obj).data('url');
    $('#file-mkdir-modal .modal-footer .btn-primary:last').data('url',url);
    $('#file-mkdir-modal .modal-header h3').html(App.i18n.CREATE_DIR);
    $('#file-mkdir-modal').niftyModal('show',{afterOpen:function(modal){
        $('#file-mkdir-input').val('');
        $('#file-mkdir-input').focus();
    }});
}

function filePlay(obj,playlist){
    if(playlist==null) playlist='a[playable]';
    var url=$(obj).data('url'),i=$(playlist).index($(obj)),
    fileName=$(obj).data('name'),
    mime=$(obj).data('mime'),
    player,id='file-play-video',
    headTitle=$('#file-play-modal .modal-header h3');
    headTitle.html(App.i18n.PLAY+': '+fileName);
    var body=$('#file-play-modal .modal-body');
    body.css({'padding':'0','text-align':'center'});
    $(obj).css({'color':'yellow'});
    $('#file-play-modal').niftyModal('show',{afterOpen:function(modal){
        /*
        if(String(url.split('.').pop()).toLowerCase()=='ts'){
            //transferTS(url,'#file-play-video');
            url=BACKEND_URL+'/ts2m3u8?ts='+encodeURIComponent(url);
            mime='application/x-mpegURL';
        }
        */
        $('#file-play-video source').attr('src',url).attr('type',mime);
        $(window).trigger('resize');
        player = videojs(id, null, function(){
	        this.on('ended',function(){
	    	    i++;
	            if(i >= $(playlist).length) i = 0;
                var ve=$(playlist).eq(i);
                ve.css({'color':'yellow'});
	            this.src({type: ve.data('mime'), src: ve.data('url')});
	            this.play();
                headTitle.html(App.i18n.PLAY+': '+ve.data('name'));
	        });
        });
        player.play();
    },afterClose:function(){
        if(!player) return;
        var c=$('<video-js id="file-play-video" class="vjs-default-skin" width=500 height=500 controls><source src="" type=""></video-js>');
        player.dispose();
        player=null;
        body.html(c);
    }});
}

function codeMirrorChangeMode(editor,val) {
  var m, mode, spec;
  if (m = /.+\.([^.]+)$/.exec(val)) {
    var info = CodeMirror.findModeByExtension(m[1]);
    if (info) {
      mode = info.mode;
      spec = info.mime;
    }
  } else if (/\//.test(val)) {
    var info = CodeMirror.findModeByMIME(val);
    if (info) {
      mode = info.mode;
      spec = val;
    }
  } else {
    mode = spec = val;
  }
  if (mode) {
    editor.setOption("mode", spec);
    CodeMirror.autoLoadMode(editor, mode);
  } else {
    console.log("Could not find a mode corresponding to " + val);
  }
}
Dropzone.autoDiscover=false;
CodeMirror.modeURL = ASSETS_URL+"/js/editor/markdown/lib/codemirror/mode/%N/%N.js";
function dropzoneResizeHeight(isZip){
  return function(){
    var el=isZip?'#multi-upload-zip-modal':'#multi-upload-modal';
    App.resizeModalHeight(el);
  }
}
$(function(){
    initDropzone($.extend({
        chunking:true,
        parallelChunkUploads:true,
        retryChunksLimit:3,
        retryChunks:true,
        maxFilesize:1024 // 文件最大尺寸(MB)
    },window.dropzoneOptions||{}));
    dropzone=$('#multi-upload-dropzone').get(0).dropzone;
    dropzoneZIP=$('#multi-upload-zip-dropzone').length>0?$('#multi-upload-zip-dropzone').get(0).dropzone:null;
    dropzone.on('addedfiles',dropzoneResizeHeight(false));
    if(dropzoneZIP)dropzoneZIP.on('addedfiles',dropzoneResizeHeight(true));
    initCodeMirrorEditor();
	$('#uploadBtn').off().on('click',function(event){
		$('#multi-upload-modal').niftyModal('show',{
            closeOnClickOverlay:false,
            afterClose:function(modal){
                dropzone.removeAllFiles();
                refreshList();
            }
        });
    });
    if($('#uploadZipBtn').length>0){
        $('#uploadZipBtn').off().on('click',function(event){
            $('#multi-upload-zip-modal').niftyModal('show',{
                closeOnClickOverlay:false,
                afterClose:function(modal){
                    if(dropzoneZIP)dropzoneZIP.removeAllFiles();
                    refreshList();
                }
            });
        });
    }
    $(window).off().on('resize',function(){
        $('#file-edit-modal,#file-play-modal').css({height:$(window).height(),width:'100%','max-width':'100%',left:0,top:0,transform:'none'});
        $('#file-edit-form,#file-play-video').css({height:$(window).height()-150,width:'100%','max-width':'100%',overflow:'auto'});
        dropzoneResizeHeight(false)();
        if(dropzoneZIP)dropzoneResizeHeight(true)();
    });
    $(window).trigger('resize');
    $('#file-rename-modal .modal-footer .btn-primary:last').off().on('click',function(){
        var url=$(this).data('url');
        App.loading('show');
        $.post(url,{name:$('#file-rename-input').val()},function(r){
            App.loading('hide');
            if(r.Code!=1)return App.message({title: App.i18n.SYS_INFO, text: r.Info},false);
            App.message({title: App.i18n.SYS_INFO, text: App.i18n.SAVE_SUCCEED},false);
            refreshList();
        },'json');
    });
    $('#file-mkdir-modal .modal-footer .btn-primary:last').off().on('click',function(){
        var url=$(this).data('url');
        App.loading('show');
        $.post(url,{name:$('#file-mkdir-input').val()},function(r){
            App.loading('hide');
            if(r.Code!=1)return App.message({title: App.i18n.SYS_INFO, text: r.Info},false);
            App.message({title: App.i18n.SYS_INFO, text: App.i18n.CREATE_SUCCEED},false);
            refreshList();
        },'json');
    });
    $('#query-current-path').on('keyup',function(){
        var q=$(this).val();
        if(q==''){
            $('#tbody-content').children('tr:not(:visible)').show();
            var disabledBoxies=$('#tbody-content input[type=checkbox][name="path[]"]:disabled');
            if(disabledBoxies.length>0)$('#checkedAll').prop('checked',false);
            disabledBoxies.prop('disabled',false);
            return;
        }
        $('#tbody-content').children('tr:not([item*="'+q+'"])').hide().find('input[type=checkbox][name="path[]"]').prop('disabled',true);
        $('#tbody-content').children('tr[item*="'+q+'"]:not(:visible)').show().find('input[type=checkbox][name="path[]"]:disabled').prop('disabled',false);
        $('#tbody-content input[type=checkbox][name="path[]"]:disabled:checked').prop('checked',false);
        $('#checkedAll').prop('checked',$('#tbody-content tr[item]:visible input[type=checkbox][name="path[]"]:checked').length==$('#tbody-content tr[item]:visible input[type=checkbox][name="path[]"]'));
        if(event.keyCode==13){
            var tr=$('#tbody-content').children('tr:visible');
            if(tr.length==1){
                var a=tr.children('td:first').children('a:first');
                var url=a.attr('href');
                window.location=url;
                return;
            }
        }
    }).focus();
    $('#btn-query-current-path').on('click',function(){
        $('#query-current-path').trigger('keyup');
    });
    App.float('#tbody-content img.previewable');
    resetCheckedbox();
    App.attachCheckedAll('#checkedAll','#tbody-content input[type=checkbox][name="path[]"]');
});