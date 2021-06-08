FROM python:3.8-slim
RUN apt-get -y update && apt-get -y install git && apt-get clean && \
    pip install  --disable-pip-version-check --no-cache-dir -U wheel pip && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /usr/share/man/?? /usr/share/man/??_*

COPY requirements.txt /requirements.txt
RUN pip install --disable-pip-version-check --no-cache-dir dockerfile-sec

ENTRYPOINT ["dockerfile-sec"]
