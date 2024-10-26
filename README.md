# Random Episode

Warning: This is in _very_ early development and is currently very basic and
will crash on error. Not recommended for everyday use yet.

This is a simple self-hosted web app for picking you a random epsiode of your
favourite TV show to watch.

The idea came from the fact that I'm incredibly indecisive and would spend longer
picking an episode than the episode would last.

It makes use of [TMDB](https://www.themoviedb.org/) for getting show and episode
information.

## Installation

There is a docker image published at `ghcr.io/danmharris/random-episode`.

You will need the following environment variables:

| Name          | Example        | Description                             |
| ------------- | -------------- | --------------------------------------- |
| PORT          | 3000           | Port to listen on                       |
| TMDB_TOKEN    | n/a            | Token to use to look up show data       |
| POSTGRES_USER | random_episode | User of postgresql db to connect to     |
| POSTGRES_PASS | Password123!   | Password of postgresql db to connect to |
| POSTGRES_HOST | 127.0.0.1      | Database host                           |
| POSTGRES_PORT | 5432           | Port of database host                   |
| POSTGRES_DB   | random_episode | Database name                           |

Note there are currently no defaults and all variables are required.

There is also a `docker-compose.yml` file which has a basic setup of PostgreSQL
and pre-populates the `POSTGRES_` variables for you.
