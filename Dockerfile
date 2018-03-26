#############################################
# Image: patientcoeng/halyard
# Description: Provides a service to support
# horizontal pod autoscaling
#############################################

FROM golang:1.9 as gobuilder
WORKDIR /src/go/halyard
COPY . .
RUN go get -u github.com/kardianos/govendor && \
    govendor sync && \
    go build .

FROM alpine:3.7
COPY --from=gobuilder /src/go/halyard/halyard /halyard
COPY svcinit.sh /svcinit.sh
ENTRYPOINT [ "/svcinit.sh" ]
