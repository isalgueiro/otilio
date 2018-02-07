FROM ubuntu:16.04

MAINTAINER Tudor Golubenco <tudor@elastic.co>

# install fpm
RUN \
    apt-get update && \
    apt-get install -y --no-install-recommends \
        # https://github.com/elastic/beats/pull/6302/files
        autoconf build-essential libffi-dev ruby-dev rpm zip dos2unix libgmp3-dev

RUN gem install fpm
