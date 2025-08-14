### Start Ollama locally (Llama inference)

```bash
# Ensure Ollama is running, then pull/run a model
ollama ps
ollama run llama3.2:3b-instruct-fp16 --keepalive 60m
# Exit the interactive session when ready
/bye
```

### Start LlamaStack with Ollama backend (uses Llama locally)

```bash
# Example: expose on 8321 and point to local Ollama
export LLAMA_STACK_PORT=8321
export LLAMA_STACK_MODEL=llama3.2:3b-instruct-fp16

docker run -it --rm --network host \
  llamastack/distribution-ollama:0.2.0 \
  --port $LLAMA_STACK_PORT \
  --env INFERENCE_MODEL=$LLAMA_STACK_MODEL \
  --env OLLAMA_URL=http://localhost:11434
```

### Start LlamaStack with Gemini backend

```bash
# Requires GEMINI_API_KEY. Uses config at hack/gemini.yaml
docker run -it --rm \
  -v ./hack/gemini.yaml:/app/gemini.yaml:z \
  -v ${SQLITE_STORE_DIR:-~/.llama/distributions/gemini}:/data \
  -e GEMINI_API_KEY=$GEMINI_API_KEY \
  -e SQLITE_STORE_DIR=/data \
  -p 8321:8321 \
  docker.io/llamastack/distribution-starter \
  --config /app/gemini.yaml
```

### Use llama-stack-client for inference

```bash
llama-stack-client models list

llama-stack-client inference chat-completion \
  --message "hello, what model are you?" \
  --model-id "meta-llama/Llama-3.2-3B-Instruct"
```

Notes:
- Both examples expose LlamaStack on port 8321. Adjust if you change `LLAMA_STACK_PORT`.
- For the Ollama backend, ensure the model in `LLAMA_STACK_MODEL` is available in your local Ollama.

### Usage (local scripts)

```bash
# Ingest knowledge base into LlamaStack vector DB
python ingester.py

# Search, stats, validate, reset examples
python ingester.py --search "registry push failed" --max-results 3
python ingester.py --stats
python ingester.py --validate
python ingester.py --reset --force

# Run debug assistant test scenarios
python test_debug_assistant.py

# Interactive mode
python test_debug_assistant.py --interactive

# Single error test
python test_debug_assistant.py --single "Build step failed with exit code 1"
```
