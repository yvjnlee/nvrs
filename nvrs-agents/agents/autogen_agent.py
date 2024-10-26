import os
import sys
from autogen import AssistantAgent, UserProxyAgent
import requests

grpc_client_path = os.path.abspath(os.path.join(os.path.dirname(__file__), '../grpc_client'))
sys.path.append(grpc_client_path)

# Import gRPC client methods
from grpc_client import update_status, submit_task
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# Configure the API key for GPT-4 model (if needed)
llm_config = {"model": "gpt-4", "api_key": os.getenv("OPENAI_API_KEY")}

class CustomAutoGenAgent:
    def __init__(self, name, role, base_url=None):
        self.name = name
        self.role = role
        self.base_url = base_url or "http://localhost:8080"
        self.agent_id = None
        self.token = None
        self.assistant = AssistantAgent("assistant", llm_config=llm_config)
        self.user_proxy = UserProxyAgent("user_proxy")

    def register(self):
        """Register the agent with the backend via REST API."""
        url = f"{self.base_url}/agents/register"
        headers = {"Content-Type": "application/json"}
        data = {"name": self.name, "role": self.role}
        response = requests.post(url, json=data, headers=headers)

        if response.status_code == 200:
            response_data = response.json()
            self.token = response_data.get("token")
            self.agent_id = response_data.get("agent_id")
            print(f"Agent '{self.name}' registered successfully. ID: {self.agent_id}")
        else:
            print("Failed to register agent.")

    def perform_task(self, task_description):
        """Submit a task using AutoGen gRPC client."""
        if self.agent_id:
            print(f"Submitting task: {task_description}")
            submit_task(self.agent_id, task_description)
        else:
            print("Agent is not registered.")

    def update_status(self, status):
        """Update agent status using AutoGen."""
        if self.agent_id:
            print(f"Updating status to: {status}")
            update_status(self.agent_id, status)
        else:
            print("Agent is not registered.")

    def start_conversation(self):
        """Start a conversation between agents using AutoGen."""
        self.user_proxy.initiate_chat(
            self.assistant,
            message="Analyze the given data and generate insights."
        )


# Testing the agent
if __name__ == "__main__":
    agent = CustomAutoGenAgent("AutoAgent001", "DataAnalyzer")
    agent.register()  # Register with the backend
    agent.update_status("active")  # Update status via gRPC
    agent.perform_task("Analyze financial data")
    agent.start_conversation()  # Start a conversation using AutoGen
