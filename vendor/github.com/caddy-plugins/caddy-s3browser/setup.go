package s3browser

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/caddy"
	"github.com/admpub/caddy/caddyhttp/httpserver"
	"github.com/minio/minio-go/v6"
	md2html "github.com/russross/blackfriday"
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
		//buf, _ := json.MarshalIndent(b.Fs, ``, `  `)
		//fmt.Println(string(buf))
	}
	updating = false
	if err != nil {
		return err
	}
	var duration time.Duration
	if len(b.Config.Refresh) == 0 {
		b.Config.Refresh = "5m"
	}
	duration, err = time.ParseDuration(b.Config.Refresh)
	if err != nil {
		fmt.Println("error parsing refresh, falling back to default of 5 minutes")
		duration = 5 * time.Minute
	}
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
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

	tpl, err := template.New("listing").Parse(DefaultTemplate(b.Config))
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
	var (
		minioClient *minio.Client
		err         error
	)
	if len(b.Config.Region) == 0 {
		minioClient, err = minio.New(b.Config.Endpoint, b.Config.Key, b.Config.Secret, b.Config.Secure)
	} else {
		minioClient, err = minio.NewWithRegion(b.Config.Endpoint, b.Config.Key, b.Config.Secret, b.Config.Secure, b.Config.Region)
	}
	if err != nil {
		return fs, err
	}
	if !b.Config.Secure {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		minioClient.SetCustomTransport(tr)
	}
	parseMarkdown := func(objectName string) (r string) {
		objectName = strings.TrimPrefix(objectName, `/`)
		f, err := minioClient.GetObject(b.Config.Bucket, objectName, minio.GetObjectOptions{})
		if err != nil {
			if b.Config.Debug && !strings.Contains(err.Error(), ` key does not exist`) {
				fmt.Println(objectName+`:`, err)
			}
			return
		}
		buf, err := ioutil.ReadAll(f)
		if err != nil {
			if b.Config.Debug {
				fmt.Println(objectName+`:`, err)
			}
		} else {
			buf = md2html.MarkdownCommon(buf)
			r = string(buf)
		}
		f.Close()
		return
	}
	findObjects := func(prefix string) {
		doneCh := make(chan struct{})
		defer close(doneCh)
		objectCh := minioClient.ListObjects(b.Config.Bucket, prefix, true, doneCh)

		for obj := range objectCh {
			if obj.Err != nil {
				continue
			}

			dir, file := path.Split(obj.Key)
			if len(dir) > 0 && dir[:0] != "/" {
				dir = "/" + dir
			}
			if len(dir) == 0 {
				dir = "/" // if dir is empty, then set to root
			}
			// Note: dir should start & end with / now
			folders := getFolders(dir)
			if len(folders) < 3 {
				// files are in root
				// less than three bc "/" split becomes ["",""]
				// Do nothing as file will get added below & root already exists
			} else {
				// TODO: loop through folders and ensure they are in the tree
				// make sure to add folder to parent as well
				foldersLen := len(folders)
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
							Folders: []Folder{Folder{Name: folder}},
						}
					}
					// check if folder itself exists
					if _, ok := fs[folder]; !ok {
						// create parent
						fs[folder] = Directory{
							Path: folder,
						}
						tmp := fs[parent]
						tmp.Folders = append(fs[parent].Folders, Folder{Name: folder})
						fs[parent] = tmp
					}
				}
			}

			// STEP Two
			// add file to directory
			tempFile := File{Name: file, Bytes: obj.Size, Date: obj.LastModified, Folder: joinFolders(folders)}
			fsCopy := fs[tempFile.Folder]
			fsCopy.Path = tempFile.Folder
			fsCopy.Files = append(fsCopy.Files, tempFile) // adding file list of files
			if file == `README.md` {
				objectName := path.Join(fsCopy.Path, `README.md`)
				fsCopy.README = parseMarkdown(objectName)
			}
			fs[tempFile.Folder] = fsCopy
		} // end looping through all the files
	}
	for _, prefix := range b.Config.prefixes {
		findObjects(prefix)
	}
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
		case "region":
			b.Config.Region, err = StringArg(c)
		case "prefix":
			b.Config.Prefix, err = StringArg(c)
			saved := map[string]struct{}{}
			for _, prefix := range strings.Split(b.Config.Prefix, "|") {
				prefix = strings.TrimSpace(prefix)
				prefix = strings.TrimPrefix(prefix, `/`)
				if _, ok := saved[prefix]; !ok {
					b.Config.prefixes = append(b.Config.prefixes, prefix)
					saved[prefix] = struct{}{}
				}
			}
		case "cdnurl":
			b.Config.CDNURL, err = StringArg(c)
			if len(b.Config.CDNURL) > 0 {
				b.Config.CDNURL = strings.TrimSuffix(b.Config.CDNURL, `/`)
			}
		case "secure":
			b.Config.Secure, err = BoolArg(c)
		case "refresh":
			b.Config.Refresh, err = StringArg(c)
		case "debug":
			b.Config.Debug, err = BoolArg(c)
		case "csscdn":
			b.Config.CSSCDN, err = StringArg(c)
		default:
			return c.Errf("Unknown s3browser arg: %s", c.Val())
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// StringArg Assert only one arg and return it
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
