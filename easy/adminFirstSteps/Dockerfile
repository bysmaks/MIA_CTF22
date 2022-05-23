FROM ubuntu:18.04

RUN apt update && apt install  openssh-server -y
RUN useradd server -m && \
    echo server:server | chpasswd
RUN echo 'server:SAdfsf4f!' | chpasswd
RUN echo "CTF{my_first_ssh_connection!}" > /tmp/flag
RUN chown root:root /tmp/flag && chmod -w /tmp/flag
RUN service ssh start

CMD ["/usr/sbin/sshd","-D"]

