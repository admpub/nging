9t
==============================

[![Build Status](https://travis-ci.org/gongo/9t.svg?branch=master)](https://travis-ci.org/gongo/9t)
[![Coverage Status](https://coveralls.io/repos/gongo/9t/badge.svg?branch=master)](https://coveralls.io/r/gongo/9t?branch=master)

9t (nine-tailed fox in Japanese) is a multi-file tailer (like `tail -f a.log b.log ...`).

Usage
------------------------------

```
$ 9t file1 [file2 ...]
```

### Demo

![Demo](./images/9t.gif)

1. Preparation for demo

    ```sh
    $ yukari() { echo '世界一かわいいよ!!' }
    $ while :; do       yukari >> tamura-yukari.log ; sleep 0.2 ; done
    $ while :; do echo $RANDOM >> random.log        ; sleep 3   ; done
    $ while :; do         date >>      d.log        ; sleep 1   ; done
    ```

1. Run

    ```
    $ 9t tamura-yukari.log random.log d.log
    ```

Installation
------------------------------

```
$ go get github.com/admpub/9t/cmd/9t
```

Motivation
------------------------------

So far, Multiple file display can be even `tail -f`.

![Demo](./images/tailf.gif)

But, I wanted to see in a similar format as the `heroku logs --tail`.

```
app[web.1]: foo bar baz
app[worker.1]: pizza pizza
app[web.1]: foo bar baz
app[web.2]: just do eat..soso..
.
.
```

License
------------------------------

MIT License
