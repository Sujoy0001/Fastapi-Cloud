import httpx
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import os

app = FastAPI(title="Luhn Checksum Gateway")

# Get the Go service URL from an environment variable (set in docker-compose)
GO_SERVICE_URL = os.getenv("GO_SERVICE_URL", "http://checksum-engine:8080")

# Pydantic model for request body
class PartialNumberRequest(BaseModel):
    partial_number: str

# Pydantic model for Go's response
class GoChecksumResponse(BaseModel):
    checksum_digit: str

# Pydantic model for final API response
class FullNumberResponse(BaseModel):
    full_number: str

@app.post("/api/v1/checksum", response_model=FullNumberResponse)
async def get_full_number(request: PartialNumberRequest):
    """
    Receives 15 digits, delegates the calculation to the Go service, 
    and returns the full 16-digit card number.
    """
    partial_num = request.partial_number
    
    # 1. Validation (FastAPI's strength)
    if len(partial_num) != 15 or not partial_num.isdigit():
        raise HTTPException(
            status_code=400, 
            detail="Input must be exactly 15 numerical digits."
        )

    # 2. Call Go Service (Synchronous internal API call)
    try:
        async with httpx.AsyncClient() as client:
            go_response = await client.post(
                f"{GO_SERVICE_URL}/calculate-checksum",
                json={"partial_number": partial_num}
            )
            go_response.raise_for_status() # Raise exception for 4xx/5xx status codes
            
            checksum_data = GoChecksumResponse(**go_response.json())

    except httpx.HTTPStatusError as e:
        # Handle errors returned by the Go service
        raise HTTPException(
            status_code=500, 
            detail=f"Go service error: {e.response.json().get('error', 'Unknown error')}"
        )
    except httpx.RequestError:
        # Handle connection errors if Go service is down
        raise HTTPException(
            status_code=503, 
            detail="Could not connect to the Go Checksum Engine."
        )

    # 3. Aggregation and Response
    full_number = partial_num + checksum_data.checksum_digit
    
    return FullNumberResponse(full_number=full_number)