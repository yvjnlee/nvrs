# Use an official Python image as a base
FROM python:3.9-slim

# Set the working directory inside the container
WORKDIR /app

# Copy the requirements file and install dependencies
COPY requirements.txt .
RUN pip install -r requirements.txt

# Copy the rest of the application files
COPY . .

# Expose the port your application listens on (if applicable)
EXPOSE 5000

# Command to run the Python application
CMD ["python", "agent.py"]
