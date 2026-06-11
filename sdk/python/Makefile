.PHONY: typecheck typecheck-strict lint test install-dev help

PYTHON ?= python3
MYPY ?= mypy

help:
	@echo "agentapi Python SDK — make targets"
	@echo "  install-dev    Install dev dependencies (mypy, ruff, pytest)"
	@echo "  typecheck      Run mypy in non-strict mode"
	@echo "  typecheck-strict  Run mypy --strict ./  (CI gate)"
	@echo "  lint           Run ruff"
	@echo "  test           Run pytest"

install-dev:
	$(PYTHON) -m pip install -e ".[dev]"

typecheck:
	$(MYPY) agentapi/

typecheck-strict:
	$(MYPY) --strict ./

lint:
	ruff check .

test:
	$(PYTHON) -m pytest
