import os
import sys
import requests

# Add the absolute path of the grpc_client directory to sys.path
grpc_client_path = os.path.abspath(os.path.join(os.path.dirname(__file__), '../grpc_client'))
sys.path.append(grpc_client_path)

# Import gRPC client methods
from grpc_client import update_status, submit_task

class Agent:
    def __init__(self, name, role, base_url=None):
        self.name = name
        self.role = role
        # Use an environment variable for base_url; fallback to localhost for local testing
        self.base_url = base_url or os.getenv("GATEWAY_HOST", "http://localhost:8080")
        self.token = None        # JWT token, assigned after registration
        self.agent_id = None     # Agent ID, assigned after registration

    def register(self):
        """Register the agent via REST API to get an ID and token."""
        url = f"{self.base_url}/agents/register"
        headers = {"Content-Type": "application/json"}
        data = {"name": self.name, "role": self.role}
        response = requests.post(url, json=data, headers=headers)
        
        if response.status_code == 200:
            response_data = response.json()
            self.token = response_data.get("token")
            self.agent_id = response_data.get("agent_id")  # Store agent ID

            if self.token and self.agent_id:
                print(f"Agent '{self.name}' registered successfully. ID: {self.agent_id}, Token: {self.token}")
            else:
                print("Registration successful, but token or agent ID missing in response.")
        else:
            print("Registration failed:", response.json())

    def submit_task(self, task_description):
        """Submit a task using gRPC."""
        if self.agent_id and self.token:
            print(f"Submitting task: {task_description} via gRPC...")
            submit_task(self.agent_id, task_description)
        else:
            print("Agent not registered. Cannot submit task.")

    def update_status(self, status):
        """Update the agent's status using gRPC."""
        if self.agent_id and self.token:
            print(f"Updating status to: {status} via gRPC...")
            update_status(self.agent_id, status)
        else:
            print("Agent not registered. Cannot update status.")


# Testing
if __name__ == "__main__":
    # No hardcoding; use an environment variable for the base URL
    base_url = os.getenv("GATEWAY_HOST", "http://localhost:8080")
    agent = Agent("Agent007", "Spy", base_url)
    
    agent.register()            # Register via REST
    agent.submit_task("Analyze mission data")  # Submit task via gRPC
    agent.update_status("working")  # Update status via gRPC
