package main

import "github.com/jackc/pgx/v5/pgxpool"

var Pool = pool()

type Notifier struct {
    enqueue chan int
}

type Worker struct {
    done    chan *Worker
    payload []Event 
}

func pool() func() *pgxpool.Pool {
    var pool *pgxpool.Pool

    return func() *pgxpool.Pool {
        if pool != nil {
            return pool
        }

        cfg, e := pgxpool.ParseConfig()

        return pool
    }
}

func notifier() *Notifier {
    return &Notifier{
        make(chan int),
    } 
}

func worker() *Worker {
    return &Worker{
        make(chan *Worker),
        make([]Event, 1),
    }
}

func (n *Notifier) run() {

}

func (w *Worker) run() {

}
