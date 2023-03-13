FROM python:3.7.8-alpine

COPY . /app
WORKDIR /app

# GCC for Alpine Linux in Docker
RUN apk add build-base


# Dependencies
RUN pip install pipenv && \
    pipenv requirements > requirements.txt && \
    pip install -r requirements.txt


ENTRYPOINT ["python"]
CMD ["app/app.py"]