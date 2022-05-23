FROM	node:10.16.0-jessie

ARG http_proxy
ARG https_proxy
ARG HTTP_PROXY
ARG HTTPS_PROXY

WORKDIR /usr/src/wait

RUN 	wget https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh
RUN 	chmod u+x wait-for-it.sh
RUN	chmod 777 /usr/src/wait/wait-for-it.sh

WORKDIR /usr/src/app

COPY 	package*.json ./

RUN 	npm cache verify

RUN 	npm install
RUN 	npm install -g nodemon

COPY 	. .

RUN useradd server -m && \
    echo server:sdkflmdskofdo! | chpasswd
RUN chown -R server:server /usr/src/app

USER server


EXPOSE 4000

CMD [ "npm", "start" ]
