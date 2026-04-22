NGING_BUILDER_PATH=./tool/nging-builder
# * build *
if [ "$1" != "docker" ];then
cd $NGING_BUILDER_PATH && ./run.sh linux_386,linux_amd64,linux_arm64 && cd ../../
buikdkit=`docker images | grep moby/buildkit | awk '{print $1}' | head -n 1`
if [ "$buikdkit" = "" ];then
    docker buildx create --name container-builder --driver docker-container --use --bootstrap
fi
fi
docker buildx build . \
    --platform linux/386,linux/amd64,linux/arm64 \
    -t "admpub/nging-dockermgr:latest" \
    --build-arg VERSION=$(grep NgingVersion $NGING_BUILDER_PATH/builder.conf | sed 's/NgingVersion[ ]*:[ ]*//g' | sed 's/"//g' | sed 's/ //g') \
    --push
