# NoBots
Caddy Server plugin to protect your website against web crawlers and bots.

## Usage
The directive for the Caddyfile is really simple. First, you have to place the bomb path next to the `nobots` keyword, for example `bomb.gz` in the example below.

Then you can specify user agent either as strings or regular expresions. When using regular expresions you must add the `regexp` keyword in front of the regex.

Caddyfile example:

```
nobots "bomb.gz" {
  "Googlebot/2.1 (+http://www.googlebot.com/bot.html)"
  "DuckDuckBot"
  regexp "^[Bb]ot"
  regexp "bingbot"
}
```

There is another keyword that is useful in case you want to allow crawlers and bots navigate through specific parts of your website. The keyword is `public` and its values are regular expresions, so you can use it as following:

```
nobots "bomb.gz" {
  "Googlebot/2.1 (+http://www.googlebot.com/bot.html)"
  public "^/public"
  public "^/[a-z]{,5}/public"
}
```

The above example will send the bot to all URIs except those that match with `/public` and `[a-z]{,5}/public`.

NOTE: By default all URIs.


## How to create a bomb
The bomb is not provided within the plugin so you have to create one. In Linux it is really easy, you can use the following commands.

```
dd if=/dev/zero bs=1M count=1024 | gzip > 1G.gzip
dd if=/dev/zero bs=1M count=10240 | gzip > 10G.gzip
dd if=/dev/zero bs=1M count=1048576 | gzip > 1T.gzip
```

To optimize the final bomb you may compress the parts several times:

```
cat 10G.gzip | gzip > 10G.gzipx2
cat 1T.gzip | gzip | gzip | gzip > 1T.gzipx4
 ```
*NOTE*: The extension `.gzipx2` or `.gzipx4` is only to highlight how many times the file was compressed.


