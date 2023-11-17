(function () {
var _interval = null, ws;
var idCPU='#dockerContainerCPUStats',idMem='#dockerContainerMemStats',idNetRecv='#netRecv',idNetSend='#netSend',idDiskRead='#diskRead',idDiskWrite='#diskRead',idTopTmpl='tmplTopRankTable',idTop='#topRank';
function connectWS(onopen) {
    if (ws) {
        if (onopen != null) onopen();
        return false;
    }
    var url = App.wsURL(BACKEND_URL);
    ws = new WebSocket(url + "/docker/base/container/stats/" + containerId);
    ws.onopen = function (evt) {
        if (onopen != null) onopen();
    };
    ws.onclose = function (evt) {
        ws = null;
    };
    ws.onmessage = function (evt) {
        if($(idCPU).length<1){
          ws.close();clear();
          return;
        }
        //console.dir(evt.data);
        var info=JSON.parse(evt.data);
        updateCPU(info);updateMem(info);updateNet(info);updateDisk(info);
    };
    ws.onerror = function (evt) {
        if ('data' in evt) console.log(evt.data);
        clear();
    };
}
function tick() {
    connectWS(function () {
        var ping = "ping";
        try {
            ws.send(ping);
        } catch (error) {
            clear();
            console.error(error.message);
        }
    });
}
function clear() {
    if (_interval) {
        clearInterval(_interval);
        _interval = null;
    }
}
function initCPU(){
    $(idCPU).easyPieChart({
        easing: 'easeOutBounce',
        barColor: '#69c',
        trackColor: '#ace',
        scaleColor: false,
        lineWidth: 20,
        trackWidth: 16,
        lineCap: 'butt',
        onStep: function(from, to, percent) {
            $(this.el).find('.percent').text(percent.toFixed(5)+'%');
        }
    });
}
function fixNaN(v) {
    if(typeof v == 'undefined') return 0;
    v=Number(v);
    return isNaN(v)?0:v;
}
function updateCPU(info){
    var cpuStats=info.cpu_stats;
    var prevStats=info.precpu_stats;
    var perCPUUsage=cpuStats.cpu_usage.percpu_usage||0;
    var cpuPercent = 0.0;
    var cpuDelta = fixNaN(cpuStats.cpu_usage.total_usage) - fixNaN(prevStats.cpu_usage.total_usage);
    var systemDelta = fixNaN(cpuStats.system_cpu_usage) - fixNaN(prevStats.system_cpu_usage);
    //console.log(cpuDelta,systemDelta)
    if (cpuDelta > 0.0 && cpuDelta > 0.0) {
        cpuPercent = (cpuDelta / systemDelta) * fixNaN(perCPUUsage) * 100.0;
    }
    $(idCPU).data('easyPieChart').update(cpuPercent);
}
function initMem(){
    $(idMem).easyPieChart({
        easing: 'easeOutBounce',
        barColor: '#69c',
        trackColor: '#ace',
        scaleColor: false,
        lineWidth: 20,
        trackWidth: 16,
        lineCap: 'butt',
        onStep: function(from, to, percent) {
            $(this.el).find('.percent').text(percent.toFixed(5)+'%');
        }
    });
}
function updateMem(info){
    var stats=info.memory_stats;
    var percent=0;
    if(stats.limit!=0){
        percent=stats.usage/stats.limit;
    }
    $(idMem).data('easyPieChart').update(percent*100);
}
function updateDisk(info){
    var stats=info.blkio_stats.io_service_bytes_recursive;
    var read=0,write=0;
    //console.dir(stats)
    for(var k in stats){
        var v=stats[k];
        switch(v.op){
            case 'read':
                read+=v.value;break;
            case 'write':
                write+=v.value;break;
        }
    }
    $(idDiskRead).text(App.formatBytes(read));
    $(idDiskWrite).text(App.formatBytes(write));
}
function updateNet(info){
    var stats=info.networks;
    var recv=0,send=0;
    for(var k in stats){
        var v=stats[k];
        recv+=v.rx_bytes;
        send+=v.tx_bytes;
    }
    $(idNetRecv).text(App.formatBytes(recv));
    $(idNetSend).text(App.formatBytes(send));
}
function getTopList(){
    $.get(BACKEND_URL+'/docker/base/container/top/'+containerId,{},function(r){
        if(r.Code!=1) {
            clear();
            return App.message({text:r.Info,type:'error'});
        }
        $(idTop).html(template(idTopTmpl,r.Data));
    },'json');
}
$(function(){
    clear();
    initCPU();initMem();
    getTopList();
    tick();
    _interval=window.setInterval(function(){
        //tick();
        getTopList();
    },2000);
})
})(containerId);