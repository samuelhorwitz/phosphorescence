const prepareNextTick = Symbol('prepareNextTick');
const tick = Symbol('tick');
const nextTick = Symbol('nextTick');

export default class Queue {
    #logger;
    #destroyed;
    #queue;
    #tickQueue;

    constructor(logger) {
        this.#logger = logger.addPrefix('Queue');
        this.#logger.debug('Queue created');
        this.#destroyed = false;
        this.#queue = [];
        this.#tickQueue = [];
        this[prepareNextTick]();
    }

    [prepareNextTick]() {
        this.#logger.debug('Preparing next tick');
        let resolver;
        let promise = new Promise(resolve => resolver = resolve);
        this.#tickQueue.push({promise, resolver});
    }

    [tick]() {
        this.#logger.debug('Tick');
        let {resolver} = this.#tickQueue.shift();
        resolver();
    }

    [nextTick]() {
        let {promise} = this.#tickQueue[0];
        return promise;
    }

    push(x) {
        if (this.#destroyed) {
            throw new Error('Queue destroyed');
        }
        this.#logger.debug('Push', x);
        this.#queue.push(x);
        this[prepareNextTick]();
        if (this.#queue.length == 1) {
            this.#logger.debug('Has elements again');
            this[tick]();
        }
    }

    async shift() {
        if (this.#destroyed) {
            throw new Error('Queue destroyed');
        }
        this.#logger.debug('Wants to shift');
        if (this.#queue.length == 0) {
            this.#logger.debug('Waits to shift because queue has no elements');
            await this[nextTick]();
            if (this.#destroyed) {
                throw new Error('Queue destroyed');
            }
        }
        this.#logger.debug('Shift', this.#queue[0]);
        return this.#queue.shift();
    }

    destroy() {
        if (this.#destroyed) {
            throw new Error('Queue already destroyed');
        }
        this.#logger.debug('Destroyed');
        this.#destroyed = true;
        this.#queue.length = 0;
        this[tick]();
        this.#tickQueue.length = 0;
    }
}