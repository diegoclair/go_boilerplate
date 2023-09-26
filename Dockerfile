FROM golang:latest

# Add Maintainer info
LABEL maintainer="Rodrigues Diego <diego93rodrigues@gmail.com>"

# Set the current working directory inside the container 
WORKDIR /app

# Copy the source from the current directory to the working Directory inside the container 
COPY . /app

# Add docker-compose-wait tool 
ENV WAIT_VERSION 2.9.0
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait

EXPOSE 5000

#This is used to run the application with live reload
RUN go install github.com/githubnemo/CompileDaemon@latest

ENTRYPOINT /wait && CompileDaemon --build="go build -buildvcs=false -o myapp ." --command="./myapp"
