FROM 172.18.127.68:80/support/platform_base_tools:cuda11.8-cudnn8-ubuntu22.04-sshpass-password-ibtools

RUN wget https://repo.anaconda.com/miniconda/Miniconda3-py311_23.11.0-2-Linux-x86_64.sh
RUN bash Miniconda3-py311_23.11.0-2-Linux-x86_64.sh -p /opt/miniconda -b
RUN rm Miniconda3-py311_23.11.0-2-Linux-x86_64.sh

ENV PATH=/opt/miniconda/bin:${PATH}
