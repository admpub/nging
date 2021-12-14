pwd=$GOPATH/src/github.com/nging-plugins
dir=$pwd
filelist=`ls $dir`
for file in $filelist; do
    if test -d $pwd/$file; then
	    echo "push: $pwd/$file"
        cd $pwd/$file
        git add .
        git commit -m update
        git push
    fi
done