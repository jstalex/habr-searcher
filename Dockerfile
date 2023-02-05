FROM golang:1.19

ENV TokenForHabrSearcher 5623528912:AAGD6aPPzi0xZCWtwK3nJdYf4eIoNgHnULw

RUN mkdir /home/habr-searcher

COPY . /home/habr-searcher 

WORKDIR /home/habr-searcher

CMD ["make", "run"]

