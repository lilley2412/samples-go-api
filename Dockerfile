FROM scratch

ADD main /app/main

ENTRYPOINT [ "/app/main" ]