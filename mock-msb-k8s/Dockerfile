FROM python:3.8-alpine

COPY . .

WORKDIR /app

RUN pip3 install -r ../requirements.txt

EXPOSE 5003

CMD ["python","app.py"]