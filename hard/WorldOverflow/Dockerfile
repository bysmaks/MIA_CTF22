FROM ubuntu
RUN apt update -y && apt install -y tcpd xinetd
RUN useradd user
COPY src/world_overflow_config /etc/xinetd.d/world_overflow
COPY src/world_overflow /world_overflow
COPY src/flag.txt /flag.txt
RUN chmod 777 world_overflow
EXPOSE 8888
CMD xinetd -dontfork
