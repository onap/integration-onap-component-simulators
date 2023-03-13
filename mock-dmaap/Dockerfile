FROM python:3
COPY . /app
WORKDIR /app
COPY ./requirements.txt ./
RUN pip install -r ./requirements.txt
ENV FLASK_APP=app/app.py
CMD ["python", "app/app.py"]