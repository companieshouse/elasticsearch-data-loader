FROM 416670754337.dkr.ecr.eu-west-2.amazonaws.com/ci-corretto-runtime-21:latest

RUN dnf install -y awscli-2-2.32.1

WORKDIR /opt
COPY /app .
COPY docker_start.sh .

RUN chmod +x ./docker_start.sh

CMD ["./docker_start.sh"]
