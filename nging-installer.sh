# via curl:
# sudo sh -c "$(curl -fsSL https://raw.githubusercontent.com/admpub/nging/master/nging-installer.sh)"

# via wget:
# sudo sh -c "$(wget https://raw.githubusercontent.com/admpub/nging/master/nging-installer.sh -O -)"
# or
# sudo wget https://raw.githubusercontent.com/admpub/nging/master/nging-installer.sh -O ./nging-installer.sh && chmod +x ./nging-installer.sh && ./nging-installer.sh

osname=`uname -s`
arch=`uname -m`
version="3.0.3"
url="https://img.nging.coscms.com/nging/v${version}/"

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
    "armv7l"|"armv7") 
        arch="arm-7"
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

binname="nging"
filename="nging_${osname}_$arch"
filefullname="$filename.zip"

exitOnFailure() {
    echo "❌ command failed"
    exit 1
}

install() {

    wget -c "${url}$filefullname" -O ./$filefullname || exitOnFailure

    unzip $filefullname -d ./$filename || exitOnFailure 

    cp -R ./$filename/$filename/* ./$filename || exitOnFailure
    rm -rf "./$filename/$filename"

    rm $filefullname
    chmod +x ./$filename/$binname || exitOnFailure
    cd ./$filename
    #./$binname

    ./$binname service install || exitOnFailure
    ./$binname service start || exitOnFailure
    echo "🎉 Congratulations! Installed successfully."
}

upgrade() {
    # 停止服务
    cd ./$filename
    ./$binname service stop
    cd ../

    wget -c "${url}$filefullname" -O ./$filefullname || exitOnFailure

    unzip $filefullname -d ./$filename || exitOnFailure 

    cp -R ./$filename/$filename/* ./$filename || exitOnFailure
    rm -rf "./$filename/$filename"

    rm $filefullname
    chmod +x ./$filename/$binname || exitOnFailure
    cd ./$filename
    #./$binname

    # 再次启动服务
    ./$binname service start
    echo "🎉 Congratulations! Successfully upgraded."
}

uninstall() {
    cd ./$filename
    ./$binname service stop || exitOnFailure
    ./$binname service uninstall || exitOnFailure
    echo "🎉 Congratulations! Successfully uninstalled."
    cd ../
    rm -r ./$filename || exitOnFailure
    echo "🎉 Congratulations! File deleted successfully."
}


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
