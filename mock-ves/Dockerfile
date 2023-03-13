FROM python:3.8-alpine

COPY . /app
WORKDIR /app

# GCC for Alpine Linux in Docker
RUN apk add build-base


# Dependencies
RUN pip install -r requirements.txt

EXPOSE 30417

ENTRYPOINT ["python"]

CMD ["app/app.py"]