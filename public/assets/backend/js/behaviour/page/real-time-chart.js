(function(){
var ws,idElem='#CPU-Usage',idNetElem='#Net-Usage',idNetPacket='#NetPacket-Usage'; 
var _chartCPU,_chartNet,_chartNetPacket,options = {
  series: {
    lines: {
      show: true,
      lineWidth: 2, 
      fill: true,
      fillColor: {
        colors: [{
          opacity: 0.25
        }, {
          opacity: 0.25
        }]
      } 
    },
    points: {
      show: false
    },
    shadowSize: 2
  },
  legend:{
    show: true
  },
  grid: {
    labelMargin: 10,
    axisMargin: 500,
    hoverable: true,
    clickable: true,
    tickColor: "rgba(0,0,0,0.15)",
    borderWidth: 0
  },
  colors: ["#B450B2", "#4A8CF7", "#52e136"],
  xaxis: {
    mode: "time",
    tickSize: [2, "second"],
    tickFormatter: function (v, axis) {
      var date = new Date(v);
      if (date.getSeconds() % 10 == 0) {
          var hours = date.getHours() < 10 ? "0" + date.getHours() : date.getHours();
          var minutes = date.getMinutes() < 10 ? "0" + date.getMinutes() : date.getMinutes();
          var seconds = date.getSeconds() < 10 ? "0" + date.getSeconds() : date.getSeconds();

          return hours + ":" + minutes + ":" + seconds;
      } 
      return "";
    }
  },
  yaxis: {
    min: 0,
    max: 100,        
    tickSize: 10,
    tickFormatter: function (v, axis) {
      if (v % 10 == 0) return v + "%";
      return "";
    }
  }
},netOptions = {
  series: {
    lines: {
      show: true,
      lineWidth: 2, 
      fill: true,
      fillColor: {
        colors: [{
          opacity: 0.25
        }, {
          opacity: 0.25
        }]
      } 
    },
    points: {
      show: false
    },
    shadowSize: 2
  },
  legend:{
    show: true
  },
  grid: {
    labelMargin: 10,
    axisMargin: 500,
    hoverable: true,
    clickable: true,
    tickColor: "rgba(0,0,0,0.15)",
    borderWidth: 0
  },
  colors: ["#B450B2", "#4A8CF7", "#52e136"],
  xaxis: {
    mode: "time",
    tickSize: [2, "second"],
    tickFormatter: function (v, axis) {
      var date = new Date(v);
      if (date.getSeconds() % 10 == 0) {
          var hours = date.getHours() < 10 ? "0" + date.getHours() : date.getHours();
          var minutes = date.getMinutes() < 10 ? "0" + date.getMinutes() : date.getMinutes();
          var seconds = date.getSeconds() < 10 ? "0" + date.getSeconds() : date.getSeconds();

          return hours + ":" + minutes + ":" + seconds;
      } 
      return "";
    }
  },
  yaxis: {
    tickFormatter: function (v, axis) {
      return App.formatBytes(v);
    }
  }
},netPacketOptions = {
  series: {
    lines: {
      show: true,
      lineWidth: 2, 
      fill: true,
      fillColor: {
        colors: [{
          opacity: 0.25
        }, {
          opacity: 0.25
        }]
      } 
    },
    points: {
      show: false
    },
    shadowSize: 2
  },
  legend:{
    show: true
  },
  grid: {
    labelMargin: 10,
    axisMargin: 500,
    hoverable: true,
    clickable: true,
    tickColor: "rgba(0,0,0,0.15)",
    borderWidth: 0
  },
  colors: ["#B450B2", "#4A8CF7", "#52e136"],
  xaxis: {
    mode: "time",
    tickSize: [2, "second"],
    tickFormatter: function (v, axis) {
      var date = new Date(v);
      if (date.getSeconds() % 10 == 0) {
          var hours = date.getHours() < 10 ? "0" + date.getHours() : date.getHours();
          var minutes = date.getMinutes() < 10 ? "0" + date.getMinutes() : date.getMinutes();
          var seconds = date.getSeconds() < 10 ? "0" + date.getSeconds() : date.getSeconds();

          return hours + ":" + minutes + ":" + seconds;
      } 
      return "";
    }
  }
};

// === Net ====
function getNetData(info) {
  return [{
    data: info.Net.BytesRecv,
    //color: '#0f0',
    label: App.i18n.chart.DOWNLOAD_SPEED
  },{
    data: info.Net.BytesSent, 
    //color: '#00f', 
    label: App.i18n.chart.UPLOAD_SPEED
  }];
}
function updateNet(data) {
  _chartNet.setData(data);
  // Since the axes don't change, we don't need to call plot.setupGrid()
  _chartNet.setupGrid();
  _chartNet.draw();
}
function chartNet(info) {
  data=getNetData(info)
  _chartNet=$(idNetElem).data('plot');
  if(!_chartNet) return initNetData(data);
  updateNet(data);
}
function initNetData(data){
  _chartNet = $.plot($(idNetElem), data, netOptions);
  $(idNetElem).data('plot',_chartNet);
}

// === Net-Packet ====
function getNetPacketData(info) {
  return [{
    data: info.Net.PacketsRecv,
    //color: '#0f0',
    label: App.i18n.chart.RECV_PACKETS
  },{
    data: info.Net.PacketsSent, 
    //color: '#00f', 
    label: App.i18n.chart.SENT_PACKETS
  }];
}
function updateNetPacket(data) {
  _chartNetPacket.setData(data);
  // Since the axes don't change, we don't need to call plot.setupGrid()
  _chartNetPacket.setupGrid();
  _chartNetPacket.draw();
}
function chartNetPacket(info) {
  data=getNetPacketData(info)
  _chartNetPacket=$(idNetPacket).data('plot');
  if(!_chartNetPacket) return initNetPacketData(data);
  updateNetPacket(data);
}
function initNetPacketData(data){
  _chartNetPacket = $.plot($(idNetPacket), data, netPacketOptions);
  $(idNetPacket).data('plot',_chartNetPacket);
}

// === CPU ====
function getData(info) {
  return [{
    data: info.CPU,
    //color: '#0f0',
    label: App.i18n.chart.CPU_USAGE
  },{
    data: info.Mem, 
    //color: '#00f', 
    label: App.i18n.chart.MEMORY_USAGE
  }];
}
function update(data) {
  _chartCPU.setData(data);
  // Since the axes don't change, we don't need to call plot.setupGrid()
  _chartCPU.setupGrid();
  _chartCPU.draw();
}
function chartCPU(info) {
  data=getData(info)
  _chartCPU=$(idElem).data('plot');
  if(!_chartCPU) return initData(data);
  update(data);
}
function initData(data){
  _chartCPU = $.plot($(idElem), data, options);
  $(idElem).data('plot',_chartCPU);
}

function connectWS(onopen){
	if (ws) {
		if(onopen!=null)onopen();
		return false;
	}
	var url=App.wsURL(BACKEND_URL);
	ws = new WebSocket(url+"/server/dynamic");
	ws.onopen = function(evt) {
		if(onopen!=null)onopen();
	};
	ws.onclose = function(evt) {
	  ws = null;
  };
	ws.onmessage = function(evt) {
    console.dir(evt.data);
    var info=JSON.parse(evt.data);
    chartCPU(info);
    chartNet(info);
    chartNetPacket(info);
	};
}
function tick(){
    connectWS(function(){
      ws.send("ping");
    });
    if($(idElem).length<1){
      clearInterval(window._interval);
      window._interval=null;
    }
}
$(function(){
  if(typeof(window._interval)!='undefined' && window._interval){
    clearInterval(window._interval);
  }
  tick();
  window._interval=window.setInterval(tick,2000);
});
})();