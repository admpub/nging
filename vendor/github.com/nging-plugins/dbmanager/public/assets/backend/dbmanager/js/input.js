function fontWidth(v){
    var t=$('<span style="visibility:hidden;white-space:nowrap;font-size:15px">'+v+'</span>');
    $('body').append(t);
    var w=t[0].offsetWidth;
    t.remove();
    return w;
}
function textInputFloat(a,title,alwayShow,offsetY){
    var v = $(a).val(), w = $(a).width();
    if(v.length<1||!$(a).is(':visible'))return;
    var sizeW=fontWidth(v)+4,minW=v.length*15;
    if(sizeW<minW)sizeW=minW;
    if(sizeW<w)return;
    if(title)sizeW+=80;
    if(!alwayShow)$(a).hide();
    if(!offsetY)offsetY=0;
    var id=$(a).attr('name').replace(/[\[\]]+/g,'_')+'_textinput';
    var te=$(a).next('#'+id);
    if(te.length>0){
        te.show();
        te.find('input').val(v);
        te.find('input').focus();
        return;
    }
    var width=200;
    if(sizeW>width){
        width=sizeW;
        if(width>$(window).width()){
            width='100%';
        }else{
            width=width+'px';
        }
    }
    te=$('<span id="'+id+'" class="input-group floatup-input-layer" style="position:absolute;z-index:10;width:'+width+';box-shadow:1px 1px 5px #333;margin-top:'+offsetY+'px">'+(title?'<span class="input-group-addon">'+title+'</span>':'')+'<input class="form-control" value="'+v+'"></span>');
    var sb=function(event){
        var tv=$(this).val();
        $(a).attr('value',tv).val(tv);
        te.hide();
        if(!alwayShow)$(a).show();
    };
    te.find('input').on('blur',sb).on('keyup',function(event){
        if(event.keyCode!=13)return;
        sb.call(this,event);
    });
    $(a).after(te);
    te.find('input').focus();
    if(te.offset().left+te.width()>$(window).width()){
        te.css('right',0);
    }
}