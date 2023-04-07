FROM python:3.11-bullseye
WORKDIR /app
RUN apt-get update && apt-get install -y curl python3-dev libpq-dev libpq5
COPY . .
RUN pip install -r requirements.txt

CMD ["python3", "main.py"]
