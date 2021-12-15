i18n
====

i18n support for golang applications.  Supports message translation with
placeholders and plurals, locale-specific string sorting, and number/currency
formatting.

Copyright 2014 Vubeology, Inc.

Released under the MIT License (see LICENSE).

Usage
-----

Read the full documentation here: http://godoc.org/github.com/vube/i18n

License
-------

As stated about this golang package is released under the MIT License (see
LICENSE).

### Third Party Package Licenses

This i18n package makes use of third party packages in addition to the golang
standard library and supplemental libraries. This package however does not
modify or redistribute any third party package material.

While this i18n package is released under the MIT License, you must ensure that
your use of this package also complies with the licenses under which each third
party dependency is released.

#### launchpad.net/gocheck

This i18n package makes use of the launchpad.net/gocheck package, released under
a Simplified BSD License. For specific license details, refer directly to the
the gocheck package.

#### gopkg.in/yaml.v1

This i18n package makes use of the gopkg.in/yaml.v1 package, released under
the LGPLv3 License. For specific license details, refer directly to the yaml
package.

Call for open source help!
--------------------------

We could use some help.  We do however have some guidelines if you want to
contribute to our package.

### For supplementing locale data:

If you have locale rules data that we are missing, we welcome all additional
rules data in our standard yaml format.  Please include comments in the yaml
on how you sourced the data - AKA, you are a native speaker of a language, you
got it from XYZ website, a professional translator provided the data, etc.

When supplementing locale data, you may add a locale who's language uses a set
of plural rules that this package does not support.  In this case, please add
an additional plural rule function to the plurals.go file and add it to the
plural rule map.

There is a unit test in rules_test.go that checks loading every single locale
in this package. If you are adding a brand new locale to the list, please add it
to this unit test.

### For fixing bugs:

If you find a bug that you'd like to fix, please include a unit test that
validates your work.  This test should fail without the fix you provide and pass
with the fix you provide.

### For new features:

If you've decided to either tackle a feature on our wish list or you have a
feature that you need in order to use our package, please provide a minimum of
80% unit test coverge over the code written for this new feature.

### A note on unit testing:

We use the launchpad.net/gocheck package in our unit tests. We ask that you also
use this package for tests that you write, for consistency sake.

