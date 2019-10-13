function transferTS(tsURL,videoWrapper){
    if(!videoWrapper) videoWrapper='#video-wrapper';
    var $ = document.querySelector.bind(document);
    var vjsParsed,
        video, 
        mediaSource;
    // 定义通用的事件回调处理函数，只做打印事件类型
    function logevent (event) {
      console.log(event);
    }

    // ajax
    let xhr = new XMLHttpRequest();
    xhr.open('GET', tsURL);
    // 接收的是 video/mp2t 二进制数据，Blob类型也可以，但arraybuffer类型方便后续直接处理 
    xhr.responseType = "arraybuffer";
    xhr.send();
    xhr.onreadystatechange = function () {
      if (xhr.readyState ==4) {
        if (xhr.status == 200) {
          transferFormat(xhr.response);
        } else {
          console.log('error');
        }
      }
    }
    
    function transferFormat (data) {
      // 将源数据从ArrayBuffer格式保存为可操作的Uint8Array格式
      // https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/ArrayBuffer
      var segment = new Uint8Array(data); 
      var combined = false;
      // 接收无音频ts文件，OutputType设置为'video'，带音频ts设置为'combined'
      var outputType = 'video';
      var remuxedSegments = [];
      var remuxedBytesLength = 0;
      var remuxedInitSegment = null;

      // remux选项默认为true，将源数据的音频视频混合为mp4，设为false则不混合
      var transmuxer = new muxjs.mp4.Transmuxer({remux: false});
      
      // 监听data事件，开始转换流
      transmuxer.on('data', function(event) {
        console.log(event);
        if (event.type === outputType) {
          remuxedSegments.push(event);
          remuxedBytesLength += event.data.byteLength;
          remuxedInitSegment = event.initSegment;
        }
      });
      // 监听转换完成事件，拼接最后结果并传入MediaSource
      transmuxer.on('done', function () {
        var offset = 0;
        var bytes = new Uint8Array(remuxedInitSegment.byteLength + remuxedBytesLength)
        bytes.set(remuxedInitSegment, offset);
        offset += remuxedInitSegment.byteLength;

        for (var j = 0, i = offset; j < remuxedSegments.length; j++) {
          bytes.set(remuxedSegments[j].data, i);
          i += remuxedSegments[j].byteLength;
        }
        remuxedSegments = [];
        remuxedBytesLength = 0;
        // 解析出转换后的mp4相关信息，与最终转换结果无关
        vjsParsed = muxjs.mp4.tools.inspect(bytes);
        console.log('transmuxed', vjsParsed);

        prepareSourceBuffer(combined, outputType, bytes);
      });
      // push方法可能会触发'data'事件，因此要在事件注册完成后调用
      transmuxer.push(segment); // 传入源二进制数据，分割为m2ts包，依次调用上图中的流程
      // flush的调用会直接触发'done'事件，因此要事件注册完成后调用
      transmuxer.flush(); // 将所有数据从缓存区清出来
    }

    function prepareSourceBuffer (combined, outputType, bytes) {
      var buffer;
      video = document.createElement('video');
      video.controls = true;
      // MediaSource Web API: https://developer.mozilla.org/zh-CN/docs/Web/API/MediaSource
      mediaSource = new MediaSource(); 
      video.src = URL.createObjectURL(mediaSource);
    
      $(videoWrapper).appendChild(video); // 将H5 video元素添加到对应DOM节点下
    
      // 转换后mp4的音频格式 视频格式
      var codecsArray = ["avc1.64001f", "mp4a.40.5"];
    
      mediaSource.addEventListener('sourceopen', function () {
        // MediaSource 实例默认的duration属性为NaN
        mediaSource.duration = 0;
        // 转换为带音频、视频的mp4
        if (combined) {
          buffer = mediaSource.addSourceBuffer('video/mp4;codecs="' + 'avc1.64001f,mp4a.40.5' + '"');
        } else if (outputType === 'video') {
          // 转换为只含视频的mp4
          buffer = mediaSource.addSourceBuffer('video/mp4;codecs="' + codecsArray[0] + '"');
        } else if (outputType === 'audio') {
          // 转换为只含音频的mp4
          buffer = mediaSource.addSourceBuffer('audio/mp4;codecs="' + (codecsArray[1] ||codecsArray[0]) + '"');
        }
    
        buffer.addEventListener('updatestart', logevent);
        buffer.addEventListener('updateend', logevent);
        buffer.addEventListener('error', logevent);
        video.addEventListener('error', logevent);
        // mp4 buffer 准备完毕，传入转换后的数据
        // 将 bytes 放入 MediaSource 创建的sourceBuffer中
        // https://developer.mozilla.org/en-US/docs/Web/API/SourceBuffer/appendBuffer
        buffer.appendBuffer(bytes);
        // 自动播放
        // video.play();
      });
    };
}