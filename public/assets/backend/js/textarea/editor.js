(function(a){
    var TextAreaEditor = {
        __calcBookmark: function(bookmark) {
            return (bookmark.charCodeAt(0)-1)+(bookmark.charCodeAt(3)-1)*65536+(bookmark.charCodeAt(2)-1);
        },
        __getSelectPos: function(editor, end) {
            if (!editor) return;
            if (typeof editor.selectionStart != "undefined")
                return end ? editor.selectionEnd : editor.selectionStart;
            if (!editor.createTextRange) return;
            editor.focus ();
            var range = document.selection.createRange().duplicate();
            if (!end) range.collapse(true)
            range.setEndPoint("StartToEnd", range);
            var start = document.body.createTextRange();
            start.moveToElementText(editor);
            var start = this.__calcBookmark(range.getBookmark()) - this.__calcBookmark(start.getBookmark());
            return start;
        },
        getSelectStart: function(editor) {
            return this.__getSelectPos(editor);
        },
        getSelectEnd: function(editor) {
            return this.__getSelectPos(editor, true);
        },
        getSelectRange: function(editor) {
            return [this.getSelectStart(editor), this.getSelectEnd(editor)];
        },
        setSelectRange: function(editor, range) {
            if (!editor) return;
            if (range[0] > range[1]) return;
            editor.focus();
            if (editor.setSelectionRange) {
                editor.setSelectionRange(range[0], range[1]);
            } else if (editor.createTextRange) {
                var textRange = editor.createTextRange();
                textRange.collapse(true);
                textRange.moveEnd("character", range[1]);
                textRange.moveStart("character", range[0]);
                textRange.select();
            }
        },
        textPos: function(editor, pos) {
            if (!editor) return;
            if (!editor.createTextRange) return pos;
            var value = editor.value;
            for (var i = 0; i <= pos; i++) if (value.charAt(i) == "\n") pos++;
            return pos;
        },
        cartPos: function(editor, pos) {
            if (!editor) return;
            if (!editor.createTextRange) return pos;
            var value = editor.value;
            var j = 0;
            for (var i = 0; i <= pos; i++) if (value.charAt(i) == "\n") j++;
            return pos - j;
        },
        selectEx: function(editor, startPattern, endPattern, include) {
            if (!editor) return;
            startPattern = "" + startPattern;
            endPattern = "" + endPattern;
            var range = this.getSelectRange(editor);
            var value = editor.value;
            var textRange = [this.textPos(editor, range[0]), this.textPos(editor, range[1])];
            var startStr = value.substr(0, textRange[0]);
            var endStr = value.substring(textRange[1], value.length);
            var i = startStr.lastIndexOf(startPattern);
            if (i < 0) return;
            var j = endStr.indexOf(endPattern);
            if (j < 0) return;
            j += textRange[1];
            if (include)
                j += endPattern.length;
            else i += startPattern.length;
            this.setSelectRange(editor, [this.cartPos(editor, i), this.cartPos(editor, j)]);
        },
        setSelectText: function (editor, value, posStart, posEnd) {
            if (!editor) return;
            editor.focus();
            if (document.selection) {
				var sel = document.selection.createRange();
				sel.text = value;
				sel.select();
            } else if (typeof editor.selectionStart != "undefined") { // firefox
                var str = editor.value;
                var start = editor.selectionStart;
                var scroll = editor.scrollTop;
                editor.value = str.substr(0, start) + value +
                    str.substring(editor.selectionEnd, str.length);
                editor.selectionStart = start + value.length;
                editor.selectionEnd = start + value.length;
                editor.scrollTop = scroll;
            }
            if (posStart==null) return;
            var range=this.getSelectRange(editor);
            if (range[0]==range[1]) range[0]=range[1]-value.length;
            if (posStart==null || posStart>=value.length || posStart<0) posStart=0;
            if (posEnd!=null && posEnd<value.length) {
              if (posEnd<0) posEnd=value.length+posEnd;
              if (posEnd<posStart) posEnd=value.length;
            } else {
              posEnd=value.length;
            }
            range[1]=range[0]+posEnd;
            range[0]=range[0]+posStart;
            this.setSelectRange(editor, range);
        }
    };
    a.TextAreaEditor=TextAreaEditor;
})(window);
