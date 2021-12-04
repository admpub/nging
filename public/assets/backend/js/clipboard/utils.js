function attachCopy(elem,options){
    var clipboard = new ClipboardJS(elem, options||{}); 
    clipboard.on('success', function (e) { 
        e.clearSelection(); 
        var targetName=$(elem).data('target-name')||'';
        if(typeof App != 'undefined'){
            showTooltip(e.trigger, App.t('%s已复制成功',targetName), true); 
        }else{
            alert(targetName+'已复制成功');
        }
    });
    clipboard.on('error', function (e) { 
        showTooltip(e.trigger, fallbackMessage(e.action)); 
    });
    function showTooltip(e, msg, succeed) {
        if(typeof App != 'undefined'){
            App.message({text: msg, type: succeed?'success':'error'});
        }else{
            alert(msg);
        }
    }
    function fallbackMessage(action) {
        var actionMsg = ''; 
        var actionKey = (action === 'cut' ? 'X' : 'C'); 
        if (/iPhone|iPad/i.test(navigator.userAgent)) { 
            actionMsg = 'No support :('; 
        } else if (/Mac/i.test(navigator.userAgent)) { 
            actionMsg = 'Press ⌘-' + actionKey + ' to ' + action;
        } else {
            actionMsg = 'Press Ctrl-' + actionKey + ' to ' + action;
        }
        return actionMsg;
    }
    return clipboard;
}