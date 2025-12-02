# ai_api.py  ← Save as this name (one file only)

from fastapi import FastAPI
from fastapi.responses import StreamingResponse
from pydantic import BaseModel
import os
from groq import Groq
import asyncio

app = FastAPI(title="My Fast AI API")

# Put your free Groq API key here (get it instantly at https://console.groq.com/keys)
client = Groq(api_key="your-groq-api-key-here")  # ← CHANGE THIS!

class Message(BaseModel):
    role: str
    content: str

class ChatRequest(BaseModel):
    messages: list[Message]
    stream: bool = True  # change to False if you don't want streaming

@app.post("/chat")
async def chat(request: ChatRequest):
    # Real streaming from Groq (super fast: Llama 70B at 500+ tokens/sec)
    def stream_response():
        stream = client.chat.completions.create(
            model="llama-3.3-70b-versatile",  # or mixtral-8x7b-32768, gemma2-9b-it
            messages=[{"role": m.role, "content": m.content} for m in request.messages],
            temperature=0.7,
            max_tokens=2048,
            stream=True
        )
        for chunk in stream:
            if chunk.choices[0].delta.content:
                yield chunk.choices[0].delta.content

    if request.stream:
        return StreamingResponse(stream_response(), media_type="text/event-stream")
    else:
        # Non-streaming version
        response = client.chat.completions.create(
            model="llama-3.3-70b-versatile",
            messages=[{"role": m.role, "content": m.content} for m in request.messages],
            temperature=0.7,
        )
        return {"reply": response.choices[0].message.content}

@app.get("/")
def home():
    return {"message": "Your AI is alive! Send POST to /chat"}

# Run with: python ai_api.py
if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)