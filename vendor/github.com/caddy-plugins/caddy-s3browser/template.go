package s3browser

type DepCSS []string

func (d DepCSS) String() string {
	var s string
	for _, css := range d {
		s += `<link rel="stylesheet" href="` + css + `">`
	}
	return s
}

var (
	Cloudflare = DepCSS{
		`//cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/3.3.6/css/bootstrap.min.css`,
		`//cdnjs.cloudflare.com/ajax/libs/flat-ui/2.3.0/css/flat-ui.min.css`,
	}

	BootCDN = DepCSS{
		`//cdn.bootcdn.net/ajax/libs/flat-ui/2.3.0/css/vendor/bootstrap/css/bootstrap.min.css`,
		`//cdn.bootcdn.net/ajax/libs/flat-ui/2.3.0/css/flat-ui.min.css`,
	}

	Depencies = BootCDN
)

var DefaultTemplate = func() string {
	return `<!DOCTYPE html>
<html>
	<head>
		<title>{{ .ReadableName }} | S3 Browser</title>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		` + Depencies.String() + `
		<style>
			body {
				cursor: default;
			}
			.navbar {
				margin-bottom: 20px;
			}
			.credits {
				padding-left: 15px;
				padding-right: 15px;
			}
			h1 {
				font-size: 20px;
				margin: 0;
			}
			th .glyphicon {
				font-size: 15px;
			}
			table .icon {
				width: 30px;
			}
		</style>
    <!-- template source from https://raw.githubusercontent.com/dockhippie/caddy/master/rootfs/etc/caddy/browse.tmpl -->
	</head>
	<body>
		<nav class="navbar navbar-inverse navbar-static-top">
			<div class="container-fluid">
				<div class="navbar-header">
					<a class="navbar-brand" href="/">
						S3 Browser
					</a>
				</div>

				<div class="navbar-text navbar-right hidden-xs credits">
					Powered by <a href="https://github.com/admpub/nging">Nging</a>
				</div>
			</div>
		</nav>

		<div class="container-fluid">
			<ol class="breadcrumb">
				<li>
					<a href="/"><span class="glyphicon glyphicon-home" aria-hidden="true"></span></a>
				</li>
				{{- range .Breadcrumbs -}}
					<li>
						<a href="/{{ html .Link }}">
							{{ html .ReadableName }}
						</a>
					</li>
				{{- end -}}
			</ol>

			<div class="panel panel-default">
				<table class="table table-hover table-striped">
					<thead>
						<tr>
							<th class="icon"></th>
							<th class="name">
								Name
							</th>
							<th class="size col-sm-2">
								Size
							</th>
							<th class="modified col-sm-2">
								Modified
							</th>
						</tr>
					</thead>

					<tbody>
						{{- if ne .Path "/" -}}
							<tr>
								<td>
									<span class="glyphicon glyphicon-arrow-left" aria-hidden="true"></span>
								</td>
								<td>
									<a href="..">
										Go up
									</a>
								</td>
								<td>
									&mdash;
								</td>
								<td>
									&mdash;
								</td>
							</tr>
						{{- end -}}
						{{- range .Folders -}}
							<tr>
								<td class="icon">
									<span class="glyphicon glyphicon-folder-close" aria-hidden="true"></span>
								</td>
								<td class="name">
									<a href="{{ html .Name }}">
										{{- .ReadableName -}}
									</a>
								</td>
								<td class="size">
									&mdash;
								</td>
								<td class="modified">
									&mdash;
								</td>
							</tr>
						{{- end -}}
						{{- range .Files -}}
							{{- if ne .Name "" -}}
							<tr>
								<td class="icon">
									<span class="glyphicon glyphicon-file" aria-hidden="true"></span>
								</td>
								<td class="name">
									<a href="{{ $.CDNURL }}{{ .Folder }}{{ html .Name }}">
										{{- .Name -}}
									</a>
								</td>
								<td class="size">
									{{- .HumanSize -}}
								</td>
								<td class="modified">
									{{- .HumanModTime "01/02/2006 03:04:05 PM" -}}
								</td>
							</tr>
							{{- end -}}
						{{- end -}}
					</tbody>
				</table>
			</div>
		</div>
	</body>
</html>
`
}
