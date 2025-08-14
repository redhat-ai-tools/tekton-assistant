ollama ps
ollama run llama3.2:3b-instruct-fp16 --keepalive 60m
/bye
start llamastack: 
docker run -it --network host llamastack/distribution-ollama:0.2.0 --port $LLAMA_STACK_PORT --env INFERENCE_MODEL=$LLAMA_STACK_MODEL --env OLLAMA_URL=http://localhost:11434

llama stack client: 
llama-stack-client inference chat-completion --message "hello, what model are you?" --model-id "meta-llama/Llama-3.2-3B-Instruct"


------
using gemini
docker run -it --rm \
  -v ./gemini.yaml:/app/gemini.yaml:z \
  -v ${SQLITE_STORE_DIR:-~/.llama/distributions/gemini}:/data \
  -e GEMINI_API_KEY=$GEMINI_API_KEY \
  -e SQLITE_STORE_DIR=/data \
  -p 8321:8321 \
  docker.io/llamastack/distribution-starter \
  --config /app/gemini.yaml

llama-stack-client models list

llama-stack-client inference chat-completion --message "hello, what model are you"

