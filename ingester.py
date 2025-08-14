#!/usr/bin/env python3
"""
Tekton Knowledge Base Vector Database Ingester

This script ingests TaskRun failure knowledge base entries from kb.json
into LlamaStack vector database for semantic search and retrieval.
"""

import json
import os
import time
import logging
from pathlib import Path
from typing import List, Dict, Any, Optional
from dataclasses import dataclass, asdict
from datetime import datetime
import uuid
from termcolor import cprint

try:
    from llama_stack_client import LlamaStackClient
except ImportError as e:
    print(f"Missing required dependencies. Please install with:")
    print("pip install llama-stack-client termcolor")
    exit(1)


@dataclass
class IngestionMetrics:
    """Metrics for tracking ingestion performance."""
    total_entries: int = 0
    successful_ingestions: int = 0
    failed_ingestions: int = 0
    start_time: float = 0
    end_time: float = 0
    processing_time: float = 0
    error_types: Dict[str, int] = None
    severity_counts: Dict[str, int] = None
    
    def __post_init__(self):
        if self.error_types is None:
            self.error_types = {}
        if self.severity_counts is None:
            self.severity_counts = {}
    
    @property
    def success_rate(self) -> float:
        """Calculate success rate percentage."""
        if self.total_entries == 0:
            return 0.0
        return (self.successful_ingestions / self.total_entries) * 100
    
    @property
    def entries_per_second(self) -> float:
        """Calculate processing speed."""
        if self.processing_time == 0:
            return 0.0
        return self.successful_ingestions / self.processing_time
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert metrics to dictionary for logging/reporting."""
        return asdict(self)


@dataclass
class KnowledgeEntry:
    """Structured representation of a knowledge base entry."""
    id: str
    error: str
    context: str
    solution: str
    reference: str
    source: str
    metadata: Dict[str, Any]
    combined_text: str
    error_type: str
    severity: str


class TektonKnowledgeIngester:
    """
    Tekton Knowledge Base Ingester using LlamaStack.
    
    Handles ingestion of structured error-solution pairs into vector database
    with comprehensive logging, metrics, and error handling.
    """
    
    def __init__(self, 
                 kb_file: str = "kb.json",
                 endpoint_url: str = "http://localhost:8321",
                 vector_db_id: str = "tekton_errors_db",
                 embedding_model: str = "text-embedding-004",
                 embedding_dimension: int = 384,
                 chunk_size: int = 512,
                 batch_size: int = 50,
                 log_level: str = "INFO"):
        """Initialize the Tekton Knowledge Ingester."""
        self.kb_file = Path(kb_file)
        self.endpoint_url = endpoint_url
        self.vector_db_id = vector_db_id
        self.embedding_model = embedding_model
        self.embedding_dimension = embedding_dimension
        self.chunk_size = chunk_size
        self.batch_size = batch_size
        
        # Setup logging
        self._setup_logging(log_level)
        
        # Initialize metrics
        self.metrics = IngestionMetrics()
        
        # Initialize LlamaStack client
        self._initialize_client()
        
        # Setup vector database
        self._setup_vector_db()
    
    def _setup_logging(self, log_level: str) -> None:
        """Setup structured logging with colors."""
        logging.basicConfig(
            level=getattr(logging, log_level.upper()),
            format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
            handlers=[
                logging.StreamHandler(),
                logging.FileHandler('tekton_ingester.log')
            ]
        )
        self.logger = logging.getLogger(__name__)
        self.logger.info("Logging initialized")
    
    def _initialize_client(self) -> None:
        """Initialize LlamaStack client with error handling."""
        try:
            self.client = LlamaStackClient(base_url=self.endpoint_url)
            self.logger.info(f"LlamaStack client initialized: {self.endpoint_url}")
            cprint(f"✅ Connected to LlamaStack: {self.endpoint_url}", "green")
        except Exception as e:
            self.logger.error(f"Failed to initialize LlamaStack client: {e}")
            cprint(f"❌ Failed to connect to LlamaStack: {e}", "red")
            raise
    
    def vector_db_exists(self) -> bool:
        """Check if vector database exists."""
        try:
            vector_dbs = self.client.vector_dbs.list()
            exists = any(db.identifier == self.vector_db_id for db in vector_dbs)
            self.logger.info(f"Vector DB {self.vector_db_id} exists: {exists}")
            return exists
        except Exception as e:
            self.logger.warning(f"Could not check vector DB existence: {e}")
            return False
    
    def _setup_vector_db(self) -> None:
        """Setup vector database with proper configuration."""
        try:
            if self.vector_db_exists():
                cprint(f"📊 Using existing vector DB: {self.vector_db_id}", "blue")
                self.logger.info(f"Using existing vector database: {self.vector_db_id}")
            else:
                cprint(f"🔧 Creating new vector DB: {self.vector_db_id}", "yellow")
                self.logger.info(f"Creating new vector database: {self.vector_db_id}")
                
                self.client.vector_dbs.register(
                    vector_db_id=self.vector_db_id,
                    embedding_model=self.embedding_model,
                    embedding_dimension=self.embedding_dimension,
                    provider_id=os.environ.get('VECTOR_STORE_PROVIDER', 'faiss')
                )
                cprint(f"✅ Vector DB created successfully", "green")
                self.logger.info("Vector database created successfully")
                
        except Exception as e:
            if "already exists" in str(e).lower():
                self.logger.info("Vector database already exists, continuing...")
                cprint(f"📊 Vector DB already exists: {self.vector_db_id}", "blue")
            else:
                self.logger.error(f"Failed to setup vector database: {e}")
                cprint(f"❌ Vector DB setup failed: {e}", "red")
                raise
    
    def infer_error_type(self, error: str) -> str:
        """Categorize error types based on error message patterns."""
        error_lower = error.lower()
        
        # OOM and resource issues
        if any(keyword in error_lower for keyword in ["oom", "exit code 137", "memory"]):
            return "OOM"
        # Registry and image issues  
        elif any(keyword in error_lower for keyword in ["registry", "push", "pull", "image"]):
            return "Registry_Push_Failure"
        # Build failures
        elif any(keyword in error_lower for keyword in ["build", "compile", "make"]):
            return "Build_Failure"
        # Version compatibility
        elif any(keyword in error_lower for keyword in ["version", "compatibility", "requires"]):
            return "Version_Compatibility"
        # Permission issues
        elif any(keyword in error_lower for keyword in ["permission", "denied", "unauthorized"]):
            return "Permission_Denied"
        # Network and connectivity
        elif any(keyword in error_lower for keyword in ["network", "timeout", "connection"]):
            return "Network_Issue"
        # Missing files/resources
        elif any(keyword in error_lower for keyword in ["missing", "not found", "no such file"]):
            return "Missing_Files"
        
        return "Other"
    
    def infer_severity(self, error: str, error_type: str) -> str:
        """Determine error severity based on impact and error type."""
        error_lower = error.lower()
        
        # Critical errors that block execution
        if any(keyword in error_lower for keyword in ["oom", "crashloopbackoff", "panic"]):
            return "critical"
        # High severity for build/registry failures
        elif error_type in ["Build_Failure", "Registry_Push_Failure", "Version_Compatibility"]:
            return "high"
        # Medium for permission and missing file issues
        elif error_type in ["Permission_Denied", "Missing_Files"]:
            return "medium"
        # Low for warnings and non-blocking issues
        elif any(keyword in error_lower for keyword in ["warning", "skipping"]):
            return "low"
        
        return "medium"  # Default
    
    def load_knowledge_base(self) -> List[KnowledgeEntry]:
        """Load and parse the knowledge base JSON file."""
        if not self.kb_file.exists():
            raise FileNotFoundError(f"Knowledge base file not found: {self.kb_file}")
        
        self.logger.info(f"Loading knowledge base from: {self.kb_file}")
        cprint(f"📖 Loading knowledge base: {self.kb_file}", "blue")
        
        with open(self.kb_file, 'r', encoding='utf-8') as f:
            kb_data = json.load(f)
        
        entries = []
        self.metrics.total_entries = len(kb_data)
        
        for idx, entry in enumerate(kb_data):
            try:
                # Create combined searchable text
                combined_text = self._create_combined_text(entry)
                
                # Generate unique ID
                entry_id = self._generate_entry_id(entry, idx)
                
                # Infer error classification
                error_type = self.infer_error_type(entry.get("error", ""))
                severity = self.infer_severity(entry.get("error", ""), error_type)
                
                # Update metrics
                self.metrics.error_types[error_type] = self.metrics.error_types.get(error_type, 0) + 1
                self.metrics.severity_counts[severity] = self.metrics.severity_counts.get(severity, 0) + 1
                
                kb_entry = KnowledgeEntry(
                    id=entry_id,
                    error=entry.get("error", ""),
                    context=entry.get("context", ""),
                    solution=entry.get("solution", ""),
                    reference=entry.get("reference", ""),
                    source=entry.get("source", ""),
                    metadata=entry.get("metadata", {}),
                    combined_text=combined_text,
                    error_type=error_type,
                    severity=severity
                )
                entries.append(kb_entry)
                
            except Exception as e:
                self.logger.error(f"Error processing entry {idx}: {e}")
                self.metrics.failed_ingestions += 1
        
        self.logger.info(f"Loaded {len(entries)} knowledge base entries")
        cprint(f"✅ Loaded {len(entries)} entries from knowledge base", "green")
        
        # Cache entries for fallback search
        self._cached_entries = entries
        
        return entries
    
    def _create_combined_text(self, entry: Dict[str, Any]) -> str:
        """Create combined searchable text optimized for LLM processing."""
        parts = []
        
        # Main content with clear structure
        if entry.get("error"):
            parts.append(f"ERROR: {entry['error']}")
        if entry.get("context"):
            parts.append(f"CONTEXT: {entry['context']}")
        if entry.get("solution"):
            parts.append(f"SOLUTION: {entry['solution']}")
        
        # Metadata for better search
        metadata = entry.get("metadata", {})
        if metadata.get("task_name"):
            parts.append(f"TEKTON_TASK: {metadata['task_name']}")
        if metadata.get("failed_step"):
            parts.append(f"FAILED_STEP: {metadata['failed_step']}")
        if metadata.get("repository"):
            parts.append(f"REPOSITORY: {metadata['repository']}")
        if metadata.get("related_errors"):
            parts.append(f"RELATED_ERRORS: {', '.join(metadata['related_errors'])}")
        
        return "\n".join(parts)
    
    def _generate_entry_id(self, entry: Dict[str, Any], idx: int) -> str:
        """Generate a unique, deterministic ID for the entry."""
        metadata = entry.get("metadata", {})
        
        # Use error type + timestamp if available
        if metadata.get("error_type") and metadata.get("timestamp"):
            error_type = metadata["error_type"].lower().replace("_", "-")
            timestamp = metadata["timestamp"][:10]  # YYYY-MM-DD
            return f"{error_type}-{timestamp}-{idx:03d}"
        
        # Use hash of error message for uniqueness
        error_hash = str(hash(entry.get("error", "")))[-6:]
        return f"kb-entry-{idx:03d}-{error_hash}"
    
    def ingest_entries(self, entries: List[KnowledgeEntry]) -> IngestionMetrics:
        """Ingest knowledge base entries into vector database."""
        if not entries:
            cprint("⚠️  No entries to ingest", "yellow")
            return self.metrics
        
        self.metrics.start_time = time.time()
        
        cprint(f"🚀 Starting ingestion of {len(entries)} entries...", "blue")
        self.logger.info(f"Starting ingestion of {len(entries)} entries")
        
        try:
            # Process in batches for better performance
            for i in range(0, len(entries), self.batch_size):
                batch = entries[i:i + self.batch_size]
                batch_num = i // self.batch_size + 1
                total_batches = (len(entries) + self.batch_size - 1) // self.batch_size
                
                self.logger.info(f"Processing batch {batch_num}/{total_batches} ({len(batch)} entries)")
                cprint(f"📦 Processing batch {batch_num}/{total_batches}", "cyan")
                
                try:
                    # Prepare documents for LlamaStack
                    documents = []
                    for entry in batch:
                        doc_metadata = {
                            "error": entry.error[:200],  # Truncate for metadata
                            "error_type": entry.error_type,
                            "severity": entry.severity,
                            "source": entry.source,
                            "reference": entry.reference,
                            "task_name": entry.metadata.get("task_name", ""),
                            "repository": entry.metadata.get("repository", ""),
                            "timestamp": int(time.time())
                        }
                        
                        documents.append({
                            "document_id": entry.id,
                            "content": entry.combined_text,
                            "metadata": doc_metadata
                        })
                    
                    # Insert batch into vector database
                    try:
                        # Try RAG tool insert
                        self.client.tool_runtime.rag_tool.insert(
                            documents=documents,
                            vector_db_id=self.vector_db_id,
                            chunk_size_in_tokens=self.chunk_size
                        )
                    except Exception as rag_error:
                        self.logger.warning(f"RAG tool insert failed: {rag_error}, trying vector_dbs.insert")
                        
                        # Fallback: Try direct vector database insert
                        self.client.vector_dbs.insert(
                            vector_db_id=self.vector_db_id,
                            documents=documents
                        )
                    
                    self.metrics.successful_ingestions += len(batch)
                    cprint(f"✅ Batch {batch_num} ingested successfully", "green")
                    
                except Exception as e:
                    self.logger.error(f"Failed to ingest batch {batch_num}: {e}")
                    cprint(f"❌ Batch {batch_num} failed: {e}", "red")
                    self.metrics.failed_ingestions += len(batch)
            
            self.metrics.end_time = time.time()
            self.metrics.processing_time = self.metrics.end_time - self.metrics.start_time
            
            # Log final metrics
            self._log_metrics()
            
            return self.metrics
            
        except Exception as e:
            self.logger.error(f"Ingestion failed: {e}")
            cprint(f"💥 Ingestion failed: {e}", "red")
            raise
    
    def _log_metrics(self) -> None:
        """Log comprehensive ingestion metrics."""
        metrics = self.metrics
        
        cprint("\n📊 Ingestion Metrics", "blue", attrs=['bold'])
        cprint("=" * 50, "blue")
        
        # Basic metrics
        print(f"Total entries:        {metrics.total_entries}")
        print(f"Successful:          {metrics.successful_ingestions}")
        print(f"Failed:              {metrics.failed_ingestions}")
        print(f"Success rate:        {metrics.success_rate:.1f}%")
        print(f"Processing time:     {metrics.processing_time:.2f}s")
        print(f"Speed:               {metrics.entries_per_second:.1f} entries/sec")
        
        # Error type distribution
        if metrics.error_types:
            cprint("\n🏷️  Error Types:", "cyan")
            for error_type, count in metrics.error_types.items():
                print(f"  {error_type}: {count}")
        
        # Severity distribution  
        if metrics.severity_counts:
            cprint("\n⚠️  Severity Levels:", "yellow")
            for severity, count in metrics.severity_counts.items():
                print(f"  {severity}: {count}")
        
        self.logger.info(f"Ingestion completed: {metrics.to_dict()}")
    
    def search(self, query: str, max_results: int = 5) -> Dict[str, Any]:
        """Search the knowledge base for relevant entries."""
        try:
            cprint(f"🔍 Searching: '{query}'", "blue")
            self.logger.info(f"Searching for: '{query}'")
            
            # Query vector database using the correct LlamaStack API
            try:
                # Try the RAG tool query method
                results = self.client.tool_runtime.rag_tool.query(
                    vector_db_ids=[self.vector_db_id],  # Use vector_db_ids (plural)
                    query=query,  # Back to query parameter
                    max_chunks=max_results
                )
                
                # Format results
                formatted_results = {
                    "query": query,
                    "results": {
                        "documents": [chunk.content for chunk in results.chunks],
                        "metadatas": [chunk.metadata for chunk in results.chunks],
                        "scores": [chunk.score for chunk in results.chunks]
                    },
                    "count": len(results.chunks)
                }
                
            except Exception as api_error:
                self.logger.warning(f"RAG tool query failed: {api_error}, trying vector_dbs.query")
                
                # Fallback: Try direct vector database query
                try:
                    results = self.client.vector_dbs.query(
                        vector_db_id=self.vector_db_id,
                        query_text=query,
                        k=max_results
                    )
                    
                    # Format results for vector_dbs API
                    formatted_results = {
                        "query": query,
                        "results": {
                            "documents": [doc.content for doc in results.documents] if hasattr(results, 'documents') else [],
                            "metadatas": [doc.metadata for doc in results.documents] if hasattr(results, 'documents') else [],
                            "scores": [doc.score for doc in results.documents] if hasattr(results, 'documents') else []
                        },
                        "count": len(results.documents) if hasattr(results, 'documents') else 0
                    }
                    
                except Exception as fallback_error:
                    self.logger.warning(f"Vector DB query also failed: {fallback_error}, using simple search")
                    
                    # Final fallback: Simple text matching (load entries if needed)
                    if not hasattr(self, '_cached_entries'):
                        self.logger.info("Loading knowledge base for fallback search")
                        try:
                            self._cached_entries = self.load_knowledge_base()
                        except Exception as load_error:
                            self.logger.error(f"Failed to load knowledge base for fallback: {load_error}")
                            self._cached_entries = []
                    
                    if self._cached_entries:
                        formatted_results = self._simple_text_search(query, max_results)
                    else:
                        formatted_results = {
                            "query": query,
                            "results": {"documents": [], "metadatas": [], "scores": []},
                            "count": 0
                        }
            
            cprint(f"✅ Found {formatted_results['count']} results", "green")
            return formatted_results
            
        except Exception as e:
            self.logger.error(f"Search failed: {e}")
            cprint(f"❌ Search failed: {e}", "red")
            return {"query": query, "results": {"documents": [], "metadatas": [], "scores": []}, "count": 0}
    
    def _simple_text_search(self, query: str, max_results: int) -> Dict[str, Any]:
        """Fallback simple text search when API methods fail."""
        if not hasattr(self, '_cached_entries'):
            return {"query": query, "results": {"documents": [], "metadatas": [], "scores": []}, "count": 0}
        
        query_lower = query.lower()
        matches = []
        
        for entry in self._cached_entries:
            score = 0
            text_lower = entry.combined_text.lower()
            
            # Simple scoring based on keyword presence
            for word in query_lower.split():
                if word in text_lower:
                    score += text_lower.count(word)
            
            if score > 0:
                matches.append((entry, score))
        
        # Sort by score and limit results
        matches.sort(key=lambda x: x[1], reverse=True)
        matches = matches[:max_results]
        
        # Format results
        documents = [entry.combined_text for entry, _ in matches]
        metadatas = [
            {
                "error": entry.error,
                "error_type": entry.error_type,
                "severity": entry.severity,
                "solution": entry.solution,
                "source": entry.source,
                "reference": entry.reference
            }
            for entry, _ in matches
        ]
        scores = [min(score / 10.0, 1.0) for _, score in matches]  # Normalize to 0-1
        
        return {
            "query": query,
            "results": {
                "documents": documents,
                "metadatas": metadatas,
                "scores": scores
            },
            "count": len(matches)
        }
    
    def get_stats(self) -> Dict[str, Any]:
        """Get comprehensive statistics about the knowledge base."""
        try:
            # Get vector database info
            vector_dbs = self.client.vector_dbs.list()
            our_db = next((db for db in vector_dbs if db.identifier == self.vector_db_id), None)
            
            stats = {
                "vector_db_id": self.vector_db_id,
                "endpoint_url": self.endpoint_url,
                "embedding_model": self.embedding_model,
                "embedding_dimension": self.embedding_dimension,
                "chunk_size": self.chunk_size,
                "batch_size": self.batch_size,
                "ingestion_metrics": self.metrics.to_dict() if self.metrics else {},
                "vector_db_exists": our_db is not None,
                "timestamp": datetime.now().isoformat()
            }
            
            return stats
            
        except Exception as e:
            self.logger.error(f"Failed to get stats: {e}")
            return {"error": str(e)}


def run(action: str = "ingest", **kwargs) -> Dict[str, Any]:
    """
    Programmatic interface to run ingester operations.
    
    Args:
        action: Operation to perform ('ingest', 'search', 'stats', 'validate')
        **kwargs: Configuration parameters
        
    Returns:
        Dict containing operation results and metadata
    """
    # Default configuration
    base_config = {
        "kb_file": "kb.json",
        "endpoint_url": "http://localhost:8321", 
        "vector_db_id": "tekton_errors_db",
        "embedding_model": "text-embedding-004",
        "embedding_dimension": 384,
        "chunk_size": 512,
        "batch_size": 50,
        "log_level": "INFO"
    }
    
    # Valid constructor parameters for TektonKnowledgeIngester
    constructor_params = {
        "kb_file", "endpoint_url", "vector_db_id", "embedding_model", 
        "embedding_dimension", "chunk_size", "batch_size", "log_level"
    }
    
    # Filter constructor parameters from kwargs
    ingester_config = base_config.copy()
    action_params = {}
    
    for key, value in kwargs.items():
        if key in constructor_params:
            ingester_config[key] = value
        else:
            action_params[key] = value
    
    try:
        # Initialize ingester with only valid constructor parameters
        ingester = TektonKnowledgeIngester(**ingester_config)
        
        if action == "ingest":
            # Load and ingest knowledge base
            entries = ingester.load_knowledge_base()
            metrics = ingester.ingest_entries(entries)
            
            return {
                "status": "success",
                "action": "ingest",
                "metrics": metrics.to_dict(),
                "message": f"Ingested {metrics.successful_ingestions} entries successfully"
            }
            
        elif action == "search":
            query = action_params.get("query", "")
            max_results = action_params.get("max_results", 5)
            
            if not query:
                return {
                    "status": "error", 
                    "action": "search",
                    "message": "Query parameter required for search"
                }
            
            results = ingester.search(query, max_results)
            return {
                "status": "success",
                "action": "search", 
                "results": results,
                "message": f"Found {results['count']} results for query: '{query}'"
            }
            
        elif action == "stats":
            stats = ingester.get_stats()
            return {
                "status": "success",
                "action": "stats",
                "stats": stats,
                "message": "Statistics retrieved successfully"
            }
            
        elif action == "validate":
            # Validate knowledge base format and vector DB connectivity
            entries = ingester.load_knowledge_base()
            db_exists = ingester.vector_db_exists()
            
            validation_results = {
                "kb_file_exists": ingester.kb_file.exists(),
                "kb_entries_count": len(entries),
                "vector_db_exists": db_exists,
                "vector_db_accessible": True,  # If we got this far, it's accessible
                "endpoint_reachable": True
            }
            
            return {
                "status": "success",
                "action": "validate",
                "validation": validation_results,
                "message": f"Validation completed - {len(entries)} entries found, DB exists: {db_exists}"
            }
            
        elif action == "reset":
            # Reset/clear vector database
            if ingester.vector_db_exists():
                try:
                    # Try to unregister the vector database
                    ingester.client.vector_dbs.unregister(vector_db_id=ingester.vector_db_id)
                    ingester.logger.info(f"Unregistered vector database: {ingester.vector_db_id}")
                    
                    # Recreate the vector database
                    ingester._setup_vector_db()
                    
                    return {
                        "status": "success",
                        "action": "reset",
                        "message": f"Vector database '{ingester.vector_db_id}' has been reset"
                    }
                    
                except Exception as reset_error:
                    ingester.logger.warning(f"Failed to unregister vector DB: {reset_error}")
                    
                    # Fallback: Try to clear documents if unregister fails
                    try:
                        # Get all document IDs and delete them
                        docs = ingester.client.vector_dbs.get_documents(vector_db_id=ingester.vector_db_id)
                        if docs:
                            for doc in docs:
                                try:
                                    ingester.client.vector_dbs.delete_document(
                                        vector_db_id=ingester.vector_db_id,
                                        document_id=doc.document_id
                                    )
                                except Exception as doc_delete_error:
                                    ingester.logger.warning(f"Failed to delete document {doc.document_id}: {doc_delete_error}")
                        
                        return {
                            "status": "success",
                            "action": "reset",
                            "message": f"Vector database '{ingester.vector_db_id}' documents cleared"
                        }
                        
                    except Exception as clear_error:
                        return {
                            "status": "error",
                            "action": "reset",
                            "error": str(clear_error),
                            "message": f"Failed to reset vector database: {clear_error}"
                        }
            else:
                return {
                    "status": "warning",
                    "action": "reset", 
                    "message": f"Vector database '{ingester.vector_db_id}' does not exist"
                }
        
        else:
            return {
                "status": "error",
                "action": action,
                "message": f"Unknown action: {action}. Available: ingest, search, stats, validate, reset"
            }
            
    except Exception as e:
        return {
            "status": "error",
            "action": action,
            "error": str(e),
            "message": f"Operation failed: {e}"
        }


def main():
    """Enhanced CLI entry point with comprehensive argument handling."""
    import argparse
    import sys
    
    parser = argparse.ArgumentParser(
        description="Tekton Knowledge Base Ingester - Ingest TaskRun failures into LlamaStack",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # Basic ingestion
  python ingester.py
  
  # Search for solutions
  python ingester.py --search "registry push failed" --max-results 3
  
  # Show statistics
  python ingester.py --stats
  
  # Validate setup
  python ingester.py --validate
  
  # Reset vector database
  python ingester.py --reset
  
  # Custom configuration
  python ingester.py --endpoint http://localhost:8080 --vector-db my_errors_db
  
  # Quiet mode with JSON output
  python ingester.py --search "build failed" --quiet --json-output
        """
    )
    
    # Action arguments (mutually exclusive)
    action_group = parser.add_mutually_exclusive_group()
    action_group.add_argument("--search", metavar="QUERY", help="Search for solutions using query")
    action_group.add_argument("--stats", action="store_true", help="Show database statistics")
    action_group.add_argument("--validate", action="store_true", help="Validate setup and connectivity")
    action_group.add_argument("--reset", action="store_true", help="Reset/clear vector database")
    
    # Configuration arguments
    config_group = parser.add_argument_group("Configuration")
    config_group.add_argument("--kb-file", default="kb.json", help="Knowledge base JSON file (default: kb.json)")
    config_group.add_argument("--endpoint", default="http://localhost:8321", help="LlamaStack endpoint URL")
    config_group.add_argument("--vector-db", default="tekton_errors_db", help="Vector database identifier")
    config_group.add_argument("--embedding-model", default="text-embedding-004", help="Embedding model name")
    config_group.add_argument("--embedding-dim", type=int, default=384, help="Embedding dimension")
    config_group.add_argument("--chunk-size", type=int, default=512, help="Chunk size in tokens")
    config_group.add_argument("--batch-size", type=int, default=50, help="Batch size for ingestion")
    
    # Search specific arguments
    search_group = parser.add_argument_group("Search Options")
    search_group.add_argument("--max-results", type=int, default=5, help="Maximum search results")
    search_group.add_argument("--min-score", type=float, help="Minimum similarity score filter")
    
    # Output and logging arguments
    output_group = parser.add_argument_group("Output Options")
    output_group.add_argument("--log-level", default="INFO", choices=["DEBUG", "INFO", "WARNING", "ERROR"], 
                             help="Logging level")
    output_group.add_argument("--quiet", action="store_true", help="Suppress colored output")
    output_group.add_argument("--json-output", action="store_true", help="Output results in JSON format")
    output_group.add_argument("--output-file", help="Save results to file")
    
    # Advanced options
    advanced_group = parser.add_argument_group("Advanced Options")
    advanced_group.add_argument("--force", action="store_true", help="Force operations without confirmation")
    advanced_group.add_argument("--dry-run", action="store_true", help="Show what would be done without executing")
    advanced_group.add_argument("--config-file", help="Load configuration from JSON file")
    
    args = parser.parse_args()
    
    # Load configuration from file if specified
    config = {}
    if args.config_file:
        try:
            with open(args.config_file, 'r') as f:
                config = json.load(f)
        except Exception as e:
            cprint(f"❌ Failed to load config file: {e}", "red")
            sys.exit(1)
    
    # Build base configuration for ingester constructor
    base_config = {
        "kb_file": args.kb_file,
        "endpoint_url": args.endpoint,
        "vector_db_id": args.vector_db,
        "embedding_model": args.embedding_model,
        "embedding_dimension": args.embedding_dim,
        "chunk_size": args.chunk_size,
        "batch_size": args.batch_size,
        "log_level": args.log_level
    }
    
    # Update base config with any loaded config file
    base_config.update(config)
    
    # Determine action and build action-specific config
    action_config = base_config.copy()
    
    if args.search:
        action = "search"
        action_config["query"] = args.search
        action_config["max_results"] = args.max_results
        if args.min_score:
            action_config["min_score"] = args.min_score
    elif args.stats:
        action = "stats"
    elif args.validate:
        action = "validate"
    elif args.reset:
        action = "reset"
        if not args.force:
            confirm = input(f"⚠️  This will delete vector database '{args.vector_db}'. Continue? [y/N]: ")
            if confirm.lower() != 'y':
                cprint("Operation cancelled", "yellow")
                sys.exit(0)
    else:
        action = "ingest"
    
    # Dry run mode
    if args.dry_run:
        cprint(f"🔍 DRY RUN: Would execute action '{action}' with config:", "cyan")
        for key, value in action_config.items():
            print(f"  {key}: {value}")
        sys.exit(0)
    
    try:
        # Execute action
        result = run(action, **action_config)
        
        # Handle output
        if args.json_output:
            import json
            output = json.dumps(result, indent=2, default=str)
        else:
            output = _format_result(result, args.quiet)
        
        # Save to file if specified
        if args.output_file:
            with open(args.output_file, 'w') as f:
                f.write(output)
            if not args.quiet:
                cprint(f"📁 Results saved to: {args.output_file}", "green")
        else:
            print(output)
        
        # Exit with appropriate code
        if result["status"] == "error":
            sys.exit(1)
        elif result["status"] == "warning":
            sys.exit(2)
        else:
            sys.exit(0)
    
    except KeyboardInterrupt:
        cprint("\n⚠️  Operation cancelled by user", "yellow")
        sys.exit(130)
    except Exception as e:
        if args.json_output:
            error_result = {"status": "error", "error": str(e), "action": action}
            print(json.dumps(error_result, indent=2))
        else:
            cprint(f"💥 Unexpected error: {e}", "red")
        sys.exit(1)


def _format_result(result: Dict[str, Any], quiet: bool = False) -> str:
    """Format result for human-readable output."""
    status = result["status"]
    action = result["action"]
    
    if quiet:
        return result.get("message", "")
    
    # Status emoji
    status_emoji = {
        "success": "✅",
        "error": "❌", 
        "warning": "⚠️"
    }
    
    output = []
    output.append(f"{status_emoji.get(status, '❓')} {action.title()} {status.title()}")
    output.append("=" * 50)
    
    if result.get("message"):
        output.append(f"Message: {result['message']}")
    
    # Action-specific formatting
    if action == "search" and result.get("results"):
        search_results = result["results"]
        output.append(f"\n🔍 Search Results for: '{search_results['query']}'")
        output.append(f"Found {search_results['count']} results:")
        
        for i, (doc, metadata, score) in enumerate(zip(
            search_results["results"]["documents"],
            search_results["results"]["metadatas"], 
            search_results["results"]["scores"]
        )):
            output.append(f"\n#{i+1} (Score: {score:.3f})")
            output.append(f"Error Type: {metadata.get('error_type', 'N/A')}")
            output.append(f"Severity: {metadata.get('severity', 'N/A')}")
            output.append(f"Error: {metadata.get('error', 'N/A')[:100]}...")
            if i < len(search_results["results"]["documents"]) - 1:
                output.append("-" * 30)
    
    elif action == "stats" and result.get("stats"):
        stats = result["stats"]
        output.append(f"\n📊 Database Statistics:")
        for key, value in stats.items():
            if key != "ingestion_metrics":
                output.append(f"{key}: {value}")
        
        if stats.get("ingestion_metrics"):
            metrics = stats["ingestion_metrics"]
            output.append(f"\n📈 Ingestion Metrics:")
            output.append(f"Total entries: {metrics.get('total_entries', 0)}")
            output.append(f"Success rate: {metrics.get('success_rate', 0):.1f}%")
            output.append(f"Processing speed: {metrics.get('entries_per_second', 0):.1f} entries/sec")
    
    elif action == "validate" and result.get("validation"):
        validation = result["validation"]
        output.append(f"\n🔍 Validation Results:")
        for key, value in validation.items():
            status_icon = "✅" if value else "❌"
            output.append(f"{status_icon} {key}: {value}")
    
    elif action == "ingest" and result.get("metrics"):
        metrics = result["metrics"]
        output.append(f"\n📊 Ingestion Results:")
        output.append(f"Total entries: {metrics.get('total_entries', 0)}")
        output.append(f"Successful: {metrics.get('successful_ingestions', 0)}")
        output.append(f"Failed: {metrics.get('failed_ingestions', 0)}")
        output.append(f"Success rate: {metrics.get('success_rate', 0):.1f}%")
        output.append(f"Processing time: {metrics.get('processing_time', 0):.2f}s")
        output.append(f"Speed: {metrics.get('entries_per_second', 0):.1f} entries/sec")
    
    if result.get("error"):
        output.append(f"\n❌ Error: {result['error']}")
    
    return "\n".join(output)


if __name__ == "__main__":
    main()
