# Editor.js 文档 / Editor.js Documentation

## 概述 / Overview

`editor.js` 是 Nging 系统的编辑器管理模块，提供了多种富文本和代码编辑器的集成支持，包括 TinyMCE、EditorMD、CodeMirror、Markdown-it 等，以及文件上传、图片裁剪、日期选择等功能。

---

## 目录 / Table of Contents

1. [编辑器类型 / Editor Types](#编辑器类型--editor-types)
2. [富文本编辑器 / Rich Text Editors](#富文本编辑器--rich-text-editors)
3. [Markdown 编辑器 / Markdown Editors](#markdown-编辑器--markdown-editors)
4. [代码编辑器 / Code Editors](#代码编辑器--code-editors)
5. [文件上传 / File Upload](#文件上传--file-upload)
6. [表单组件 / Form Components](#表单组件--form-components)
7. [工具函数 / Utility Functions](#工具函数--utility-functions)

---

## 编辑器类型 / Editor Types

### App.editor.browsingFileURL / 文件浏览 URL / File Browse URL

默认的文件浏览器 URL。

```javascript
App.editor.browsingFileURL = BASE_URL + '/user/file/finder'
```

---

## 富文本编辑器 / Rich Text Editors

### App.editor.tinymce(elem, uploadUrl, options, useSimpleToolbar) / TinyMCE 编辑器 / TinyMCE Editor

初始化 TinyMCE 富文本编辑器。

**语法 / Syntax:**
```javascript
App.editor.tinymce(elem, uploadUrl, options, useSimpleToolbar)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `elem` | Selector/Element | 文本域元素 / Textarea element |
| `uploadUrl` | String | 上传 URL，默认使用 `action` 属性 / Upload URL, defaults to `action` attribute |
| `options` | Object | TinyMCE 配置选项 / TinyMCE config options |
| `useSimpleToolbar` | Boolean | 是否使用简化工具栏 / Whether to use simple toolbar |

**默认配置 / Default Configuration:**

```javascript
{
    height: 500,
    menubar: true,
    language: App.langTag('_'),
    plugins: [
        'print', 'preview', 'paste', 'importcss', 'searchreplace', 
        'autolink', 'autosave', 'save', 'directionality', 'code', 
        'visualblocks', 'visualchars', 'fullscreen', 'image', 'link', 
        'media', 'template', 'codesample', 'table', 'charmap', 'hr', 
        'pagebreak', 'nonbreaking', 'anchor', 'toc', 'insertdatetime', 
        'advlist', 'lists', 'wordcount', 'imagetools', 'textpattern', 
        'noneditable', 'charmap', 'quickbars', 'emoticons'
    ],
    toolbar: 'undo redo | bold italic underline strikethrough | fontselect fontsizeselect formatselect | alignleft aligncenter alignright alignjustify | outdent indent | numlist bullist | forecolor backcolor removeformat | pagebreak | charmap emoticons | fullscreen preview save print | insertfile image media template link anchor codesample | ltr rtl | customDateButton',
    toolbar_sticky: true,
    autosave_ask_before_unload: false,
    autosave_interval: "30s",
    autosave_prefix: "{path}{query}-{id}-",
    autosave_restore_when_empty: true,
    autosave_retention: "2m",
    image_advtab: true,
    image_caption: true,
    relative_urls: false,
    image_title: true,
    quickbars_selection_toolbar: 'bold italic | quicklink h2 h3 blockquote quicktable',
    toolbar_drawer: 'sliding',
    contextmenu: 'link table',
    templates: [...]
}
```

**示例 / Example:**

```javascript
// 基本用法 / Basic usage
App.editor.tinymce('#content');

// 带上传配置 / With upload config
App.editor.tinymce('#content', '/api/upload');

// 自定义配置 / Custom config
App.editor.tinymce('#content', '/api/upload', {
    height: 600,
    toolbar: 'undo redo | bold italic | link image'
});

// 使用简化工具栏 / Use simple toolbar
App.editor.tinymce('#content', '/api/upload', {}, true);

// 批量初始化 / Batch initialize
App.editor.tinymces('textarea.editor');
```

### App.editor.finderDialog(remoteURL, callback, zIndex) / 文件选择对话框 / File Picker Dialog

打开文件选择对话框。

**语法 / Syntax:**
```javascript
App.editor.finderDialog(remoteURL, callback, zIndex)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 默认值 / Default | 描述 / Description |
|-------------------|-------------|-----------------|-------------------|
| `remoteURL` | String | - | 文件浏览器 URL / File browser URL |
| `callback` | Function | - | 回调函数，接收文件列表 / Callback function, receives file list |
| `zIndex` | Number | `2000` | 对话框层级 / Dialog z-index |

**示例 / Example:**

```javascript
App.editor.finderDialog(
    '/file/finder?filetype=image&multiple=1',
    function(fileList, infoList) {
        console.log('Selected files:', fileList);
        console.log('File info:', infoList);
        
        // fileList: ['http://example.com/image1.jpg', ...]
        // infoList: [{name: 'image1.jpg', size: 1024}, ...]
    }
);
```

---

## Markdown 编辑器 / Markdown Editors

### App.editor.markdown(elem, uploadUrl, options) / EditorMD 编辑器 / EditorMD Editor

初始化 EditorMD Markdown 编辑器。

**语法 / Syntax:**
```javascript
App.editor.markdown(elem, uploadUrl, options)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `elem` | Selector/Element | 文本域元素 / Textarea element |
| `uploadUrl` | String | 上传 URL / Upload URL |
| `options` | Object | EditorMD 配置选项 / EditorMD config options |

**默认配置 / Default Configuration:**

```javascript
{
    width: "100%",
    height: container.height(),
    path: ASSETS_URL + '/js/editor/markdown/lib/',
    markdown: $(elem).val(),
    placeholder: $(elem).attr('placeholder') || '',
    codeFold: true,
    saveHTMLToTextarea: true,
    searchReplace: true,
    watch: true,
    htmlDecode: "style,script,iframe|on*",
    emoji: true,
    taskList: true,
    tocm: true,
    tex: true,
    flowChart: true,
    sequenceDiagram: true,
    imageUpload: true,
    imageFormats: ["jpg", "jpeg", "gif", "png", "bmp"],
    imageUploadURL: uploadUrl,
    crossDomainUpload: true
}
```

**示例 / Example:**

```javascript
// 基本用法 / Basic usage
App.editor.markdown('#content');

// 带上传配置 / With upload config
App.editor.markdown('#content', '/api/upload');

// 自定义配置 / Custom config
App.editor.markdown('#content', '/api/upload', {
    height: 600,
    theme: 'dark',
    previewTheme: 'dark',
    editorTheme: 'ambiance'
});

// 批量初始化 / Batch initialize
App.editor.markdowns('textarea.markdown-editor');
```

### App.editor.markdownToHTML(elem, markdownData, options) / Markdown 转 HTML / Markdown to HTML

将 Markdown 内容渲染为 HTML。

**语法 / Syntax:**
```javascript
App.editor.markdownToHTML(elem, markdownData, options)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `elem` | Selector/Element | 容器元素 / Container element |
| `markdownData` | String | Markdown 内容 / Markdown content |
| `options` | Object | 渲染选项 / Render options |

**示例 / Example:**

```javascript
// 从文本域渲染 / Render from textarea
App.editor.markdownToHTML('#markdown-container');

// 从字符串渲染 / Render from string
App.editor.markdownToHTML('#markdown-container', '# Hello World\n\nThis is a **test**.');

// 自定义选项 / Custom options
App.editor.markdownToHTML('#markdown-container', null, {
    toc: true,
    taskList: true,
    tex: false,
    flowChart: true
});
```

### App.editor.md / Markdown 编辑器别名 / Markdown Editor Alias

`App.editor.md` 是 `App.editor.markdown` 的别名。

```javascript
App.editor.md('#content');
```

### App.editor.markdownItToHTML(box, isContainer) / Markdown-it 渲染 / Markdown-it Render

使用 Markdown-it 渲染 Markdown 内容。

**语法 / Syntax:**
```javascript
App.editor.markdownItToHTML(box, isContainer)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 默认值 / Default | 描述 / Description |
|-------------------|-------------|-----------------|-------------------|
| `box` | Selector/Element | - | 容器元素 / Container element |
| `isContainer` | Boolean | `true` | 是否为容器 / Whether it's a container |

---

## 代码编辑器 / Code Editors

### App.editor.codemirror(elem, options, loadLangType) / CodeMirror 编辑器 / CodeMirror Editor

初始化 CodeMirror 代码编辑器。

**语法 / Syntax:**
```javascript
App.editor.codemirror(elem, options, loadLangType)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `elem` | Selector/Element | 文本域或 div 元素 / Textarea or div element |
| `options` | Object | CodeMirror 配置选项 / CodeMirror config options |
| `loadLangType` | String/Object | 加载的语言模式 / Language mode to load |

**默认配置 / Default Configuration:**

```javascript
{
    lineNumbers: true,
    lineWrapping: true,
    gutters: ["CodeMirror-linenumbers", "CodeMirror-foldgutter"],
    autoCloseTags: true,
    autoCloseBrackets: true,
    showTrailingSpace: true,
    indentWithTabs: true,
    matchBrackets: true,
    styleActiveLine: true,
    styleSelectedText: true,
    highlightSelectionMatches: true,
    smartIndent: true,
    mode: "text/x-csrc",
    width: null,
    height: null,
    hintOptions: {completeSingle: false}
}
```

**示例 / Example:**

```javascript
// 基本用法 / Basic usage
App.editor.codemirror('#code-editor');

// 指定语言模式 / Specify language mode
App.editor.codemirror('#code-editor', {
    mode: 'text/javascript',
    theme: 'ambiance'
});

// 加载特定语言 / Load specific language
App.editor.codemirror('#code-editor', null, 'javascript');
```

### App.editor.codeHighlight(elem) / 代码高亮 / Code Highlight

为代码块应用语法高亮。

**语法 / Syntax:**
```javascript
App.editor.codeHighlight(elem)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 默认值 / Default | 描述 / Description |
|-------------------|-------------|-----------------|-------------------|
| `elem` | Selector/Object | `'pre[class^=language-]'` | 代码块元素 / Code block element |

**示例 / Example:**

```javascript
// 高亮所有代码块 / Highlight all code blocks
App.editor.codeHighlight();

// 高亮指定元素 / Highlight specific element
App.editor.codeHighlight('#code-block');
```

---

## 文件上传 / File Upload

### App.editor.fileInput(elem, options, successCallback, errorCallback) / 文件输入 / File Input

初始化文件输入组件（包括文件选择器和预览功能）。

**语法 / Syntax:**
```javascript
App.editor.fileInput(elem, options, successCallback, errorCallback)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `elem` | Selector | 容器选择器 / Container selector |
| `options` | Object | 上传选项 / Upload options |
| `successCallback` | Function | 成功回调 / Success callback |
| `errorCallback` | Function | 错误回调 / Error callback |

**数据属性 / Data Attributes:**

- `data-toggle="finder"` - 文件选择器 / File picker
- `data-finder-url` - 文件浏览器 URL / File browser URL
- `data-input` - 输入框选择器 / Input selector
- `data-preview-btn` - 预览按钮选择器 / Preview button selector
- `data-preview-img` - 预览图片选择器 / Preview image selector
- `data-toggle="uploadPreviewer"` - 上传预览 / Upload preview
- `data-upload-url` - 上传 URL / Upload URL
- `data-preview-container` - 预览容器 / Preview container

**示例 / Example:**

```html
<!-- 文件选择器 / File picker -->
<div class="input-group">
    <input type="text" name="avatar" class="form-control" />
    <span class="input-group-btn">
        <button class="btn btn-default" data-toggle="finder" data-finder-url="/file/finder">
            <i class="fa fa-folder-open"></i>
        </button>
    </span>
</div>

<!-- 上传预览 / Upload preview -->
<div>
    <button class="btn btn-primary" data-toggle="uploadPreviewer" data-upload-url="/api/upload">
        <i class="fa fa-upload"></i> 上传
    </button>
    <div class="preview-container"></div>
</div>

<script>
App.editor.fileInput();
</script>
```

### App.editor.uploadPreviewer(elem, options) / 上传预览 / Upload Previewer

初始化上传预览组件。

**语法 / Syntax:**
```javascript
App.editor.uploadPreviewer(elem, options)
```

### App.editor.dropzone(elem, options, onSuccss, onError, onRemove) / Dropzone 上传 / Dropzone Upload

初始化 Dropzone 文件上传组件。

**语法 / Syntax:**
```javascript
App.editor.dropzone(elem, options, onSuccss, onError, onRemove)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `elem` | Selector | Dropzone 容器 / Dropzone container |
| `options` | Object | Dropzone 配置选项 / Dropzone config options |
| `onSuccss` | Function | 成功回调 / Success callback |
| `onError` | Function | 错误回调 / Error callback |
| `onRemove` | Function | 移除回调 / Remove callback |

**示例 / Example:**

```javascript
// 基本用法 / Basic usage
App.editor.dropzone('#my-dropzone', {
    url: '/api/upload',
    maxFilesize: 2,  // MB
    acceptedFiles: 'image/*'
});

// 支持签名 URL / Support signed URL
App.editor.dropzone('#my-dropzone', {
    getSignedPutURL: function(file, cb, done) {
        // 获取签名上传 URL
        $.post('/api/sign-url', {name: file.name}, function(res) {
            cb(res.URL);
        });
    }
});

// 带回调 / With callbacks
App.editor.dropzone('#my-dropzone', {
    url: '/api/upload'
}, function(file, resp) {
    console.log('Upload success:', resp);
}, function(file, error) {
    console.error('Upload error:', error);
}, function(file) {
    console.log('File removed:', file);
});
```

---

## 表单组件 / Form Components

### App.editor.select2() / Select2 初始化 / Select2 Initialize

初始化 Select2 组件。

**语法 / Syntax:**
```javascript
App.editor.select2()
```

### App.editor.selectPage(elem, options, loaded) / 选择页组件 / Select Page Component

初始化选择页组件（类似 Select2，支持分页）。

**语法 / Syntax:**
```javascript
App.editor.selectPage(elem, options, loaded)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `elem` | Selector | 元素选择器 / Element selector |
| `options` | Object | 配置选项 / Config options |
| `loaded` | Boolean | 是否已加载 / Whether already loaded |

**配置选项 / Config Options:**

```javascript
{
    showField: 'name',        // 显示字段 / Display field
    keyField: 'id',          // 键字段 / Key field
    data: [],                 // 数据或 URL / Data or URL
    params: function(){return {}},
    eAjaxSuccess: function(d){...},
    eSelect: function (data) {},
    eClear: function () {}
}
```

**示例 / Example:**

```javascript
// 基本用法 / Basic usage
App.editor.selectPage('#user-select', {
    url: '/api/users',
    showField: 'name',
    keyField: 'id'
});

// 批量初始化 / Batch initialize
App.editor.selectPages(
    '#user-select',
    {url: '/api/users'},
    '#category-select',
    {url: '/api/categories'}
);
```

### App.editor.cascadeSelect(elem, selectedIds, url) / 级联选择 / Cascade Select

初始化级联选择组件。

**语法 / Syntax:**
```javascript
App.editor.cascadeSelect(elem, selectedIds, url)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `elem` | Selector | 选择器元素 / Selector element |
| `selectedIds` | Array | 选中的 ID 列表 / Selected ID list |
| `url` | String | 数据源 URL / Data source URL |

### App.editor.dateRangePicker(rangeElem, options) / 日期范围选择器 / Date Range Picker

初始化日期范围选择器。

**语法 / Syntax:**
```javascript
App.editor.dateRangePicker(rangeElem, options)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `rangeElem` | Selector | 输入元素 / Input element |
| `options` | Object | 配置选项 / Config options |

**示例 / Example:**

```javascript
App.editor.dateRangePicker('#daterangepicker', {
    format: 'YYYY-MM-DD',
    separator: ' - ',
    locale: {
        applyLabel: '确定',
        cancelLabel: '取消',
        fromLabel: '起始',
        toLabel: '结束'
    }
});
```

### App.editor.datePicker(elem, options) / 日期选择器 / Date Picker

初始化日期选择器。

**语法 / Syntax:**
```javascript
App.editor.datePicker(elem, options)
```

### App.editor.inputmask(elem, options) / 输入掩码 / Input Mask

初始化输入掩码组件。

**语法 / Syntax:**
```javascript
App.editor.inputmask(elem, options)
```

**示例 / Example:**

```javascript
// 电话号码 / Phone number
App.editor.inputmask('#phone', {
    mask: '999-9999-9999'
});

// 日期 / Date
App.editor.inputmask('#date', {
    mask: '9999-99-99'
});

// 邮箱 / Email
App.editor.inputmask('#email', {
    mask: '*{1,20}[.*{1,20}]@*{1,20}.*{2,6}[.*{1,2}]'
});
```

### App.editor.clipboard(elem, options) / 剪贴板 / Clipboard

初始化剪贴板复制功能。

**语法 / Syntax:**
```javascript
App.editor.clipboard(elem, options)
```

---

## 编辑器切换 / Editor Switching

### App.editor.switch(editorName, texta, cancelFn, tips) / 编辑器切换 / Editor Switch

在不同类型的编辑器之间切换。

**语法 / Syntax:**
```javascript
App.editor.switch(editorName, texta, cancelFn, tips)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `editorName` | String | 目标编辑器类型: 'tinymce', 'editormd', 'vditor', 'text' / Target editor type |
| `texta` | Element | 文本域元素 / Textarea element |
| `cancelFn` | Function | 取消回调 / Cancel callback |
| `tips` | Boolean | 是否显示提示 / Whether to show tips |

**编辑器类型映射 / Editor Type Mapping:**

| 编辑器名 / Editor Name | 类型 / Type | 描述 / Description |
|---------------------|-----------|-------------------|
| 'tinymce' | html | HTML 编辑器 / HTML editor |
| 'ckeditor' | html | HTML 编辑器 / HTML editor |
| 'ueditor' | html | HTML 编辑器 / HTML editor |
| 'editormd' | markdown | Markdown 编辑器 / Markdown editor |
| 'markdown' | markdown | Markdown 编辑器 / Markdown editor |
| 'vditor' | markdown | Markdown 编辑器 / Markdown editor |
| 'text' | text | 纯文本编辑器 / Plain text editor |

### App.editor.switcher(swicherElem, contentElem, defaultEditorName) / 编辑器切换器 / Editor Switcher

初始化编辑器切换器。

**语法 / Syntax:**
```javascript
App.editor.switcher(swicherElem, contentElem, defaultEditorName)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `swicherElem` | Selector | 切换按钮元素 / Switch button element |
| `contentElem` | Selector | 内容元素 / Content element |
| `defaultEditorName` | String | 默认编辑器名 / Default editor name |

**示例 / Example:**

```html
<!-- 切换器 / Switcher -->
<div class="btn-group" id="editor-switcher">
    <button type="button" class="btn btn-default" value="text">纯文本</button>
    <button type="button" class="btn btn-default" value="markdown">Markdown</button>
    <button type="button" class="btn btn-default" value="tinymce">富文本</button>
</div>

<!-- 编辑器容器 / Editor container -->
<textarea id="content" data-editor-name="tinymce"></textarea>

<script>
App.editor.switcher('#editor-switcher', '#content', 'tinymce');
</script>
```

---

## 工具函数 / Utility Functions

### App.editor.loadingOverlay(options) / 加载遮罩 / Loading Overlay

创建加载遮罩层。

**语法 / Syntax:**
```javascript
App.editor.loadingOverlay(options)
```

**示例 / Example:**

```javascript
var loader = App.editor.loadingOverlay({
    zIndex: 9999,
    image: '<i class="fa fa-spinner fa-spin"></i>',
    text: 'Loading...'
});

loader.show();
loader.hide();
```

### App.editor.dialog(options) / 对话框 / Dialog

创建 Bootstrap 对话框。

**语法 / Syntax:**
```javascript
App.editor.dialog(options)
```

**示例 / Example:**

```javascript
App.editor.dialog({
    title: 'Confirm',
    message: 'Are you sure you want to delete?',
    buttons: [{
        label: 'Cancel',
        action: function(dialogItself){
            dialogItself.close();
        }
    }, {
        label: 'OK',
        cssClass: 'btn-primary',
        action: function(dialogItself){
            // 执行操作
            dialogItself.close();
        }
    }]
});
```

### App.editor.popup(elem, options, callback) / 弹出层 / Popup

初始化弹出层（图片预览等）。

**语法 / Syntax:**
```javascript
App.editor.popup(elem, options, callback)
```

**示例 / Example:**

```javascript
// 图片预览 / Image preview
App.editor.popup('.image-zoom', {
    type: 'image',
    closeBtnInside: false,
    mainClass: 'mfp-with-zoom',
    zoom: {
        enabled: true,
        duration: 300,
        easing: 'ease-in-out'
    }
});
```

### App.editor.galleryPopup(elem, options, callback) / 图片库弹出 / Gallery Popup

初始化图片库弹出层（支持多图浏览）。

**语法 / Syntax:**
```javascript
App.editor.galleryPopup(elem, options, callback)
```

**示例 / Example:**

```javascript
App.editor.galleryPopup('.gallery-item', {
    closeBtnInside: false,
    zoom: {opener: null},
    gallery: {
        enabled: true, 
        navigateByImgClick: true
    }
});
```

---

## 多语言编辑器 / Multilingual Editor

### App.editor.multilingualContentEditor(formContainer, contentElem, uploadUrl, helpBlock, beforeSwitchCallback) / 多语言内容编辑器 / Multilingual Content Editor

初始化多语言内容编辑器。

**语法 / Syntax:**
```javascript
App.editor.multilingualContentEditor(formContainer, contentElem, uploadUrl, helpBlock, beforeSwitchCallback)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `formContainer` | Selector | 表单容器 / Form container |
| `contentElem` | Selector | 内容元素 / Content element |
| `uploadUrl` | String | 上传 URL / Upload URL |
| `helpBlock` | String | 帮助块 HTML / Help block HTML |
| `beforeSwitchCallback` | Function | 切换前回调 / Callback before switch |

**示例 / Example:**

```html
<div id="form-container">
    <div data-editor-name="tinymce" id="content-en">
        <textarea name="content" class="form-control"></textarea>
    </div>
    <div data-editor-name="tinymce" id="content-zh">
        <textarea name="content" class="form-control"></textarea>
    </div>
</div>

<script>
App.editor.multilingualContentEditor('#form-container', 'textarea', '/api/upload');
</script>
```

---

## 图片裁剪 / Image Cropping

### App.editor.cropImage(uploadURL, thumbnailElem, originalElem, type, width, height) / 图片裁剪 / Image Cropping

初始化图片裁剪组件。

**语法 / Syntax:**
```javascript
App.editor.cropImage(uploadURL, thumbnailElem, originalElem, type, width, height)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `uploadURL` | String | 上传 URL / Upload URL |
| `thumbnailElem` | Selector | 缩略图元素 / Thumbnail element |
| `originalElem` | Selector | 原图元素 / Original image element |
| `type` | String | 裁剪类型 / Crop type |
| `width` | Number | 目标宽度 / Target width |
| `height` | Number | 目标高度 / Target height |

---

## 工具函数 / Utility Functions

### App.utils.elemToId(elem) / 元素转 ID / Element to ID

获取或生成元素的 ID。

**语法 / Syntax:**
```javascript
App.utils.elemToId(elem)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `elem` | Selector/Object | 元素选择器或对象 / Element selector or object |

**返回值 / Returns:** 元素的 ID 字符串 / Element ID string

### App.utils.unixtime() / Unix 时间戳 / Unix Timestamp

获取当前 Unix 时间戳。

**语法 / Syntax:**
```javascript
App.utils.unixtime()
```

**返回值 / Returns:** 毫秒级时间戳 / Millisecond timestamp

**示例 / Example:**

```javascript
var timestamp = App.utils.unixtime();
// 1640995200000
```

---

## 完整示例 / Complete Examples

### 示例 1: 带上传的富文本编辑器 / Example 1: Rich Text Editor with Upload

```javascript
// HTML
<textarea id="content" name="content" class="form-control" action="/api/upload"></textarea>

// JavaScript
App.editor.tinymce('#content', '/api/upload', {
    height: 500,
    plugins: 'image media link table code',
    toolbar: 'undo redo | bold italic | link image media | table | code'
});
```

### 示例 2: Markdown 编辑器 / Example 2: Markdown Editor

```javascript
// HTML
<textarea id="markdown-content" class="form-control" action="/api/upload"></textarea>

// JavaScript
App.editor.markdown('#markdown-content', '/api/upload', {
    height: 600,
    theme: 'dark',
    previewTheme: 'dark'
});
```

### 示例 3: 代码编辑器 / Example 3: Code Editor

```javascript
// HTML
<textarea id="code-editor"># Write your code here</textarea>

// JavaScript
App.editor.codemirror('#code-editor', {
    mode: 'text/javascript',
    lineNumbers: true,
    theme: 'monokai'
});
```

### 示例 4: 文件上传 / Example 4: File Upload

```html
<div id="file-upload">
    <button class="btn btn-primary" data-toggle="uploadPreviewer" data-upload-url="/api/upload">
        <i class="fa fa-cloud-upload"></i> 上传文件
    </button>
    <div class="preview-container"></div>
</div>

<script>
App.editor.fileInput('#file-upload', {}, function(fileURL) {
    console.log('Upload success:', fileURL);
    App.message({text: '上传成功', type: 'success'});
}, function(error) {
    console.error('Upload error:', error);
    App.message({text: '上传失败', type: 'error'});
});
</script>
```

### 示例 5: 编辑器切换 / Example 5: Editor Switching

```html
<div class="btn-group" id="editor-switcher">
    <button type="button" class="btn btn-default" value="text">纯文本</button>
    <button type="button" class="btn btn-default" value="markdown">Markdown</button>
    <button type="button" class="btn btn-default" value="tinymce">富文本</button>
</div>

<textarea id="content" data-editor-name="tinymce" action="/api/upload"></textarea>

<script>
App.editor.switcher('#editor-switcher', '#content', 'tinymce');
</script>
```

### 示例 6: 日期选择器 / Example 6: Date Picker

```javascript
// 日期范围 / Date range
<input id="daterange" type="text" class="form-control" />

<script>
App.editor.dateRangePicker('#daterange', {
    format: 'YYYY-MM-DD',
    locale: {
        applyLabel: '确定',
        cancelLabel: '取消',
        fromLabel: '起始',
        toLabel: '结束',
        daysOfWeek: ['日', '一', '二', '三', '四', '五', '六'],
        monthNames: ['一月', '二月', '三月', '四月', '五月', '六月', '七月', '八月', '九月', '十月', '十一月', '十二月']
    }
});
</script>
```

---

## 最佳实践 / Best Practices

### 1. 编辑器初始化 / Editor Initialization

```javascript
// 推荐：使用 data 属性配置 / Recommended: Use data attributes
<textarea id="content" 
          data-editor-name="tinymce" 
          action="/api/upload"
          data-tinymce-options='{"height": 500}'></textarea>

<script>
App.editor.tinymce('#content');
</script>
```

### 2. 错误处理 / Error Handling

```javascript
App.editor.tinymce('#content', '/api/upload', {}, true)
    .then(function(editor) {
        console.log('Editor initialized');
    })
    .catch(function(error) {
        console.error('Editor initialization failed:', error);
        App.message({text: '编辑器初始化失败', type: 'error'});
    });
```

### 3. 文件上传验证 / File Upload Validation

```javascript
App.editor.dropzone('#upload-zone', {
    url: '/api/upload',
    maxFilesize: 2,  // MB
    acceptedFiles: 'image/*',
    maxFiles: 5
}, function(file, resp) {
    if (resp.error) {
        App.message({text: resp.error, type: 'error'});
        return;
    }
    App.message({text: '上传成功', type: 'success'});
});
```

### 4. 多语言内容编辑 / Multilingual Content Editing

```javascript
// 确保所有语言标签页都初始化编辑器
$('.langset textarea[data-editor-name]').each(function() {
    var lang = $(this).closest('.langset').find('.nav-tabs .active a').data('lang');
    var uploadUrl = '/api/upload?lang=' + lang;
    
    App.editor.tinymce(this, uploadUrl, {
        height: 400,
        plugins: 'basic'
    });
});
```

---

## 浏览器兼容性 / Browser Compatibility

- **Chrome/Edge**: 完全支持 / Full support
- **Firefox**: 完全支持 / Full support
- **Safari**: 完全支持 / Full support
- **IE11+**: 部分支持（建议使用现代浏览器）/ Partial support (modern browsers recommended)

---

## 依赖项 / Dependencies

### 必需依赖 / Required Dependencies

- **jQuery**: >= 1.9.0
- **TinyMCE**: >= 5.0
- **EditorMD**: >= 1.5
- **CodeMirror**: >= 5.0
- **Bootstrap**: >= 3.3.7

### 可选依赖 / Optional Dependencies

- **jQuery UI**: 用于拖拽等功能 / For drag & drop features
- **Dropzone**: 文件上传 / File upload
- **Bootstrap Dialog**: 对话框 / Dialogs
- **Magnific Popup**: 弹出层 / Popups
- **Select2**: 增强选择框 / Enhanced select
- **InputMask**: 输入掩码 / Input mask

---

## 加载器配置 / Loader Configuration

依赖 `/public/assets/backend/js/behaviour/general.js` 和 `/public/assets/backend/js/loader/loader.js`  
Depends on `/public/assets/backend/js/behaviour/general.js` and `/public/assets/backend/js/loader/loader.js` 

### App.loader.libs / 库文件配置 / Library Files Configuration

```javascript
App.loader.libs = {
    editormd: ['#editor/markdown/css/editormd.min.css', '#editor/markdown/editormd.min.js'],
    codemirror: ['#editor/markdown/lib/codemirror/codemirror.min.css', ...],
    tinymce: ['#editor/tinymce/custom.css', '#editor/tinymce/tinymce.min.js', ...],
    dropzone: ['#jquery.ui/css/dropzone.min.css', '#dropzone/dropzone.min.js'],
    select2: ['#jquery.select2/select2.css', '#jquery.select2/select2.min.js'],
    dateRangePicker: ['#daterangepicker/daterangepicker.min.css', ...],
    magnificPopup: ['#magnific-popup/magnific-popup.min.css', ...],
    inputmask: ['#inputmask/inputmask.min.js', '#inputmask/jquery.inputmask.min.js'],
    clipboard: ['#clipboard/clipboard.min.js', '#clipboard/utils.js']
    // ... 更多库 / more libraries
};
```

