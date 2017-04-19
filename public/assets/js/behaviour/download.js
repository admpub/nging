function FormatProgressBar(cellValue) {
    var intVal = parseInt(cellValue);
    var cellHtml = '<div class="progress"><div class="progress-bar" style="width: ' + intVal + '%"></div></div>'
    return cellHtml;
}
function FormatByte(cellValue) {
    var intVal = parseInt(cellValue);
    var ras = " B"
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
    var cellHtml = (intVal).toFixed(1) + ras;
    return cellHtml;
}
function StateIcon(state) {
    var c,t;
    switch(state){
        case 'Completed':
        c='ok-circle text-success';
        break;
        case 'Running':
        c='play-circle text-info';
        break;
        case 'Stopped':
        c='ban-circle text-warning';
        break;
        case 'Failed':
        c='remove-circle text-danger'
        break;
        default:
        state='Stopped';
        c='ban-circle text-warning';
    }
    t=states[state];
    return '<span class="glyphicon glyphicon-'+c+'" title="'+t+'"></span>';
}
function FormatSpeedByte(cellValue) {
    var intVal = parseInt(cellValue);
    var ras = " B/s"
    if (intVal > 1024) {
        intVal /= 1024
        ras = " KB/s"
    }
    if (intVal > 1024) {
        intVal /= 1024
        ras = " MB/s"
    }
    if (intVal > 1024) {
        intVal /= 1024
        ras = " GB/s"
    }
    if (intVal > 1024) {
        intVal /= 1024
        ras = " TB/s"
    }
    var cellHtml = (intVal).toFixed(1) + ras;
    return cellHtml;
}
var downloadWS,downloadWSInterval,downloadAPIPrefix='/download';
function connectSockJS(onopen,onmessage){
	if (downloadWS) {
		if(onopen!=null)onopen();
		return false;
	}
	downloadWS = new SockJS(downloadAPIPrefix+'/progress');
	downloadWS.onopen    = function(){
		if(onopen!=null)onopen();
	};
	downloadWS.onclose   = function(){
        downloadWS = null;
    };
	downloadWS.onmessage = function(msg){
		if(onmessage!=null)onmessage(msg.data);
	};
}

function sockJSConnect(){
    var tmpl = $('#tr-template').html();
    connectSockJS(function(){
        downloadWS.send("progress");
    },function(r){
        var rows=JSON.parse(r);
        var total = 100*rows.length, finished = 0;
        var content = '';
        var checkedAll = $('#fileTable .allCheck').prop('checked');
        for (var i = 0; i < rows.length; i++){
            var v = rows[i];
            finished=finished+v.Progress;
            if($('#id-'+v.Id).length>0){
                $('#downed-'+v.Id).text(FormatByte(v.Downloaded));
                $('#percent-'+v.Id).text(v.Progress);
                $('#speed-'+v.Id).text(FormatSpeedByte(v.Speed));
                $('#progress-'+v.Id).html(FormatProgressBar(v.Progress));
                $('#state-'+v.Id).html(StateIcon(v.State));
                continue;
            }
            var tmplCopy=tmpl;
            for(var j in v){
                var re=new RegExp('\\{'+j+'\\}','g');
                var vl=v[j];
                switch(j){
                    case 'Downloaded':
                    vl=FormatByte(vl);
                    break;
                    case 'Size':
                    vl=FormatByte(vl);
                    break;
                    case 'Speed':
                    vl=FormatSpeedByte(vl);
                    break;
                    case 'FileName':
                    vl='<span id="state-'+v.Id+'">'+StateIcon(v.State)+'</span> '+vl;
                    break;
                    case 'Progress':
                    var re2=new RegExp('\\{Percent\\}');
                    tmplCopy=tmplCopy.replace(re2,vl);
                    vl=FormatProgressBar(vl);
                    break;
                }
                tmplCopy=tmplCopy.replace(re,vl);
            }
            if(checkedAll){
                tmplCopy=$(tmplCopy);
                tmplCopy.find('.idCheck').prop('checked',true);
            }
            $('#fileList').append(tmplCopy);
        }
        if(downloadWSInterval && total<=finished){
            window.clearInterval(downloadWSInterval);
            downloadWSInterval=null;
        }
    });
}

function sockJSRead(){
    if(downloadWSInterval)return;
    if(!downloadWS){
        sockJSConnect();
    }else{
        downloadWS.send("progress");
    }
    downloadWSInterval=setInterval(function(){
        if(!downloadWS){
            window.clearInterval(downloadWSInterval);
            return;
        }
        downloadWS.send("progress");
    }, 2000);
}

function reqJSON(url,data,callback) {
    loading(false);
    var opt={
        contentType: "application/json; charset=utf-8",
        url: downloadAPIPrefix+url,
        type: "POST",
        dataType: "json"
    };
    if(data) opt.data = JSON.stringify(data);
    $.ajax(opt).error(function(jsonData) {
        loading(true);
        alert(jsonData);
    }).success(function(jsonData){
        loading(true);
        if(jsonData.Code!=1) {
            App.message({text: jsonData.Info,type:'error'},false);
            return;
        }
        if(callback)callback();
        sockJSRead();
    });
}

function reqForm(url,data,callback) {
    loading(false);
    var opt={
        url: downloadAPIPrefix+url,
        type: "POST",
        dataType: "json"
    };
    if(data) opt.data = data;
    $.ajax(opt).error(function(jsonData) {
        loading(true);
        alert(jsonData);
    }).success(function(jsonData){
        loading(true);
        if(jsonData.Code!=1) {
            App.message({text: jsonData.Info,type:'error'},false);
            return;
        }
        if(callback)callback();
        sockJSRead();
    });
}

function AddDownload() {
    var req = {
        PartCount: parseInt($("#part_count_id").val()),
        FilePath: $("#save_path_id").val(),
        Url: $("#url_id").val()
    };
    reqJSON("/add_task",req);
}
function checkedIds(){
    var ids = [];
    $('#fileTable .idCheck:checked').each(function(){
        ids.push(parseInt($(this).val()));
    });
    return ids;
}
function RemoveDownload() {
    var req = {id:checkedIds()};
    reqForm("/remove_task",req,function(){
        for(var i=0;i<req.id.length;i++){
            $('#id-'+req.id[i]).parent('tr').remove();
        }
    });
}
function StartDownload() {
    var req = {id:checkedIds()};
    reqForm("/start_task",req);
}
function StopDownload() {
    var req = {id:checkedIds()};
    reqForm("/stop_task",req);
}
function StartAllDownload() {
    reqJSON("/start_all_task");
}
function StopAllDownload() {
    reqJSON("/stop_all_task");
}
function OnChangeUrl() {
    var filename = $("#url_id").val().split('/').pop()
    $("#save_path_id").val(filename)
}
function loading(close){
    App.loading(close?'hide':'show');
}

$(function(){
    sockJSRead();
    $('body').on('click','#fileTable .allCheck',function(){
        $('#fileTable .idCheck').prop('checked',$(this).prop('checked'));
    });
});