errorutil
=========

Errorutil is a small go package to help show syntax errors in for example JSON documents.

It was forked from [Camlistore](http://camlistore.org) to make a smaller dependency.


Example
-------

An example of how to use the package to show errors when decoding with
[encoding/json](http://golang.org/pkg/encoding/json/).

    if err = decoder.Decode(&objmap); err != nil {
            extra := ""

            // if it's a syntax error, add more information
            if serr, ok := err.(*json.SyntaxError); ok {
                    if _, serr := fh.Seek(0, os.SEEK_SET); serr != nil {
                            log.Fatalf("seek error: %v", serr)
                    }
                    line, col, highlight := errorutil.HighlightBytePosition(fh, serr.Offset)
                    extra = fmt.Sprintf(":\nError at line %d, column %d (file offset %d):\n%s",
                            line, col, serr.Offset, highlight)
            }

            return nil, fmt.Errorf("error parsing JSON object in config file %s%s\n%v",
                    fh.Name(), extra, err)
    }


License
-------

This package is licesed under the Apache License, version 2.0. It was developed
by Brad Fitzpatrick as part of the Camlistore project.
