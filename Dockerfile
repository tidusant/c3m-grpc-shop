FROM alpine
# Add Maintainer Info
LABEL maintainer="Duy Ha <duyhph@gmail.com>"
# Set the Current Working Directory inside the container
WORKDIR /app
# Copy exec file and config
COPY shop ./

# Build the Go app
#RUN go build -o main .

# Expose port to the outside world
#EXPOSE 9191
# Run the executable
CMD ["./shop"]