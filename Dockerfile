FROM golang:1.16-alpine

RUN apk add --update --no-cache \
    make \
    git \
    curl

WORKDIR /go/src/github.com/gandarez/changelog-action

COPY . .

# build
RUN make build-linux

# apply permissions
RUN chmod a+x ./build/linux/amd64/changelog

# symbolic link
RUN ln -s /go/src/github.com/gandarez/changelog-action/build/linux/amd64/changelog /bin/

# Specify the container's entrypoint as the action
ENTRYPOINT ["/bin/changelog"]
