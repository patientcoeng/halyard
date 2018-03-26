#############################################
# Image: patientcoeng/halyard
# Description: Provides a service to support
# horizontal pod autoscaling
#############################################

FROM golang:1.9 as gobuilder
WORKDIR /go/src/github.com/patientcoeng/halyard
COPY . .
RUN go get -u github.com/kardianos/govendor && \
    govendor sync && \
    go build .

FROM alpine:3.7
COPY --from=gobuilder /go/src/github.com/patientcoeng/halyard/halyard /halyard
COPY svcinit.sh /svcinit.sh
ENTRYPOINT [ "/svcinit.sh" ]
