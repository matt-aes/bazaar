# small base image for Python containers
FROM python:3.7-alpine as base

# Load current directory and subdirectories
WORKDIR /app
COPY . /app

# Install required packages
RUN pip3 --no-cache-dir install -r requirements.txt

# Expose port 5002 and launch the application
EXPOSE 5002
ENTRYPOINT ["python3", "service.py"]
