FROM 416670754337.dkr.ecr.eu-west-2.amazonaws.com/ci-corretto-runtime-21:latest

WORKDIR /opt
COPY /app .
COPY docker_start.sh .
COPY /app/run-elastic-search.sh .

RUN chmod +x ./docker_start.sh ./run-elastic-search.sh

CMD ["./docker_start.sh"]
