import os
import sys
import logging
from autogen.agentchat import ConversableAgent, UserProxyAgent
import requests
from dotenv import load_dotenv

# Define the path for grpc_client
grpc_client_path = os.path.abspath(os.path.join(os.path.dirname(__file__), '../grpc_client'))
sys.path.append(grpc_client_path)
from grpc_client import update_status, submit_task

# Load environment variables
load_dotenv()

# Set GPT-4 API key configuration
llm_config = {"model": "gpt-4", "api_key": os.getenv("OPENAI_API_KEY")}

# Initialize logging
logging.basicConfig(level=logging.INFO)

class CustomAgent(ConversableAgent):
    def __init__(self, name, role, base_url=None):
        super().__init__(name=name, description=role, llm_config=llm_config)
        self.role = role
        self.base_url = base_url or "http://localhost:8080"
        self.agent_id = None
        self.token = None

        # Initialize user proxy agent for human interaction
        self.user_proxy = UserProxyAgent(
            "user_proxy",
            llm_config=False,
            human_input_mode="ALWAYS",
            code_execution_config={"use_docker": False}
        )

    def register(self):
        """Register agent with backend using REST API."""
        try:
            url = f"{self.base_url}/agents/register"
            headers = {"Content-Type": "application/json"}
            data = {"name": self.name, "role": self.role}
            response = requests.post(url, json=data, headers=headers)
            response.raise_for_status()

            response_data = response.json()
            self.token = response_data.get("token")
            self.agent_id = response_data.get("agent_id")

            if not self.token or not self.agent_id:
                raise ValueError("Registration succeeded, but missing token or agent_id.")

            logging.info(f"Agent '{self.name}' registered successfully. ID: {self.agent_id}")

        except requests.RequestException as e:
            logging.error(f"Failed to register agent due to network error: {e}")
            raise
        except ValueError as ve:
            logging.error(f"Agent registration error: {ve}")
            raise

    def perform_task(self, task_description):
        """Submit a task using gRPC client."""
        if self.agent_id:
            logging.info(f"Submitting task: {task_description}")
            try:
                submit_task(self.agent_id, task_description)
            except Exception as e:
                logging.error(f"Task submission failed: {e}")
        else:
            logging.warning("Agent is not registered, task cannot be submitted.")

    def update_status(self, status):
        """Update agent status using gRPC client."""
        if self.agent_id:
            logging.info(f"Updating status to: {status}")
            try:
                update_status(self.agent_id, status)
            except Exception as e:
                logging.error(f"Status update failed: {e}")
        else:
            logging.warning("Agent is not registered, status cannot be updated.")

    def start_conversation(self, recipient=None):
        """Start conversation with user or another agent."""
        if recipient is None:
            logging.info("Conversation started. Type 'exit' to end.")
            while True:
                user_message = input("You: ")
                if user_message.lower() == 'exit':
                    logging.info("Ending conversation.")
                    break

                try:
                    result = self.user_proxy.initiate_chat(self, message=user_message)
                    if isinstance(result, list) and result and isinstance(result[0], str):
                        logging.info(f"{self.name}: {result[0]}")
                    else:
                        logging.error(f"Unexpected result format or empty response: {result}")
                        logging.info(f"{self.name}: (No valid response received)")
                except IndexError:
                    logging.error("Result from initiate_chat was empty.")
                except Exception as e:
                    logging.error(f"Unexpected error during conversation: {e}")
        else:
            initial_message = "Hello, let's start the task."
            self.send(initial_message, recipient)

# Testing the CustomAgent
if __name__ == "__main__":
    agent = CustomAgent("AutoAgent001", "DataAnalyzer")
    try:
        agent.register()
        agent.update_status("active")
        agent.perform_task("Analyze financial data")
        agent.start_conversation()
    except Exception as e:
        logging.error(f"An error occurred in main execution: {e}")
