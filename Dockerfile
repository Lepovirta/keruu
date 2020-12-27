FROM gcr.io/distroless/base-debian10
COPY keruu /usr/local/bin/keruu
ENTRYPOINT ["keruu"]
