"""
agent-lab: A lightweight AI Agent framework implementing the ReAct pattern.

Classes exported at package level:
    Agent   — the core reasoning/acting loop
    Tool    — base class / descriptor for agent tools
"""

from .agent import Agent
from .tools import Tool

__all__ = ["Agent", "Tool"]
