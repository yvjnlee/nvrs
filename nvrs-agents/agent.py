import requests

class Agent:
    def __init__(self, name, role, base_url):
        self.name = name
        self.role = role
        self.base_url = base_url
        self.token = None        # JWT token, assigned after registration
        self.agent_id = None      # Agent ID, assigned after registration

    def register(self):
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
        url = f"{self.base_url}/api/tasks"
        headers = {
            "Authorization": f"Bearer {self.token}",
            "Content-Type": "application/json"
        }
        data = {"agent_id": self.agent_id, "task": task_description}  # Use agent_id
        response = requests.post(url, json=data, headers=headers)
        
        if response.status_code == 200:
            print("Task submitted successfully")
        else:
            print("Task submission failed:", response.json())

    def update_status(self, status):
        url = f"{self.base_url}/api/agents/status"
        headers = {
            "Authorization": f"Bearer {self.token}",
            "Content-Type": "application/json"
        }
        data = {"agent_id": self.agent_id, "status": status}  # Use agent_id
        response = requests.post(url, json=data, headers=headers)
        
        if response.status_code == 200:
            print("Status updated successfully")
        else:
            print("Status update failed:", response.json())


# Testing
if __name__ == "__main__":
    base_url = "http://localhost:8080"  # Replace with your server URL
    agent = Agent("Agent007", "Spy", base_url)
    
    agent.register()
    agent.submit_task("Analyze mission data")
    agent.update_status("working")
