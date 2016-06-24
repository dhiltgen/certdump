FROM scratch
ADD certdump /bin/
ENTRYPOINT ["/bin/certdump"]
CMD ["--scan"]

