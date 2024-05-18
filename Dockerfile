FROM python:3.11-slim

WORKDIR /usr/src/app

COPY backend/requirements.txt /usr/src/app
# RUN apt-get update && apt-get install -y %LIBRARIES% && rm -rf /var/lib/apt/lists/*
RUN pip install --no-cache-dir -r /usr/src/app/requirements.txt && rm /usr/src/app/requirements.txt

COPY backend/*.py /usr/src/app/
COPY frontend/build /usr/src/app/static

CMD [ "python", "/usr/src/app/main.py" ]
