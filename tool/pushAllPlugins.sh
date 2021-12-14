pwd=../../../nging-plugins
dir=.
filelist=`ls $dir`
for file in $filelist; do
	    echo $pwd/$file
    if test -d $file; then
	    echo $pwd/$file
        cd $pwd/$file
        git add .
        git commit -m update
        git push
        cd $pwd
    fi
done