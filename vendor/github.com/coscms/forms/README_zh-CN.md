forms
==========

[![GoDoc](https://godoc.org/github.com/coscms/forms?status.png)](http://godoc.org/github.com/coscms/forms)


简介
===========

用`forms`可以方便的创建HTML表单，支持指定表单模板，在使用时您不用写HTML代码就可以便捷的创建出符合个人需求的表单。你可以根据自身需要创建一个表单实例然后添加、填充、自定义表单字段。

`forms`可以很容易的集成到你的应用中，只需要通过调用form对象的Render方法即可渲染出表单的HTML代码:
	
	tmpl.Execute(buf, map[string]interface{}{"form": form})

在模板中使用:

	{{ if .form }}{{ .form.Render }}{{ end }}


安装方式
============

使用以下代码安装forms包:

	go get github.com/coscms/forms

表单模板
=====

本包中已经内置了两套表单模板：base和bootstrap3。它们代表两种不同的风格。除此之外，您也可以创建自己的表单模板。

入门指引
------------

创建一个form实例，并指定表单模板为base，表单的提交方式为POST，表单的提交网址为“/action.html”：
	
	form := NewForm("base", POST, "/action.html")


现在，我们可以通过form实例来自定义表单的属性。form实例中的每一个方法都会返回form指针，因此你可以非常方便的采用链式调用来多次执行方法。下面，我们来定义HTML标签`<form>`中的id、class和style属性：

	form.SetID("TestForm").AddClass("form").AddCSS("border", "auto")

添加其它表单字段也很容易，比如，我们要添加一个name属性值为“text_field”的文本输入框：

	form.Elements(fields.TextField("text_field"))

Elements方法可以添加`<fieldset></fieldset>`包围起来的表单字段或单独的表单字段，然后根据在Elements方法中的添加顺序来依次显示表单元素。 我们也可以通过name属性值来删除某个表单元素:

	form.RemoveElement("text_field")

典型用法:

	form := NewForm("base", POST, "/action.html").Elements(
		fields.TextField("text_field").SetLabel("Username"),
		FieldSet("psw_fieldset",
			fields.PasswordField("psw1").AddClass("password_class").SetLabel("Password 1"),
			fields.PasswordField("psw2").AddClass("password_class").SetLabel("Password 2"),
			),
		fields.SubmitButton("btn1", "Submit"),
		)

表单验证：

	
	type User struct {
		Username 	string	`valid:"Required;AlphaDash;MaxSize(30)"`
		Password1 	string	`valid:"Required"`
		Password2	string	`valid:"Required"`
	}

	u := &User{}
	form.SetModel(u) //必须设置数据模型
	valid, passed := form.valid()
	if !passed { 
		//验证未通过时的操作代码
	}
	_ = valid


表单验证的详细用法请访问: [https://github.com/webx-top/validation](https://github.com/webx-top/validation)

调用 `form.Render()` 返回如下表单：
	
	<form method="POST" action="/action.html">
		<label>Username</label>
		<input type="text" name="text_field">
		<fieldset>
			<label>Password 1</label>
			<input type="password" name="psw1" class="password_class ">
			<label>Password 2</label>
			<input type="password" name="psw2" class="password_class ">
		</fieldset>
		<button type="submit" name="btn1">Submit</button>
	</form>

从model实例创建表单
-------------------

我们可以通过model实例来自动创建表单，免除了手动一个个添加表单字段的麻烦。本forms包会根据model实例中的属性字段及其类型来依次填充到表单中作为表单字段。
model实例内属性字段类型与表单字段对应关系如下：

* string: TextField （文本输入框）
* bool: Checkbox （复选框）
* time.Time: DatetimeField （日期输入框）
* int: NumberField （数字输入框）
* struct: 递归解析

也可以通过添加tag到model实例中的属性字段内来定义表单字段的类型。  
不带tag的代码：
	
	type User struct {
		Username 	string
		Password1 	string
		Password2	string
	}

	u := &User{}

	form := NewFormFromModel(u, "bootstrap3", POST, "/action.html")

验证表单数据：

	valid, passed := form.valid()
	if !passed { 
		// validation does not pass
	}
	_ = valid

	form.Render()

生成的HTML代码如下:

	<form method="POST" action="/action.html">
		<label>Username</label>
		<input type="text" name="Username">
		<label>Password1</label>
		<input type="text" name="Password1">
		<label>Password2</label>
		<input type="text" name="Password2">
		<button type="submit" name="submit">Submit</button>
	</form>

默认会添加一个提交按钮。

注意：我们可以象前面介绍的那样添加、修改或删除某一个表单字段。

When creating a form from a model instance, field names are created by appending the field name to the baseline; the baseline is empty for single level structs but is crafted when nested structs are found: in this case it becomes the field name followed by a dot.
So for example, if the struct is:

	type A struct {
		field1 	int
		field2 	int
	}

	type B struct {
		field0 	int
		struct1	A
	}

The final form will contain fields "field0", "struct1.field1" and "struct1.field2".

Tags
----

Struct tags can be used to slightly modify automatic form creation. 下面列出的这些tag会被解析:

* form_options: 可以包含如下关键词，同时使用多个关键词时，用分号（;）隔开
	- -: 跳过此字段, 不转为HTML表单字段
	- checked: 针对Checkbox，默认选中
	- multiple: 指定select为允许多选
* form_widget: 指定表单部件类型。支持以下类型：
	- text
	- hidden
	- textarea
	- password
	- select
	- datetime
	- date
	- time
	- number
	- range
	- radio
	- checkbox
	- static (简单的静态文本)
* form_fieldset: 定义fieldset标题文字
* form_sort: 排序编号 (按升序排列, 编号从0开始，范围为0 ~ 总数-1)
* form_choices: select或radio输入字段的选项
	- radio/checkbox 范例(格式: id|value): 1|选项一|2|选项二|3|选项三
	- select 范例(格式: group|id|value): 组1|A|选项A|组1|B|选项B 
		- "" 组名为空白时，默认将不渲染`<optgroup></optgroup>`。
* form_max: 允许的最大值 (用于number、range、datetime、date 和 time 类型输入框)
* form_min: 允许的最小值 (用于number、range、datetime、date 和 time 类型输入框)
* form_step: 步进值 (用于range输入字段)
* form_rows: 行数 (用于textarea)
* form_cols: 列数 (用于textarea)
* form_value: 输入字段的默认值
* form_label: label内容

例如：

	type User struct {
		Username 	string
		Password1 	string 	`form_widget:"password" form_label:"Password 1"`
		Password2	string 	`form_widget:"password" form_label:"Password 2"`
		SkipThis	int 	`form_options:"-"`
	}

	u := User{}

	form := NewFormFromModel(u, "bootstrap3", POST, "/action.html")
	form.Render()

它们最后会翻译成以下代码：

 	<form method="POST" action="/action.html">
		<label>Username</label>
		<input type="text" name="Username">
		<label>Password 1</label>
		<input type="password" name="Password1">
		<label>Password 2</label>
		<input type="password" name="Password2">
		<button type="submit" name="submit">Submit</button>
	</form>

Fields
======

Field objects in `forms` implement the `fields.FieldInterface` which exposes methods to edit classes, parameters, tags and CSS styles.
See the [documentation](http://godoc.org/github.com/coscms/forms) for details.

Most of the field widgets have already been created and integrate with Bootstrap. It is possible, however, to define custom widgets to render fields by simply assigning an object implementing the widgets.WidgetInterface to the Widget field.

Also, error messages can be added to fields via the `AddError(err)` method: in a Bootstrap environment they will be correctly rendered.

Text fields
-----------

This category includes text, password, textarea and hidden fields. They are all instantiated by providing the name, except the TextAreaField which also requires a dimension in terms of rows and columns.

	f0 := fields.TextField("text")
	f1 := fields.PasswordField("password")
	f2 := fields.HiddenField("hidden")
	f3 := fields.TextAreaField("textarea", 30, 50)

Option fields
-------------

This category includes checkbox, select and radio button fields.
Checkbox field requires a name and a set of options to populate the field. The options are just a set of InputChoice (ID-Value pairs) objects:

	opts := []fields.InputChoice{
		fields.InputChoice{ID:"A", Val:"Option A"},
		fields.InputChoice{ID:"B", Val:"Option B"},
	}
	f := fields.CheckboxField("checkbox", opts)
	f.AddSelected("A", "B")

Radio buttons, instead, require a name and a set of options to populate the field. The options are just a set of InputChoice (ID-Value pairs) objects:

	opts := []fields.InputChoice{
		fields.InputChoice{ID:"A", Val:"Option A"},
		fields.InputChoice{ID:"B", Val:"Option B"},
	}
	f := fields.RadioField("radio", opts)

Select fields, on the other hand, allow option grouping. This can be achieved by passing a `map[string][]InputChoice` in which keys are groups containing choices given as values; the default (empty) group is "", which is not translated into any `<optgroup></optgroup>` element.

	opts := map[string][]fields.InputChoice{
		"": []fields.InputChoice{fields.InputChoice{"A", "Option A"}},
		"group1": []fields.InputChoice{
			fields.InputChoice{ID:"B", Val:"Option B"},
			fields.InputChoice{ID:"C", Val:"Option C"},
		}
	}
	f := fields.SelectField("select", opts)

Select fields can allow multiple choices. To enable this option simply call the `MultipleChoice()` method on the field and provide the selected choices via `AddSelected(...string)`:

	f.MultipleChoice()
	f.AddSelected("A", "B")

Number fields
-------------

Number and range fields are included.
Number field only require a name to be instantiated; minimum and maximum values can optionally be set by adding `min` and `max` parameters respectively.

	f := fields.NumberField("number")
	f.SetParam("min", "1")

Range fields, on the other hand, require both minimum and maximum values (plus the identifier). The optional "step" value is set via `SetParam`.

	f := fields.RangeField("range", 1, 10)
	f.SetParam("step", "2")


Datetime fields
---------------

Datetime, date and time input fields are defined in `go-form-it`.

	f0 := fields.DatetimeField("datetime")
	f1 := fields.DateField("date")
	f2 := fields.TimeField("time")

Values can be set via `SetValue` method; there's no input validation but format strings are provided to ensure the correct time-to-string conversion.

	t := time.Now()
	f0.SetValue(t.Format(fields.DATETIME_FORMAT))
	f1.SetValue(t.Format(fields.DATE_FORMAT))
	f2.SetValue(t.Format(fields.TIME_FORMAT))

Buttons
-------

Buttons can be created calling either the `Button`, `SubmitButton` or `ResetButton` constructor methods and providing a text identifier and the content of the button itself.

	btn0 := fields.Button("btn", "Click me!")


License
=======

`forms` is released under the MIT license. See [LICENSE](https://github.com/coscms/forms/blob/master/LICENSE).