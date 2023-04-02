FROM kldtks/edge:edge-builder-goa AS builder

COPY ./ /src/

RUN rm -rf src/gen
RUN rm -rf src/cmd

RUN cd /src && go mod tidy && export PATH=$PATH:/go/bin && goa gen admin-panel/design && goa example admin-panel/design && go mod tidy && go build ./cmd/api_service

FROM bitnami/kubectl

USER root

RUN mkdir app
COPY --from=builder /src/api_service /app/api_service
COPY --from=builder /src/setup /setup
COPY --from=builder /src/init_data /init_data


RUN chmod 777 /app/api_service

ENTRYPOINT [ "/app/api_service" ]
