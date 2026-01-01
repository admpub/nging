# Modal.js 文档 / Modal.js Documentation

## 概述 / Overview

`modal.js` 是一个用于管理模态框（Modal）功能的 JavaScript 模块，提供了模态框表单初始化、多语言支持、分页、表单字段操作等功能。

---

## 目录 / Table of Contents

1. [核心函数 / Core Functions](#核心函数--core-functions)
2. [表单字段操作 / Form Field Operations](#表单字段操作--form-field-operations)
3. [多语言支持 / Multilingual Support](#多语言支持--multilingual-support)
4. [全局 API / Global API](#全局-api--global-api)

---

## 核心函数 / Core Functions

### initModalBody / 初始化模态框主体

初始化模态框的主体内容，包括标签页切换和分页功能。

**语法 / Syntax:**
```javascript
App.initModalBody($md, ajaxData, onloadCallback, showFooterTabIndexies)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `$md` | jQuery Object | 模态框的 jQuery 对象 / Modal jQuery object |
| `ajaxData` | Object | AJAX 请求数据 / AJAX request data |
| `onloadCallback` | Function | 加载回调函数 / Load callback function |
| `showFooterTabIndexies` | Array | 显示底部的标签索引数组 / Array of tab indexes to show footer |

**功能说明 / Description:**

- 设置标签页点击事件，控制底部按钮的显示/隐藏
- 初始化模态框内的分页表格
- 默认显示第一个标签页的底部按钮

**示例 / Example:**
```javascript
var $modal = $('#myModal');
var ajaxData = { type: 'list' };
var onloadCallback = {
    'ajaxList': function($table) {
        console.log('Table loaded');
    },
    'switchPage': function(page) {
        console.log('Page switched to:', page);
    }
};

App.initModalBody($modal, ajaxData, onloadCallback, [0, 2]);
```

---

### initModalBodyPagination / 初始化模态框分页

初始化模态框内的分页功能。

**语法 / Syntax:**
```javascript
App.initModalBodyPagination(that, ajaxData, onloadCallback)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `that` | jQuery Object | 表格的 jQuery 对象 / Table jQuery object |
| `ajaxData` | Object | AJAX 请求数据 / AJAX request data |
| `onloadCallback` | Function | 加载回调函数 / Load callback function |

**功能说明 / Description:**

- 初始化表格的分页功能
- 支持自定义页面切换回调
- 支持自定义分页数据

**示例 / Example:**
```javascript
var $table = $('table[data-page-size]');
var ajaxData = { category: 'products' };
var onloadCallback = {
    'switchPage': function(page) {
        console.log('Switching to page:', page);
    }
};

App.initModalBodyPagination($table, ajaxData, onloadCallback);
```

---

### initModalForm / 初始化模态框表单

初始化模态框的表单，支持多语言和自定义回调。

**语法 / Syntax:**
```javascript
App.initModalForm(button, modal, fields, afterOpenCallback, onSubmitCallback, multilingualFieldPrefix)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `button` | jQuery Object | 触发按钮 / Trigger button |
| `modal` | jQuery Object | 模态框对象 / Modal object |
| `fields` | Array/Object | 表单字段数组或配置对象 / Form fields array or config object |
| `afterOpenCallback` | Function | 打开后回调 / Callback after opening |
| `onSubmitCallback` | Function | 提交回调 / Submit callback |
| `multilingualFieldPrefix` | String | 多语言字段前缀 / Multilingual field prefix |

**参数对象模式 / Object Parameter Mode:**

```javascript
App.initModalForm(button, modal, {
    fields: ['field1', 'field2'],
    afterOpenCallback: function(button, modal, context) { /* ... */ },
    onSubmitCallback: function(button, modal, values) { /* ... */ },
    multilingualFieldPrefix: 'Language'
});
```

**功能说明 / Description:**

- 设置模态框标题
- 重置表单
- 处理多语言标签页
- 设置表单字段值
- 处理表单提交

**示例 / Example:**

```javascript
// 基本用法
$('#editButton').on('click', function() {
    var $this = $(this);
    App.initModalForm(
        $this,
        $('#editModal'),
        ['name', 'email', 'phone'],
        function(button, modal, context) {
            // 打开后的回调
            context.formData = {
                name: 'John Doe',
                email: 'john@example.com'
            };
        },
        function(button, modal, values) {
            // 提交回调
            console.log('Form values:', values);
            $.post('/api/update', values, function(res) {
                if(res.code === 0) {
                    App.message({title: 'Success', text: 'Updated successfully'});
                    modal.modal('hide');
                }
            });
        }
    );
});

// 多语言表单
$('#multiLangButton').on('click', function() {
    App.initModalForm(
        $(this),
        $('#multiLangModal'),
        {
            fields: ['title', 'description'],
            multilingualFieldPrefix: 'Language',
            afterOpenCallback: function(button, modal, context) {
                // 自定义字段值获取
                return function(fieldName) {
                    return context.formData[fieldName] || '';
                };
            },
            onSubmitCallback: function(button, modal, values) {
                console.log('Multilingual values:', values);
                // values.langDefault: 默认语言
                // values.data: 所有语言的数据
                // values[fieldName]: 默认语言的字段值
            }
        }
    );
});
```

---

### updateMultilingualFormByModal / 通过模态框更新多语言表单

根据模态框的值更新父表单的多语言字段。

**语法 / Syntax:**
```javascript
App.updateMultilingualFormByModal(parent, values, prefixNames, nameFixer, parentForDefaultLang, callback)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `parent` | jQuery Object | 父表单容器 / Parent form container |
| `values` | Object | 模态框返回的值 / Values from modal |
| `prefixNames` | Array | 前缀名称数组 / Array of prefix names |
| `nameFixer` | Function | 字段名称修正函数 / Field name fixer function |
| `parentForDefaultLang` | jQuery Object | 默认语言表单容器 / Default language form container |
| `callback` | Function | 字段更新回调 / Field update callback |

**功能说明 / Description:**

- 为非默认语言创建隐藏字段
- 更新默认语言的字段值
- 处理强制翻译标记
- 支持自定义字段名处理

**示例 / Example:**

```javascript
var modalValues = {
    data: {
        'Language[zh-CN][title]': '中文标题',
        'Language[zh-CN][content]': '中文内容',
        'Language[en-US][title]': 'English Title',
        'Language[en-US][content]': 'English Content'
    },
    langDefault: 'en-US',
    multilingual: true,
    title: 'English Title',
    content: 'English Content'
};

App.updateMultilingualFormByModal(
    $('#mainForm'),
    modalValues,
    ['Translation'],
    function(fieldName) {
        // 修正字段名
        return fieldName.replace(/_/, '-');
    },
    $('#defaultLangForm'),
    function(field, name, value) {
        console.log('Updated:', name, '=', value);
    }
);
```

---

## 表单字段操作 / Form Field Operations

### setFormFieldValue / 设置表单字段值

设置表单字段的值，支持多种输入类型。

**语法 / Syntax:**
```javascript
App.setFormFieldValue(input, val)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `input` | jQuery Object | 输入元素的 jQuery 对象 / Input jQuery object |
| `val` | Any | 要设置的值 / Value to set |

**支持的元素类型 / Supported Element Types:**

- **Checkbox**: 多选框，根据值设置选中状态
- **Radio**: 单选框，根据值设置选中状态
- **Select**: 下拉框，根据值设置选中项
- **Text/Hidden/其他**: 文本框和隐藏字段，直接设置值

**示例 / Example:**

```javascript
// 设置文本框
App.setFormFieldValue($('#name'), 'John Doe');

// 设置复选框
App.setFormFieldValue($('input[name="hobbies[]"]'), ['reading', 'gaming']);

// 设置单选框
App.setFormFieldValue($('input[name="gender"]'), 'male');

// 设置下拉框
App.setFormFieldValue($('#country'), 'US');
```

---

### getFormFieldValue / 获取表单字段值

获取表单字段的值，支持多种输入类型。

**语法 / Syntax:**
```javascript
App.getFormFieldValue(input)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `input` | jQuery Object | 输入元素的 jQuery 对象 / Input jQuery object |

**返回值 / Returns:**

- **Checkbox (带 [])**: 返回选中的值数组
- **Checkbox (不带 [])**: 返回最后一个选中的值
- **Radio**: 返回选中的值
- **其他**: 返回字段的值

**示例 / Example:**

```javascript
// 获取文本框值
var name = App.getFormFieldValue($('#name'));

// 获取多选框值数组
var hobbies = App.getFormFieldValue($('input[name="hobbies[]"]'));
// 结果: ['reading', 'gaming']

// 获取单选框值
var gender = App.getFormFieldValue($('input[name="gender"]'));

// 获取下拉框值
var country = App.getFormFieldValue($('#country'));
```

---

### parseLangFieldName / 解析语言字段名

从多语言字段名中解析出字段名。

**语法 / Syntax:**
```javascript
App.parseLangFieldName(name, prefix)
```

**参数 / Parameters:**

| 参数名 / Parameter | 类型 / Type | 描述 / Description |
|-------------------|-------------|-------------------|
| `name` | String | 完整的字段名 / Full field name |
| `prefix` | String | 字段前缀 / Field prefix |

**示例 / Example:**

```javascript
var fieldName = App.parseLangFieldName('Language[zh-CN][title]', 'Language[zh-CN]');
// 返回: 'title'

var fieldName = App.parseLangFieldName('Language[en][description]', 'Language[en]');
// 返回: 'description'
```

---

## 多语言支持 / Multilingual Support

### 多语言表单处理流程 / Multilingual Form Processing Flow

1. **检测多语言表单** / Detect multilingual form
   ```javascript
   var multilingual = modal.find('.langset').length > 0;
   ```

2. **获取默认语言** / Get default language
   ```javascript
   var langDefault = firstTabLi.children('a').data('lang');
   ```

3. **构建字段列表** / Build field list
   - 默认语言字段: `Language[en][fieldName]`
   - 其他语言字段: `Language[zh-CN][fieldName]`

4. **字段值获取** / Field value retrieval
   - 默认语言: 直接使用字段名
   - 其他语言: 使用带语言前缀的完整字段名

### 多语言字段名格式 / Multilingual Field Name Format

```
Language[languageCode][fieldName]
```

**示例 / Examples:**

```javascript
Language[en][title]          // 英文标题
Language[zh-CN][title]       // 简体中文标题
Language[ja][description]     // 日文描述
```

---

## 全局 API / Global API

### 暴露的方法 / Exposed Methods

所有主要函数都通过 `App` 对象暴露给全局作用域：

```javascript
// 初始化模态框主体
App.initModalBody($md, ajaxData, onloadCallback, showFooterTabIndexies);

// 初始化模态框分页
App.initModalBodyPagination(that, ajaxData, onloadCallback);

// 初始化模态框表单
App.initModalForm(button, modal, fields, afterOpenCallback, onSubmitCallback, multilingualFieldPrefix);

// 更新多语言表单
App.updateMultilingualFormByModal(parent, values, prefixNames, nameFixer, parentForDefaultLang, callback);

// 获取表单字段值
App.getFormFieldValue(input);

// 设置表单字段值
App.setFormFieldValue(input, val);
```

---

## 完整示例 / Complete Examples

### 示例 1: 基本编辑模态框 / Example 1: Basic Edit Modal

```javascript
$('#editButton').on('click', function() {
    var $btn = $(this);
    var $modal = $('#editModal');

    // 设置模态框标题
    $btn.data('modal-title', 'Edit User');

    // 设置表单数据
    $btn.data('form-data', {
        name: 'John Doe',
        email: 'john@example.com',
        status: 1
    });

    // 初始化模态框表单
    App.initModalForm(
        $btn,
        $modal,
        ['name', 'email', 'status'],
        function(button, modal, context) {
            console.log('Modal opened');
        },
        function(button, modal, values) {
            // 提交表单
            $.post('/user/update', values, function(res) {
                if(res.code === 0) {
                    App.message({title: 'Success', text: 'User updated'});
                    modal.modal('hide');
                    // 刷新列表
                    window.location.reload();
                }
            });
        }
    );

    // 显示模态框
    $modal.modal('show');
});
```

### 示例 2: 多语言表单 / Example 2: Multilingual Form

```javascript
$('#addTranslate').on('click', function() {
    var $btn = $(this);
    var $modal = $('#translateModal');

    App.initModalForm(
        $btn,
        $modal,
        {
            fields: ['title', 'content'],
            multilingualFieldPrefix: 'Language',
            afterOpenCallback: function(button, modal, context) {
                // 根据按钮数据填充表单
                var rowData = button.data('row-data');
                if(rowData) {
                    return function(fieldName) {
                        return rowData[fieldName] || '';
                    };
                }
            },
            onSubmitCallback: function(button, modal, values) {
                // 更新主表单
                App.updateMultilingualFormByModal(
                    $('#mainForm'),
                    values,
                    ['Translation'],
                    null,
                    $('#defaultLangFields'),
                    function(field, name, value) {
                        console.log('Field updated:', name, value);
                    }
                );

                modal.modal('hide');
            }
        }
    );

    $modal.modal('show');
});
```

### 示例 3: 带分页的模态框 / Example 3: Modal with Pagination

```javascript
$('#selectItem').on('click', function() {
    var $modal = $('#selectModal');

    // 加载内容
    $modal.find('.modal-body').load('/items/list', function() {
        // 初始化模态框主体
        App.initModalBody(
            $modal,
            { type: 'select' },
            {
                'ajaxList': function($table) {
                    console.log('Table loaded');
                },
                'switchPage': function(page) {
                    console.log('Page:', page);
                }
            },
            [0] // 只在第一个标签显示底部按钮
        );
    });

    $modal.modal('show');
});
```

### 示例 4: 动态字段处理 / Example 4: Dynamic Field Handling

```javascript
$('#dynamicForm').on('click', function() {
    var $btn = $(this);
    var $modal = $('#dynamicModal');

    // 动态获取字段列表
    var fields = [];

    App.initModalForm(
        $btn,
        $modal,
        {
            afterOpenCallback: function(button, modal, context) {
                // 从服务器获取字段配置
                return function(fieldName) {
                    // 返回字段值
                    return button.data(fieldName);
                };
            },
            onSubmitCallback: function(button, modal, values) {
                console.log('Submit values:', values);

                // 处理特殊字段
                if(values.multilingual) {
                    // 多语言处理
                    var translations = {};
                    for(var lang in values.data) {
                        if(lang !== values.langDefault) {
                            translations[lang] = values.data[lang];
                        }
                    }
                    console.log('Translations:', translations);
                }
            }
        }
    );

    $modal.modal('show');
});
```

---

## 最佳实践 / Best Practices

### 1. 数据属性使用 / Data Attributes Usage

使用 HTML5 data 属性存储模态框相关数据：

```html
<button id="editBtn"
        data-modal-title="Edit User"
        data-form-id="123"
        data-name="John Doe">
    Edit
</button>
```

```javascript
$('#editBtn').on('click', function() {
    var $btn = $(this);
    var formData = {
        id: $btn.data('form-id'),
        name: $btn.data('name')
    };
    // ...
});
```

### 2. 回调函数命名 / Callback Function Naming

使用有意义的回调函数名：

```javascript
var callbacks = {
    'userList': function($table) {
        // 用户列表加载完成
    },
    'productList': function($table) {
        // 产品列表加载完成
    },
    'switchPage': function(page) {
        // 页面切换
    },
    'pageData': function($table) {
        return $table.data();
    }
};

App.initModalBody($modal, {}, callbacks, [0]);
```

### 3. 错误处理 / Error Handling

在回调函数中添加错误处理：

```javascript
onSubmitCallback: function(button, modal, values) {
    $.post('/api/save', values)
        .done(function(res) {
            if(res.code === 0) {
                App.message({title: 'Success', text: 'Saved'});
                modal.modal('hide');
            } else {
                App.message({title: 'Error', text: res.error || 'Save failed'});
            }
        })
        .fail(function() {
            App.message({title: 'Error', text: 'Network error'});
        });
}
```

### 4. 表单验证 / Form Validation

在提交前进行验证：

```javascript
onSubmitCallback: function(button, modal, values) {
    // 验证必填字段
    if(!values.name || !values.name.trim()) {
        App.message({title: 'Warning', text: 'Name is required'});
        return;
    }

    // 验证邮箱格式
    if(values.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(values.email)) {
        App.message({title: 'Warning', text: 'Invalid email format'});
        return;
    }

    // 提交表单
    // ...
}
```

### 5. 多语言字段管理 / Multilingual Field Management

```javascript
// 统一管理多语言字段前缀
var MULTILINGUAL_PREFIX = 'Language';

// 获取默认语言
var getDefaultLang = function() {
    return $('body').data('default-lang') || 'en';
};

// 构建多语言字段名
var buildLangFieldName = function(lang, fieldName) {
    return MULTILINGUAL_PREFIX + '[' + lang + '][' + fieldName + ']';
};

// 使用
var titleField = buildLangFieldName('zh-CN', 'title');
// 结果: 'Language[zh-CN][title]'
```

---

## 常见问题 / FAQ

### Q1: 如何在模态框中处理复杂数据?

**A:** 使用 `afterOpenCallback` 和自定义数据处理逻辑：

```javascript
afterOpenCallback: function(button, modal, context) {
    var data = button.data('complex-data');
    
    return function(fieldName) {
        if(fieldName === 'tags') {
            // 处理数组数据
            return data.tags ? data.tags.join(',') : '';
        }
        return data[fieldName];
    };
}
```

### Q2: 如何动态控制底部按钮显示?

**A:** 使用 `showFooterTabIndexies` 参数：

```javascript
// 只在索引为 0 和 2 的标签页显示底部按钮
App.initModalBody($modal, {}, null, [0, 2]);
```

### Q3: 多语言表单如何处理默认语言?

**A:** 系统会自动使用第一个标签页的语言作为默认语言：

```javascript
// 获取默认语言
var langDefault = firstTabLi.children('a').data('lang');

// 字段值存储在两个地方:
// 1. values[fieldName] - 默认语言值
// 2. values.data[Language[lang][fieldName]] - 所有语言值
```

### Q4: 如何自定义分页行为?

**A:** 在 `onloadCallback` 中提供 `switchPage` 函数：

```javascript
var callbacks = {
    'switchPage': function(page) {
        // 自定义分页逻辑
        var params = { page: page, size: 20 };
        $.get('/custom/list', params, function(res) {
            // 更新表格内容
        });
    }
};

App.initModalBodyPagination($table, {}, callbacks);
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
- **Bootstrap Modal**: 模态框 UI 组件 / Modal UI component
- **App**: 全局应用对象 (包含 message 等方法) / Global app object (includes message method, etc.)

---

## 更新日志 / Changelog

### Version 1.0.0
- 初始版本 / Initial version
- 支持基本模态框功能 / Basic modal support
- 多语言表单支持 / Multilingual form support
- 分页功能 / Pagination support
