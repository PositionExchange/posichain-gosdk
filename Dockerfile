FROM golang:1.16.3

RUN apt clean && apt update && apt install -y libgmp-dev libssl-dev make gcc g++ dnsutils

WORKDIR /workspace
RUN git clone https://github.com/PositionExchange/mcl.git
RUN git clone https://github.com/PositionExchange/bls.git
RUN cd bls && make -j8 BLS_SWAP_G=1
RUN cp ./bls/lib/libbls384_256.so /usr/local/lib
RUN cp ./mcl/lib/libmcl.so /usr/local/lib
RUN echo "/usr/local/lib" >> /etc/ld.so.conf
RUN ldconfig

WORKDIR /workspace/sdk
COPY . .
RUN make
ENTRYPOINT ./psc
