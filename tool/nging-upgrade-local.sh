userDir=/home/admpub

if [ "$1" != "" ]; then
    userDir="$1"
fi

osname=`uname -s`
arch=`uname -m`

case "$arch" in
    "x86_64") 
        arch="amd64"
        ;;
    "i386"|"i686") 
        arch="386"
        ;;
    "aarch64_be"|"aarch64"|"armv8b"|"armv8l"|"armv8"|"arm64") 
        arch="arm64"
        ;;
    "armv7") 
        arch="arm-7"
        ;;
    "armv7l") 
        arch="arm-6"
        ;;
    "armv6") 
        arch="arm-6"
        ;;
    "armv5"|"arm") 
        arch="arm-5"
        ;;
    *)
        echo "Unsupported Arch:${arch}"
        exit 1
        ;;
esac

case $osname in
    "Darwin") 
        osname="darwin"
        ;;
    "Linux") 
        osname="linux"
        ;; 
    *)
        echo "Unsupported System:${osname}"
        exit 1
        ;;
esac

osArch=${osname}_${arch}

cd $userDir/nging
./nging service stop
cd $userDir
tar -zxvf $userDir/go/src/github.com/admpub/nging/dist/nging_$osArch.tar.gz -C $userDir/nging
cp -R ./nging/nging_$osArch/* ./nging
rm -rf ./nging/nging_$osArch
cd $userDir/nging
./nging service start
#cd $userDir
