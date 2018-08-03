FROM debian:jessie
RUN apt-get update && apt-get install -y apt-transport-https && apt install -y curl
ADD host_age /
ADD ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["./host_age"]
