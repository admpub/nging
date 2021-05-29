forms
==========

[![GoDoc](https://godoc.org/github.com/coscms/forms?status.png)](http://godoc.org/github.com/coscms/forms)

Description
===========

`forms` makes form creation and handling easy. It allows the creation of form without having to write HTML code or bother to make the code Bootstrap compatible.
You can just create your form instance and add / populate / customize fields based on your needs. Or you can let `forms` do that for you starting from any object instance.

To integrate `forms` forms into your application simply pass the form object to the template and call its Render method.
In your code:

```go
tmpl.Execute(buf, map[string]interface{}{"form": form})
```

In your template:

```html
{{ if .form }}{{ .form.Render }}{{ end }}
```

Installation
============

To install this package simply:

```go
go get github.com/coscms/forms
```

Forms
=====

There are two predefined styles for forms: base HTML forms and Bootstrap forms: they have different structures and predefined classes.
Style aside, forms can be created from scratch or starting from a base instance.

From scratch
------------

You can create a form instance by simply deciding its style and providing its method and action:

```go
form := NewWithConfig(&config.Config{
    Theme:"base", 
    Method:POST, 
    Action:"/action.html",
})
```

Now that you have a form instance you can customize it by adding classes, parameters, CSS values or id. Each method returns a pointer to the same form, so multiple calls can be chained:

```go
form.SetID("TestForm").AddClass("form").AddCSS("border", "auto")
```

Obviously, elements can be added as well:

```go
form.Elements(fields.TextField("text_field"))
```

Elements can be either FieldSets or Fields: the formers are simply collections of fields translated into a `<fieldset></fieldset>` element.
Elements are added in order, and they are displayed in the exact same order. Note that single elements can be removed from a form referencing them by name:

```go
form.RemoveElement("text_field")
```

Typical usage looks like this:

```go
form := NewWithConfig(&config.Config{
    Theme:"base", 
    Method:POST, 
    Action:"/action.html",
}).Elements(
    fields.TextField("text_field").SetLabel("Username"),
    FieldSet("psw_fieldset",
        fields.PasswordField("psw1").AddClass("password_class").SetLabel("Password 1"),
        fields.PasswordField("psw2").AddClass("password_class").SetLabel("Password 2"),
        ),
    fields.SubmitButton("btn1", "Submit"),
    )
```

validation:

```go
type User struct {
    Username     string    `valid:"Required;AlphaDash;MaxSize(30)"`
    Password1     string    `valid:"Required"`
    Password2    string    `valid:"Required"`
}
u := &User{}
form.SetModel(u) //Must set model
err := form.valid()
if err != nil { 
    // validation does not pass
}
```

Details about validation, please visit: https://github.com/webx-top/validation

A call to `form.Render()` returns the following form:

```html
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
```

From model instance
-------------------

Instead of manually creating a form, it can be automatically created from an existing model instance: the package will try to infer the field types based on the instance fields and fill them accordingly.
Default type-to-field mapping is as follows:

* string: TextField
* bool: Checkbox
* time.Time: DatetimeField
* int: NumberField
* struct: recursively parse

You can customize field behaviors by adding tags to instance fields.
Without tags this code:

```go
type User struct {
    Username     string
    Password1     string
    Password2    string
}
u := &User{}
form := NewFromModel(u, &config.Config{
    Theme:"bootstrap3", 
    Method:POST, 
    Action:"/action.html",
})
```

validation:

```go
err := form.valid()
if err != nil { 
    // validation does not pass
}
form.Render()
```

would yield this HTML form:

```html
<form method="POST" action="/action.html">
    <label>Username</label>
    <input type="text" name="Username">
    <label>Password1</label>
    <input type="text" name="Password1">
    <label>Password2</label>
    <input type="text" name="Password2">
    <button type="submit" name="submit">Submit</button>
</form>
```

A submit button is added by default.

Notice that the form is still editable and fields can be added, modified or removed like before.

When creating a form from a model instance, field names are created by appending the field name to the baseline; the baseline is empty for single level structs but is crafted when nested structs are found: in this case it becomes the field name followed by a dot.
So for example, if the struct is:

```go
type A struct {
    field1     int
    field2     int
}
type B struct {
    field0     int
    struct1    A
}
```

The final form will contain fields "field0", "struct1.field1" and "struct1.field2".

Tags
----

Struct tags can be used to slightly modify automatic form creation. In particular the following tags are parsed:

* form_options: can contain the following keywords separated by Semicolon (;)

    - -: skip field, do not convert to HTML field
    - checked: for Checkbox fields, check by default
    - multiple: for select fields, allows multiple choices

* form_widget: override custom widget with one of the following

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
    - static (simple text)

* form_fieldset: define fieldset name
* form_sort: sort number (asc, 0 ~ total-1)
* form_choices: defines options for select and radio input fields
    - radio/checkbox example(format: id|value): 1|Option One|2|Option 2|3|Option 3
    - select example(format: group|id|value): G1|A|Option A|G1|B|Option B
        - "" group is the default one and does not trigger a `<optgroup></optgroup>` rendering.
* form_max: max value (number, range, datetime, date and time fields)
* form_min: min value (number, range, datetime, date and time fields)
* form_step: step value (range field)
* form_rows: number of rows (textarea field)
* form_cols: number of columns (textarea field)
* form_value: input field value (used if field is empty)
* form_label: label for input field

The code would therefore be better like this:

```go
type User struct {
    Username     string
    Password1     string     `form_widget:"password" form_label:"Password 1"`
    Password2    string     `form_widget:"password" form_label:"Password 2"`
    SkipThis    int     `form_options:"-"`
}
u := User{}
form := NewFromModel(u, &config.Config{
    Theme:"bootstrap3", 
    Method:POST, 
    Action:"/action.html",
})
form.Render()
```

which translates into:

```html
<form method="POST" action="/action.html">
    <label>Username</label>
    <input type="text" name="Username">
    <label>Password 1</label>
    <input type="password" name="Password1">
    <label>Password 2</label>
    <input type="password" name="Password2">
    <button type="submit" name="submit">Submit</button>
</form>
```

Fields
======

Field objects in `forms` implement the `fields.FieldInterface` which exposes methods to edit classes, parameters, tags and CSS styles.
See the [documentation](http://godoc.org/github.com/coscms/forms) for details.

Most of the field widgets have already been created and integrate with Bootstrap. It is possible, however, to define custom widgets to render fields by simply assigning an object implementing the widgets.WidgetInterface to the Widget field.

Also, error messages can be added to fields via the `AddError(err)` method: in a Bootstrap environment they will be correctly rendered.

Text fields
-----------

This category includes text, password, textarea and hidden fields. They are all instantiated by providing the name, except the TextAreaField which also requires a dimension in terms of rows and columns.

```go
f0 := fields.TextField("text")
f1 := fields.PasswordField("password")
f2 := fields.HiddenField("hidden")
f3 := fields.TextAreaField("textarea", 30, 50)
```

Option fields
-------------

This category includes checkbox, select and radio button fields.
Checkbox field requires a name and a set of options to populate the field. The options are just a set of InputChoice (ID-Value pairs) objects:

```go
opts := []fields.InputChoice{
    fields.InputChoice{ID:"A", Val:"Option A"},
    fields.InputChoice{ID:"B", Val:"Option B"},
}
f := fields.CheckboxField("checkbox", opts)
f.AddSelected("A", "B")
```

Radio buttons, instead, require a name and a set of options to populate the field. The options are just a set of InputChoice (ID-Value pairs) objects:

```go
opts := []fields.InputChoice{
    fields.InputChoice{ID:"A", Val:"Option A"},
    fields.InputChoice{ID:"B", Val:"Option B"},
}
f := fields.RadioField("radio", opts)
```

Select fields, on the other hand, allow option grouping. This can be achieved by passing a `map[string][]InputChoice` in which keys are groups containing choices given as values; the default (empty) group is "", which is not translated into any `<optgroup></optgroup>` element.

```go
opts := map[string][]fields.InputChoice{
    "": []fields.InputChoice{fields.InputChoice{"A", "Option A"}},
    "group1": []fields.InputChoice{
        fields.InputChoice{ID:"B", Val:"Option B"},
        fields.InputChoice{ID:"C", Val:"Option C"},
    }
}
f := fields.SelectField("select", opts)
```

Select fields can allow multiple choices. To enable this option simply call the `MultipleChoice()` method on the field and provide the selected choices via `AddSelected(...string)`:

```go
f.MultipleChoice()
f.AddSelected("A", "B")
```

Number fields
-------------

Number and range fields are included.
Number field only require a name to be instantiated; minimum and maximum values can optionally be set by adding `min` and `max` parameters respectively.

```go
f := fields.NumberField("number")
f.SetParam("min", "1")
```

Range fields, on the other hand, require both minimum and maximum values (plus the identifier). The optional "step" value is set via `SetParam`.

```go
f := fields.RangeField("range", 1, 10)
f.SetParam("step", "2")
```

Datetime fields
---------------

Datetime, date and time input fields are defined in `forms`.

```go
f0 := fields.DatetimeField("datetime")
f1 := fields.DateField("date")
f2 := fields.TimeField("time")
```

Values can be set via `SetValue` method; there's no input validation but format strings are provided to ensure the correct time-to-string conversion.

```go
t := time.Now()
f0.SetValue(t.Format(fields.DATETIME_FORMAT))
f1.SetValue(t.Format(fields.DATE_FORMAT))
f2.SetValue(t.Format(fields.TIME_FORMAT))
```

Buttons
-------

Buttons can be created calling either the `Button`, `SubmitButton` or `ResetButton` constructor methods and providing a text identifier and the content of the button itself.

```go
btn0 := fields.Button("btn", "Click me!")
```

License
=======

`forms` is released under the MIT license. See [LICENSE](https://github.com/coscms/forms/blob/master/LICENSE).