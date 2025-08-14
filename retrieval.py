#!/usr/bin/env python3
"""
Tekton Debug Assistant - Core RAG Agent

This module provides the core TektonDebugAssistant class for debugging Tekton pipeline errors
using LlamaStack RAG capabilities.

Prerequisites:
1. Run `python ingester.py --reset && python ingester.py` first to setup the knowledge base
2. Ensure LlamaStack is running on http://localhost:8321
3. Optional: Set INFERENCE_MODEL environment variable (default: gemini-1.5-flash)

Note: This implementation uses proper session management for conversation continuity.

Usage as module:
    from retrieval import TektonDebugAssistant
    
    assistant = TektonDebugAssistant()
    assistant.setup()
    response = assistant.debug_pipeline_error("Build failed with exit code 1")

Usage as script:
    python retrieval.py --test   # Test knowledge base connection
    python retrieval.py --daemon # Run as background service (future)
"""

import os
import json
import uuid
from typing import List, Dict, Any
from llama_stack_client import LlamaStackClient
from llama_stack_client.lib.agents.agent import Agent
from llama_stack_client.lib.agents.event_logger import EventLogger as AgentEventLogger

class TektonDebugAssistant:
    def __init__(self, llamastack_endpoint: str = "http://localhost:8321", 
                 vector_db_id: str = "tekton_errors_db",
                 debug_mode: bool = False):
        self.client = LlamaStackClient(base_url=llamastack_endpoint)
        self.vector_db_id = vector_db_id
        self.debug_mode = debug_mode
        self.rag_agent = None
        self.session_id = None
        print(f"🔧 Initialized Tekton Debug Assistant")
        print(f"   Endpoint: {llamastack_endpoint}")
        print(f"   Vector DB: {vector_db_id}")
        print("📋 Note: Assuming vector database is already setup via ingester.py")
         
    def test_knowledge_base_connection(self) -> bool:
        """Quick test to verify knowledge base is accessible"""
        try:
            # Try to list vector databases to see if our DB exists
            vector_dbs = self.client.vector_dbs.list()
            
            # Check if our vector DB exists
            db_exists = any(db.identifier == self.vector_db_id for db in vector_dbs)
            
            if db_exists:
                print(f"✅ Vector database '{self.vector_db_id}' is ready")
                return True
            else:
                print(f"❌ Vector database '{self.vector_db_id}' not found")
                print("🔧 Please run: python ingester.py --reset && python ingester.py")
                return False
                
        except Exception as e:
            print(f"❌ Cannot connect to vector database: {e}")
            print("🔧 Make sure LlamaStack is running and knowledge base is setup")
            return False
    
    def create_rag_agent(self, vector_store_id: str):
        """Create RAG agent specialized for Tekton debugging"""
        
        system_prompt = """You are a Tekton pipeline debugging expert with access to a comprehensive knowledge base of common Tekton issues and solutions.

IMPORTANT: You have access to a knowledge_search tool that searches the 'tekton_errors_db' vector database. Always use this tool to search the knowledge base for relevant information before providing solutions.

Your process:
1. **Search First**: Use the knowledge_search tool to find relevant documented solutions in the tekton_errors_db knowledge base
2. **Analyze**: Based on search results, explain what the error means
3. **Provide Solutions**: Give clear, actionable debugging steps from the knowledge base
4. **Include Examples**: Show specific YAML configurations and commands

When a user asks about a Tekton error:
- ALWAYS start by searching the knowledge base using the knowledge_search tool
- Use the exact information, references, and solutions found in the search results
- If the search finds relevant entries, base your entire response on that information
- Include the exact reference URLs found in the knowledge base

Response Structure:
1. **Problem Analysis**: Explain the error based on knowledge base findings
2. **Root Cause**: Identify the cause using knowledge base context
3. **Solution Steps**: Provide the exact solution from knowledge base, with details
4. **Code Examples**: Include YAML configurations from knowledge base or based on its guidance
5. **Verification**: Steps to confirm the fix worked
6. **References**: Use the exact reference URLs from knowledge base search results

Guidelines:
- ALWAYS use the knowledge_search tool for every user query
- Base responses primarily on knowledge base search results
- Use exact reference URLs from the knowledge base (e.g., https://tekton.dev/docs/memory-limits)
- Include error types, severity levels, and metadata from search results
- If search returns results, never say "no knowledge base entry found"

Keep responses practical and actionable, prioritizing knowledge base information."""

        self.rag_agent = Agent(
            self.client,
            model=os.environ.get("INFERENCE_MODEL", "gemini-1.5-flash"),
            instructions=system_prompt,
            enable_session_persistence=True,
            tools=[
                {
                    "name": "builtin::rag",
                    "args": {
                        "vector_db_ids": [vector_store_id],
                    },
                }
            ],
        )
        
        print("🤖 RAG agent created successfully")
        return self.rag_agent
    
    def start_debug_session(self):
        """Start a new debugging session"""
        try:
            # Create session with a unique name
            session_name = f"tekton-debug-{uuid.uuid4().hex[:8]}"
            self.session_id = self.rag_agent.create_session(session_name)
            print(f"🔧 Started debug session: {self.session_id}")
            return self.session_id
        except Exception as e:
            print(f"❌ Session creation failed: {e}")
            raise
    
    def debug_pipeline_error(self, error_message: str) -> str:
        """Debug a specific pipeline error"""
        if not self.rag_agent or not self.session_id:
            raise ValueError("Agent not initialized. Call setup() first.")
        
        print(f"\n🔍 User> {error_message}")
        
        try:
            # Use the agent with proper session management
            response = self.rag_agent.create_turn(
                messages=[{"role": "user", "content": error_message}],
                session_id=self.session_id,
            )
            
            # Log the interaction using AgentEventLogger
            full_response = ""
            rag_tool_used = False
            rag_results_count = 0
            
            for log in AgentEventLogger().log(response):
                # In debug mode, show everything
                if self.debug_mode:
                    log.print()
                    
                    # Enhanced RAG monitoring for debug mode
                    log_str = str(log)
                    
                    # Check for tool calls
                    if hasattr(log, 'tool_call'):
                        tool_call_str = str(log.tool_call)
                        if 'knowledge_search' in tool_call_str or 'rag' in tool_call_str.lower():
                            rag_tool_used = True
                            print(f"🔍 RAG TOOL CALLED: {log.tool_call}")
                    
                    # Check for tool responses
                    if hasattr(log, 'tool_response'):
                        tool_response_str = str(log.tool_response)
                        if any(term in tool_response_str.lower() for term in ['vector', 'knowledge', 'search', 'tekton']):
                            print(f"📋 RAG RESPONSE: {log.tool_response}")
                            if 'results' in tool_response_str or 'documents' in tool_response_str:
                                rag_results_count += 1
                    
                    # Check log content for RAG indicators
                    if 'knowledge' in log_str.lower() or 'vector' in log_str.lower():
                        print(f"📊 RAG INDICATOR: {log_str[:100]}...")
                else:
                    # In production mode, only show the final content
                    if hasattr(log, 'content') and log.content:
                        print(log.content, end='')
                
                # Always collect the full response for return
                if hasattr(log, 'content'):
                    full_response += str(log.content)
            
            # Summary of RAG usage (only in debug mode)
            if self.debug_mode:
                print(f"\n📊 RAG USAGE SUMMARY:")
                print(f"   RAG Tool Called: {'✅ YES' if rag_tool_used else '❌ NO'}")
                print(f"   RAG Results: {rag_results_count}")
                if not rag_tool_used:
                    print("   ⚠️  Agent may be using general knowledge instead of KB")
                print("-" * 40)
            
            return full_response
            
        except Exception as e:
            print(f"❌ Agent interaction failed: {e}")
            # Fallback: return simple error response
            return f"Error processing request: {error_message}. Details: {str(e)}"
    
    def setup(self):
        """Complete setup: agent creation (assumes vector DB already exists)"""
        self.create_rag_agent(self.vector_db_id)
        self.start_debug_session()
        print("🚀 Tekton Debug Assistant ready!")

def test_connection():
    """Quick test mode to verify knowledge base connectivity"""
    print("🧪 Testing Knowledge Base Connection...")
    assistant = TektonDebugAssistant()
    
    if assistant.test_knowledge_base_connection():
        print("\n✅ Knowledge base test passed!")
        print("🚀 Ready to run debug assistant")
        print("💡 Use test_debug_assistant.py for full testing scenarios")
        print("💡 Set INFERENCE_MODEL environment variable if needed (default: gemini-1.5-flash)")
    else:
        print("\n❌ Knowledge base test failed!")
        print("📋 Setup required before using retrieval assistant")

if __name__ == "__main__":
    import sys
    
    if len(sys.argv) > 1 and sys.argv[1] == "--test":
        test_connection()
    elif len(sys.argv) > 1 and sys.argv[1] == "--daemon":
        print("🚀 Daemon mode not yet implemented")
        print("💡 Use test_debug_assistant.py for testing scenarios")
    else:
        print("Tekton Debug Assistant - Core Module")
        print("Usage:")
        print("  python retrieval.py --test              # Test connection")
        print("  python test_debug_assistant.py          # Run test scenarios")
        print("  python test_debug_assistant.py --interactive  # Interactive mode")
        print("  from retrieval import TektonDebugAssistant    # Use as module")
