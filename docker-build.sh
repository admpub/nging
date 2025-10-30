NGING_BUILDER_PATH=./tool/nging-builder
# * build *
cd $NGING_BUILDER_PATH && ./run.sh linux_amd64 && cd ../../
docker build . -t "admpub/nging-dockermgr:latest" --build-arg VERSION=$(grep NgingVersion $NGING_BUILDER_PATH/builder.conf | sed 's/NgingVersion[ ]*:[ ]*//g' | sed 's/"//g' | sed 's/ //g')