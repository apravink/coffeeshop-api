FROM scratch
ADD main /
EXPOSE 8081
CMD ["/main"]
