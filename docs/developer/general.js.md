# General.js 文档 / General.js Documentation

## 概述 / Overview

`general.js` 是 Nging 系统的核心 JavaScript 模块，提供了丰富的前端功能，包括 UI 组件初始化、AJAX 请求、消息通知、表单处理、导航管理等功能。

---

## 目录 / Table of Contents

1. [配置 / Configuration](#配置--configuration)
2. [初始化 / Initialization](#初始化--initialization)
3. [国际化 / Internationalization](#国际化--internationalization)
4. [UI 组件 / UI Components](#ui-组件--ui-components)
5. [导航管理 / Navigation Management](#导航管理--navigation-management)
6. [AJAX 相关 / AJAX Related](#ajax-相关--ajax-related)
7. [消息通知 / Message Notifications](#消息通知--message-notifications)
8. [表单操作 / Form Operations](#表单操作--form-operations)
9. [工具函数 / Utility Functions](#工具函数--utility-functions)

---

## 配置 / Configuration

### 基础配置 / Basic Configuration

```javascript
var config = {
    tooltip: true,        // 启用提示框 / Enable tooltips
    popover: true,        // 启用弹出框 / Enable popovers
    nanoScroller: true,   // 启用滚动条 / Enable scroller
    nestableLists: true,  // 启用可嵌套列表 / Enable nestable lists
    hiddenElements: true,  // 启用隐藏元素绑定 / Enable hidden elements binding
    bootstrapSwitch: true, // 启用开关组件 / Enable switch component
    dateTime: true,       // 启用日期时间选择器 / Enable datetime picker
    select2: true,        // 启用 Select2 / Enable Select2
    tags: true,           // 启用标签输入 / Enable tags input
    slider: true          // 启用滑块 / Enable slider
};
```

---

## 初始化 / Initialization

### App.init(options) / 应用初始化

初始化应用的所有 UI 组件和功能。

**语法 / Syntax:**
```javascript
App.init(options)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `options` | Object | 配置选项，会合并到默认配置中 / Config options to merge with defaults |

**示例 / Example:**
```javascript
App.init({
    tooltip: true,
    popover: true,
    nanoScroller: true
});
```

---

## 国际化 / Internationalization

### App.i18n / 国际化文本 / Internationalization Text

内置的国际化文本字典。

**可用键 / Available Keys:**

```javascript
{
    SYS_INFO: 'System Information',
    UPLOAD_ERR: 'Upload Error',
    PLEASE_SELECT_FOR_OPERATE: 'Please select the item you want to operate',
    PLEASE_SELECT_FOR_REMOVE: 'Please select the item you want to delete',
    CONFIRM_REMOVE: 'Are you sure you want to delete them?',
    SELECTED_ITEMS: 'You have selected %d items',
    SUCCESS: 'The operation was successful',
    FAILURE: 'Operation failed',
    UPLOADING: 'File uploading, please wait...',
    UPLOAD_SUCCEED: 'Upload successfully',
    BUTTON_UPLOAD: 'Upload'
}
```

### App.t(key, ...args) / 翻译函数 / Translation Function

获取翻译文本，支持格式化参数。

**语法 / Syntax:**
```javascript
App.t(key, ...args)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `key` | String | 翻译键 / Translation key |
| `...args` | Any | 格式化参数 / Format arguments |

**返回值 / Returns:** 翻译后的字符串 / Translated string

**示例 / Example:**
```javascript
// 简单翻译
var text = App.t('SYS_INFO');

// 带参数的翻译
var message = App.t('SELECTED_ITEMS', 5);
// 结果 / Result: "You have selected 5 items"

// 使用 # 语法指定格式化模板
var msg = App.t('#You have selected %d items', 10);
```

### App.langInfo() / 获取语言信息 / Get Language Information

获取当前语言的语言信息对象。

**语法 / Syntax:**
```javascript
App.langInfo()
```

**返回值 / Returns:**
```javascript
{
    encoding: 'zh',    // 语言代码 / Language code
    country: 'CN'      // 国家代码 / Country code
}
```

### App.langTag(seperator) / 获取语言标签 / Get Language Tag

获取语言标签字符串，如 `zh-CN`、`en-US`。

**语法 / Syntax:**
```javascript
App.langTag(seperator)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 默认值 / Default | 描述 / Description |
|-------------------|-------------|-----------------|-------------------|
| `seperator` | String | `-` | 分隔符 / Separator |

**示例 / Example:**
```javascript
App.langTag()        // 'zh-CN'
App.langTag('_')      // 'zh_CN'
```

---

## UI 组件 / UI Components

### App.initLeftNav() / 初始化左侧导航 / Initialize Left Navigation

初始化左侧垂直导航菜单。

**语法 / Syntax:**
```javascript
App.initLeftNav()
```

**功能 / Features:**
- 展开/折叠子菜单 / Expand/collapse submenus
- 侧边栏折叠时显示悬浮菜单 / Show floating menu when sidebar collapsed
- 响应式菜单 / Responsive menu

### App.initTool() / 初始化工具菜单 / Initialize Tool Menu

初始化悬浮工具菜单。

**语法 / Syntax:**
```javascript
App.initTool()
```

### App.toggleSideBar() / 切换侧边栏 / Toggle Sidebar

切换侧边栏的展开/折叠状态。

**语法 / Syntax:**
```javascript
App.toggleSideBar()
```

### App.wizard() / 初始化向导 / Initialize Wizard

初始化表单向导组件。

**语法 / Syntax:**
```javascript
App.wizard()
```

### App.pageAside() / 页面侧边栏 / Page Aside

初始化页面侧边栏的折叠/展开功能。

**语法 / Syntax:**
```javascript
App.pageAside()
```

### App.tableReponsiveInit() / 表格响应式初始化 / Table Responsive Init

初始化表格的响应式布局。

**语法 / Syntax:**
```javascript
App.tableReponsiveInit()
```

### App.returnToTopButton() / 返回顶部按钮 / Return to Top Button

创建返回顶部的浮动按钮。

**语法 / Syntax:**
```javascript
App.returnToTopButton()
```

---

## 导航管理 / Navigation Management

### App.markNavByURL(url) / 根据 URL 标记导航 / Mark Navigation by URL

根据当前 URL 自动标记激活的导航项。

**语法 / Syntax:**
```javascript
App.markNavByURL(url)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `url` | String | URL 路径，默认为当前页面 URL / URL path, defaults to current page URL |

### App.markNav(curNavA, position) / 标记导航 / Mark Navigation

标记导航项为激活状态。

**语法 / Syntax:**
```javascript
App.markNav(curNavA, position)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `curNavA` | jQuery Object | 导航链接元素 / Navigation link element |
| `position` | String | 位置: `'left'` 或 `'top'` / Position: `'left'` or `'top'` |

### App.unmarkNav(curNavA, position) / 取消标记导航 / Unmark Navigation

取消导航的激活状态。

**语法 / Syntax:**
```javascript
App.unmarkNav(curNavA, position)
```

---

## AJAX 相关 / AJAX Related

### App.attachAjaxURL(elem) / 附加 AJAX URL 事件 / Attach AJAX URL Events

为元素附加 AJAX 请求事件处理。

**语法 / Syntax:**
```javascript
App.attachAjaxURL(elem)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 默认值 / Default | 描述 / Description |
|-------------------|-------------|-----------------|-------------------|
| `elem` | Selector/Element | `document` | 容器元素 / Container element |

**数据属性 / Data Attributes:**
- `data-ajax-url` - 请求 URL / Request URL
- `data-ajax-method` - 请求方法 / Request method (get/post)
- `data-ajax-params` - 请求参数 / Request parameters
- `data-ajax-confirm` - 确认消息 / Confirm message
- `data-ajax-accept` - 响应类型 / Response type (json/html)
- `data-ajax-target` - 响应目标元素 / Response target element
- `data-ajax-callback` - 回调函数名 / Callback function name
- `data-ajax-onsuccess` - 成功回调 / Success callback
- `data-ajax-reload` - 成功后是否重载页面 / Reload after success

**示例 / Example:**
```html
<a href="#" 
   data-ajax-url="/api/delete" 
   data-ajax-method="post"
   data-ajax-confirm="Are you sure?"
   data-ajax-accept="json">
   Delete
</a>
```

### App.startLazyload(elem) / 启动懒加载 / Start Lazy Load

为元素启动懒加载功能。

**语法 / Syntax:**
```javascript
App.startLazyload(elem)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 默认值 / Default | 描述 / Description |
|-------------------|-------------|-----------------|-------------------|
| `elem` | Selector/Element | `document` | 容器元素 / Container element |

**数据属性 / Data Attributes:**
- `lazyload-url` - 加载 URL / Load URL
- `lazyload-method` - 请求方法 / Request method
- `lazyload-params` - 请求参数 / Request parameters
- `lazyload-accept` - 响应类型 / Response type
- `lazyload-target` - 目标元素 / Target element

### App.postFormData(form, postData, success, error, accept) / 提交表单数据 / Post Form Data

通过 AJAX 提交表单数据。

**语法 / Syntax:**
```javascript
App.postFormData(form, postData, success, error, accept)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `form` | Selector/Element | 表单元素 / Form element |
| `postData` | Object | 额外数据 / Additional data |
| `success` | Function | 成功回调 / Success callback |
| `error` | Function | 错误回调 / Error callback |
| `accept` | String | 响应类型，默认 'json' / Response type, default 'json' |

### App.attachPjax(elem, callbacks, timeout) / 附加 PJAX / Attach PJAX

为链接附加 PJAX 无刷新加载功能。

**语法 / Syntax:**
```javascript
App.attachPjax(elem, callbacks, timeout)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 默认值 / Default | 描述 / Description |
|-------------------|-------------|-----------------|-------------------|
| `elem` | Selector | `'a'` | 链接选择器 / Link selector |
| `callbacks` | Object | `{}` | 回调函数 / Callback functions |
| `timeout` | Number | `5000` | 超时时间(毫秒) / Timeout in ms |

**回调选项 / Callback Options:**
- `onclick` - 点击时 / On click
- `onsend` - 发送请求时 / On send
- `oncomplete` - 完成时 / On complete
- `ontimeout` - 超时时 / On timeout
- `onstart` - 开始时 / On start
- `onend` - 结束时 / On end

---

## 消息通知 / Message Notifications

### App.message(options, sticky) / 显示消息 / Show Message

显示消息通知。

**语法 / Syntax:**
```javascript
App.message(options, sticky)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `options` | String/Object | 消息内容或配置对象 / Message content or config object |
| `sticky` | Boolean | 是否保持显示 / Whether to keep showing |

**配置选项 / Options:**

```javascript
{
    title: 'Title',           // 标题 / Title
    text: 'Message',          // 内容 / Text
    image: 'avatar.png',      // 头像图片 / Avatar image
    class_name: 'success',     // 样式: success/info/warning/error/danger/clean/primary/dark
    sticky: false,            // 是否保持显示 / Whether to keep showing
    time: 5000,              // 显示时间(毫秒) / Display time in ms
    after_open: function(){}   // 打开后回调 / Callback after open
}
```

**示例 / Example:**
```javascript
// 简单消息
App.message('Operation successful');

// 完整配置
App.message({
    title: 'Success',
    text: 'Data saved successfully',
    type: 'success',
    time: 3000
});

// 保持显示的消息
App.message('Please wait...', true);

// 清除所有消息
App.message('clear');

// 清除指定消息
App.message('remove', messageId);
```

### App.notifyListen() / 监听通知 / Listen to Notifications

启动实时通知监听（WebSocket 或 SSE）。

**语法 / Syntax:**
```javascript
App.notifyListen()
```

**功能 / Features:**
- 自动选择 WebSocket 或 Server-Sent Events
- 支持多种通知模式: notify、element、modal
- 自动重连 / Auto reconnect

---

## 表单操作 / Form Operations

### App.checkedAll(ctrl, target) / 全选 / Check All

控制一组复选框的全选/取消全选。

**语法 / Syntax:**
```javascript
App.checkedAll(ctrl, target)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `ctrl` | Element | 控制复选框 / Control checkbox |
| `target` | Selector | 目标复选框组 / Target checkbox group |

**示例 / Example:**
```javascript
// HTML
<input type="checkbox" id="checkAll" />
<input type="checkbox" name="items[]" value="1" />
<input type="checkbox" name="items[]" value="2" />

// JS
App.checkedAll('#checkAll', 'input[name="items[]"]');
```

### App.attachCheckedAll(ctrl, target, showNumElem) / 附加全选事件 / Attach Check All Event

为全选复选框附加事件处理。

**语法 / Syntax:**
```javascript
App.attachCheckedAll(ctrl, target, showNumElem)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `ctrl` | Selector | 全选复选框 / Check all checkbox |
| `target` | Selector | 目标复选框组 / Target checkbox group |
| `showNumElem` | Selector | 显示选中数量的元素 / Element to show selected count |

### App.passwordInputShowPassword(container) / 显示密码输入 / Show Password Input

为密码输入框添加显示/隐藏密码功能。

**语法 / Syntax:**
```javascript
App.passwordInputShowPassword(container)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 默认值 / Default | 描述 / Description |
|-------------------|-------------|-----------------|-------------------|
| `container` | Selector | `null` | 容器元素 / Container element |

**示例 / Example:**
```html
<div class="input-group">
    <input type="password" name="password" />
    <span class="input-group-btn">
        <a class="show-password" data-target="input[name='password']">
            <i class="fa fa-eye"></i> 显示
        </a>
    </span>
</div>

<script>
App.passwordInputShowPassword();
</script>
```

### App.editableSortNumber(container, url, callback) / 可编辑排序号 / Editable Sort Number

使表格中的排序号可编辑。

**语法 / Syntax:**
```javascript
App.editableSortNumber(container, url, callback)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `container` | Selector | 表格容器 / Table container |
| `url` | String | 更新 URL / Update URL |
| `callback` | Function | 回调函数 / Callback function |

### App.opSelected(elem, postField, removeURL, callback, confirmMsg, unselectedMsg) / 操作选中项 / Operate Selected Items

对选中的项目执行操作（如删除）。

**语法 / Syntax:**
```javascript
App.opSelected(elem, postField, removeURL, callback, confirmMsg, unselectedMsg)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `elem` | Selector | 复选框选择器 / Checkbox selector |
| `postField` | String | 提交字段名 / Submit field name |
| `removeURL` | String | 请求 URL / Request URL |
| `callback` | Function | 回调函数 / Callback function |
| `confirmMsg` | String | 确认消息 / Confirm message |
| `unselectedMsg` | String | 未选择提示 / Unselected message |

**示例 / Example:**
```javascript
App.opSelected(
    'input[name="id[]"]',
    'id',
    '/api/delete',
    function() {
        console.log('Deleted successfully');
    },
    'Are you sure you want to delete these items?',
    'Please select items to delete'
);
```

---

## 工具函数 / Utility Functions

### App.getJQueryObject(a) / 获取 jQuery 对象 / Get jQuery Object

确保返回 jQuery 对象。

**语法 / Syntax:**
```javascript
App.getJQueryObject(a)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `a` | Selector/jQuery Object | 元素选择器或 jQuery 对象 / Element selector or jQuery object |

**返回值 / Returns:** jQuery 对象 / jQuery object

### App.htmlEncode(value) / HTML 编码 / HTML Encode

对文本进行 HTML 实体编码。

**语法 / Syntax:**
```javascript
App.htmlEncode(value)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `value` | String | 要编码的文本 / Text to encode |

**返回值 / Returns:** 编码后的字符串 / Encoded string

**编码映射 / Encoding Mapping:**
```javascript
{
    "&": "&amp;",
    "<": "&lt;",
    ">": "&gt;",
    " ": "&nbsp;",
    "'": "&#39;",
    '"': "&quot;"
}
```

### App.htmlDecode(value) / HTML 解码 / HTML Decode

对 HTML 实体进行解码。

**语法 / Syntax:**
```javascript
App.htmlDecode(value)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `value` | String | 要解码的文本 / Text to decode |

**返回值 / Returns:** 解码后的字符串 / Decoded string

### App.text2html(text, noescape) / 文本转 HTML / Text to HTML

将纯文本转换为 HTML，处理换行符。

**语法 / Syntax:**
```javascript
App.text2html(text, noescape)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `text` | String | 文本内容 / Text content |
| `noescape` | Boolean | 是否跳过转义 / Whether to skip escaping |

### App.textNl2br(text) / 换行符转换 / Newline to Break

将换行符转换为 `<br>` 标签。

**语法 / Syntax:**
```javascript
App.textNl2br(text)
```

**转换映射 / Conversion Mapping:**
```javascript
{
    "\n": '<br />',
    "  ": '&nbsp; ',
    "\t": '&nbsp; &nbsp; '
}
```

### App.trimSpace(text) / 去除空格 / Trim Spaces

去除字符串首尾的空格。

**语法 / Syntax:**
```javascript
App.trimSpace(text)
```

### App.randomString(len) / 生成随机字符串 / Generate Random String

生成指定长度的随机字符串。

**语法 / Syntax:**
```javascript
App.randomString(len)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 默认值 / Default | 描述 / Description |
|-------------------|-------------|-----------------|-------------------|
| `len` | Number | `32` | 字符串长度 / String length |

**返回值 / Returns:** 随机字符串 / Random string

**注意 / Note:** 已排除容易混淆的字符 oOLl,9gq,Vv,Uu,I1

### App.formatBytes(bytes, precision) / 格式化字节数 / Format Bytes

将字节数格式化为易读的文件大小。

**语法 / Syntax:**
```javascript
App.formatBytes(bytes, precision)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 默认值 / Default | 描述 / Description |
|-------------------|-------------|-----------------|-------------------|
| `bytes` | Number | - | 字节数 / Bytes |
| `precision` | Number | `2` | 小数位数 / Decimal places |

**示例 / Example:**
```javascript
App.formatBytes(1024)      // "1.00 KB"
App.formatBytes(1048576)   // "1.00 MB"
App.formatBytes(1073741824) // "1.00 GB"
```

### App.replaceURLParam(name, value, url) / 替换 URL 参数 / Replace URL Parameter

替换或添加 URL 参数。

**语法 / Syntax:**
```javascript
App.replaceURLParam(name, value, url)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `name` | String | 参数名 / Parameter name |
| `value` | String | 参数值 / Parameter value |
| `url` | String | URL，默认为当前页面 URL / URL, defaults to current page URL |

**返回值 / Returns:** 修改后的 URL / Modified URL

**示例 / Example:**
```javascript
App.replaceURLParam('page', 2, '/list?page=1')
// 返回 / Returns: "/list?page=2"
```

### App.insertAtCursor(myField, myValue, posStart, posEnd) / 在光标位置插入 / Insert at Cursor

在文本域的光标位置插入文本。

**语法 / Syntax:**
```javascript
App.insertAtCursor(myField, myValue, posStart, posEnd)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `myField` | Element | 文本域元素 / Textarea element |
| `myValue` | String | 要插入的文本 / Text to insert |
| `posStart` | Number | 起始位置 / Start position |
| `posEnd` | Number | 结束位置 / End position |

### App.loading(op) / 加载状态 / Loading State

显示或隐藏加载指示器。

**语法 / Syntax:**
```javascript
App.loading(op)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `op` | String | 操作: `'show'` 或 `'hide'` / Operation: 'show' or 'hide' |

**示例 / Example:**
```javascript
App.loading('show');  // 显示加载 / Show loading
App.loading('hide');  // 隐藏加载 / Hide loading
```

### App.getClientID(type) / 获取客户端 ID / Get Client ID

获取客户端通知 ID。

**语法 / Syntax:**
```javascript
App.getClientID(type)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 默认值 / Default | 描述 / Description |
|-------------------|-------------|-----------------|-------------------|
| `type` | String | `'notify'` | ID 类型 / ID type |

**返回值 / Returns:** 客户端 ID / Client ID

### App.setClientID(data, type) / 设置客户端 ID / Set Client ID

将客户端 ID 添加到请求数据中。

**语法 / Syntax:**
```javascript
App.setClientID(data, type)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `data` | Object/String | 请求数据 / Request data |
| `type` | String | ID 类型 / ID type |

**返回值 / Returns:** 添加了 ID 的数据 / Data with ID added

---

## 布局相关 / Layout Related

### App.topFloat(elems, top, autoWith) / 顶部浮动 / Top Float

使元素在滚动时固定在顶部。

**语法 / Syntax:**
```javascript
App.topFloat(elems, top, autoWith)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `elems` | Selector | 元素选择器 / Element selector |
| `top` | Number | 顶部偏移量 / Top offset |
| `autoWith` | Boolean | 是否自动设置宽度 / Whether to auto set width |

### App.bottomFloat(elems, top, autoWith) / 底部浮动 / Bottom Float

使元素在滚动时固定在底部。

**语法 / Syntax:**
```javascript
App.bottomFloat(elems, top, autoWith)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `elems` | Selector | 元素选择器 / Element selector |
| `top` | Number | 底部偏移量 / Bottom offset |
| `autoWith` | Boolean | 是否自动设置宽度 / Whether to auto set width |

### App.autoFixedThead(prefix) / 自动固定表头 / Auto Fixed Table Header

自动使表格的表头在滚动时固定。

**语法 / Syntax:**
```javascript
App.autoFixedThead(prefix)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 默认值 / Default | 描述 / Description |
|-------------------|-------------|-----------------|-------------------|
| `prefix` | String | `''` | 选择器前缀 / Selector prefix |

**示例 / Example:**
```html
<table class="table">
    <thead class="auto-fixed">
        <tr>
            <th>Name</th>
            <th>Email</th>
        </tr>
    </thead>
    <tbody>
        <!-- 数据行 / Data rows -->
    </tbody>
</table>

<script>
App.autoFixedThead();
</script>
```

---

## WebSocket 相关 / WebSocket Related

### App.wsURL(url) / WebSocket URL / WebSocket URL

构建 WebSocket URL，自动处理协议和主机。

**语法 / Syntax:**
```javascript
App.wsURL(url)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `url` | String | 相对或绝对 URL / Relative or absolute URL |

**返回值 / Returns:** 完整的 WebSocket URL / Complete WebSocket URL

**示例 / Example:**
```javascript
App.wsURL('/ws/chat')           // "ws://example.com/ws/chat"
App.wsURL('wss://example.com/ws/chat')  // "wss://example.com/ws/chat"
```

### App.websocket(showmsg, url, onopen, onclose) / WebSocket 连接 / WebSocket Connection

创建 WebSocket 连接。

**语法 / Syntax:**
```javascript
App.websocket(showmsg, url, onopen, onclose)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `showmsg` | Function | 消息处理函数 / Message handler |
| `url` | String | WebSocket URL / WebSocket URL |
| `onopen` | Function/Object | 打开回调 / Open callback |
| `onclose` | Function | 关闭回调 / Close callback |

**返回值 / Returns:** WebSocket 对象 / WebSocket object

---

## 警告框相关 / Alert Related

### App.alertBlock(content, title, type) / 警告框 / Alert Block

生成警告框 HTML。

**语法 / Syntax:**
```javascript
App.alertBlock(content, title, type)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `content` | String | 内容 / Content |
| `title` | String | 标题 / Title |
| `type` | String | 类型: 'success', 'info', 'warn', 'error' / Type |

**示例 / Example:**
```javascript
App.alertBlock('Operation successful', 'Success', 'success');
```

### App.alertBlockx(content, title, type) / 白色警告框 / White Alert Block

生成白色圆角警告框 HTML。

**语法 / Syntax:**
```javascript
App.alertBlockx(content, title, type)
```

---

## 其他功能 / Other Features

### App.switchLang(lang) / 切换语言 / Switch Language

切换应用语言。

**语法 / Syntax:**
```javascript
App.switchLang(lang)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `lang` | String | 语言代码 / Language code (e.g., 'zh-CN', 'en-US') |

### App.reportBug(url) / 报告错误 / Report Bug

发送错误报告到服务器。

**语法 / Syntax:**
```javascript
App.reportBug(url)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `url` | String | 报告 URL / Report URL |

### App.pushState(data, title, url) / 推入历史 / Push State

向浏览器历史记录推入新状态。

**语法 / Syntax:**
```javascript
App.pushState(data, title, url)
```

### App.replaceState(data, title, url) / 替换历史 / Replace State

替换当前浏览器历史记录状态。

**语法 / Syntax:**
```javascript
App.replaceState(data, title, url)
```

### App.parseBool(b) / 解析布尔值 / Parse Boolean

将字符串解析为布尔值。

**语法 / Syntax:**
```javascript
App.parseBool(b)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `b` | Any | 要解析的值 / Value to parse |

**返回值 / Returns:** 布尔值 / Boolean value

**视为 false 的值 / Values considered false:**
- `'0'`, `'false'`, `'n'`, `'no'`, `'off'`, `null`

**示例 / Example:**
```javascript
App.parseBool('true')    // true
App.parseBool('false')   // false
App.parseBool('yes')     // true
App.parseBool('no')      // false
```

---

## 最佳实践 / Best Practices

### 1. 使用国际化 / Use Internationalization

```javascript
// 推荐 / Recommended
App.message({ text: App.t('SAVE_SUCCESS'), type: 'success' });

// 避免 / Avoid
App.message({ text: '保存成功', type: 'success' });
```

### 2. AJAX 请求处理 / AJAX Request Handling

```javascript
// 使用 data-ajax-url 属性 / Use data-ajax-url attribute
<a href="#"
   data-ajax-url="/api/delete"
   data-ajax-method="post"
   data-ajax-confirm="Are you sure?"
   data-ajax-accept="json">
   Delete
</a>

<script>
App.attachAjaxURL();
</script>
```

### 3. 表单验证 / Form Validation

```javascript
// 全选框使用 / Check all checkbox usage
<input type="checkbox" id="checkAll" />
<input type="checkbox" name="items[]" value="1" />
<input type="checkbox" name="items[]" value="2" />

<script>
App.attachCheckedAll('#checkAll', 'input[name="items[]"]', '#selectedCount');
</script>
```

### 4. 响应式表头 / Responsive Table Header

```html
<thead class="auto-fixed">
    <tr>
        <th>Column 1</th>
        <th>Column 2</th>
    </tr>
</thead>

<script>
App.autoFixedThead();
</script>
```

---

## 浏览器兼容性 / Browser Compatibility

- **Chrome/Edge**: 完全支持 / Full support
- **Firefox**: 完全支持 / Full support
- **Safari**: 完全支持 / Full support
- **IE11+**: 支持 (需要 polyfills) / Supported (requires polyfills)

---

## 依赖项 / Dependencies

- **jQuery**: 核心 DOM 操作 / Core DOM manipulation
- **Bootstrap**: UI 组件 / UI components
- **jQuery NanoScroller**: 自定义滚动条 / Custom scrollbar
- **Bootstrap Switch**: 开关组件 / Switch component
- **Select2**: 增强选择框 / Enhanced select
- **Bootstrap Datepicker**: 日期选择器 / Date picker
- **Typeahead**: 自动完成 / Autocomplete
- **Gritter**: 消息通知 / Message notifications
- **NProgress**: 进度条 / Progress bar
