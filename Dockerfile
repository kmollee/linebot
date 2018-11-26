
#golang #docker
# By using a scratch base image, we save about ~5MB over Alpine base images and we ship with a smaller attack surface.
# This is the first stage, for building things that will be required by the
# final stage (notably the binary)
FROM golang as builder
# Copy in just the go.mod and go.sum files, and download the dependencies. By
# doing this before copying in the other dependencies, the Docker build cache
# can skip these steps so long as neither of these two files change.
ADD app /app
RUN cd /app && go mod download
# Build the Go app with CGO_ENABLED=0 so we use the pure-Go implementations for
# things like DNS resolution (so we don't build a binary that depends on system
# libraries)
RUN cd /app && CGO_ENABLED=0 go build -o /goapp
# Create a "nobody" non-root user for the next image by crafting an /etc/passwd
# file that the next image can copy in. This is necessary since the next image
# is based on scratch, which doesn't have adduser, cat, echo, or even sh.
# RUN echo "nobody:x:65534:65534:nodoby:/:" > /etc_passwd
# The second and final stage
FROM scratch

ENV HOME /app
WORKDIR /app
# Copy the binary from the builder stage
COPY --from=builder  /goapp .
# Copy the certs from the builder stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Copy the /etc_passwd file we created in the builder stage into /etc/passwd in
# the target stage. This creates a new non-root user as a security best
# practice.
# COPY --from=0 /etc_passwd /etc/passwd
# Run as the new non-root by default
# USER nobody

EXPOSE 80
# ENTRYPOINT [ "/app" ]

CMD [ "./goapp" ]