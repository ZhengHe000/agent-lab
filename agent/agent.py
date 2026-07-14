"""Core agent implementing the ReAct (Reason + Act) loop.

The loop follows this pattern until a Final Answer is produced or
the maximum number of steps is reached:

    1. Build a prompt that includes the tools description, conversation
       history, and the current question.
    2. Call the LLM and parse its response for a Thought / Action /
       Action Input triplet **or** a Final Answer.
    3. If an action is requested, execute the corresponding tool and
       append the observation to the history.
    4. Repeat.
"""

from __future__ import annotations

import re
import textwrap
from dataclasses import dataclass, field
from typing import Dict, List, Optional, Tuple

from .llm import LLMBase, MockLLM
from .tools import DEFAULT_TOOLS, Tool


@dataclass
class Step:
    """A single step in the agent's reasoning trace."""

    thought: str = ""
    action: str = ""
    action_input: str = ""
    observation: str = ""
    final_answer: str = ""

    @property
    def is_final(self) -> bool:
        return bool(self.final_answer)

    def __str__(self) -> str:
        parts: List[str] = []
        if self.thought:
            parts.append(f"Thought: {self.thought}")
        if self.action:
            parts.append(f"Action: {self.action}")
            parts.append(f"Action Input: {self.action_input}")
        if self.observation:
            parts.append(f"Observation: {self.observation}")
        if self.final_answer:
            parts.append(f"Final Answer: {self.final_answer}")
        return "\n".join(parts)


class Agent:
    """An AI agent that reasons and acts using registered tools.

    Args:
        llm:       Language model backend (defaults to :class:`MockLLM`).
        tools:     List of :class:`Tool` objects. Defaults to
                   :data:`DEFAULT_TOOLS`.
        max_steps: Maximum ReAct iterations before giving up (default: 10).
        verbose:   Print each step to stdout when *True* (default: False).
    """

    def __init__(
        self,
        llm: Optional[LLMBase] = None,
        tools: Optional[List[Tool]] = None,
        max_steps: int = 10,
        verbose: bool = False,
    ) -> None:
        self._llm: LLMBase = llm or MockLLM()
        self._tools: Dict[str, Tool] = {
            t.name: t for t in (tools if tools is not None else DEFAULT_TOOLS)
        }
        self._max_steps = max_steps
        self._verbose = verbose

    # ------------------------------------------------------------------
    # Public API
    # ------------------------------------------------------------------

    def run(self, question: str) -> str:
        """Run the agent on *question* and return the final answer string.

        Args:
            question: The user's natural-language question or task.

        Returns:
            The agent's final answer as a plain string.
        """
        history: List[Step] = []

        for step_idx in range(self._max_steps):
            prompt = self._build_prompt(question, history)
            raw = self._llm.complete(prompt)
            step = self._parse_response(raw)

            if self._verbose:
                print(f"\n--- Step {step_idx + 1} ---")
                print(step)

            if step.action and step.action in self._tools:
                step.observation = self._tools[step.action].run(step.action_input)
                if self._verbose:
                    print(f"Observation: {step.observation}")
            elif step.action and step.action not in self._tools:
                step.observation = (
                    f"Unknown tool '{step.action}'. "
                    f"Available tools: {', '.join(self._tools)}"
                )

            history.append(step)

            if step.is_final:
                return step.final_answer

        # Fallback: return the last observation or a timeout message
        if history and history[-1].observation:
            return history[-1].observation
        return "Agent reached the maximum number of steps without a final answer."

    def add_tool(self, tool: Tool) -> None:
        """Register a new tool with the agent."""
        self._tools[tool.name] = tool

    @property
    def tools(self) -> Dict[str, Tool]:
        """Read-only view of registered tools."""
        return dict(self._tools)

    # ------------------------------------------------------------------
    # Internal helpers
    # ------------------------------------------------------------------

    def _build_prompt(self, question: str, history: List[Step]) -> str:
        tool_descriptions = "\n".join(
            f"- {name}: {tool.description}"
            for name, tool in self._tools.items()
        )

        system_block = textwrap.dedent(f"""
            You are a helpful AI assistant (智能体). Answer the user's question by
            reasoning step-by-step and using the available tools when necessary.

            Available tools:
            {tool_descriptions}

            Use the following format:
                Thought: <your reasoning>
                Action: <tool name>
                Action Input: <tool argument>
                Observation: <tool result>
                ... (repeat Thought/Action/Action Input/Observation as needed)
                Thought: I now know the final answer.
                Final Answer: <your final answer>

            If you can answer directly without a tool, write:
                Thought: <your reasoning>
                Final Answer: <your answer>
        """).strip()

        history_block = ""
        for step in history:
            history_block += f"\n{step}\n"

        return f"{system_block}\n\nQuestion: {question}{history_block}"

    @staticmethod
    def _parse_response(text: str) -> Step:
        """Parse the LLM output into a :class:`Step`."""
        step = Step()

        thought_match = re.search(r"Thought:\s*(.+?)(?=\n(?:Action|Final Answer)|$)", text, re.S)
        if thought_match:
            step.thought = thought_match.group(1).strip()

        action_match = re.search(r"Action:\s*(.+)", text)
        if action_match:
            step.action = action_match.group(1).strip()

        action_input_match = re.search(r"Action Input:\s*(.+?)(?=\nObservation|\Z)", text, re.S)
        if action_input_match:
            step.action_input = action_input_match.group(1).strip()

        final_match = re.search(r"Final Answer:\s*(.+)", text, re.S)
        if final_match:
            step.final_answer = final_match.group(1).strip()

        return step
