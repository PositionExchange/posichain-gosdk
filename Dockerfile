FROM golang:1.16.3

RUN apt clean && apt update && apt install -y libgmp-dev libssl-dev make gcc g++ dnsutils
RUN export GOPATH="/go"
WORKDIR /go/src/github.com/PositionExchange
RUN git clone https://github.com/PositionExchange/mcl.git
RUN git clone https://github.com/PositionExchange/bls.git
RUN cd bls && make -j8 BLS_SWAP_G=1
RUN cp ./bls/lib/libbls384_256.so /usr/local/lib
RUN cp ./mcl/lib/libmcl.so /usr/local/lib
RUN echo "/usr/local/lib" >> /etc/ld.so.conf
RUN ldconfig

WORKDIR /go/src/github.com/PositionExchange/posichain-gosdk
COPY . .
RUN make
ENTRYPOINT ["./psc"]
