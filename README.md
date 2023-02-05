# Habr-searcher

This app can help you with tracking habr posts on any tags

## Installation
Clone repo

    $ git clone https://github.com/jstalex/habr-searcher
Change directory

    $ cd habr-searcher
Set your Telegram API token

    $ export TokenForHabrSearcher="<your token>"
Build & run

    $ make run 


Or use Docker

Build

    $ docker build -t habr-searcher .

Run

    $ docker run -d -e TokenForHabrSearcher="<your token>" habr-searcher
