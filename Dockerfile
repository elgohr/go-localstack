FROM localstack/localstack:%s
RUN echo "#!/bin/bash" > timeout-entrypoint.sh \
    && echo "timeout -s SIGKILL %d docker-entrypoint.sh" >> timeout-entrypoint.sh \
    && chmod +x timeout-entrypoint.sh
ENTRYPOINT ["./timeout-entrypoint.sh"]

