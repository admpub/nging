# Caddy s3browser

## Example config
```
dl.example.com {
	s3browser {
		key ADDKEYHERE
		secret ADDSECRETHERE
		bucket ADDBUCKETHERE
		endpoint nyc3.digitaloceanspaces.com
		secure true
		refresh 5m
		debug false
	}
	proxy / https://examplebucket.nyc3.digitaloceanspaces.com {
		header_upstream Host examplebucket.nyc3.digitaloceanspaces.com
	}
}
```

This will provide directory listing for an S3 bucket (you are able to use minio, or other S3 providers). To serve files via Caddy as well you'll need to use the `proxy` directive as well. The server must be able to have public access to the files in the bucket.

Note: For performance reasons, the file listing is fetched once every 5minutes to reduce load on S3 (or S3 equivalent).

## Prior Art
* This is based on the [Browse plugin](https://github.com/mholt/caddy/tree/master/caddyhttp/browse) that is built into Caddy
* The template is based on the [browse template](https://github.com/dockhippie/caddy/blob/master/rootfs/etc/caddy/browse.tmpl) from Webhippie
* [s3server](https://github.com/jessfraz/s3server) from jessfraz
* [pretty-s3-index-html](https://github.com/nolanlawson/pretty-s3-index-html) by Nolan Lawson
