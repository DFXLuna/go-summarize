## go-summarize
Video summarization pipeline, currently only supports transcription.

Something that sets this project apart from others is that semantic chunking is done manually, rather than just throwing everything into a vectordb and calling it a day.


## Observations
- Whisper peaks at ~11GB VRAM when transcribing a long file
