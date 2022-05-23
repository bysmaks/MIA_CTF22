FROM node

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY . /usr/src/app/
RUN npm install
RUN useradd server -m && \
    echo server:sdkflmdskofdo! | chpasswd
RUN chown -R server:server /usr/src/app

USER server
EXPOSE 9090
CMD npm run start
