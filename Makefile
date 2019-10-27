build:
	@docker build -t geekshub:latest .

run:
	@docker run -it geekshub:latest /bin/bash
