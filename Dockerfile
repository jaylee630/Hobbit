FROM 192.168.154.250/library/golang:1.9.1-glide
MAINTAINER Sean.Wang <sean.wang@ucloud.cn>

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

EXPOSE 9091

ADD . $GOPATH/src/github.com/reposkeeper/golang-admin-basic/
WORKDIR $GOPATH/src/github.com/reposkeeper/golang-admin-basic/

RUN make build
CMD cd ./build && ./golang-admin-basic
