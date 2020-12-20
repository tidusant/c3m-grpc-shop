FROM alpine
# Add Maintainer Info
LABEL maintainer="Duy Ha <duyhph@gmail.com>"
# Set the Current Working Directory inside the container
WORKDIR /app
# Copy exec file and config
COPY main ./

# Run the executable
CMD ["./main"]