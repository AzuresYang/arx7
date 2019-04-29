# 指定基本镜像为alpine,一个5M大小的轻量级别的linux系统
FROM alpine:latest

# 创建目录，拷贝arx执行app文件
RUN mkdir /arxapp
RUN mkdir /arxapp/log
ADD arx /arxapp/
RUN chmod +x /arxapp/arx

# 启动容器后执行的脚本
# ENTRYPOINT exec nohup /arxapp/arx spider --port=31001 > spider.log 2>&1 &
ENTRYPOINT exec  /arxapp/arx spider --port=31001 --output=/arxapp/log/spider.log
