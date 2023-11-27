# AI-Cheats

AI Cheats is a backend-plugin solution for cracking platform-specific tests (JetIQ in particular) with OpenAI.

### How To Run

You'll need 3 environment variables set:
1. `DATABASE_URL` - The connection URL to your postgres database
2. `PORT` - The port you want the server to run on
3. `OPENAI_KEY` - Your OpenAI API key

The server is originally designed to be easily deployable on [fl0](https://fl0.com), but you can run it locally
yourself if you have the database running and variables set.

Also make sure you run migration scripts in `main/db` folder.

### How To Use
After successfully booting up the server, you can send an any request to `/:id` with an example body:
```text
What nationality was Frédéric Chopin?
1. German
2. French
3. Polish
4. Italian
```

Before doing so, however, make sure you have a user with the same ID in the database, otherwise you'll get a 404, and
also that user must have `active_till` column set to any future date value.

### Stuff that may be improved
- [ ] IP-based usage locking
- [ ] Fix Dockerfile
- [ ] Add tests
