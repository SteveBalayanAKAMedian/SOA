import grpc
import uuid
from google.protobuf.timestamp_pb2 import Timestamp
import task_pb2 as task_service_pb2
import task_pb2_grpc as task_service_pb2_grpc
from datetime import datetime

def create_task(stub, user_id, title, description):
    task = task_service_pb2.CreateTaskRequest(
        user_id=user_id,
        title=title,
        description=description
    )
    response = stub.CreateTask(task)
    print("CreateTask response:", response)
    return response.task

def update_task(stub, task_id, user_id, title, description, status):
    task = task_service_pb2.UpdateTaskRequest(
        id=task_id,
        user_id=user_id,
        title=title,
        description=description,
        status=status
    )
    response = stub.UpdateTask(task)
    print("UpdateTask response:", response)
    return response.task

def delete_task(stub, task_id, user_id):
    request = task_service_pb2.DeleteTaskRequest(id=task_id, user_id=user_id)
    response = stub.DeleteTask(request)
    print("DeleteTask response:", response)
    return response.success

def get_task(stub, task_id):
    request = task_service_pb2.GetTaskRequest(id=task_id)
    response = stub.GetTask(request)
    print("GetTask response:", response)
    return response.task

def list_tasks(stub, page, page_size):
    request = task_service_pb2.ListTasksRequest(page=page, page_size=page_size)
    response = stub.ListTasks(request)
    print("ListTasks response:", response)
    return response.tasks

def main():
    # Replace with your gRPC server address
    grpc_address = "localhost:50051"
    
    # Create a channel and a stub
    channel = grpc.insecure_channel(grpc_address)
    stub = task_service_pb2_grpc.TaskServiceStub(channel)
    
    # Test data
    user_id = str(uuid.uuid4())  # Generate a valid UUID for user_id
    title = "Test Task"
    description = "This is a test task"
    status = "pending"
    
    # Create a task
    create_response = create_task(stub, user_id, title, description)
    
    # Get the task ID from the response
    task_id = create_response.id
    
    # Update the task
    update_response = update_task(stub, task_id, user_id, "Updated Task", "Updated description", "completed")
    
    # Get the task by ID
    get_response = get_task(stub, task_id)
    
    # List tasks with pagination
    list_response = list_tasks(stub, page=1, page_size=10)
    
    # Delete the task
    delete_response = delete_task(stub, task_id, user_id)

if __name__ == "__main__":
    main()
