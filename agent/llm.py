"""LLM abstraction layer.

Provides:
    LLMBase   — abstract interface every backend must implement.
    MockLLM   — rule-based backend that works without any API key.
    OpenAILLM — thin wrapper around the OpenAI chat-completions API
                (requires ``openai`` package and ``OPENAI_API_KEY``).
"""

from __future__ import annotations

import os
import re
from abc import ABC, abstractmethod


class LLMBase(ABC):
    """Abstract language model interface."""

    @abstractmethod
    def complete(self, prompt: str) -> str:
        """Return a completion string for the given prompt."""


# ---------------------------------------------------------------------------
# Mock (rule-based) LLM — no API key required
# ---------------------------------------------------------------------------

class MockLLM(LLMBase):
    """A rule-based mock LLM useful for testing and offline demos.

    Parses simple keywords in the user question and emits the
    corresponding Thought / Action / Final Answer pattern.
    """

    def complete(self, prompt: str) -> str:
        # If the prompt already contains a real observation (after the question),
        # wrap up with a Final Answer.
        if self._has_real_observation(prompt):
            last_obs = self._extract_last_observation(prompt)
            return (
                "Thought: I now have the observation and can give a final answer.\n"
                f"Final Answer: {last_obs}\n"
            )

        question = self._extract_question(prompt)
        q_lower = question.lower()

        if any(word in q_lower for word in ("time", "date", "now", "today", "当前时间", "现在")):
            return (
                "Thought: The user wants to know the current time.\n"
                "Action: current_time\n"
                "Action Input: now\n"
            )

        if any(word in q_lower for word in ("read", "open", "file", "读取", "文件")):
            words = question.split()
            path = next((w for w in words if os.sep in w or "." in w), "unknown_file.txt")
            return (
                f"Thought: The user wants to read a file.\n"
                f"Action: read_file\n"
                f"Action Input: {path}\n"
            )

        if any(word in q_lower for word in ("calculate", "compute", "math", "计算", "+", "-", "*")):
            # Strip common keyword prefixes and keep only the mathematical expression
            expr = re.sub(
                r"^(?:calculate|compute|math|what is|eval|evaluate)\s*",
                "",
                question,
                flags=re.IGNORECASE,
            ).strip() or "1+1"
            return (
                f"Thought: This is a math problem. I will use the calculator.\n"
                f"Action: calculator\n"
                f"Action Input: {expr}\n"
            )

        if any(word in q_lower for word in ("search", "find", "look up", "搜索", "查找", "查询")):
            return (
                f"Thought: I should search the web for relevant information.\n"
                f"Action: web_search\n"
                f"Action Input: {question}\n"
            )

        # Default: answer directly without tools
        return (
            f"Thought: I can answer this directly.\n"
            f"Final Answer: I received your question: \"{question}\". "
            "To enable full AI capabilities, set OPENAI_API_KEY and use OpenAILLM.\n"
        )

    @staticmethod
    def _has_real_observation(prompt: str) -> bool:
        """Return True if the prompt contains an actual Observation: after the Question line."""
        question_idx = prompt.find("Question:")
        if question_idx == -1:
            return False
        after_question = prompt[question_idx:]
        return "Observation:" in after_question

    @staticmethod
    def _extract_last_observation(prompt: str) -> str:
        """Return the last Observation value found in the prompt."""
        matches = re.findall(r"Observation:\s*(.+?)(?=\nThought:|\Z)", prompt, re.S)
        if matches:
            return matches[-1].strip()
        return "done"

    @staticmethod
    def _extract_question(prompt: str) -> str:
        """Pull the most recent user question from the prompt."""
        lines = prompt.strip().splitlines()
        for line in reversed(lines):
            if line.startswith("Question:"):
                return line[len("Question:"):].strip()
        return lines[-1].strip() if lines else ""


# ---------------------------------------------------------------------------
# OpenAI backend
# ---------------------------------------------------------------------------

class OpenAILLM(LLMBase):
    """Wrapper around the OpenAI chat-completions API.

    Requires:
        - ``pip install openai``
        - Environment variable ``OPENAI_API_KEY`` set to a valid key.

    Args:
        model:       OpenAI model identifier (default: ``gpt-4o-mini``).
        temperature: Sampling temperature (default: 0).
    """

    def __init__(self, model: str = "gpt-4o-mini", temperature: float = 0.0) -> None:
        try:
            import openai  # noqa: PLC0415
        except ImportError as exc:
            raise ImportError(
                "OpenAI package not installed. Run: pip install openai"
            ) from exc

        api_key = os.environ.get("OPENAI_API_KEY")
        if not api_key:
            raise EnvironmentError(
                "OPENAI_API_KEY environment variable is not set."
            )

        self._client = openai.OpenAI(api_key=api_key)
        self._model = model
        self._temperature = temperature

    def complete(self, prompt: str) -> str:
        response = self._client.chat.completions.create(
            model=self._model,
            temperature=self._temperature,
            messages=[{"role": "user", "content": prompt}],
        )
        return response.choices[0].message.content or ""


# ---------------------------------------------------------------------------
# Factory helper
# ---------------------------------------------------------------------------

def create_llm(backend: str = "mock", **kwargs: object) -> LLMBase:
    """Create and return an LLM instance by name.

    Args:
        backend: ``"mock"`` (default) or ``"openai"``.
        **kwargs: Forwarded to the backend constructor.

    Returns:
        An :class:`LLMBase` instance.
    """
    backends = {
        "mock": MockLLM,
        "openai": OpenAILLM,
    }
    if backend not in backends:
        raise ValueError(f"Unknown backend '{backend}'. Choose from: {list(backends)}")
    return backends[backend](**kwargs)
