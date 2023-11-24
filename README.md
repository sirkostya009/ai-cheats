# AI-Cheats

AI Cheats is a backend-plugin solution for cracking platform-specific tests (JetIQ in particular) with OpenAI.

### How To Run:

You'll need 3 environment variables set:
1. `DATABASE_URL` - The connection URL to your postgres database
2. `PORT` - The port you want the server to run on
3. `OPENAI_KEY` - Your OpenAI API key

The server is originally designed to be easily deployable on [fl0](https://fl0.com), but you can run it locally
yourself if you have the database running and variables set.

### Stuff that may be improved:
- [ ] IP-based usage locking
