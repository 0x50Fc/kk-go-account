FROM registry-internal.cn-hangzhou.aliyuncs.com/kk/kk-golang

RUN mkdir github.com

RUN mkdir github.com/hailongz

COPY . github.com/hailongz/kk-go-account

RUN go get -u github.com/hailongz/kk-go

RUN go get -u github.com/hailongz/kk-go-db

RUN go get -u github.com/hailongz/kk-go-task

RUN go install github.com/hailongz/kk-go-account

ENV KK_ADDR 127.0.0.1:87

ENV KK_DB_URL root:123456@/db

ENV KK_DB_PREFIX account_

ENV KK_NAME kk.account.

CMD kk-go-account $KK_NAME $KK_ADDR $KK_DB_URL $KK_DB_PREFIX
