# Locale detection for caddy

## Configuration schema

    locale <availableLocales...> {
      detect <methods...>
      cookie <cookie name>
      available <availableLocales...>
      path <path scope>
    }

A `method` can be currently `cookie` or `header`. If `cookie` is added, `cookie name` defines from which cookie the
locale is read. The `header` method extracts the locales from the `Accept-Language` header. The first `availableLocale`
is also the default locale, which is picked if none of the locales from the detection methods is in `availableLocales`.

The defaults are: `methods = [header]`,  `cookie name = locale`, `path scope = /`.

## Example

    locale en de {
      detect cookie header
    }

    rewrite {
      ext /
      to index.{>Detected-Locale}.html index.html
    }

    header / Vary "Cookie, Accept-Language"
