FROM golang:1.21.6

LABEL maintainer="Zyro Kamson <zyro.kamson@gmail.com>"

# Move to working directory (/usr/src/app).
WORKDIR /usr/src/app

# Install Air for hot reloading.
RUN go install github.com/cosmtrek/air@latest

# Copy the code into the container.
COPY . .
RUN go mod tidy