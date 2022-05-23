FROM node

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY . /usr/src/app/
RUN npm install

RUN useradd server -m && \
    echo server:sdkflmdskofdo! | chpasswd
RUN chown -R server:server /usr/src/app
RUN echo "CTF{so_easy_peasy}" > /tmp/flag.txt
RUN chmod o+r /tmp/flag.txt

USER server
EXPOSE 1337
CMD node index.js
