#FROM golang:1.21 as builder
#
#RUN go install github.com/cespare/reflex@v0.3.1
#
#WORKDIR /app
#COPY reflex.conf /
#
#COPY start.sh /
#
#ENTRYPOINT ["reflex", "-c", "/reflex.conf"]

FROM golang:1.21

WORKDIR /app
COPY . .

# Set the entry point
ENTRYPOINT ["reflex", "-c", "/reflex.conf"]