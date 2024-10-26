import grpc
import agent_pb2
import agent_pb2_grpc


def update_status(agent_id, status):
    # Connect to the gRPC server at localhost:50051
    with grpc.insecure_channel('localhost:50051') as channel:
        # Create a client stub to interact with the AgentService
        stub = agent_pb2_grpc.AgentServiceStub(channel)
        
        # Create and send the UpdateStatus request
        request = agent_pb2.StatusRequest(agent_id=agent_id, status=status)
        response = stub.UpdateStatus(request)
        print(f"UpdateStatus Response: {response.message}")


def submit_task(agent_id, task):
    # Connect to the gRPC server at localhost:50051
    with grpc.insecure_channel('localhost:50051') as channel:
        # Create a client stub to interact with the AgentService
        stub = agent_pb2_grpc.AgentServiceStub(channel)
        
        # Create and send the SubmitTask request
        request = agent_pb2.TaskRequest(agent_id=agent_id, task=task)
        response = stub.SubmitTask(request)
        print(f"SubmitTask Response: {response.message}")


# Example usage
if __name__ == "__main__":
    # Test UpdateStatus
    update_status(agent_id=1, status="active")

    # Test SubmitTask
    submit_task(agent_id=1, task="Analyze data")
