"""Built-in tools available to the agent."""

from __future__ import annotations

import math
import datetime
import os
from dataclasses import dataclass, field
from typing import Callable, Dict


@dataclass
class Tool:
    """Descriptor for a single agent tool.

    Attributes:
        name:        Unique identifier used in action strings.
        description: Human-readable explanation shown to the LLM.
        func:        Callable that receives a single string argument and
                     returns a string observation.
    """

    name: str
    description: str
    func: Callable[[str], str]

    def run(self, argument: str) -> str:
        """Execute the tool and return the observation string."""
        try:
            return str(self.func(argument))
        except Exception as exc:  # noqa: BLE001
            return f"Error: {exc}"


# ---------------------------------------------------------------------------
# Default tool implementations
# ---------------------------------------------------------------------------

def _calculator(expression: str) -> str:
    """Evaluate a safe mathematical expression."""
    allowed_names: Dict[str, object] = {
        k: v for k, v in math.__dict__.items() if not k.startswith("_")
    }
    allowed_names["abs"] = abs
    allowed_names["round"] = round
    try:
        result = eval(expression, {"__builtins__": {}}, allowed_names)  # noqa: S307
        return str(result)
    except Exception as exc:  # noqa: BLE001
        return f"Invalid expression: {exc}"


def _current_time(_: str) -> str:
    """Return the current date and time."""
    now = datetime.datetime.now()
    return now.strftime("%Y-%m-%d %H:%M:%S")


def _read_file(path: str) -> str:
    """Read the contents of a local file."""
    path = path.strip()
    if not os.path.isfile(path):
        return f"File not found: {path}"
    with open(path, encoding="utf-8") as fh:
        return fh.read()


def _write_file(argument: str) -> str:
    """Write text to a file.  Argument format: ``<path>|<content>``."""
    if "|" not in argument:
        return "Usage: <path>|<content>"
    path, _, content = argument.partition("|")
    path = path.strip()
    with open(path, "w", encoding="utf-8") as fh:
        fh.write(content)
    return f"Written {len(content)} characters to '{path}'."


def _web_search(query: str) -> str:
    """Stub web-search tool (returns a placeholder response).

    Replace this implementation with a real search API integration
    (e.g. SerpAPI, Bing Search, DuckDuckGo) when an API key is available.
    """
    return (
        f"[Web Search stub] Query received: '{query}'. "
        "Configure a real search API to get live results."
    )


# ---------------------------------------------------------------------------
# Registry of all default tools
# ---------------------------------------------------------------------------

DEFAULT_TOOLS: list[Tool] = [
    Tool(
        name="calculator",
        description=(
            "Evaluate a mathematical expression. "
            "Input: a Python-compatible math expression, e.g. '2 ** 10' or 'sqrt(144)'."
        ),
        func=_calculator,
    ),
    Tool(
        name="current_time",
        description="Return the current local date and time. Input: ignored.",
        func=_current_time,
    ),
    Tool(
        name="read_file",
        description="Read the contents of a local file. Input: absolute or relative file path.",
        func=_read_file,
    ),
    Tool(
        name="write_file",
        description=(
            "Write text content to a local file. "
            "Input format: '<path>|<content>' where | separates the file path from the content."
        ),
        func=_write_file,
    ),
    Tool(
        name="web_search",
        description="Search the web for information. Input: search query string.",
        func=_web_search,
    ),
]
