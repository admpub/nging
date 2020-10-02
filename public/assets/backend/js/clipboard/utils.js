function attachCopy(elem,options){
    var clipboard = new ClipboardJS(elem, options||{}); 
    clipboard.on('success', function (e) { 
        e.clearSelection(); 
        var targetName=$(elem).data('target-name')||'';
        showTooltip(e.trigger, App.t('%s已复制成功',targetName), true); 
    });
    clipboard.on('error', function (e) { 
        showTooltip(e.trigger, fallbackMessage(e.action)); 
    });
    function showTooltip(e, msg, succeed) {
        App.message({text: msg, type: succeed?'success':'error'});
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