service Nging stop
cd ~/

#local
cp /home/www/eget_download/nging/latest/nging_linux_amd64.zip ./

#remote
#wget https://dl.eget.io/nging/latest/nging_linux_amd64.zip

unzip nging_linux_amd64.zip
rm nging_linux_amd64.zip
service Nging start
