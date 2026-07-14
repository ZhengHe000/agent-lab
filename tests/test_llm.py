"""Tests for agent/llm.py"""

import pytest

from agent.llm import MockLLM, create_llm


class TestMockLLM:
    def setup_method(self):
        self.llm = MockLLM()

    def _ask(self, question: str) -> str:
        return self.llm.complete(f"Question: {question}")

    def test_time_question_triggers_current_time(self):
        result = self._ask("What time is it now?")
        assert "current_time" in result

    def test_math_question_triggers_calculator(self):
        result = self._ask("calculate 2 + 2")
        assert "calculator" in result

    def test_search_question_triggers_web_search(self):
        result = self._ask("search for Python tutorials")
        assert "web_search" in result

    def test_file_question_triggers_read_file(self):
        result = self._ask("read file /tmp/test.txt")
        assert "read_file" in result

    def test_default_returns_final_answer(self):
        result = self._ask("What is the meaning of life?")
        assert "Final Answer" in result

    def test_empty_prompt_handled(self):
        result = self.llm.complete("")
        assert isinstance(result, str)


class TestCreateLLM:
    def test_mock_backend(self):
        llm = create_llm("mock")
        assert isinstance(llm, MockLLM)

    def test_unknown_backend_raises(self):
        with pytest.raises(ValueError, match="Unknown backend"):
            create_llm("unknown_backend")
