FROM 172.18.127.68:80/base-images/ubuntu22_cuda118_python311:v1.1

ENV PATH=/root/miniconda3/bin:/usr/local/cuda-11.8/bin:$PATH LD_LIBRARY_PATH=/usr/local/cuda-11.8/lib64:$LD_LIBRARY_PATH

RUN echo 'root:adminADMIN123' | chpasswd

RUN apt install sshpass -y