package s3browser

import (
	"crypto/tls"
	"fmt"
	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddyhttp/httpserver"
	"github.com/minio/minio-go/v6"
	"html/template"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

var (
	updating bool
)

func init() {
	caddy.RegisterPlugin("s3browser", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

// setup configures a new S3BROWSER middleware instance.
func setup(c *caddy.Controller) error {
	var err error
	cfg := httpserver.GetConfig(c)

	b := &Browse{}
	if err = parse(b, c); err != nil {
		return err
	}
	if b.Config.Debug {
		fmt.Println("Config:")
		fmt.Println(b.Config)
	}
	updating = true
	if b.Config.Debug {
		fmt.Println("Fetching Files..")
	}
	b.Fs, err = getFiles(b)
	if b.Config.Debug {
		fmt.Println("Files...")
		fmt.Println(b.Fs)
	}
	updating = false
	if err != nil {
		return err
	}
	var duration time.Duration
	if b.Config.Refresh == "" {
		b.Config.Refresh = "5m"
	}
	duration, err = time.ParseDuration(b.Config.Refresh)
	if err != nil {
		fmt.Println("error parsing refresh, falling back to default of 5 minutes")
		duration = 5 * time.Minute
	}
	ticker := time.NewTicker(duration)
	go func() {
		// create more indexes every X minutes based off interval
		for range ticker.C {
			if !updating {
				if b.Config.Debug {
					fmt.Println("Updating Files..")
				}
				if b.Fs, err = getFiles(b); err != nil {
					fmt.Println(err)
					updating = false
				}
			}
		}
	}()

	tpl, err := template.New("listing").Parse(defaultTemplate)
	if err != nil {
		return err
	}
	b.Template = tpl

	cfg.AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		b.Next = next
		return b
	})

	return nil
}

func getFiles(b *Browse) (map[string]Directory, error) {
	updating = true
	fs := make(map[string]Directory)
	fs["/"] = Directory{
		Path: "/",
	}
	minioClient, err := minio.New(b.Config.Endpoint, b.Config.Key, b.Config.Secret, b.Config.Secure)
	if err != nil {
		return fs, err
	}

	if !b.Config.Secure {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		minioClient.SetCustomTransport(tr)
	}

	doneCh := make(chan struct{})
	defer close(doneCh)
	objectCh := minioClient.ListObjects(b.Config.Bucket, "", true, doneCh)

	for obj := range objectCh {
		if obj.Err != nil {
			continue
		}

		dir, file := path.Split(obj.Key)
		if len(dir) > 0 && dir[:0] != "/" {
			dir = "/" + dir
		}
		if dir == "" {
			dir = "/" // if dir is empty, then set to root
		}
		// Note: dir should start & end with / now

		if len(getFolders(dir)) < 3 {
			// files are in root
			// less than three bc "/" split becomes ["",""]
			// Do nothing as file will get added below & root already exists
		} else {
			// TODO: loop through folders and ensure they are in the tree
			// make sure to add folder to parent as well
			foldersLen := len(getFolders(dir))
			for i := 2; i < foldersLen; i++ {
				parent := getParent(getFolders(dir), i)
				folder := getFolder(getFolders(dir), i)
				if b.Config.Debug {
					fmt.Printf("folders: %q i: %d parent: %s folder: %s\n", getFolders(dir), i, parent, folder)
				}

				// check if parent exists
				if _, ok := fs[parent]; !ok {
					// create parent
					fs[parent] = Directory{
						Path:    parent,
						Folders: []Folder{Folder{Name: getFolder(getFolders(dir), i)}},
					}
				}
				// check if folder itself exists
				if _, ok := fs[folder]; !ok {
					// create parent
					fs[folder] = Directory{
						Path: folder,
					}
					tmp := fs[parent]
					tmp.Folders = append(fs[parent].Folders, Folder{Name: getFolder(getFolders(dir), i)})
					fs[parent] = tmp
				}
			}
		}

		// STEP Two
		// add file to directory
		tempFile := File{Name: file, Bytes: obj.Size, Date: obj.LastModified, Folder: joinFolders(getFolders(dir))}
		fsCopy := fs[joinFolders(getFolders(dir))]
		fsCopy.Path = joinFolders(getFolders(dir))
		fsCopy.Files = append(fsCopy.Files, tempFile) // adding file list of files
		fs[joinFolders(getFolders(dir))] = fsCopy
	} // end looping through all the files
	updating = false
	return fs, nil
}

func getFolders(s string) []string {
	// first and last entry should be empty
	return strings.Split(s, "/")
}

func joinFolders(s []string) string {
	return strings.Join(s, "/")
}

func getParent(s []string, i int) string {
	// trim one from end
	if i < 3 {
		return "/"
	}
	s[i-1] = ""
	return joinFolders(s[0:(i)])
}

func getFolder(s []string, i int) string {
	if i < 3 {
		s[2] = ""
		return joinFolders(s[0:3])
	}
	s[i] = ""
	return joinFolders(s[0:(i + 1)])
}

func parse(b *Browse, c *caddy.Controller) (err error) {
	c.RemainingArgs()
	b.Config = Config{}
	b.Config.Secure = true
	b.Config.Debug = false
	for c.NextBlock() {
		var err error
		switch c.Val() {
		case "key":
			b.Config.Key, err = StringArg(c)
		case "secret":
			b.Config.Secret, err = StringArg(c)
		case "endpoint":
			b.Config.Endpoint, err = StringArg(c)
		case "bucket":
			b.Config.Bucket, err = StringArg(c)
		case "secure":
			b.Config.Secure, err = BoolArg(c)
		case "refresh":
			b.Config.Refresh, err = StringArg(c)
		case "debug":
			b.Config.Debug, err = BoolArg(c)
		default:
			return c.Errf("Unknown s3browser arg: %s", c.Val())
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// Assert only one arg and return it
func StringArg(c *caddy.Controller) (string, error) {
	args := c.RemainingArgs()
	if len(args) != 1 {
		return "", c.ArgErr()
	}
	return args[0], nil
}

func BoolArg(c *caddy.Controller) (bool, error) {
	args := c.RemainingArgs()
	if len(args) != 1 {
		return true, c.ArgErr()
	}
	return strconv.ParseBool(args[0])
}

const defaultTemplate = `<!DOCTYPE html>
<html>
	<head>
		<title>{{ .ReadableName }} | S3 Browser</title>

		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">

		<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/3.3.6/css/bootstrap.min.css">
		<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/flat-ui/2.3.0/css/flat-ui.min.css">

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
					Powered by <a href="https://caddyserver.com">Caddy</a>
				</div>
			</div>
		</nav>

		<div class="container-fluid">
			<ol class="breadcrumb">
				<li>
					<a href="/"><span class="glyphicon glyphicon-home" aria-hidden="true"></span></a>
				</li>
				{{ range .Breadcrumbs }}
					<li>
						<a href="/{{ html .Link }}">
							{{ html .ReadableName }}
						</a>
					</li>
				{{ end }}
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
						{{ if ne .Path "/" }}
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
						{{ end }}
						{{ range .Folders }}
							<tr>
								<td class="icon">
									<span class="glyphicon glyphicon-folder-close" aria-hidden="true"></span>
								</td>
								<td class="name">
									<a href="{{ html .Name }}">
										{{ .ReadableName }}
									</a>
								</td>
								<td class="size">
									&mdash;
								</td>
								<td class="modified">
									&mdash;
								</td>
							</tr>
						{{ end }}
						{{ range .Files }}
							{{ if ne .Name ""}}
							<tr>
								<td class="icon">
									<span class="glyphicon glyphicon-file" aria-hidden="true"></span>
								</td>
								<td class="name">
									<a href="./{{ html .Name }}">
										{{ .Name }}
									</a>
								</td>
								<td class="size">
									{{ .HumanSize }}
								</td>
								<td class="modified">
									{{ .HumanModTime "01/02/2006 03:04:05 PM" }}
								</td>
							</tr>
							{{ end }}
						{{ end }}
					</tbody>
				</table>
			</div>
		</div>
	</body>
</html>
`
