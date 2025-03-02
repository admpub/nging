# via curl:
# sudo sh -c "$(curl -fsSL https://raw.githubusercontent.com/admpub/nging/master/nging-installer.sh)"

# via wget:
# sudo sh -c "$(wget https://raw.githubusercontent.com/admpub/nging/master/nging-installer.sh -O -)"
# or
# sudo wget https://raw.githubusercontent.com/admpub/nging/master/nging-installer.sh -O ./nging-installer.sh && chmod +x ./nging-installer.sh && ./nging-installer.sh

osname=`uname -s`
arch=`uname -m`
version="5.2.7"

if [ "$2" != "" ] && [ "$2" != "-" ]; then
    version="$2"
fi

url="https://img.nging.coscms.com/nging/v${version}/"
savedir="nging"

case "$arch" in
    "x86_64"|"amd64") 
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
        savedir="nging_linux_arm-7" # 兼容旧版
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
    "FreeBSD") 
        osname="freebsd"
        ;; 
    *)
        echo "Unsupported System:${osname}"
        exit 1
        ;;
esac

binname="nging"
filename="nging_${osname}_$arch"
filefullname="$filename.tar.gz"

exitOnFailure() {
    echo "❌ command failed"
    exit 1
}

install() {

    wget -c "${url}$filefullname" -O ./$filefullname || exitOnFailure

    mkdir ./$savedir

    tar -zxvf $filefullname -C ./$savedir || exitOnFailure
    #unzip $filefullname -d ./$filename || exitOnFailure 

    if [ -d "./$savedir/$filename" ]; then
        cp -R ./$savedir/$filename/* ./$savedir || exitOnFailure
        rm -rf "./$savedir/$filename"
    fi

    rm $filefullname
    chmod +x ./$savedir/$binname || exitOnFailure
    cd ./$savedir
    #./$binname

    ./$binname service install || exitOnFailure
    ./$binname service start || exitOnFailure
    echo "🎉 Congratulations! Installed successfully."
}

upgrade() {
    # 停止服务
    cd ./$savedir
    ./$binname service stop
    cd ../

    wget -c "${url}$filefullname" -O ./$filefullname || exitOnFailure

    mkdir ./$savedir
    
    tar -zxvf $filefullname -C ./$savedir || exitOnFailure
    #unzip $filefullname -d ./$filename || exitOnFailure 

    sleep 5s && pkill `pwd`/$savedir/$binname

    if [ -d "./$savedir/$filename" ]; then
        cp -R ./$savedir/$filename/* ./$savedir || exitOnFailure
        rm -rf "./$savedir/$filename"
    fi

    rm $filefullname
    chmod +x ./$savedir/$binname || exitOnFailure
    cd ./$savedir
    #./$binname

    # 再次启动服务
    ./$binname service start
    echo "🎉 Congratulations! Successfully upgraded."
}

uninstall() {
    cd ./$savedir
    ./$binname service stop || exitOnFailure
    ./$binname service uninstall || exitOnFailure

    cd ../
    sleep 5s && pkill `pwd`/$savedir/$binname

    echo "🎉 Congratulations! Successfully uninstalled."
    rm -r ./$savedir || exitOnFailure
    echo "🎉 Congratulations! File deleted successfully."
}

if [ "$3" != "" ]; then
    savedir="$3"
fi

case "$1" in
    "up"|"upgrade")
        upgrade
        ;;
    "un"|"uninstall")
        uninstall
        ;;
    "install")
        install
        ;;
    *)
        install
        ;;
esac
