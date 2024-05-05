FROM golang:1.21-bookworm as api

# Define working directory.
WORKDIR /go/app

# Copy the local package files to the container's workspace.
COPY ./app-bin /go/app/app-bin

# Set the environment variables
ENV ENV="GA" \
  # Time zone
  TZ="UTC" \
  # Default language
  LANGUAGE="en" \
  # Default port
  PORT="3000"


# Expose port
EXPOSE 3000

# HEALTHCHECK
HEALTHCHECK --interval=5s --timeout=3s \
  CMD ps aux | grep 'app\-bin' || exit 1

# ENTRYPOINT
# ENTRYPOINT ["/docker-entrypoint"]

# Launch the app
CMD ["/go/app/app-bin", "server"]
