#!/usr/bin/env python3
"""
Test scenarios for Tekton Debug Assistant

Usage:
    python test_debug_assistant.py                    # Run all test scenarios
    python test_debug_assistant.py --interactive      # Interactive mode
    python test_debug_assistant.py --single "error"   # Test single error
"""

import sys
from retrieval import TektonDebugAssistant

def run_test_scenarios():
    """Run predefined test scenarios"""
    # Initialize the debug assistant (production mode - clean output)
    assistant = TektonDebugAssistant(debug_mode=True)
    
    # Test knowledge base first
    if not assistant.test_knowledge_base_connection():
        print("❌ Cannot proceed without knowledge base")
        return False
    
    # Setup the system
    assistant.setup()
    
    # Test with common Tekton error scenarios
    test_prompts = [
        "My pipeline step failed with exit code 137, what does this mean?",
        
        "Getting error: imagePullBackOff on my build task",
        
        "Pod is stuck in Pending state with message: 0/3 nodes are available: 3 Insufficient cpu",
        
        "Pipeline fails with: couldn't find workspace 'source' declared in pipeline",
        
        "Build step fails with 'npm: command not found' - how do I fix this?",
        
        "Build step failed with exit code 1 in buildah task",
        
        # More complex scenario
        """I have a failing PipelineRun with this error:
        
        Step build-app in task build-task failed with exit code 1
        Log: /bin/sh: line 1: mvn: command not found
        
        My task definition uses image: alpine:latest
        What's wrong and how do I fix it?"""
    ]
    
    print("\n" + "="*60)
    print("🧪 TESTING TEKTON DEBUG ASSISTANT")
    print("="*60)
    
    for i, prompt in enumerate(test_prompts, 1):
        print(f"\n--- Test Case {i} ---")
        print(f"Query: {prompt[:60]}{'...' if len(prompt) > 60 else ''}")
        
        try:
            response = assistant.debug_pipeline_error(prompt)
            print(f"✅ Response generated successfully")
        except Exception as e:
            print(f"❌ Error: {e}")
        
        print("-" * 40)
        
        # Add a pause between tests for readability
        if i < len(test_prompts):
            user_input = input("\nPress Enter to continue (or 'q' to quit): ").strip()
            if user_input.lower() in ['q', 'quit']:
                break
    
    return True

def interactive_mode():
    """Interactive mode for debugging Tekton errors"""
    assistant = TektonDebugAssistant(debug_mode=False)
    
    # Test knowledge base first
    if not assistant.test_knowledge_base_connection():
        return
    
    # Setup the system
    assistant.setup()
    
    print("\n" + "="*60)
    print("🔧 TEKTON DEBUG ASSISTANT - Interactive Mode")
    print("="*60)
    print("Enter your Tekton error messages or 'quit' to exit")
    print("Example: 'Build step failed with exit code 1'")
    print("-" * 60)
    
    while True:
        try:
            user_input = input("\n🔍 Describe your error: ").strip()
            
            if user_input.lower() in ['quit', 'exit', 'q']:
                print("👋 Goodbye!")
                break
            
            if not user_input:
                print("Please enter an error description or 'quit' to exit")
                continue
            
            assistant.debug_pipeline_error(user_input)
            
        except KeyboardInterrupt:
            print("\n👋 Goodbye!")
            break
        except Exception as e:
            print(f"❌ Unexpected error: {e}")

def test_single_error(error_message: str):
    """Test a single error message"""
    assistant = TektonDebugAssistant(debug_mode=False)
    
    # Test knowledge base first
    if not assistant.test_knowledge_base_connection():
        return
    
    # Setup the system
    assistant.setup()
    
    print(f"\n🧪 Testing single error: {error_message}")
    print("="*60)
    
    try:
        response = assistant.debug_pipeline_error(error_message)
        print(f"✅ Response generated successfully")
        return response
    except Exception as e:
        print(f"❌ Error: {e}")
        return None

def main():
    """Main function with argument parsing"""
    if len(sys.argv) == 1:
        # Default: run test scenarios
        run_test_scenarios()
        
    elif len(sys.argv) == 2:
        if sys.argv[1] == "--interactive":
            interactive_mode()
        elif sys.argv[1] == "--help":
            print(__doc__)
        else:
            print("Usage: python test_debug_assistant.py [--interactive|--single 'error'|--help]")
            
    elif len(sys.argv) == 3 and sys.argv[1] == "--single":
        error_message = sys.argv[2]
        test_single_error(error_message)
        
    else:
        print("Usage: python test_debug_assistant.py [--interactive|--single 'error'|--help]")

if __name__ == "__main__":
    main()
